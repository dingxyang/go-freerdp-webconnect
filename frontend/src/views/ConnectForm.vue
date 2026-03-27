<template>
  <div class="connect-form">
    <h2>FreeRDP WebConnect</h2>
    <form @submit.prevent="handleConnect">
      <div class="form-group">
        <label for="host">主机地址</label>
        <input id="host" v-model="form.host" type="text" placeholder="192.168.1.100" required />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label for="port">端口</label>
          <input id="port" v-model.number="form.port" type="number" />
        </div>
        <div class="form-group">
          <label for="resolution">分辨率</label>
          <select id="resolution" v-model="form.resolution">
            <option value="1024x768">1024 x 768</option>
            <option value="1280x720">1280 x 720</option>
            <option value="1280x800">1280 x 800</option>
            <option value="1366x768">1366 x 768</option>
            <option value="1440x900">1440 x 900</option>
            <option value="1920x1080">1920 x 1080</option>
          </select>
        </div>
      </div>
      <div class="form-group">
        <label for="user">用户名</label>
        <input id="user" v-model="form.user" type="text" placeholder="administrator" />
      </div>
      <div class="form-group">
        <label for="pass">密码</label>
        <input id="pass" v-model="form.pass" type="password" />
      </div>
      <button type="submit" :disabled="connecting">
        {{ connecting ? '连接中...' : '连接' }}
      </button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
    <div class="version" v-if="version">
      App {{ version.app }} | FreeRDP {{ version.freerdp }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

const emit = defineEmits<{
  connect: [wsUrl: string, width: number, height: number]
}>()

const form = reactive({
  host: '',
  port: 3389,
  user: 'administrator',
  pass: '',
  resolution: '1024x768',
})

const connecting = ref(false)
const error = ref('')
const version = ref<{ app: string; freerdp: string } | null>(null)

onMounted(async () => {
  try {
    // @ts-ignore - Wails 运行时注入的全局对象
    version.value = await window.go.main.App.GetVersion()
  } catch {
    // 非 Wails 环境忽略
  }
})

async function handleConnect() {
  if (!form.host) {
    error.value = '请输入主机地址'
    return
  }

  connecting.value = true
  error.value = ''

  try {
    const [w, h] = form.resolution.split('x').map(Number)

    // @ts-ignore - Wails 运行时注入的全局对象
    const wsUrl: string = await window.go.main.App.Connect(
      form.host, form.user, form.pass, form.port, w, h
    )

    emit('connect', wsUrl, w, h)
  } catch (e: any) {
    error.value = e.message || '连接失败'
  } finally {
    connecting.value = false
  }
}
</script>

<style scoped>
.connect-form {
  max-width: 400px;
  margin: 80px auto;
  padding: 32px;
  background: #1a1a2e;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
}

h2 {
  text-align: center;
  color: #e0e0e0;
  margin-bottom: 24px;
  font-weight: 500;
}

.form-group {
  margin-bottom: 16px;
}

.form-row {
  display: flex;
  gap: 12px;
}
.form-row .form-group {
  flex: 1;
}

label {
  display: block;
  color: #aaa;
  font-size: 13px;
  margin-bottom: 4px;
}

input, select {
  width: 100%;
  padding: 10px 12px;
  background: #16213e;
  border: 1px solid #333;
  border-radius: 6px;
  color: #e0e0e0;
  font-size: 14px;
  box-sizing: border-box;
}
input:focus, select:focus {
  outline: none;
  border-color: #4a9eff;
}

button {
  width: 100%;
  padding: 12px;
  background: #4a9eff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  cursor: pointer;
  margin-top: 8px;
}
button:hover:not(:disabled) {
  background: #3a8eef;
}
button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error {
  color: #ff6b6b;
  font-size: 13px;
  margin-top: 8px;
  text-align: center;
}

.version {
  text-align: center;
  color: #666;
  font-size: 11px;
  margin-top: 20px;
}
</style>

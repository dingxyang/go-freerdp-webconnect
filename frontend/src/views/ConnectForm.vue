<template>
  <div class="connect-form">
    <h2>FreeRDP WebConnect</h2>
    <form @submit.prevent="handleConnect">
      <div class="tabs" role="tablist" aria-label="连接参数分类">
        <button
          type="button"
          role="tab"
          class="tab-btn"
          :class="{ active: activeTab === 'basic' }"
          :aria-selected="activeTab === 'basic'"
          @click="activeTab = 'basic'"
        >
          基础参数
        </button>
        <button
          type="button"
          role="tab"
          class="tab-btn"
          :class="{ active: activeTab === 'advanced' }"
          :aria-selected="activeTab === 'advanced'"
          @click="activeTab = 'advanced'"
        >
          高级参数
        </button>
      </div>

      <div v-show="activeTab === 'basic'" class="tab-panel" role="tabpanel">
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
      </div>

      <div v-show="activeTab === 'advanced'" class="tab-panel advanced-panel" role="tabpanel">
        <div class="advanced-grid">
          <div class="form-group">
            <label for="perf">性能</label>
            <select id="perf" v-model.number="form.perf">
              <option :value="0">局域网</option>
              <option :value="1">宽带</option>
              <option :value="2">调制解调器</option>
            </select>
          </div>
          <div class="form-group">
            <label for="fntlm">强制 NTLM 认证</label>
            <select id="fntlm" v-model.number="form.fntlm">
              <option :value="0">禁用</option>
              <option :value="1">NTLM v1</option>
              <option :value="2">NTLM v2</option>
            </select>
          </div>
        </div>
        <div class="check-grid">
          <label class="check-item"><input v-model="form.nowallp" type="checkbox" /> 禁用壁纸</label>
          <label class="check-item"><input v-model="form.nowdrag" type="checkbox" /> 禁用窗口全拖动</label>
          <label class="check-item"><input v-model="form.nomani" type="checkbox" /> 禁用菜单动画</label>
          <label class="check-item"><input v-model="form.notheme" type="checkbox" /> 禁用主题</label>
          <label class="check-item"><input v-model="form.nonla" type="checkbox" :disabled="form.notls" /> 禁用网络级别身份验证 (NLA)</label>
          <label class="check-item"><input v-model="form.notls" type="checkbox" /> 禁用 TLS</label>
        </div>
      </div>

      <button type="submit" class="submit-btn" :disabled="connecting">
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
import { ref, reactive, onMounted, watch } from 'vue'

const emit = defineEmits<{
  connect: [wsUrl: string, width: number, height: number]
}>()

const form = reactive({
  host: '',
  port: 3389,
  user: 'administrator',
  pass: '',
  resolution: '1024x768',
  perf: 0,
  fntlm: 0,
  nowallp: false,
  nowdrag: false,
  nomani: false,
  notheme: false,
  nonla: false,
  notls: false,
})

const connecting = ref(false)
const error = ref('')
const version = ref<{ app: string; freerdp: string } | null>(null)
const activeTab = ref<'basic' | 'advanced'>('basic')

onMounted(async () => {
  try {
    // @ts-ignore - Wails 运行时注入的全局对象
    version.value = await window.go.main.App.GetVersion()
  } catch {
    // 非 Wails 环境忽略
  }
})

watch(() => form.notls, (v) => {
  if (v) {
    form.nonla = true
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
      form.host, form.user, form.pass, form.port, w, h,
      form.perf, form.fntlm,
      form.nowallp, form.nowdrag, form.nomani, form.notheme, form.nonla, form.notls
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
  width: min(420px, calc(100% - 24px));
  margin: 24px auto;
  padding: 28px 28px 24px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(207, 10, 44, 0.12);
  border-radius: 18px;
  box-shadow: 0 24px 60px rgba(15, 23, 42, 0.12);
  backdrop-filter: blur(14px);
  max-height: calc(100vh - 24px);
  overflow-y: auto;
  scrollbar-gutter: stable;
}

.connect-form::-webkit-scrollbar {
  width: 8px;
}

.connect-form::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.6);
  border-radius: 999px;
}

form {
  padding-bottom: 2px;
}

h2 {
  text-align: center;
  color: #111827;
  margin-bottom: 24px;
  font-weight: 600;
  letter-spacing: 0.01em;
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

.tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-bottom: 14px;
}

.tab-btn {
  border: 1px solid #d6deea;
  border-radius: 10px;
  min-height: 40px;
  background: #f8fafc;
  color: #475569;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  border-color: #c2d2e5;
  background: #ffffff;
}

.tab-btn.active {
  color: #9f1239;
  border-color: rgba(207, 10, 44, 0.35);
  background: rgba(207, 10, 44, 0.08);
  box-shadow: 0 6px 14px rgba(207, 10, 44, 0.12);
}

.tab-panel {
  margin: 6px 0 14px;
}

.advanced-panel {
  padding: 10px 12px 12px;
  border: 1px solid #d9e2ee;
  border-radius: 12px;
  background: #f8fafc;
}

.advanced-grid {
  display: flex;
  gap: 12px;
  margin-top: 12px;
}

.advanced-grid .form-group {
  flex: 1;
  margin-bottom: 10px;
}

.check-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px 12px;
}

.check-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 13px;
  color: #475569;
}

.check-item input[type='checkbox'] {
  width: 16px;
  height: 16px;
}

label {
  display: block;
  color: #6b7280;
  font-size: 13px;
  margin-bottom: 6px;
}

input, select {
  width: 100%;
  padding: 10px 12px;
  min-height: 42px;
  background: #f8fafc;
  border: 1px solid #d6deea;
  border-radius: 10px;
  color: #111827;
  font-size: 14px;
  line-height: 1.4;
  box-sizing: border-box;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}
input:focus, select:focus {
  outline: none;
  border-color: #cf0a2c;
  box-shadow: 0 0 0 4px rgba(207, 10, 44, 0.12);
  background: #ffffff;
}

input::placeholder {
  color: #9ca3af;
}

select {
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  padding-right: 42px;
  background-image:
    linear-gradient(45deg, transparent 50%, #64748b 50%),
    linear-gradient(135deg, #64748b 50%, transparent 50%);
  background-position:
    calc(100% - 18px) calc(50% - 1px),
    calc(100% - 12px) calc(50% - 1px);
  background-size: 6px 6px, 6px 6px;
  background-repeat: no-repeat;
  cursor: pointer;
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #cf0a2c 0%, #f43f5e 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 15px;
  cursor: pointer;
  margin-top: 8px;
  box-shadow: 0 14px 28px rgba(207, 10, 44, 0.24);
  transition: transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
}
.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 18px 32px rgba(207, 10, 44, 0.28);
  filter: saturate(1.05);
}
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  box-shadow: none;
}

.error {
  color: #dc2626;
  font-size: 13px;
  margin-top: 8px;
  text-align: center;
}

.version {
  text-align: center;
  color: #94a3b8;
  font-size: 11px;
  margin-top: 20px;
}

@media (max-width: 520px) {
  .connect-form {
    width: calc(100% - 12px);
    margin: 6px auto;
    max-height: calc(100vh - 12px);
    padding: 20px;
    border-radius: 14px;
  }
  .form-row,
  .advanced-grid {
    flex-direction: column;
    gap: 8px;
  }
  .check-grid {
    grid-template-columns: 1fr;
  }
}
</style>

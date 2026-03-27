<template>
  <div class="rdp-desktop">
    <div class="toolbar">
      <span class="title">{{ host }} - 远程桌面</span>
      <button class="disconnect-btn" @click="handleDisconnect">断开连接</button>
    </div>
    <div class="canvas-wrapper">
      <canvas ref="canvasRef" :width="width" :height="height"></canvas>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RDPClient } from '../rdp'

const props = defineProps<{
  wsUrl: string
  width: number
  height: number
  host: string
}>()

const emit = defineEmits<{
  disconnected: []
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)
const client = new RDPClient()

client.on('disconnected', () => {
  emit('disconnected')
})

client.on('error', (msg) => {
  console.error('RDP Error:', msg)
  emit('disconnected')
})

onMounted(() => {
  if (canvasRef.value) {
    client.connect(props.wsUrl, canvasRef.value)
  }
})

onUnmounted(() => {
  client.disconnect()
})

function handleDisconnect() {
  client.disconnect()
}
</script>

<style scoped>
.rdp-desktop {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: linear-gradient(180deg, #f8fafc 0%, #eef2f7 100%);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.92);
  border-bottom: 1px solid #dbe2ea;
  box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
  backdrop-filter: blur(14px);
  flex-shrink: 0;
}

.title {
  color: #475569;
  font-size: 13px;
  font-weight: 500;
}

.disconnect-btn {
  padding: 6px 16px;
  background: #cf0a2c;
  color: white;
  border: none;
  border-radius: 999px;
  font-size: 13px;
  cursor: pointer;
  box-shadow: 0 10px 20px rgba(207, 10, 44, 0.2);
  transition: transform 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}
.disconnect-btn:hover {
  background: #b50c2b;
  transform: translateY(-1px);
  box-shadow: 0 14px 24px rgba(181, 12, 43, 0.24);
}

.canvas-wrapper {
  flex: 1;
  overflow: auto;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 24px;
}

canvas {
  display: block;
  cursor: default;
  border-radius: 14px;
  box-shadow: 0 24px 48px rgba(15, 23, 42, 0.14);
}
</style>

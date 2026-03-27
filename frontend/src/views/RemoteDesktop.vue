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
  background: #0a0a0a;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 16px;
  background: #1a1a2e;
  border-bottom: 1px solid #333;
  flex-shrink: 0;
}

.title {
  color: #aaa;
  font-size: 13px;
}

.disconnect-btn {
  padding: 4px 16px;
  background: #e74c3c;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
}
.disconnect-btn:hover {
  background: #c0392b;
}

.canvas-wrapper {
  flex: 1;
  overflow: auto;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 0;
}

canvas {
  display: block;
  cursor: default;
}
</style>

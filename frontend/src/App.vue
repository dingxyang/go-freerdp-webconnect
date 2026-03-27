<template>
  <RemoteDesktop
    v-if="connected"
    :ws-url="wsUrl"
    :width="rdpWidth"
    :height="rdpHeight"
    :host="rdpHost"
    @disconnected="onDisconnected"
  />
  <ConnectForm v-else @connect="onConnect" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ConnectForm from './views/ConnectForm.vue'
import RemoteDesktop from './views/RemoteDesktop.vue'

const connected = ref(false)
const wsUrl = ref('')
const rdpWidth = ref(1024)
const rdpHeight = ref(768)
const rdpHost = ref('')

function onConnect(url: string, width: number, height: number) {
  wsUrl.value = url
  rdpWidth.value = width
  rdpHeight.value = height
  // 从 URL 提取 host 参数用于标题显示
  try {
    const u = new URL(url)
    rdpHost.value = u.searchParams.get('host') || '远程桌面'
  } catch {
    rdpHost.value = '远程桌面'
  }
  connected.value = true
}

function onDisconnected() {
  connected.value = false
  wsUrl.value = ''
}
</script>

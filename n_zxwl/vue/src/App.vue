<template>
  <meta name="referrer" content="unsafe-url">
  <RouterView />
</template>
<script setup>
import { useUserStore } from '@/utils/auth.js';
import { initWebSocket, closeWebSocket } from '@/utils/wsUtil.js';
import {ref,onMounted,onUnmounted,computed,onBeforeUnmount} from "vue";
import {preloadStatus} from '@/router/index.js'

const {getUser,deleteUser} = useUserStore()

const preloadStats = computed(() => ({
  success: preloadStatus.success.size,
  failed: preloadStatus.failed.size,
  pending: preloadStatus.pending.size
}))

onMounted(() => {
  // 应用启动时，如果用户已登录，自动连接 WebSocket
  const userLocal = getUser()
  if (userLocal && (userLocal.id !== null || userLocal.email !== ''))
    initWebSocket(userLocal.id,userLocal.email);
})
onUnmounted(() => closeWebSocket()) // 应用卸载时关闭 WebSocket
onBeforeUnmount(() => deleteUser()) // 删除用户信息
</script>

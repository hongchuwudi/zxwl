<template>
  <div class="chat-container">
    <!-- 顶部标题栏 -->
    <div class="chat-header">
      <el-button type="text" @click="goBack" class="back-btn">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <el-image class="header-logo" :src="logo" fit="cover" />
      <h2 class="header-title">{{ name }}</h2>
      <div class="header-status">
        <el-badge v-if="unreadCount > 0" :value="unreadCount" :max="99" />
      </div>
    </div>

    <!-- 消息列表 -->
    <div ref="messageContainer" class="message-list">
      <div
          v-for="(message, index) in messages"
          :key="index"
          class="message-item"
          :class="{ 'self': isSelf(message.email) }"
      >
        <el-avatar class="message-avatar" :src="message.picture || defaultAvatar" />
        <div class="message-content">
          <div class="message-info">
            <span class="message-name">{{ message.name }}</span>
            <span class="message-time">{{ formatTime(message.created_at) }}</span>
          </div>
          <div class="message-bubble">
            {{ message.content }}
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <el-input
          v-model="inputMessage"
          placeholder="输入消息..."
          @keyup.enter="sendMessage"
          class="message-input"
          resize="none"
          type="textarea"
          :autosize="{ minRows: 1, maxRows: 4 }"
      />
      <el-button
          type="primary"
          @click="sendMessage"
          class="send-btn"
          :disabled="!inputMessage.trim()"
      >
        发送
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
//const defaultAvatar = ref('https://example.com/default-avatar.png')

// 路由参数
const schoolId = parseInt(route.params.id)
const name = route.query.name
const logo = route.query.logo

// 响应式数据
const messages = ref([])
const inputMessage = ref('')
const messageContainer = ref(null)
const unreadCount = ref(0)
const ws = ref(null)
const capi = axios.create({
  baseURL: 'http://localhost:8792', // 根据实际情况修改
});

// 初始化WebSocket
const initWebSocket = () => {
  const userEmail = localStorage.getItem('userEmail')
  if (!userEmail) {
    ElMessage.error('请先登录')
    router.push('/login')
    return
  }

  ws.value = new WebSocket(`ws://localhost:8792/ws`)

  ws.value.onopen = () => {
    console.log('WebSocket连接已建立');
    // 加入聊天室
    ws.value.send(JSON.stringify({
      action: 'join',
      schoolID: schoolId,
      email: userEmail
    }))
  }

  ws.value.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      console.log("收到消息:", data)

      if (data.action === 'history') {
        // 处理null和无效数据
        const rawMessages = data.messages || []
        messages.value = rawMessages
            .filter(m => m && m.content) // 过滤无效消息
            .map(m => ({
              school_id: schoolId,
              email: m.email || '',
              content: m.content,
              name: m.name || '匿名用户',
              picture: m.picture || defaultAvatar.value,
              created_at: m.created_at ? new Date(m.created_at) : new Date()
            }))
            .reverse()
        scrollToBottom()
      } else if (data.action === 'message') {
        handleNewMessage(data)
      }
    } catch (error) {
      console.error("消息处理错误:", error)
    }
  }

  const handleNewMessage = (data) => {
    const newMsg = {
      school_id: data.schoolID,
      email: data.email,
      content: data.content,
      name: data.name || '匿名用户',
      picture: data.picture || defaultAvatar.value,
      created_at: data.time ? new Date(data.time) : new Date()
    }

    // 去重逻辑
    const existingIndex = messages.value.findIndex(m =>
        m.content === newMsg.content &&
        m.email === newMsg.email &&
        Math.abs(m.created_at - newMsg.created_at) < 5000
    )

    if (existingIndex > -1) {
      messages.value[existingIndex] = newMsg
    } else {
      messages.value.push(newMsg)
    }
    scrollToBottom()
  }

  ws.value.onerror = (error) => {
    console.error('WebSocket错误:', error);
    ElMessage.error('连接异常，请刷新页面');
  }


  ws.value.onclose = () => {
    console.log('WebSocket连接关闭')
  }
}

// 发送消息
const sendMessage = async () => {
  if (!inputMessage.value.trim()) return
  const userEmail = localStorage.getItem('userEmail')

  // 临时消息（带发送状态）
  const tempMessage = {
    school_id: schoolId,
    email: userEmail,
    content: inputMessage.value.trim(),
    created_at: new Date(),
    name: '我',
    picture: defaultAvatar.value
  }

  messages.value.push(tempMessage)
  scrollToBottom()

  try {
    await capi.post('/chat/send', {
      school_id: schoolId,
      email: userEmail,
      content: inputMessage.value.trim()
    })
    inputMessage.value = ''
  } catch (error) {
    ElMessage.error('发送失败')
    // 标记失败状态
    const index = messages.value.findIndex(m => m === tempMessage)
    if (index > -1) {
      messages.value[index].name = '发送失败'
    }
  }
}

// 检查是否是自己发送的消息
const isSelf = (email) => {
  return email === localStorage.getItem('userEmail')
}


// 格式化时间
const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

// 返回按钮
const goBack = () => {
  router.go(-1)
}

// 标记已读
const readHandler = async () => {
  try {
    // 使用 axios 实例 capi 发送请求
    await capi.post('/chat/mark-read', {
      email: localStorage.getItem('userEmail'),
      school_id: schoolId
    })
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

// 组件挂载时
onMounted(async () => {
  readHandler()
  initWebSocket()
})
onMounted(async () => {
  const logData = {
    "email": localStorage.getItem('userEmail'),
    "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
    "operation": "用户进入聊天室聊天"
  };
  const logResponse = await axios.post("gapi/log", logData, {
    headers: {
      "Content-Type": "application/json"
    }
  });
})
// 组件卸载时
onUnmounted(() => {
  if (ws.value) {
    ws.value.close()
  }
})
</script>

<style lang="scss" scoped>
.chat-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f6fa;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 16px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);

  .back-btn {
    margin-right: 12px;
    font-size: 20px;
  }

  .header-logo {
    width: 40px;
    height: 40px;
    border-radius: 8px;
    margin-right: 12px;
  }

  .header-title {
    margin: 0;
    font-size: 18px;
    flex-grow: 1;
  }
}

.message-list {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background: #f5f6fa;
}

.message-item {
  display: flex;
  margin-bottom: 20px;

  &.self {
    flex-direction: row-reverse;

    .message-content {
      align-items: flex-end;
    }

    .message-bubble {
      background: #409eff;
      color: white;
      border-radius: 12px 12px 0 12px;
    }
  }
}

.message-avatar {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
}

.message-content {
  margin: 0 12px;
  max-width: 70%;
  display: flex;
  flex-direction: column;
}

.message-info {
  margin-bottom: 4px;
  display: flex;
  align-items: center;
}

.message-name {
  font-size: 12px;
  color: #606266;
  margin-right: 8px;
}

.message-time {
  font-size: 12px;
  color: #909399;
}

.message-bubble {
  padding: 12px;
  background: white;
  border-radius: 12px 12px 12px 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  line-height: 1.5;
  word-break: break-word;
}

.input-area {
  display: flex;
  padding: 16px;
  background: white;
  border-top: 1px solid #ebeef5;

  .message-input {
    flex-grow: 1;
    margin-right: 12px;

    ::v-deep(.el-textarea__inner) {
      border-radius: 20px;
      padding: 8px 16px;
    }
  }

  .send-btn {
    border-radius: 20px;
    padding: 8px 24px;
  }
}

@media (max-width: 768px) {
  .message-content {
    max-width: 85%;
  }

  .message-name {
    font-size: 10px;
  }

  .message-time {
    font-size: 10px;
  }

  .message-bubble {
    padding: 8px 12px;
    font-size: 14px;
  }
}
</style>
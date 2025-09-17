<template>
  <div class="qq9pc-container">
    <a-page-header
        title="家庭共享"
        @back="() => $router.go(-1)"
        class="page-header"
    />

    <div class="chat-layout">
      <!-- 左侧好友栏 -->
      <div class="friend-sidebar">
        <div class="friend-search">
          <a-input-search
              v-model:value="searchKeyword"
              placeholder="搜索好友"
          />
          <a-button type="primary" @click="showAddFriendModal" class="add-friend-btn">
            <user-add-outlined />
            添加
          </a-button>
        </div>

        <a-tabs v-model:activeKey="activeFriendTab" class="friend-tabs">
          <a-tab-pane key="friends" tab="家庭列表">
            <div class="friend-list" >
              <div
                  v-for="friend in filteredFriends"
                  :key="friend.id"
                  :class="['friend-item', { active: activeFriend?.id === friend.id }]"
                  @click="selectFriend(friend)"
              >
                <a-avatar :size="40" :src="friend.user_a.avatarUrl" class="friend-avatar">
                  {{ friend.user_a.username?.charAt(0) || friend.user_a.username?.charAt(0) }}
                </a-avatar>
                <div class="friend-info">
                  <div class="friend-name"> {{friend.nickname || friend.user_a.displayName ||  friend.user_a.username}}</div>
                  <div class="friend-status">
                    <a-tag :color="friend.user_a.isOnline ? 'green' : 'default'" size="small">
                      {{ friend.user_a.isOnline ? '在线' : '离线' }}
                    </a-tag>
                    <span v-if="getUnreadCount(friend.id) > 0" class="unread-badge">
                      {{ getUnreadCount(friend.id) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </a-tab-pane>

          <a-tab-pane key="requests" >
            <template #tab>
      <span class="request-tab">
        好友请求
        <span v-if="unreadRequestCount > 0" class="unread-request-badge">
          {{ unreadRequestCount }}
        </span>
      </span>
            </template>
            <div class="request-list">
              <div v-for="request in pendingRequests" :key="request.id" class="request-item">
                <a-avatar :size="40" :src="request.from_user.avatarUrl">
                  {{ request.from_user.displayName?.charAt(0) || request.from_user.username?.charAt(0) }}
                </a-avatar>
                <div class="request-info">
                  <div class="request-name">{{ request.from_user.displayName || request.from_user.username }}</div>
                  <div class="request-message" v-if="request.request_message">
                    {{ request.request_message }}
                  </div>
                  <div class="request-actions">
                    <a-button size="small" type="primary" @click="acceptFriendRequest(request)">
                      接受
                    </a-button>
                    <a-button size="small" @click="rejectFriendRequest(request)">
                      拒绝
                    </a-button>
                  </div>
                </div>
              </div>
            </div>
          </a-tab-pane>
        </a-tabs>
      </div>

      <!-- 右侧聊天内容 -->
      <div class="chat-content">
        <div v-if="activeFriend" class="chat-header">
          <a-avatar :size="40" :src="activeFriend.user_a.avatarUrl" class="chat-avatar">
            {{ activeFriend.user_a.displayName?.charAt(0) || activeFriend.user_a.username?.charAt(0) }}
          </a-avatar>
          <div class="chat-user-info">
            <div class="chat-user-name">{{ activeFriend.user_a.displayName || activeFriend.user_a.username }}</div>
            <div class="chat-user-status">
              <a-tag :color="activeFriend.user_a.isOnline ? 'green' : 'default'" size="small">
                {{ activeFriend.user_a.isOnline ? '在线' : '离线' }}
              </a-tag>
            </div>
          </div>
          <div class="chat-actions">
            <a-button type="link" @click="deleteFriend(activeFriend.user_a)" danger>
              <delete-outlined />
            </a-button>
          </div>
        </div>

        <div v-else class="chat-placeholder">
          <a-empty description="请选择好友开始聊天" />
        </div>

        <div v-if="activeFriend" class="chat-messages">
          <div class="messages-container" ref="messagesContainer">
            <div
                v-for="message in chatMessages"
                :key="message.id"
                :class="['message-item', { own: message.sender_id === currentUser.id }]"
            >
              <!-- 自己发送的消息显示在右侧 -->
              <div v-if="message.sender_id === currentUser.id" class="message-row own-message">
                <div class="message-time">{{ formatMessageTime(message.created_at) }}</div>
                <div class="message-bubble">
                  <div class="message-text">{{ message.content }}</div>
                </div>
                <a-avatar
                    :size="32"
                    :src="currentUser.avatarUrl"
                    class="message-avatar"
                />
              </div>

              <!-- 对方发送的消息显示在左侧 -->
              <div v-else class="message-row other-message">
                <a-avatar
                    :size="32"
                    :src="activeFriend.user_a.avatarUrl"
                    class="message-avatar"
                />
                <div class="message-bubble">
                  <div class="message-text">{{ message.content }}</div>
                </div>
                <div class="message-time">{{ formatMessageTime(message.created_at) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 打字状态指示器 -->
        <div v-if="isFriendTyping " class="typing-indicator">
          <div class="typing-dots">
            <span></span>
            <span></span>
            <span></span>
          </div>
          <span class="typing-text">
    {{ (activeFriend.user_a.displayName || activeFriend.user_a.username) }} 正在输入...
  </span>
        </div>

        <div v-if="activeFriend" class="chat-input">
          <a-input
              v-model:value="messageText"
              placeholder="输入消息..."
              @pressEnter="sendMessage"
              :disabled="sending"
          >
            <template #suffix>
              <a-button
                  type="primary"
                  @click="sendMessage"
                  :loading="sending"
                  :disabled="!messageText.trim()"
              >
                发送
              </a-button>
            </template>
          </a-input>
        </div>
      </div>
    </div>

    <!-- 添加好友模态框 -->
    <a-modal
        v-model:open="addFriendModalVisible"
        title="添加好友"
        @ok="handleAddFriend"
        :confirm-loading="addingFriend"
    >
      <a-form layout="vertical">
        <a-form-item label="家人ID">
          <a-input
              v-model:value="addFriendId"
              placeholder="请输入家人ID"
          />
        </a-form-item>
        <a-form-item label="家人称呼">
          <a-input
              v-model:value="addSalutation"
              placeholder="请输入家人称呼"
          />
        </a-form-item>
        <a-form-item label="打招呼信息">
          <a-textarea
              v-model:value="addFriendMessage"
              placeholder="请输入打招呼信息"
              :rows="3"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted,onUnmounted,  computed, nextTick, watch } from 'vue'
import { message } from 'ant-design-vue'
import { UserAddOutlined, DeleteOutlined } from '@ant-design/icons-vue'
import { useUserStore } from "@/utils/auth.js"
import axios from "axios"
import { getWebSocket } from '@/utils/wsUtil.js' // 导入WebSocket工具

const { getUser } = useUserStore()
const currentUser = getUser()

// 添加WebSocket相关状态
const websocket = ref(null)     // ws连接线
const onlineStatusMap = ref({}) // 存储好友在线状态 { friendId: boolean }
const typingTimeout = ref({}) // 存储清除打字状态的定时器

// 提示好友正在输入的状态
const typingCheckInterval = ref(null)
let isCurrentlyTyping = false

// 好友相关状态
const friends = ref([])
const activeFriend = ref(null)
const searchKeyword = ref('')
const activeFriendTab = ref('friends')
const addFriendModalVisible = ref(false)
const addFriendId = ref('')
const addFriendMessage = ref('')
const addSalutation = ref('')
const addingFriend = ref(false)

// 好友请求
const pendingRequests = ref([])
const unreadRequestCount = ref(0)// 添加未读请求计数状态

// 聊天相关状态
const chatMessages = ref([])
const messageText = ref('')
const sending = ref(false)
const messagesContainer = ref(null)
let typingTimer = null
const messageSound = new Audio('/iphone发送信息声.wav')
const messageGet = new Audio('/iphone接收信息声.wav')

const typingStatus = ref({})
const typingTimeouts = ref({})
// 处理接收到的打字状态消息
const handleFriendTypingStatus = (event) => {
  try {
    const { from_user_id, is_typing } = event.detail
    console.log('收到打字状态:', from_user_id, is_typing)

    if (!from_user_id) return

    const friendId = parseInt(from_user_id)
    if (isNaN(friendId)) return

    // 清除之前的定时器（如果有）
    if (typingTimeouts.value[friendId]) clearTimeout(typingTimeouts.value[friendId])

    // 更新打字状态
    typingStatus.value[friendId] = is_typing

    // 如果不是正在打字，设置定时器5秒后清除状态（防止闪烁）
    if (!is_typing) {
      typingTimeouts.value[friendId] = setTimeout(() => {
        typingStatus.value[friendId] = false
      }, 5000) // 5秒后清除打字状态
    }
  } catch (error) {
    console.error('处理好友打字状态失败:', error)
  }
}
// 计算属性：判断当前好友是否正在输入
const isFriendTyping = computed(() => {
  if (!activeFriend.value || !activeFriend.value.user_a || !activeFriend.value.user_a.id) {
    return false
  }

  const friendId = activeFriend.value.user_a.id
  return typingStatus.value[friendId] === true
})
// 发送自己的打字状态
const sendMyTypingStatus = (isTyping) => {
  if (!websocket.value || !activeFriend.value) return

  try {
    const wsMessage = {
      type: 'typing',
      data: {
        to: activeFriend.value.user_a.id,
        is_typing: isTyping
      }
    }
    websocket.value.send(wsMessage)
  } catch (error) {
    console.error('发送打字状态失败:', error)
  }
}
// 监听输入内容变化(超绝防抖)
watch(messageText, (newValue) => {
  if (!activeFriend.value) return
  if (typingTimer) clearTimeout(typingTimer)   // 清除之前的定时器

  const isTyping = newValue.trim().length > 0

  if (isTyping && !isCurrentlyTyping) {
    sendMyTypingStatus(true)     // 开始输入
    isCurrentlyTyping = true
  } else if (!isTyping && isCurrentlyTyping) {
    // 停止输入（防抖：延迟500ms确认真的停止了）
    typingTimer = setTimeout(() => {
      if (messageText.value.trim().length === 0) {
        sendMyTypingStatus(false)
        isCurrentlyTyping = false
      }
    }, 500)
  }
})

/** 好友消息 **/
// sendMessage函数，在HTTP发送成功后添加WebSocket通知
const sendMessage = async () => {
  if (!messageText.value.trim() || !activeFriend.value) return

  try {
    sending.value = true

    // 发送停止输入状态
    // 停止输入状态
    if (typingTimer) {
      clearTimeout(typingTimer)
      sendMyTypingStatus(false)
    }

    // 先在前端添加消息（乐观更新）
    const tempMessage = {
      id: Date.now(), // 临时ID
      sender_id: currentUser.id,
      content: messageText.value.trim(),
      message_type: 'text',
      created_at: new Date().toISOString()
    }

    chatMessages.value.push(tempMessage)
    scrollToBottom()
    messageSound.play().catch(() => {})

    // 使用HTTP API发送消息
    const response = await axios.post(`/gapi/user/${currentUser.id}/chat/friend/send`, {
      receiver_id: activeFriend.value.user_a.id,
      message_type: 'text',
      content: messageText.value.trim()
    })

    if (response.data.code === 200) {
      // HTTP发送成功，通过WebSocket通知后端接收消息
      if (websocket.value) {
        // 构建与后端handleChatMessage期望格式匹配的消息
        const wsMessage = {
          type: 'chat_message',
          data: { // 直接使用对象，不要JSON.stringify
            receiver_id: activeFriend.value.user_a.id,
            content: messageText.value.trim(),
            message_type: 'text'
          }
        }

        // 发送WebSocket消息
        websocket.value.send(wsMessage)
      }

      messageText.value = ''
    }
  } catch (error) {
    console.error('发送消息失败:', error)
    message.error('发送消息失败')

    // 失败时移除临时消息
    const index = chatMessages.value.findIndex(msg => msg.id === tempMessage.id)
    if (index !== -1) {
      chatMessages.value.splice(index, 1)
    }
  } finally {
    sending.value = false
  }
}
// 处理新消息
const handleNewMessage = async (event) => {
  scrollToBottom()
  console.log('handleNewMessage:', event.detail)
  const { from_user_id, content, message_type, created_at } = event.detail

  // 如果当前正在和这个好友聊天，直接添加到消息列表
  if (activeFriend.value && activeFriend.value.user_a.id === from_user_id) {
    chatMessages.value.push({
      id: Date.now(), // 临时ID
      sender_id: from_user_id,
      content: content,
      message_type: message_type,
      created_at: created_at
    })


    // 播放消息提示音
    messageGet.play().catch(() => {})
  } else {
    // 如果不是当前聊天好友，显示通知
    const friend = friends.value.find(f => f.user_a.id === from_user_id)
    if (friend) {
      showMessageNotification(friend.user_a.displayName || friend.user_a.username, content)

      // 更新未读消息计数（如果需要）
      updateUnreadCount(from_user_id)
    }
  }
}
// 获取未读消息数量
const getUnreadCount = (friendId) => {
  // 这里可以根据需要实现未读消息计数
  return 0
}


// 好友搜索计算
const filteredFriends = computed(() => {
  if (!searchKeyword.value) return friends.value
  return friends.value.filter(friend =>
      friend.user_a.displayName?.includes(searchKeyword.value) ||
      friend.user_a.username?.includes(searchKeyword.value)
  )
})
// 获取好友列表
const fetchFriends = async () => {
  try {
    const response = await axios.get(`/gapi/user/${currentUser.id}/friends`)
    if (response.data.code === 200) {
      friends.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取好友列表失败:', error)
    message.error('获取好友列表失败')
  }
}
// 获取待处理的好友请求
const fetchPendingRequests = async () => {
  try {
    const response = await axios.get(`/gapi/user/${currentUser.id}/friend-requests`)
    if (response.data.code === 200) {
      pendingRequests.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取好友请求失败:', error)
    message.error('获取好友请求失败')
  }
}
// 获取聊天记录
const fetchChatHistory = async (friendId) => {
  try {
    const response = await axios.get(`/gapi/user/${currentUser.id}/chat/friend/${friendId}/history?limit=100`)
    if (response.data.code === 200) {
      chatMessages.value = response.data.data || []
      scrollToBottom()

      // 标记消息为已读
      await axios.post(`/gapi/user/${currentUser.id}/chat/mark-read/${friendId}`)
    }
  } catch (error) {
    console.error('获取聊天记录失败:', error)
  }
}
// 选择好友
const selectFriend = (friend) => {
  activeFriend.value = friend
  fetchChatHistory(friend.user_a.id)
}
// 接受好友请求
const acceptFriendRequest = async (request) => {
  try {
    const response = await axios.post(`/gapi/user/${currentUser.id}/friend/accept/${request.id}`)
    if (response.data.code === 200) {
      message.success('好友请求已接受')
      await fetchPendingRequests()
      await fetchFriends()
    }
  } catch (error) {
    console.error('接受好友请求失败:', error)
    message.error('接受好友请求失败')
  }
}
// 拒绝好友请求
const rejectFriendRequest = async (request) => {
  try {
    // router.HandleFunc("/user/{userID}/friend/reject/{requestID}", userFriendsHandler.RejectFriendRequest).Methods("POST") // 拒绝好友请求
    const response = await axios.post(`/gapi/user/${currentUser.id}/friend/reject/${request.id}`)
    if (response.data.code === 200) {
      message.success('好友请求已拒绝')
      await fetchPendingRequests()
    }
  } catch (error) {
    console.error('拒绝好友请求失败:', error)
    message.error('拒绝好友请求失败')
  }
}
// 处理新好友请求
const handleNewRequest = (event) => {
  try {
    console.log('收到新好友请求:', event.detail)

    const requestData = event.detail
    // 将新请求添加到待处理请求列表
    pendingRequests.value.unshift(requestData)

    // 如果当前不在"好友请求"标签页，增加未读计数
    if (activeFriendTab.value !== 'requests') {
      unreadRequestCount.value++
      showNewRequestNotification(requestData)
    }

    // 播放提示音
    messageGet.play().catch(() => {})

  } catch (error) {
    console.error('处理新好友请求失败:', error)
  }
}
// 显示新请求通知
const showNewRequestNotification = (requestData) => {
  const userName = requestData.from_user?.displayName || requestData.from_user?.username || '未知用户'
  const messageText = requestData.request_message || '想添加您为好友'

  // 使用Ant Design的message组件
  message.info(`新的好友请求: ${userName}`, 5)

  // 浏览器通知
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification('新的好友请求', {
      body: `${userName}: ${messageText}`,
      icon: '/favicon.ico'
    })
  }
}
// 监听标签页切换，当切换到请求页面时清除未读计数
watch(activeFriendTab, (newTab) => {
  if (newTab === 'requests') unreadRequestCount.value = 0
})

// 添加好友
const showAddFriendModal = () => addFriendModalVisible.value = true
const handleAddFriend = async () => {
  if (!addFriendId.value) {
    message.error('请输入好友ID')
    return
  }
  try {
    addingFriend.value = true
    const response = await axios.post(`/gapi/user/${currentUser.id}/friend/request/${parseInt(addFriendId.value)}`, {
      friend_id: parseInt(addFriendId.value),
      request_message: addFriendMessage.value,
      salutation: addSalutation.value
    })

    if (response.data.code === 200) {
      message.success('好友请求已发送')
      addFriendModalVisible.value = false
      addFriendId.value = ''
      addFriendMessage.value = ''

      // 获取返回的请求ID
      const requestId = response.data?.id // 根据你的实际返回结构调整
      console.log('Request ID:', response.data?.id)
      // 通过WebSocket发送消息给后端进行转发（使用与聊天消息相同的格式）
      if (requestId && websocket.value) {
        const wsMessage = {
          type: 'friend_request', // 消息类型
          data: { // 直接使用对象，不要JSON.stringify
            id: requestId // 发送请求ID给后端
          }
        }

        // 发送WebSocket消息（与聊天消息发送方式一致）
        websocket.value.send(wsMessage)
      }

      await fetchPendingRequests()
    }
  } catch (error) {
    console.error('添加好友失败:', error)
    message.error(error.response?.data?.message || '添加好友失败')
  } finally {
    addingFriend.value = false
  }
}

// 删除好友
const deleteFriend = async (friend) => {
  try {
// router.HandleFunc("/user/{userID}/friend/{friendID}", userFriendsHandler.DeleteFriend).Methods("DELETE")//删除好友
    const response = await axios.delete(`/gapi/user/${currentUser.id}/friend/${friend.id}`)
    if (response.data.code === 200) {
      message.success('好友已删除')
      activeFriend.value = null
      chatMessages.value = []
      await fetchFriends()
    }
  } catch (error) {
    console.error('删除好友失败:', error)
    message.error('删除好友失败')
  }
}

// 处理好友状态消息
const handleFriendStatus = (event) => {
  const { user_id, is_online } = event.detail
  updateFriendOnlineStatus(user_id, is_online)
}
// 处理好友上线消息
const handleFriendOnline = (event) => {
  const { user_id, user_name } = event.detail
  updateFriendOnlineStatus(user_id, true)
  showStatusNotification(user_name, true)  // 显示上线通知
}
// 处理好友下线消息
const handleFriendOffline = (event) => {
  const { user_id, user_name } = event.detail
  updateFriendOnlineStatus(user_id, false)
  showStatusNotification(user_name, false)  // 显示下线通知
}
// 更新好友在线状态
const updateFriendOnlineStatus = (friendId, isOnline) => {
  onlineStatusMap[friendId] = isOnline
  // 更新好友列表中的在线状态
  const friendIndex = friends.value.findIndex(f => f.user_a_id === Number(friendId))
  if (friendIndex !== -1) {
    console.log('更新好友在线状态:', friendId, isOnline)
    friends.value[friendIndex].user_a.isOnline = isOnline
  }
}

// 显示状态通知
const showStatusNotification = (friendName, isOnline) => {
  const statusText = isOnline ? '上线了' : '下线了'
  const notificationText = `${friendName} ${statusText}`

  // 使用Ant Design的message组件
  message.info(notificationText, 3)

  // 或者使用浏览器通知
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification('好友状态更新', {
      body: notificationText,
      icon: '/favicon.ico'
    })
  }
}
// 显示消息通知
const showMessageNotification = (friendName, messageContent) => {
  const truncatedMessage = messageContent.length > 50
      ? messageContent.substring(0, 50) + '...'
      : messageContent

  // 使用Ant Design的message组件
  message.info(`新消息来自 ${friendName}: ${truncatedMessage}`, 5)

  // 浏览器通知
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification(`新消息来自 ${friendName}`, {
      body: truncatedMessage,
      icon: '/favicon.ico'
    })
  }
}

// 更新未读消息计数（需要你实现具体的计数逻辑）
const updateUnreadCount = (friendId) => {
  // 这里可以实现未读消息计数的逻辑
  console.log(`收到来自 ${friendId} 的新消息，更新未读计数`)
}

/** 工具函数 **/
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      const container = messagesContainer.value
      container.scrollTop = container.scrollHeight
      console.log('滚动到底部，容器高度:', container.scrollHeight, '滚动位置:', container.scrollTop)
    } else {
      console.warn('messagesContainer 未找到')
    }
  })
}
const formatMessageTime = (time) => {
  const now = new Date();
  const messageDate = new Date(time);

  // 重置时间为0点，以便比较日期
  const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const yesterdayStart = new Date(todayStart);
  yesterdayStart.setDate(yesterdayStart.getDate() - 1);

  // 计算消息日期0点的时间
  const messageDayStart = new Date(messageDate.getFullYear(), messageDate.getMonth(), messageDate.getDate());

  // 计算天数差（基于日期，而不是24小时）
  const diffDays = Math.floor((todayStart - messageDayStart) / (1000 * 60 * 60 * 24));

  if (messageDate >= todayStart) {
    // 今天
    return messageDate.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    });
  } else if (messageDate >= yesterdayStart) {
    // 昨天
    return `昨天 ${messageDate.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })}`;
  } else if (diffDays <= 7) {
    // 一周内
    return `${diffDays}天前 ${messageDate.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })}`;
  } else {
    // 超过一周
    return messageDate.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    }).replace(/\//g, '-');
  }
};
/** ws相关**/

// 设置WebSocket消息监听
const setupWebSocketListeners = () => {
  // 获取WebSocket实例
  websocket.value = getWebSocket()

  if (!websocket.value) {
    console.warn('WebSocket连接不存在')
    return
  }

  // 监听好友状态变化
  window.addEventListener('friendStatus', handleFriendStatus)
  window.addEventListener('friendOnline', handleFriendOnline)
  window.addEventListener('friendOffline', handleFriendOffline)
  window.addEventListener('newMessage', handleNewMessage)
  window.addEventListener('typingStatus', handleFriendTypingStatus) // 监听打字状态
  window.addEventListener('newRequest', handleNewRequest) // 监听新请求
}
// 移除WebSocket监听
const removeWebSocketListeners = () => {
  window.removeEventListener('friendStatus', handleFriendStatus)
  window.removeEventListener('friendOnline', handleFriendOnline)
  window.removeEventListener('friendOffline', handleFriendOffline)
  window.removeEventListener('newMessage', handleNewMessage)
  window.removeEventListener('typingStatus', handleFriendTypingStatus)
  window.removeEventListener('newRequest', handleNewRequest)
}

/** 钩子函数 **/
// 初始化
onMounted(() => {
  fetchFriends()
  fetchPendingRequests()
  setupWebSocketListeners()
})
onUnmounted(() => {
  if (typingCheckInterval.value) clearInterval(typingCheckInterval.value) // 清理定时器

  // 清理所有定时器
  Object.values(typingTimeout.value).forEach(timeout => {
    clearTimeout(timeout)
  })
  // 移除WebSocket监听
  removeWebSocketListeners()
})
// 监听消息发送，自动滚动到底部
watch(chatMessages, scrollToBottom, { deep: true })
</script>

<style scoped lang="scss">
.qq9pc-container {
  //padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}
.page-header {
  background-color: #fff;
  //margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.chat-layout {
  display: flex;
  height: calc(100vh - 120px);
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}
.friend-sidebar {
  width: 300px;
  border-right: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
}
.friend-search {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  gap: 8px;
}
.add-friend-btn {
  flex-shrink: 0;
}
::v-deep(.friend-tabs) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .ant-tabs-nav {
    margin: 0;
    padding: 0 16px;
    background: #fff;
    border-bottom: 1px solid #f0f0f0;
    flex-shrink: 0;
  }

  .ant-tabs-tab {
    padding: 12px 0;
    margin-right: 24px;
    font-weight: 500;
  }

  .ant-tabs-ink-bar {
    background: #1890ff;
  }

  .ant-tabs-content {
    flex: 1;
    overflow: hidden;
  }

  .ant-tabs-content-holder {
    flex: 1;
    overflow: hidden;
  }

  .ant-tabs-tabpane {
    padding: 0 !important;
    height: 100%;
    overflow: hidden;
  }
}
.friend-list, .request-list {
  height: 100%;
  overflow-y: auto;
  max-height: calc(100vh - 250px);
  padding: 0;
}
.friend-item, .request-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.3s ease;
  width: 100%;
  box-sizing: border-box;
}
.friend-item:hover, .request-item:hover {
  background-color: #f9f9f9;
}
.friend-item.active {
  background-color: #e6f7ff;
  border-left: 3px solid #1890ff;
}
.friend-avatar {
  margin-right: 12px;
  flex-shrink: 0;
  border: 2px solid #f0f0f0;
  transition: border-color 0.3s ease;
}
.friend-item:hover .friend-avatar,
.friend-item.active .friend-avatar {
  border-color: #1890ff;
}
.friend-info, .request-info {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}
.friend-name {
  font-weight: 500;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #262626;
  font-size: 14px;
}
.friend-status {
  display: flex;
  align-items: center;
  gap: 8px;
}
::v-deep(.status-tag) {
  margin: 0;
  font-size: 11px;
  padding: 0 6px;
  height: 20px;
  line-height: 20px;
}

.unread-badge {
  background-color: #ff4d4f;
  color: white;
  border-radius: 10px;
  padding: 2px 6px;
  font-size: 11px;
  min-width: 18px;
  height: 18px;
  text-align: center;
  line-height: 14px;
  font-weight: bold;
}

/* 未读请求小红点样式 */
.unread-request-badge {
  background-color: #ff4d4f;
  color: white;
  border-radius: 50%;
  min-width: 18px;
  height: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  margin-left: 8px;
  padding: 0 4px;
}
.request-tab {
  display: flex;
  align-items: center;
}
.request-item {
  background-color: #fafafa;
  border-left: 3px solid #ffa940;
}
.request-item:hover {
  background-color: #fff7e6;
}
.request-name {
  font-weight: 500;
  margin-bottom: 4px;
  color: #262626;
  font-size: 14px;
}
.request-message {
  color: #8c8c8c;
  font-size: 12px;
  margin-bottom: 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.request-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}
::v-deep(.request-actions .ant-btn) {
  height: 24px;
  font-size: 12px;
  padding: 0 8px;
}

.chat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.chat-header {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  background-color: #fafafa;
}
.chat-avatar {
  margin-right: 12px;
}
.chat-user-info {
  flex: 1;
}
.chat-user-name {
  font-weight: 500;
  margin-bottom: 4px;
}
.chat-actions {
  margin-left: auto;
}
.chat-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}
.chat-input {

  padding: 16px;
  border-top: 1px solid #f0f0f0;
  background-color: #fafafa;

}

/* 空状态样式 */
.tab-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #8c8c8c;
}

.tab-empty .anticon {
  font-size: 48px;
  margin-bottom: 16px;
  color: #d9d9d9;
}

/* 滚动条样式 */
.friend-list::-webkit-scrollbar,
.request-list::-webkit-scrollbar,
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.friend-list::-webkit-scrollbar-track,
.request-list::-webkit-scrollbar-track,
.chat-messages::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.friend-list::-webkit-scrollbar-thumb,
.request-list::-webkit-scrollbar-thumb,
.chat-messages::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.friend-list::-webkit-scrollbar-thumb:hover,
.request-list::-webkit-scrollbar-thumb:hover,
.chat-messages::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .chat-layout {
    flex-direction: column;
    height: auto;
  }

  .friend-sidebar {
    width: 100%;
    height: 300px;
    border-right: none;
    border-bottom: 1px solid #f0f0f0;
  }

  .friend-search {
    flex-direction: column;
  }

  .friend-item, .request-item {
    padding: 10px 12px;
  }

  .friend-avatar {
    width: 36px;
    height: 36px;
    margin-right: 10px;
  }

  .friend-name, .request-name {
    font-size: 13px;
  }

  .request-actions {
    flex-direction: column;
    gap: 4px;
  }

  ::v-deep(.request-actions .ant-btn) {
    width: 100%;
  }
}

::v-deep(.ant-tabs-nav-list){
  margin: auto;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0.5; }
}


.typing-indicator {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  color: #888;
  font-size: 12px;
}
.typing-dots {
  display: flex;
  margin-right: 8px;
}
.typing-dots span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #888;
  margin: 0 2px;
  animation: typingAnimation 1.4s infinite ease-in-out both;
}
.typing-dots span:nth-child(1) { animation-delay: -0.32s; }
.typing-dots span:nth-child(2) { animation-delay: -0.16s; }
@keyframes typingAnimation {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}


.message-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 12px;
}
.own-message {
  justify-content: flex-end;

  .message-bubble {
    background-color: #1890ff;
    color: white;
    padding: 10px 14px;
    border-radius: 18px 5px 18px 18px;
    max-width: 300px;

  }

  .message-time {
    font-size: 12px;
    color: #8c8c8c;
    align-self: center;
    margin-bottom: 4px; // 时间也稍微低一点
  }
}
.other-message {
  justify-content: flex-start;

  .message-bubble {
    background-color: #f0f0f0;
    padding: 10px 14px;
    border-radius: 5px 18px 18px 18px;
    max-width: 300px;
  }

  .message-time {
    font-size: 12px;
    color: #8c8c8c;
    align-self: center;
  }
}
.message-text {
  word-break: break-word;
  line-height: 1.4;
}
/* 容器必须设置overflow和固定高度才能滚动 */
.messages-container {
  overflow-y: auto;
  height: 444px; /* 或 max-height */
}
@keyframes messageAppear {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
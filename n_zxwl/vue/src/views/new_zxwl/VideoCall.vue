<template>
  <div class="video-call-container">
    <div class="header">
      <h1>视频通话应用</h1>
      <div class="controls">
        <button @click="toggleVideo" :class="{ disabled: !videoEnabled }">
          {{ videoEnabled ? '关闭视频' : '开启视频' }}
        </button>
        <button @click="toggleAudio" :class="{ disabled: !audioEnabled }">
          {{ audioEnabled ? '静音' : '取消静音' }}
        </button>
        <button @click="startCall" :disabled="isConnected">开始通话</button>
        <button @click="endCall" :disabled="!isConnected">结束通话</button>
      </div>
    </div>

    <div class="video-container">
      <div class="video-wrapper">
        <video ref="localVideo" autoplay muted class="video-element"></video>
        <span class="video-label">本地视频</span>
      </div>
      <div class="video-wrapper">
        <video ref="remoteVideo" autoplay class="video-element"></video>
        <span class="video-label">远程视频</span>
      </div>
    </div>

    <div class="connection-status">
      <p>状态: {{ connectionStatus }}</p>
    </div>

    <div v-if="!isConnected" class="join-form">
      <input v-model="roomId" placeholder="输入房间号" type="text" />
      <button @click="joinRoom">加入房间</button>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'

export default {
  name: 'VideoCall',
  setup() {
    const localVideo = ref(null)
    const remoteVideo = ref(null)
    const videoEnabled = ref(true)
    const audioEnabled = ref(true)
    const isConnected = ref(false)
    const connectionStatus = ref('断开连接')
    const roomId = ref('')
    const localStream = ref(null)
    const peerConnection = ref(null)
    const socket = ref(null)

    // WebRTC配置
    const configuration = {
      iceServers: [
        { urls: 'stun:stun.l.google.com:19302' },
        { urls: 'stun:stun1.l.google.com:19302' }
      ]
    }

    // 初始化媒体流
    const initStream = async () => {
      try {
        localStream.value = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true
        })
        localVideo.value.srcObject = localStream.value
      } catch (error) {
        console.error('获取媒体设备失败:', error)
        alert('无法访问摄像头或麦克风')
      }
    }

    // 创建WebSocket连接
    const connectWebSocket = () => {
      socket.value = new WebSocket('ws://localhost:8080/ws')

      socket.value.onopen = () => {
        console.log('WebSocket连接已建立')
        connectionStatus.value = '已连接信令服务器'
      }

      socket.value.onmessage = async (event) => {
        const message = JSON.parse(event.data)
        await handleSignalingMessage(message)
      }

      socket.value.onclose = () => {
        console.log('WebSocket连接已关闭')
        connectionStatus.value = '信令服务器断开连接'
      }
    }

    // 处理信令消息
    const handleSignalingMessage = async (message) => {
      switch (message.type) {
        case 'offer':
          await handleOffer(message)
          break
        case 'answer':
          await handleAnswer(message)
          break
        case 'ice-candidate':
          await handleIceCandidate(message)
          break
        case 'room-joined':
          connectionStatus.value = '已加入房间: ' + message.room
          break
        case 'error':
          alert('错误: ' + message.message)
          break
      }
    }

    // 处理Offer
    const handleOffer = async (message) => {
      if (!peerConnection.value) {
        createPeerConnection()
      }

      await peerConnection.value.setRemoteDescription(message.offer)
      const answer = await peerConnection.value.createAnswer()
      await peerConnection.value.setLocalDescription(answer)

      socket.value.send(JSON.stringify({
        type: 'answer',
        answer: answer,
        room: roomId.value
      }))
    }

    // 处理Answer
    const handleAnswer = async (message) => {
      await peerConnection.value.setRemoteDescription(message.answer)
    }

    // 处理ICE候选
    const handleIceCandidate = async (message) => {
      try {
        await peerConnection.value.addIceCandidate(message.candidate)
      } catch (error) {
        console.error('添加ICE候选失败:', error)
      }
    }

    // 创建对等连接
    const createPeerConnection = () => {
      peerConnection.value = new RTCPeerConnection(configuration)

      // 添加本地流
      localStream.value.getTracks().forEach(track => {
        peerConnection.value.addTrack(track, localStream.value)
      })

      // 处理远程流
      peerConnection.value.ontrack = (event) => {
        remoteVideo.value.srcObject = event.streams[0]
      }

      // 处理ICE候选
      peerConnection.value.onicecandidate = (event) => {
        if (event.candidate) {
          socket.value.send(JSON.stringify({
            type: 'ice-candidate',
            candidate: event.candidate,
            room: roomId.value
          }))
        }
      }

      peerConnection.value.onconnectionstatechange = () => {
        connectionStatus.value = peerConnection.value.connectionState
        isConnected.value = peerConnection.value.connectionState === 'connected'
      }
    }

    // 开始通话
    const startCall = async () => {
      if (!roomId.value) {
        alert('请先输入房间号')
        return
      }

      createPeerConnection()

      const offer = await peerConnection.value.createOffer()
      await peerConnection.value.setLocalDescription(offer)

      socket.value.send(JSON.stringify({
        type: 'offer',
        offer: offer,
        room: roomId.value
      }))
    }

    // 结束通话
    const endCall = () => {
      if (peerConnection.value) {
        peerConnection.value.close()
        peerConnection.value = null
      }

      if (socket.value) {
        socket.value.send(JSON.stringify({
          type: 'leave',
          room: roomId.value
        }))
      }

      isConnected.value = false
      connectionStatus.value = '通话已结束'

      if (remoteVideo.value.srcObject) {
        remoteVideo.value.srcObject.getTracks().forEach(track => track.stop())
      }
    }

    // 切换视频
    const toggleVideo = () => {
      if (localStream.value) {
        const videoTracks = localStream.value.getVideoTracks()
        if (videoTracks.length > 0) {
          videoEnabled.value = !videoTracks[0].enabled
          videoTracks[0].enabled = videoEnabled.value
        }
      }
    }

    // 切换音频
    const toggleAudio = () => {
      if (localStream.value) {
        const audioTracks = localStream.value.getAudioTracks()
        if (audioTracks.length > 0) {
          audioEnabled.value = !audioTracks[0].enabled
          audioTracks[0].enabled = audioEnabled.value
        }
      }
    }

    // 加入房间
    const joinRoom = () => {
      if (!roomId.value) {
        alert('请输入房间号')
        return
      }

      if (socket.value && socket.value.readyState === WebSocket.OPEN) {
        socket.value.send(JSON.stringify({
          type: 'join',
          room: roomId.value
        }))
      } else {
        alert('请先连接信令服务器')
      }
    }

    onMounted(() => {
      initStream()
      connectWebSocket()
    })

    onUnmounted(() => {
      if (socket.value) {
        socket.value.close()
      }

      if (localStream.value) {
        localStream.value.getTracks().forEach(track => track.stop())
      }

      if (peerConnection.value) {
        peerConnection.value.close()
      }
    })

    return {
      localVideo,
      remoteVideo,
      videoEnabled,
      audioEnabled,
      isConnected,
      connectionStatus,
      roomId,
      startCall,
      endCall,
      toggleVideo,
      toggleAudio,
      joinRoom
    }
  }
}
</script>

<style scoped>
.video-call-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #1a1a1a;
  color: white;
}

.header {
  padding: 1rem;
  background-color: #2d2d2d;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.controls {
  display: flex;
  gap: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  background-color: #4CAF50;
  color: white;
  cursor: pointer;
}

button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

button.disabled {
  background-color: #f44336;
}

.video-container {
  flex: 1;
  display: flex;
  padding: 1rem;
  gap: 1rem;
}

.video-wrapper {
  flex: 1;
  position: relative;
  background-color: #000;
  border-radius: 8px;
  overflow: hidden;
}

.video-element {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-label {
  position: absolute;
  bottom: 10px;
  left: 10px;
  background-color: rgba(0, 0, 0, 0.5);
  padding: 5px 10px;
  border-radius: 4px;
}

.connection-status {
  padding: 1rem;
  text-align: center;
  background-color: #2d2d2d;
}

.join-form {
  padding: 1rem;
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  background-color: #2d2d2d;
}

input {
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .video-container {
    flex-direction: column;
  }

  .header {
    flex-direction: column;
    gap: 1rem;
  }

  .controls {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>
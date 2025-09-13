<template>
  <div class="video-chat-container">
    <div class="main-content">
      <!-- 左侧视频区域 - 仅替换这部分 -->
      <div class="left-content">
        <VideoContainer :users="meetingUsers" />
      </div>

      <!-- 中间控制面板 - 保持原样 -->
      <div class="control-panel" :style="{ width: controlPanelWidth + 'px' }">
        <div class="control-buttons">
          <a-tooltip placement="right">
            <template #title>
              <span>邀请他人</span>
            </template>
            <a-button
                class="control-btn"
                type="primary"
                shape="circle"
                @click="showInviteDialog = true"
            >
              <template #icon>
                <UserAddOutlined />
              </template>
            </a-button>
          </a-tooltip>

          <a-tooltip placement="right">
            <template #title>
              <span>{{ isMuted ? '取消静音' : '静音' }}</span>
            </template>
            <a-button
                class="control-btn"
                :type="isMuted ? 'default' : 'primary'"
                shape="circle"
                @click="toggleMute"
            >
              <template #icon>
                <AudioMutedOutlined v-if="isMuted" />
                <AudioOutlined v-else />
              </template>
            </a-button>
          </a-tooltip>

          <a-tooltip placement="right">
            <template #title>
              <span>{{ isCameraOff ? '开启摄像头' : '关闭摄像头' }}</span>
            </template>
            <a-button
                class="control-btn"
                :type="isCameraOff ? 'default' : 'primary'"
                shape="circle"
                @click="toggleCamera"
            >
              <template #icon>
                <VideoCameraOutlined v-if="!isCameraOff" />
                <VideoCameraAddOutlined v-else />
              </template>
            </a-button>
          </a-tooltip>

          <a-tooltip placement="right">
            <template #title>
              <span>结束通话</span>
            </template>
            <a-button
                class="control-btn"
                type="primary"
                danger
                shape="circle"
                @click="endCall"
            >
              <template #icon>
                <PhoneOutlined />
              </template>
            </a-button>
          </a-tooltip>

          <a-tooltip placement="right">
            <template #title>
              <span>{{ showParticipants ? '隐藏参会者' : '显示参会者' }}</span>
            </template>
            <a-button
                class="control-btn"
                :type="showParticipants ? 'primary' : 'default'"
                shape="circle"
                @click="toggleParticipants"
            >
              <template #icon>
                <TeamOutlined />
              </template>
            </a-button>
          </a-tooltip>

          <!-- 摄像头选择按钮 -->
          <a-tooltip placement="right">
            <template #title>
              <span>选择摄像头</span>
            </template>
            <a-button
                class="control-btn"
                type="default"
                shape="circle"
                @click="showCameraSelector = true"
            >
              <template #icon>
                <CameraOutlined />
              </template>
            </a-button>
          </a-tooltip>
          <!-- 摄像头选择对话框 -->
          <a-modal
              v-model:open="showCameraSelector"
              title="选择摄像头"
              :footer="null"
              width="400px"
          >
            <div class="camera-selector">
              <div
                  v-for="camera in availableCameras"
                  :key="camera.deviceId"
                  class="camera-option"
                  :class="{ active: selectedCameraId === camera.deviceId }"
                  @click="switchCamera(camera.deviceId)"
              >
                <div class="camera-info">
                  <div class="camera-name">{{ camera.label || `摄像头 ${camera.deviceId.slice(0, 8)}` }}</div>
                  <div class="camera-id">{{ camera.deviceId }}</div>
                </div>
                <CheckOutlined v-if="selectedCameraId === camera.deviceId" />
              </div>
            </div>
          </a-modal>
        </div>
      </div>
    </div>

    <!-- 右侧参会者面板 - 保持原样 -->
    <div
        class="participants-panel"
        v-show="showParticipantsPanel"
        :style="{ width: participantsPanelWidth + 'px' }"
    >
      <participants-panel :users="users" />
    </div>

    <!-- 邀请对话框 - 保持原样 -->
    <invite-dialog
        v-model:open="showInviteDialog"
        @invite="handleInvite"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  UserAddOutlined,
  AudioMutedOutlined,
  AudioOutlined,
  VideoCameraOutlined,
  VideoCameraAddOutlined,
  PhoneOutlined,
  TeamOutlined,
  CameraOutlined,
  CheckOutlined
} from '@ant-design/icons-vue'

// 组件导入 - 只添加 VideoContainer，其他保持不变
import ParticipantsPanel from '../components/ParticipantsPanel.vue'
import InviteDialog from '../components/InviteDialog.vue'
import VideoContainer from '../components/VideoContainer.vue'

const router = useRouter()

const isMuted = ref(true)   // 静音状态
const isCameraOff = ref(false) // 摄像头状态
const showParticipants = ref(false) // 是否显示参会者面板
const showInviteDialog = ref(false) // 邀请对话框状态
const controlPanelWidth = ref(60)  // 控制面板宽度
const participantsPanelWidth = ref(250) // 参会者面板宽度
const availableCameras = ref([])  // 可用的摄像头列表
const selectedCameraId = ref('')  // 当前选中的摄像头ID
const showCameraSelector = ref(false) // 摄像头选择对话框状态

const localStream = ref(null) // 本地媒体流
const mediaAccessError = ref(false) // 媒体访问错误

// 用户数据 - 简化处理
const meetingUsers = ref([
  {
    id: 'local-user',
    name: '我',
    isCameraOff: false,
    isMuted: false,
    isLocal: true,
    joinTime: Date.now(),
    stream: null
  },
  {
    id: 'user2',
    name: '李四',
    isCameraOff: false,
    isMuted: true,
    isLocal: false,
    joinTime: Date.now() - 5000,
    stream: null
  },
  {
    id: 'user3',
    name: '王五',
    isCameraOff: true,
    isMuted: false,
    isLocal: false,
    joinTime: Date.now() - 3000,
    stream: null
  }
])

// 计算属性 - 保持原样
const showParticipantsPanel = computed(() => showParticipants.value)

// 控制功能 - 保持原样
const toggleMute = () => {
  isMuted.value = !isMuted.value
  if (localStream.value) {
    localStream.value.getAudioTracks().forEach(track => {
      track.enabled = !isMuted.value
    })
  }
}

const toggleCamera = () => {
  isCameraOff.value = !isCameraOff.value
  if (localStream.value) {
    localStream.value.getVideoTracks().forEach(track => {
      track.enabled = !isCameraOff.value
    })
  }
}

const endCall = () => {
  if (localStream.value) {
    localStream.value.getTracks().forEach(track => track.stop())
  }
  router.back()
}

const toggleParticipants = () => {
  showParticipants.value = !showParticipants.value
}

const handleInvite = (inviteData) => {
  console.log('邀请数据:', inviteData)
  showInviteDialog.value = false
}

// 媒体流相关 - 保持原样
const getMediaStream = async () => {
  try {
    const devices = await navigator.mediaDevices.enumerateDevices()
    const videoDevices = devices.filter(device => device.kind === 'videoinput')
    availableCameras.value = videoDevices

    const stream = await navigator.mediaDevices.getUserMedia({
      video: {
        width: { ideal: 640 },
        height: { ideal: 480 }
      },
      audio: {
        echoCancellation: true,
        noiseSuppression: true
      }
    })

    localStream.value = stream
    toggleAudio(false) // 默认静音

  } catch (error) {
    console.error('获取媒体设备失败:', error)
    mediaAccessError.value = true
  }
}

const switchCamera = async (deviceId) => {
  try {
    if (localStream.value) {
      localStream.value.getTracks().forEach(track => track.stop())
    }
    selectedCameraId.value = deviceId

    const stream = await navigator.mediaDevices.getUserMedia({
      video: {
        deviceId: { exact: deviceId },
        width: { ideal: 640 },
        height: { ideal: 480 }
      },
      audio: {
        echoCancellation: true,
        noiseSuppression: true
      }
    })

    localStream.value = stream

  } catch (error) {
    console.error('切换摄像头失败:', error)
  }
}

const toggleAudio = (enabled) => {
  if (localStream.value) {
    localStream.value.getAudioTracks().forEach(track => {
      track.enabled = enabled
    })
  }
}

// 生命周期 - 保持原样
onMounted(() => {
  getMediaStream()
})

onUnmounted(() => {
  if (localStream.value) {
    localStream.value.getTracks().forEach(track => track.stop())
  }
})
</script>

<style scoped>
.video-chat-container {
  width: 100%;
  height: 100vh;
  display: flex;
  background-color: #f5f7fa;
  overflow: hidden;
}

.main-content {
  display: flex;
  flex: 1;
  height: 100%;
}

.left-content {
  flex: 1;
  height: 100%;
  background-color: #fff;
  border-right: 1px solid #e8e8e8;
  overflow: hidden;
}

.control-panel {
  height: 100%;
  background-color: #fff;
  border-left: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;
  box-sizing: border-box;
}

.control-buttons {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.control-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.participants-panel {
  height: 100%;
  background-color: #fff;
  border-left: 1px solid #e8e8e8;
}
</style>
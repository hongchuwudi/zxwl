<template>
  <div
      class="video-element"
      :class="{
      'single-mode': layoutMode === 'single',
      'pip-mode': layoutMode === 'pip',
      'grid-mode': layoutMode === 'grid'
    }"
      @click="onVideoClick"
  >
    <div class="video-wrapper">
      <video
          v-if="!user.isCameraOff"
          ref="videoElement"
          autoplay
          :muted="user.isLocal"
          playsinline
          class="video"
      ></video>
      <div v-else class="no-video">
        <div class="user-avatar">
          {{ user.name.charAt(0) }}
        </div>
      </div>

      <div class="video-overlay">
        <div class="user-name">{{ user.name }}</div>
        <div class="status-icons">
          <span v-if="user.isMuted" class="muted-icon">
            <AudioMutedOutlined />
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { AudioMutedOutlined } from '@ant-design/icons-vue'
import { ref, onMounted, watch, computed } from 'vue'

const videoElement = ref(null)
const props = defineProps({
  user: {
    type: Object,
    required: true
  },
  stream: {
    type: MediaStream,
    default: null
  },
  layoutMode: {
    type: String,
    default: 'grid' // 'single', 'pip', 'grid'
  },
  isMainVideo: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['switch-main'])

// 监听stream变化并绑定到video元素
watch(() => props.stream, (newStream) => {
  if (videoElement.value && newStream) {
    videoElement.value.srcObject = newStream
  }
})

onMounted(() => {
  if (videoElement.value && props.stream) {
    videoElement.value.srcObject = props.stream
  }
})

const onVideoClick = () => {
  if (props.layoutMode === 'pip' && !props.isMainVideo) {
    emit('switch-main', props.user.id)
  }
}
</script>

<style scoped>
.video-element {
  position: relative;
  transition: all 0.3s ease;
}

/* 单人模式 */
.video-element.single-mode {
  width: 100%;
  height: 100%;
}

/* 画中画模式 */
.video-element.pip-mode {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 120px;
  height: 90px;
  z-index: 10;
  border: 2px solid #1890ff;
  border-radius: 6px;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.video-element.pip-mode.is-main-video {
  position: relative;
  width: 100%;
  height: 100%;
  top: 0;
  right: 0;
  border: none;
  box-shadow: none;
}

/* 网格模式 */
.video-element.grid-mode {
  width: 100%;
  height: 100%;
}

.video-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
  background-color: #000;
  border-radius: 8px;
  overflow: hidden;
}

.video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-video {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #262626;
}

.user-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: #1890ff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: bold;
}

.video-element.pip-mode .user-avatar {
  width: 40px;
  height: 40px;
  font-size: 16px;
}

.video-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 8px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.7));
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-name {
  font-size: 12px;
}

.video-element.pip-mode .user-name {
  font-size: 10px;
}

.status-icons {
  display: flex;
}

.muted-icon {
  color: #ff4d4f;
}
</style>
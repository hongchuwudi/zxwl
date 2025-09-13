<template>
  <div class="video-container">
    <div class="video-area" :class="layoutClass">
      <VideoElement
          v-for="user in displayedUsers"
          :key="user.id"
          :user="user"
          :stream="user.stream"
          :layout-mode="layoutMode"
          :is-main-video="user.id === mainVideoId"
          @switch-main="switchMainVideo"
      />
    </div>

    <!-- 分页控制（4人以上显示） -->
    <div v-if="totalPages > 1" class="pagination-controls">
      <button
          class="pagination-btn"
          @click="prevPage"
          :disabled="currentPage === 1"
      >
        <LeftOutlined />
      </button>
      <span class="pagination-info">
        第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
      </span>
      <button
          class="pagination-btn"
          @click="nextPage"
          :disabled="currentPage === totalPages"
      >
        <RightOutlined />
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { LeftOutlined, RightOutlined } from '@ant-design/icons-vue'
import VideoElement from './VideoElement.vue'

const props = defineProps({
  users: {
    type: Array,
    required: true,
    default: () => []
  }
})

const currentPage = ref(1)
const pageSize = 4
const mainVideoId = ref(null)

// 按进入时间排序
const sortedUsers = computed(() => {
  return [...props.users].sort((a, b) => a.joinTime - b.joinTime)
})

// 计算布局模式
const layoutMode = computed(() => {
  const total = sortedUsers.value.length
  if (total === 1) return 'single'
  if (total === 2) return 'pip'
  return 'grid'
})

// 计算布局类名
const layoutClass = computed(() => {
  const total = displayedUsers.value.length
  if (total === 1) return 'layout-single'
  if (total === 2) return 'layout-pip'
  if (total === 3) return 'layout-grid-3'
  return 'layout-grid-4'
})

// 总页数
const totalPages = computed(() => {
  return Math.ceil(sortedUsers.value.length / pageSize)
})

// 当前页显示的用户
const displayedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  const pageUsers = sortedUsers.value.slice(start, end)

  // 设置主视频（第一页的第一个用户）
  if (currentPage.value === 1 && pageUsers.length > 0 && !mainVideoId.value) {
    mainVideoId.value = pageUsers[0].id
  }

  return pageUsers
})

// 切换主视频（画中画模式）
const switchMainVideo = (userId) => {
  mainVideoId.value = userId
}

// 分页控制
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

// 监听用户数量变化
watch(() => props.users.length, (newCount, oldCount) => {
  if (newCount <= pageSize) {
    currentPage.value = 1
  }

  // 如果当前页的用户离开了，需要重新设置主视频
  if (newCount < oldCount) {
    const currentUserIds = displayedUsers.value.map(user => user.id)
    if (!currentUserIds.includes(mainVideoId.value)) {
      mainVideoId.value = displayedUsers.value[0]?.id || null
    }
  }
})

// 监听当前页变化，确保主视频在当前页
watch(currentPage, (newPage) => {
  if (displayedUsers.value.length > 0 && !displayedUsers.value.some(user => user.id === mainVideoId.value)) {
    mainVideoId.value = displayedUsers.value[0].id
  }
})
</script>

<style scoped>
.video-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1a1a1a;
}

.video-area {
  flex: 1;
  display: grid;
  gap: 8px;
  padding: 8px;
  position: relative;
}

/* 单人布局 - 全屏 */
.layout-single {
  grid-template-columns: 1fr;
  grid-template-rows: 1fr;
}

/* 画中画布局 - 主视频全屏，小窗口绝对定位 */
.layout-pip {
  grid-template-columns: 1fr;
  grid-template-rows: 1fr;
  position: relative;
}

/* 3人网格布局 - 2x2网格，第三个占满底部 */
.layout-grid-3 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

.layout-grid-3 ::v-deep(.video-element:nth-child(3)) {
  grid-column: span 2;
}

/* 4人网格布局 - 标准的2x2网格 */
.layout-grid-4 {
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

/* 分页控制 */
.pagination-controls {
  padding: 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  background: #2a2a2a;
  border-top: 1px solid #404040;
}

.pagination-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: #1890ff;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.pagination-btn:hover:not(:disabled) {
  background: #40a9ff;
  transform: translateY(-1px);
}

.pagination-btn:disabled {
  background: #595959;
  cursor: not-allowed;
  opacity: 0.6;
}

.pagination-info {
  color: #d9d9d9;
  font-size: 14px;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .video-area {
    gap: 4px;
    padding: 4px;
  }

  .pagination-controls {
    padding: 12px;
    gap: 12px;
  }

  .pagination-info {
    font-size: 12px;
  }
}

/* 画中画模式的特殊处理 */
.layout-pip ::v-deep(.video-element.pip-mode) {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 120px;
  height: 90px;
  z-index: 10;
}

.layout-pip ::v-deep(.video-element.pip-mode.is-main-video) {
  position: relative;
  width: 100%;
  height: 100%;
  top: 0;
  right: 0;
}

/* 确保视频元素在网格中正确显示 */
.video-area ::v-deep(.video-element) {
  width: 100%;
  height: 100%;
  min-height: 0; /* 防止网格项溢出 */
}
</style>
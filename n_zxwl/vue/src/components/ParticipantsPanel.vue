<template>
  <div class="participants-panel">
    <div class="panel-header">
      <h3>参会人员 ({{ users.length }})</h3>
    </div>
    <div class="participants-list">
      <div
          v-for="user in users"
          :key="user.id"
          class="participant-item"
      >
        <div class="avatar">
          {{ user.name.charAt(0) }}
        </div>
        <div class="user-info">
          <div class="user-name">{{ user.name }}</div>
          <div class="user-status">
            <span v-if="user.isMuted" class="muted">已静音</span>
            <span v-else class="speaking">发言中</span>
            <span v-if="user.isCameraOff" class="camera-off">摄像头关闭</span>
          </div>
        </div>
        <div class="actions">
          <a-button type="text" size="small">
            <template #icon>
              <MessageOutlined />
            </template>
          </a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { MessageOutlined } from '@ant-design/icons-vue'

defineProps({
  users: {
    type: Array,
    default: () => []
  }
})
</script>

<style scoped>
.participants-panel {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.panel-header {
  padding: 16px;
  border-bottom: 1px solid #e8e8e8;
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  color: #262626;
}

.participants-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.participant-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  transition: background-color 0.2s;
}

.participant-item:hover {
  background-color: #f5f5f5;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #1890ff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 12px;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.user-status {
  font-size: 12px;
  color: #8c8c8c;
}

.user-status .muted {
  color: #ff4d4f;
}

.user-status .speaking {
  color: #52c41a;
}

.user-status .camera-off {
  margin-left: 8px;
}

.actions {
  opacity: 0;
  transition: opacity 0.2s;
}

.participant-item:hover .actions {
  opacity: 1;
}
</style>
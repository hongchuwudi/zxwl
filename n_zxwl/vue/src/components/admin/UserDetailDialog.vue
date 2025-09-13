<script setup>
import { computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  User,
  Message,
  Calendar,
  Location,
  Phone,
  Clock,
  Male,
  Female,

} from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  user: Object
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formatGender = (gender) => {
  const genderMap = { 0: '未知', 1: '男', 2: '女' }
  return genderMap[gender] || '未知'
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

const formatLastOnline = (time) => {
  if (!time) return '从未在线'
  const now = new Date()
  const lastOnline = new Date(time)
  const diffInMinutes = Math.floor((now - lastOnline) / (1000 * 60))

  if (diffInMinutes < 1) return '刚刚'
  if (diffInMinutes < 60) return `${diffInMinutes}分钟前`
  if (diffInMinutes < 1440) return `${Math.floor(diffInMinutes / 60)}小时前`

  return lastOnline.toLocaleDateString('zh-CN')
}

const getGenderIcon = gender => gender === 1 ? Male : gender === 2 ? Female : Clock

</script>

<template>
  <el-dialog
      v-model="dialogVisible"
      :title="`用户详情 - ${user?.username || '未知用户'}`"
      width="600px"
      destroy-on-close
  >
    <div v-if="user" class="user-detail-container">
      <!-- 用户头像和信息 -->
      <div class="user-profile">
        <el-avatar
            :size="80"
            :src="user.avatarUrl"
            class="profile-avatar"
        />
        <div class="profile-info">
          <div class="username">
            {{ user.username }}
            <el-tag
                :type="user.isOnline ? 'success' : 'info'"
                size="small"
                class="online-status"
            >
              {{ user.isOnline ? '在线' : '离线' }}
            </el-tag>
          </div>
          <div class="email">
            <el-icon><Message /></el-icon>
            {{ user.email || '未设置邮箱' }}
          </div>
          <div class="display-name">
            {{ user.displayName || '未设置昵称' }}
          </div>
        </div>
      </div>

      <el-divider />

      <!-- 基本信息 -->
      <div class="detail-section">
        <h3>基本信息</h3>
        <div class="detail-grid">
          <div class="detail-item">
            <div class="detail-label">
              <el-icon><User /></el-icon>
              性别
            </div>
            <div class="detail-value">
              <el-icon v-if="user.gender === 1"><Male /></el-icon>
              <el-icon v-else-if="user.gender === 2"><Female /></el-icon>
              <el-icon v-else><Question /></el-icon>
              {{ formatGender(user.gender) }}
            </div>
          </div>

          <div class="detail-item">
            <div class="detail-label">
              <el-icon><Calendar /></el-icon>
              出生年份
            </div>
            <div class="detail-value">
              {{ user.birthYear || '未设置' }}
            </div>
          </div>

          <div class="detail-item">
            <div class="detail-label">
              <el-icon><Location /></el-icon>
              位置
            </div>
            <div class="detail-value">
              {{ user.location || '未设置位置' }}
            </div>
          </div>

          <div class="detail-item">
            <div class="detail-label">
              <el-icon><Phone /></el-icon>
              电话
            </div>
            <div class="detail-value">
              {{ user.phone || '未设置电话' }}
            </div>
          </div>
        </div>
      </div>

      <el-divider />

      <!-- 状态信息 -->
      <div class="detail-section">
        <h3>状态信息</h3>
        <div class="detail-grid">
          <div class="detail-item">
            <div class="detail-label">
              <el-icon><Clock /></el-icon>
              最后在线
            </div>
            <div class="detail-value">
              {{ formatLastOnline(user.lastOnlineTime) }}
            </div>
          </div>

          <div class="detail-item">
            <div class="detail-label">
              注册时间
            </div>
            <div class="detail-value">
              {{ formatDateTime(user.createdAt) }}
            </div>
          </div>

          <div class="detail-item">
            <div class="detail-label">
              最后更新
            </div>
            <div class="detail-value">
              {{ formatDateTime(user.updatedAt) || '-' }}
            </div>
          </div>
        </div>
      </div>

      <!-- 其他信息 -->
      <div v-if="user.bio" class="detail-section">
        <h3>个人简介</h3>
        <div class="bio-content">
          {{ user.bio }}
        </div>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">关闭</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
.user-detail-container {
  .user-profile {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;

    .profile-info {
      .username {
        font-size: 18px;
        font-weight: 600;
        margin-bottom: 8px;
        display: flex;
        align-items: center;
        gap: 8px;
      }

      .email {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #666;
        margin-bottom: 4px;
      }

      .display-name {
        color: #888;
        font-size: 14px;
      }
    }
  }

  .detail-section {
    margin-bottom: 20px;

    h3 {
      margin: 0 0 16px 0;
      color: #333;
      font-size: 16px;
    }

    .detail-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 16px;
    }

    .detail-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 8px 0;

      .detail-label {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #666;
        font-weight: 500;
      }

      .detail-value {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #333;
      }
    }

    .bio-content {
      padding: 12px;
      background-color: #f9f9f9;
      border-radius: 4px;
      line-height: 1.6;
    }
  }
}
</style>
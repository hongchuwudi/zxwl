<template>
  <div class="user-profile-container">
    <a-page-header
        title="个人信息"
        @back="() => $router.go(-1)"
        class="page-header"
    />

    <a-spin :spinning="loading">
      <div class="profile-content">
        <!-- 左侧用户信息卡片 -->
        <div class="left-column">
          <div class="user-info-section" v-if="user && user.username">
            <a-card :bordered="false" class="user-card">
              <div class="avatar-section">
                <a-avatar :size="120" class="user-avatar" :src="user.avatarUrl" v-if="user.avatarUrl"/>
                <a-avatar :size="120" class="user-avatar" v-else>
                  {{ (user && (user.displayName ? user.displayName.charAt(0) : user.username.charAt(0))) }}
                </a-avatar>
                <div class="user-name" v-if="user">
                  <h2>{{ user.displayName || user.username }}</h2>
                  <p>@{{ user.username }}</p>
                  <a-tag :color="user.isOnline ? 'green' : 'default'">
                    {{ user.isOnline ? '在线' : '离线' }}
                  </a-tag>
                </div>
                <a-skeleton v-else active :paragraph="{ rows: 2 }" />
              </div>

              <!-- 不可修改的信息区域 -->
              <div class="readonly-info">
                <div class="info-item">
                  <span class="info-label">邮箱:</span>
                  <span class="info-value">{{ user.email }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">注册时间:</span>
                  <span class="info-value">{{ formatDate(user.createdAt) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">最后更新:</span>
                  <span class="info-value">{{ formatDate(user.updatedAt) }}</span>
                </div>
                <div class="info-item" v-if="user.lastOnlineTime">
                  <span class="info-label">最后在线:</span>
                  <span class="info-value">{{ formatDate(user.lastOnlineTime) }}</span>
                </div>
              </div>
            </a-card>
          </div >
          <!-- 加载状态 -->
          <div v-else>
            <a-skeleton active :paragraph="{ rows: 4 }" />
          </div>

          <a-card class="user-info-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <user-outlined />
                <span>基本信息</span>
                <a-button
                    v-if="!isEditing"
                    type="link"
                    @click="startEditing"
                    class="edit-btn"
                >
                  编辑
                </a-button>
                <div v-else class="edit-actions">
                  <a-button type="link" @click="saveChanges" :loading="saving">保存</a-button>
                  <a-button type="link" @click="cancelEditing">取消</a-button>
                </div>
              </div>
            </template>

            <a-form layout="vertical" v-if="user">
              <a-form-item label="显示名称">
                <a-input
                    v-model:value="user.displayName"
                    :disabled="!isEditing"
                    placeholder="请输入显示名称"
                />
              </a-form-item>

              <a-form-item label="性别">
                <a-select
                    v-model:value="user.gender"
                    :disabled="!isEditing"
                    placeholder="请选择性别"
                >
                  <a-select-option :value="3" disabled="">未知</a-select-option>
                  <a-select-option :value="1">男</a-select-option>
                  <a-select-option :value="2">女</a-select-option>
                </a-select>
              </a-form-item>

              <a-form-item label="出生年份">
                <a-input-number
                    v-model:value="user.birthYear"
                    :disabled="!isEditing"
                    :min="1900"
                    :max="new Date().getFullYear()"
                    style="width: 100%"
                    placeholder="请输入出生年份"
                />
              </a-form-item>

              <a-form-item label="所在地">
                <a-input
                    v-model:value="user.location"
                    :disabled="!isEditing"
                    placeholder="请输入所在地"
                />
              </a-form-item>

              <a-form-item label="个人简介">
                <a-textarea
                    v-model:value="user.bio"
                    :disabled="!isEditing"
                    :rows="3"
                    placeholder="请输入个人简介"
                />
              </a-form-item>
            </a-form>
          </a-card>
        </div>

        <!-- 中间贪吃蛇动画 -->
        <div class="center-column">
          <a-card class="snake-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <code-outlined />
                <span>活跃度</span>
              </div>
            </template>
            <SnakeAnimation />
          </a-card>

          <!-- 志愿填报小组件 -->
          <VolunteerMiniCard />
        </div>

        <!-- 右侧设备信息 -->
        <div class="right-column">
          <a-card class="device-info-card" :bordered="false">
            <template #title>
              <div class="card-title">
                <laptop-outlined />
                <span>设备信息</span>
              </div>
            </template>

            <div v-if="user && user.deviceInfo">
              <a-collapse v-model:activeKey="activeDeviceKeys" accordion>
                <a-collapse-panel key="browser" header="浏览器信息">
                  <div class="device-detail">
                    <p><strong>名称:</strong> {{ user.deviceInfo.browser?.name }}</p>
                    <p><strong>版本:</strong> {{ user.deviceInfo.browser?.version }}</p>
                    <p><strong>用户代理:</strong> {{ user.deviceInfo.userAgent }}</p>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="os" header="操作系统">
                  <div class="device-detail">
                    <p><strong>名称:</strong> {{ user.deviceInfo.os?.name }}</p>
                    <p><strong>版本:</strong> {{ user.deviceInfo.os?.version }}</p>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="screen" header="屏幕信息">
                  <div class="device-detail">
                    <p><strong>分辨率:</strong> {{ user.deviceInfo.screen?.width }} × {{ user.deviceInfo.screen?.height }}</p>
                    <p><strong>可用分辨率:</strong> {{ user.deviceInfo.screen?.availWidth }} × {{ user.deviceInfo.screen?.availHeight }}</p>
                    <p><strong>色彩深度:</strong> {{ user.deviceInfo.screen?.colorDepth }} 位</p>
                    <p><strong>像素深度:</strong> {{ user.deviceInfo.screen?.pixelDepth }} 位</p>
                    <p><strong>方向:</strong> {{ user.deviceInfo.screen?.orientation }}</p>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="network" header="网络信息">
                  <div class="device-detail">
                    <p><strong>类型:</strong> {{ user.deviceInfo.network?.type }}</p>
                    <p><strong>RTT:</strong> {{ user.deviceInfo.network?.rtt }}ms</p>
                    <p><strong>下行速度:</strong> {{ user.deviceInfo.network?.downlink }} Mbps</p>
                    <p><strong>节省数据模式:</strong> {{ user.deviceInfo.network?.saveData ? '是' : '否' }}</p>
                    <p><strong>有效类型:</strong> {{ user.deviceInfo.network?.effectiveType }}</p>
                  </div>
                </a-collapse-panel>

                <a-collapse-panel key="system" header="系统信息">
                  <div class="device-detail">
                    <p><strong>设备类型:</strong> {{ user.deviceInfo.device }}</p>
                    <p><strong>语言:</strong> {{ user.deviceInfo.language }}</p>
                    <p><strong>时区:</strong> {{ user.deviceInfo.timezone }}</p>
                    <p><strong>检测时间:</strong> {{ formatDate(user.deviceInfo.timestamp) }}</p>
                  </div>
                </a-collapse-panel>
              </a-collapse>
            </div>

            <div v-else class="no-device-info">
              <a-empty description="暂无设备信息" />
            </div>
          </a-card>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { UserOutlined, CodeOutlined, LaptopOutlined } from '@ant-design/icons-vue'
import SnakeAnimation from '@/components/SnakeAnimation.vue'
import {useUserStore} from "@/utils/auth.js";
import VolunteerMiniCard from '@/components/VolunteerMiniCard.vue'
const {getUser } = useUserStore()
import axios from "axios";
// 用户信息
const user = ref(null)
const originalUser = ref(null)
const loading = ref(false)
const saving = ref(false)
const isEditing = ref(false)
const activeDeviceKeys = ref(['browser'])

// 获取用户信息
const fetchUserData = async () => {
  try {
    loading.value = true;
    const userLocal = getUser();
    const userEmail = userLocal.email;
    const id = userLocal.id;

    if (!userEmail) {
      message.error('未获取到用户邮箱信息');
      return;
    }

    if (!id) {
      message.error('未获取到用户ID信息');
      return;
    }

    // 使用axios调用API获取用户信息
    const response = await axios.get(`/gapi/user/${id}`);
    const result = response.data;

    if (result.error === 0) {
      user.value = result.data;
      // 保存原始数据用于取消编辑时恢复
      originalUser.value = JSON.parse(JSON.stringify(result.data));
      // 解析设备信息JSON字符串
      if (user.value.deviceInfo && typeof user.value.deviceInfo === 'string') {
        try {
          user.value.deviceInfo = JSON.parse(user.value.deviceInfo);
        } catch (e) {
          console.error('解析设备信息失败:', e);
        }
      }
    } else {
      message.error(result.message || '获取用户信息失败');
    }
  } catch (error) {
    console.error('获取用户信息出错:', error);

    // 更详细的错误处理
    if (error.response) {
      // 服务器返回了错误状态码
      message.error(`请求失败: ${error.response.status} - ${error.response.data?.message || '服务器错误'}`);
    } else if (error.request) {
      // 请求已发出但没有收到响应
      message.error('网络错误，请检查网络连接');
    } else {
      // 其他错误
      message.error('获取用户信息失败');
    }
  } finally {
    loading.value = false;
  }
}

// 开始编辑
const startEditing = () => {
  isEditing.value = true
}

// 取消编辑
const cancelEditing = () => {
  isEditing.value = false
  user.value = JSON.parse(JSON.stringify(originalUser.value))
}

// 保存更改
const saveChanges = async () => {
  try {
    saving.value = true
    console.log('保存更改:', user.value)

    const requestData = {
      username: user.value.username,
      email: user.value.email,
      displayName: user.value.displayName,  // 改为蛇形命名
      gender: Number(user.value.gender),     // 确保是数字
      birthYear: user.value.birthYear ? Number(user.value.birthYear) : 0, // 转为数字
      location: user.value.location,
      bio: user.value.bio
    }

    let result = await axios.put(`/gapi/user/update/${getUser().id}`, requestData, {
      headers: {
        'Content-Type': 'application/json'
      }
    })
    result = result.data
    if (result.error === 0) {
      message.success('用户信息更新成功')
      isEditing.value = false
      // 更新原始数据
      originalUser.value = JSON.parse(JSON.stringify(user.value))
    } else {
      message.error(result.message || '更新用户信息失败')
    }
  } catch (error) {
    console.error('更新用户信息出错:', error)
    message.error('更新用户信息失败')
  } finally {
    saving.value = false
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'

  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchUserData()
})
</script>

<style scoped>
.user-profile-container {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  background-color: #fff;
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.profile-content {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 20px;
}

.left-column, .center-column, .right-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-info-card, .snake-card, .device-info-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.edit-btn, .edit-actions {
  margin-left: auto;
}

.device-detail p {
  margin-bottom: 8px;
}

.no-device-info {
  padding: 20px 0;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .profile-content {
    grid-template-columns: 1fr 1fr;
  }

  .center-column {
    grid-column: span 2;
    order: 3;
  }
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
  }

  .center-column {
    grid-column: span 1;
  }
}
.readonly-info {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.info-item {
  display: flex;
  margin-bottom: 12px;
  align-items: flex-start;
}

.info-label {
  font-weight: 500;
  color: #666;
  min-width: 80px;
  margin-right: 12px;
}

.info-value {
  color: #333;
  word-break: break-all;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .info-item {
    flex-direction: column;
  }

  .info-label {
    margin-bottom: 4px;
    min-width: auto;
  }
}

.avatar-section {
  display: flex;
  align-items: center;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
  margin-bottom: 16px;
}

.user-avatar {
  background: #1890ff;
  color: white;
  font-size: 36px;
  font-weight: bold;
  margin-right: 16px;
  flex-shrink: 0;
}

.user-name {
  flex: 1;
  min-width: 0;
}

.user-name h2 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-name p {
  margin: 0 0 8px 0;
  color: #8c8c8c;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
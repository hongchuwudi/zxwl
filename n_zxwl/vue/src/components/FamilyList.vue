<template>
  <a-card class="family-card" :bordered="false">
    <template #title>
      <div class="card-title">
        <team-outlined />
        <span>家人列表</span>
      </div>
    </template>

    <div v-if="loading" class="loading-container">
      <a-skeleton active :paragraph="{ rows: 3 }" />
    </div>

    <div v-else-if="familyList.length > 0" class="family-list">
      <div
          v-for="member in familyList"
          :key="member.id"
          class="family-member-item"
      >
        <a-avatar
            :size="36"
            :src="member.user_a.avatarUrl"
            class="member-avatar"
        >
          {{ member.user_a.displayName?.charAt(0) || member.user_a.username?.charAt(0) }}
        </a-avatar>

        <div class="member-info">
          <div class="member-name-row">
            <span class="member-name">{{ member.nickname || member.user_a.displayName || member.user_a.username }}</span>
            <a-tag
                :color="getRelationColor(member.relation_type)"
                size="small"
                class="relation-tag"
            >
              {{ formatRelationType(member.relation_type) }}
            </a-tag>
            <a-tag
                :color="member.user_a.isOnline ? 'green' : 'default'"
                size="small"
            >
              {{ member.user_a.isOnline ? '在线' : formatLastOnline(member.user_a.lastOnlineTime) }}
            </a-tag>
          </div>

          <div v-if="member.user_a.bio" class="member-bio">
            {{ member.user_a.bio }}
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-family">
      <a-empty description="暂无家人" />
    </div>
  </a-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { TeamOutlined } from '@ant-design/icons-vue'
import { useUserStore } from "@/utils/auth.js"
import axios from "axios"

const { getUser } = useUserStore()
const currentUser = getUser()

// 家人列表状态
const familyList = ref([])
const loading = ref(false)

// 获取家人列表
const fetchFamilyList = async () => {
  try {
    loading.value = true
    // 调用查看家人列表API - 假设使用好友接口
    const response = await axios.get(`/gapi/user/${currentUser.id}/friends`)
    if (response.data.code === 200) {
      familyList.value = response.data.data || []
    }
  } catch (error) {
    console.error('获取家人列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化关系类型
const formatRelationType = (type) => {
  const typeMap = {
    'normal': '普通',
    'close': '亲密',
    'family': '家人',
    'parent': '父母',
    'child': '子女',
    'spouse': '配偶'
  }
  return typeMap[type] || type
}

// 获取关系类型对应的颜色
const getRelationColor = (type) => {
  const colorMap = {
    'normal': 'default',
    'close': 'blue',
    'family': 'purple',
    'parent': 'volcano',
    'child': 'cyan',
    'spouse': 'red'
  }
  return colorMap[type] || 'default'
}

// 格式化最后在线时间
const formatLastOnline = (time) => {
  if (!time) return '离线'

  const now = new Date()
  const lastOnline = new Date(time)
  const diffHours = Math.floor((now - lastOnline) / (1000 * 60 * 60))

  if (diffHours < 1) return '刚刚在线'
  if (diffHours < 24) return `${diffHours}小时前`

  const diffDays = Math.floor(diffHours / 24)
  return `${diffDays}天前`
}

onMounted(() => {
  fetchFamilyList()
})
</script>

<style scoped>
.family-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  margin-bottom: 0;
}

.family-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.family-member-item {
  display: flex;
  align-items: flex-start;
  padding: 8px 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.family-member-item:hover {
  background-color: #f5f5f5;
}

.member-avatar {
  margin-right: 10px;
  flex-shrink: 0;
}

.member-info {
  flex: 1;
  min-width: 0;
}

.member-name-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 2px;
  flex-wrap: wrap;
}

.member-name {
  font-weight: 500;
  font-size: 13px;
  color: #262626;
}

.relation-tag {
  font-size: 11px;
  padding: 0 5px;
  height: 18px;
  line-height: 18px;
}

.member-bio {
  font-size: 12px;
  color: #8c8c8c;
  line-height: 1.4;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.loading-container {
  padding: 12px 0;
}

.empty-family {
  padding: 12px 0;
  text-align: center;
}

:deep(.ant-tag) {
  margin: 0;
  font-size: 11px;
  padding: 0 5px;
  height: 18px;
  line-height: 18px;
}
</style>
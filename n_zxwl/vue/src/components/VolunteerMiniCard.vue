<template>
  <a-card class="volunteer-mini-card" :bordered="false">
    <template #title>
      <div class="card-title">
        <form-outlined />
        <span>志愿填报</span>
      </div>
    </template>

    <div v-if="loading" class="loading-section">
      <a-spin size="small" />
      <span>加载中...</span>
    </div>

    <div v-else-if="volunteers.length === 0" class="empty-section">
      <a-empty description="暂无志愿数据" image-style="{ height: '60px' }">
        <a-button type="primary" size="small" @click="goToVolunteer">
          去填报
        </a-button>
      </a-empty>
    </div>

    <div v-else class="volunteer-content">
      <!-- 第一行：大学名称 -->
      <div class="info-row">
        <span class="label">院校:</span>
        <span class="value schools">
          {{ schoolNames }}
        </span>
      </div>

      <!-- 第二行：专业名称 -->
      <div class="info-row">
        <span class="label">专业:</span>
        <span class="value majors">
          {{ majorNames }}
        </span>
      </div>

      <!-- 第三行：操作按钮 -->
      <div class="action-buttons">
        <a-tooltip placement="top" title="智能推荐">
          <a-button type="text" size="default " @click="goToRecommendation">
            <template #icon><bulb-outlined /></template>
          </a-button>
        </a-tooltip>

        <a-tooltip placement="top" title="搜索大学">
          <a-button type="text" size="default " @click="goToUniversitySearch">
            <template #icon><search-outlined /></template>
          </a-button>
        </a-tooltip>

        <a-tooltip placement="top" title="搜索专业">
          <a-button type="text" size="default " @click="goToMajorSearch">
            <template #icon><search-outlined /></template>
          </a-button>
        </a-tooltip>

        <a-tooltip placement="top" title="AI咨询">
          <a-button type="text" size="default " @click="goToAIChat">
            <template #icon><message-outlined /></template>
          </a-button>
        </a-tooltip>

        <a-tooltip placement="top" title="查看政策">
          <a-button type="text" size="default "  @click="goToPolicy">
            <template #icon><file-text-outlined /></template>
          </a-button>
        </a-tooltip>

        <a-tooltip placement="top" title="视频咨询">
          <a-button type="text" size="default " @click="goToVideoCall">
            <template #icon><video-camera-outlined /></template>
          </a-button>
        </a-tooltip>

        <a-tooltip placement="top" title="管理志愿">
          <a-button type="primary" size="default " @click="goToVolunteer">
            <template #icon><edit-outlined /></template>
          </a-button>
        </a-tooltip>
      </div>
    </div>
  </a-card>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  FormOutlined,
  BulbOutlined,
  SearchOutlined,
  MessageOutlined,
  FileTextOutlined,
  VideoCameraOutlined,
  EditOutlined
} from '@ant-design/icons-vue'
import axios from "axios"
import {useUserStore} from "@/utils/auth.js";
const { getUser } = useUserStore()
const router = useRouter()
const loading = ref(false)
const volunteers = ref([])

// 计算属性
const schoolNames = computed(() => {
  const schools = volunteers.value
      .filter(v => v.school_name && v.school_name.trim())
      .map(v => v.school_name)
  return schools.length > 0 ? schools.join('，') : '暂无院校'
})

const majorNames = computed(() => {
  const majors = volunteers.value
      .filter(v => v.major_name && v.major_name.trim())
      .map(v => v.major_name)
  return majors.length > 0 ? majors.join('，') : '暂无专业'
})

// 加载志愿数据
const loadVolunteers = async () => {
  loading.value = true
  try {
    const user = getUser(0)
    if (!user.id) {
      console.warn('未获取到用户ID')
      return
    }

    const response = await axios.get(`/gapi/user/${user.id}/choices`)

    if (response.data && response.data.code === 200) {
      volunteers.value = response.data.data || []
    } else if (Array.isArray(response.data)) {
      volunteers.value = response.data
    } else {
      volunteers.value = []
    }
  } catch (error) {
    console.error('加载志愿数据失败:', error)
    volunteers.value = []
  } finally {
    loading.value = false
  }
}

// 导航方法
const goToVolunteer = () => router.push('/simulate')
const goToRecommendation = () => router.push('/recommends')
const goToUniversitySearch = () => router.push('/allSchool')
const goToMajorSearch = () => router.push('/professional')
const goToAIChat = () => router.push('/aismartsel')
const goToPolicy = () => router.push('/policy')
const goToVideoCall = () => router.push('/videoCall')

onMounted(() => {
  loadVolunteers()
})
</script>

<style scoped>
.volunteer-mini-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.loading-section {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
  justify-content: center;
  color: #8c8c8c;
}

.empty-section {
  padding: 10px 0;
}

.volunteer-content {
  padding: 8px 0;
}

.info-row {
  display: flex;
  margin-bottom: 12px;
  align-items: flex-start;
}

.info-row .label {
  font-weight: 500;
  color: #666;
  min-width: 40px;
  margin-right: 12px;
  flex-shrink: 0;
}

.info-row .value {
  color: #333;
  word-break: break-all;
  line-height: 1.4;
}

.schools, .majors {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.action-buttons {
  display: flex;
  gap: 6px;
  justify-content: center;
  padding-top: 8px;
  border-top: 1px solid #f0f0f0;
  margin-top: 8px;
}

.action-buttons ::v-deep(.ant-btn) {
  height: 28px;
  width: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .info-row {
    flex-direction: column;
  }

  .info-row .label {
    margin-bottom: 4px;
    min-width: auto;
  }

  .action-buttons {
    flex-wrap: wrap;
  }
}
</style>
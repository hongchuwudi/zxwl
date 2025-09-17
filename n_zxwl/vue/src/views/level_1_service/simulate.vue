<template>
  <div class="volunteer-container">
    <a-page-header
        title="志愿填报系统"
        @back="() => $router.go(-1)"
        class="page-header"
    >
      <template #extra>
        <div class="header-buttons">
          <!-- 搜索大学按钮 -->
          <a-tooltip placement="bottom" title="搜索大学信息">
            <a-button type="default" @click="goToUniversitySearch" size="large" shape="circle">
              <template #icon><SearchOutlined /></template>
            </a-button>
          </a-tooltip>

          <!-- 搜索专业按钮 -->
          <a-tooltip placement="bottom" title="搜索专业信息">
            <a-button type="primary" @click="goToMajorSearch" size="large" shape="circle">
              <template #icon ><SearchOutlined /></template>
            </a-button>
          </a-tooltip>

          <!-- AI对话按钮 -->
          <a-tooltip placement="bottom" title="AI智能咨询">
            <a-button type="default" @click="goToAIChat" size="large" shape="circle">
              <template #icon><MessageOutlined /></template>
            </a-button>
          </a-tooltip>

          <!-- 阅读政策按钮 -->
          <a-tooltip placement="bottom" title="查看招生政策">
            <a-button type="primary" @click="goToPolicy" size="large" shape="circle">
              <template #icon><FileTextOutlined /></template>
            </a-button>
          </a-tooltip>

          <!-- 智能推荐按钮 -->
          <a-tooltip placement="bottom" title="智能志愿推荐" >
            <a-button type="text" @click="goToRecommendation" size="large" shape="circle">
              <template #icon><BulbOutlined /></template>
            </a-button>
          </a-tooltip>

          <!-- 视频咨询按钮 -->
          <a-tooltip placement="bottom" title="家庭视频商讨">
            <a-button type="primary" @click="goToVideoCall" size="large" shape="circle">
              <template #icon><VideoCameraOutlined /></template>
            </a-button>
          </a-tooltip>
        </div>
      </template>
    </a-page-header>

    <a-spin :spinning="loading">
      <div class="volunteer-content">
        <!-- 左侧用户信息 -->
        <div class="left-column">
          <a-card :bordered="false" class="user-card">
            <div class="user-info">
              <a-avatar :size="64" class="user-avatar">
                {{ userInitial }}
              </a-avatar>
              <div class="user-details">
                <h3>{{ userInfo.displayName || userInfo.username }}</h3>
                <p>{{ userInfo.email }}</p>
                <p>用户ID: {{ userInfo.id }}</p>
              </div>
            </div>
          </a-card>

          <a-card title="填报说明" :bordered="false" class="instruction-card">
            <p>1. 最多可填报45个志愿</p>
            <p>2. 每个志愿包含一个学校和一个专业</p>
            <p>3. 请为每个志愿选择优先级（0-最高，3-最低）</p>
            <p>4. 系统将按优先级顺序处理志愿</p>
            <p>5. 填报完成后可导出为Excel文件</p>
          </a-card>
          <a-card title="优先级说明" :bordered="false" class="priority-card">
            <div class="priority-list">
              <div v-for="option in priorityOptions" :key="option.value" class="priority-item">
                <span class="priority-badge" :class="`priority-${option.value}`">{{ option.value }}</span>
                <span class="priority-label">{{ option.label }}</span>
              </div>
            </div>
          </a-card>
        </div>

        <!-- 中间志愿填报区域 -->
        <div class="center-column">
          <a-card title="志愿填报表" :bordered="false" class="volunteer-card">
            <div class="table-header">
              <div class="col-order">序号</div>
              <div class="col-school">院校信息</div>
              <div class="col-major">专业信息</div>
              <div class="col-priority">优先级</div>
              <div class="col-actions">操作</div>
            </div>

            <div v-for="(volunteer, index) in sortedVolunteers" :key="volunteer.id" class="volunteer-item">
              <div class="col-order">
                <span class="order-number">{{ index + 1 }}</span>
              </div>

              <div class="col-school">
                <a-input
                    v-model:value="volunteer.school_name"
                    placeholder="请输入院校名称"
                    class="school-input"
                />
              </div>

              <div class="col-major">
                <a-input
                    v-model:value="volunteer.major_name"
                    placeholder="请输入专业名称"
                    class="major-input"
                />
              </div>

              <div class="col-priority">
                <a-select
                    v-model:value="volunteer.priority"
                    :options="priorityOptions"
                    style="width: 100%"
                />
              </div>

              <div class="col-actions">
                <a-button
                    type="danger"
                    :icon="h(DeleteOutlined)"
                    shape="circle"
                    @click="removeVolunteer(volunteer.id)"
                />
              </div>
            </div>

            <div class="add-section">
              <a-button
                  type="primary"
                  :icon="h(PlusOutlined)"
                  class="add-button"
                  :disabled="!canAddMore"
                  @click="addVolunteer"
              >
                添加志愿（{{ remainingSlots }}个剩余）
              </a-button>
              <div class="add-tips">
                最多可填报 {{ MAX_VOLUNTEERS }} 个志愿
              </div>
            </div>
          </a-card>
        </div>

        <!-- 右侧统计信息 -->
        <div class="right-column">
          <a-card title="填报统计" :bordered="false" class="stats-card">
            <div class="stats-content">
              <div class="stat-item">
                <label>已选院校：</label>
                <span class="stat-value">{{ selectedSchoolCount }}</span>
              </div>
              <div class="stat-item">
                <label>已填专业：</label>
                <span class="stat-value">{{ selectedMajorCount }}</span>
              </div>
              <div class="stat-item">
                <label>剩余名额：</label>
                <span class="stat-value">{{ remainingSlots }}</span>
              </div>
              <div class="stat-item">
                <label>最高优先级：</label>
                <span class="stat-value">{{ highestPriorityCount }}</span>
              </div>
            </div>
          </a-card>

          <a-card title="操作面板" :bordered="false" class="action-card">
            <div class="action-buttons">
              <a-button
                  type="primary"
                  :icon="h(SaveOutlined)"
                  size="large"
                  block
                  @click="saveDraft"
                  class="action-btn"
              >
                保存草稿
              </a-button>
              <a-button
                  type="default"
                  :icon="h(FileExcelOutlined)"
                  size="large"
                  block
                  @click="exportToExcel"
                  class="action-btn"
              >
                导出Excel
              </a-button>
              <a-button
                  type="dashed"
                  :icon="h(DeleteOutlined)"
                  size="large"
                  block
                  @click="resetForm"
                  class="action-btn"
              >
                重置填报
              </a-button>
            </div>
          </a-card>

        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  BulbOutlined,
  VideoCameraOutlined,
  DeleteOutlined,
  PlusOutlined,
  SaveOutlined,
  FileExcelOutlined,
  SearchOutlined,
  MessageOutlined,
  FileTextOutlined
} from '@ant-design/icons-vue'
import { useUserStore } from '@/utils/auth.js'
import * as XLSX from 'xlsx'
import axios from "axios";

const router = useRouter()
const { getUser } = useUserStore()

// 常量定义
const MAX_VOLUNTEERS = 45

// 状态变量
const loading = ref(false)
const userInfo = ref({})
const volunteers = ref([])

// 优先级选项
const priorityOptions = [
  { value: 0, label: '最高优先级' },
  { value: 1, label: '高优先级' },
  { value: 2, label: '中优先级' },
  { value: 3, label: '低优先级' }
]

// 计算属性
const userInitial = computed(() => {
  if (userInfo.value.displayName)
    return userInfo.value.displayName.charAt(0).toUpperCase()
  if (userInfo.value.username) return userInfo.value.username.charAt(0).toUpperCase()
  return 'U'
})
const remainingSlots = computed(() => MAX_VOLUNTEERS - volunteers.value.length)
const canAddMore = computed(() => remainingSlots.value > 0)
const selectedSchoolCount = computed(() => volunteers.value.filter(v => v.schoolName && v.schoolName.trim()).length)
const selectedMajorCount = computed(() => volunteers.value.filter(v => v.major && v.major.trim()).length)
const highestPriorityCount = computed(() => volunteers.value.filter(v => v.priority === 0).length)
const sortedVolunteers = computed(() => [...volunteers.value].sort((a, b) => a.priority - b.priority)) // 按优先级排序的志愿列表

// 生命周期
onMounted(async () => {
  await loadUserInfo()
  await loadVolunteers()
})

// 方法
const loadUserInfo = async () => {
  const userLocal = getUser()
  if (userLocal) {
    userInfo.value = userLocal
  } else {
    message.error('未获取到用户信息')
    router.push('/login')
  }
}

const loadVolunteers = async () => {
  loading.value = true
  try {
    const response = await axios.get(`/gapi/user/${userInfo.value.id}/choices`)

    // 检查响应结构
    console.log('API响应:', response.data)

    // 根据实际API响应结构调整
    if (response.data && response.data.code === 200) {
      volunteers.value = response.data.data || []
    } else if (Array.isArray(response.data)) {
      volunteers.value = response.data
    } else {
      volunteers.value = []
    }

  } catch (error) {
    message.error('加载志愿数据失败')
    console.error('加载志愿数据失败:', error)
    volunteers.value = [] // 确保有默认值
  } finally {
    loading.value = false // 确保loading被关闭
  }
}

const addVolunteer = () => {
  if (canAddMore.value) {
    volunteers.value.push({
      id: Date.now(),
      schoolName: '',
      major: '',
      priority: 3 // 默认低优先级
    })
  }
}

const removeVolunteer = (id) => {
  const index = volunteers.value.findIndex(v => v.id === id)
  if (index !== -1) volunteers.value.splice(index, 1)
}

const saveDraft = async () => {
  loading.value = true
  try {
    // 1. 先删除所有现有的志愿
    await axios.delete(`/gapi/user/${userInfo.value.id}/choices/all`)

    // 2. 批量保存新的志愿数据
    const payload = volunteers.value.map((item, index) => ({
      school_name: item.school_name,
      major_name: item.major_name, // 注意字段名改为 majorName
      priority: item.priority
    }))

    // 3. 逐个创建新的志愿（因为后端是单个创建的接口）
    for (const item of payload) {
      if (item.school_name.trim() === '' && item.major_name.trim() === '') {
        message.warning('请填写完整的志愿信息')

      }
      await axios.post(`/gapi/user/${userInfo.value.id}/choices`, item)
    }

    message.success('保存成功')
  } catch (error) {
    message.error('保存失败')
    console.error('保存失败:', error)
  } finally {
    loading.value = false
  }
}

const exportToExcel = () => {
  if (volunteers.value.length === 0) {
    message.warning('没有可导出的数据')
    return
  }

  // 准备数据
  const data = sortedVolunteers.value.map((vol, index) => {
    const priorityLabel = priorityOptions.find(opt => opt.value === vol.priority)?.label || '未知'

    return {
      序号: index + 1,
      院校名称: vol.schoolName,
      专业名称: vol.major,
      优先级: priorityLabel
    }
  })

  // 创建工作簿
  const worksheet = XLSX.utils.json_to_sheet(data)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, '志愿填报')

  // 生成文件名
  const fileName = `志愿填报_${new Date().toLocaleDateString().replace(/\//g, '-')}.xlsx`

  // 导出文件
  XLSX.writeFile(workbook, fileName)
  message.success('文件导出成功')
}

const resetForm = () => {
  volunteers.value = []
  message.info('已重置填报内容')
}

const goToRecommendation = () => router.push('/recommends')
const goToVideoCall = () => router.push('/user-friend')
const goToUniversitySearch = () => router.push('/allSchool')
const goToMajorSearch = () => router.push('/professional')
const goToAIChat = () => router.push('/aismartsel')
const goToPolicy = () => router.push('/policy')
</script>

<style scoped>
.volunteer-container {
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

.volunteer-content {
  display: grid;
  grid-template-columns: 300px 1fr 300px;
  gap: 20px;
}

.left-column, .center-column, .right-column {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.user-card, .instruction-card, .volunteer-card, .stats-card, .action-card, .priority-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user-avatar {
  background-color: #1890ff;
  color: white;
  font-weight: bold;
  flex-shrink: 0;
}

.user-details h3 {
  margin: 0 0 5px 0;
  font-size: 16px;
  color: #262626;
}

.user-details p {
  margin: 0;
  color: #8c8c8c;
  font-size: 14px;
}

.instruction-card p {
  margin-bottom: 10px;
  color: #595959;
}

.table-header {
  display: flex;
  background: #fafafa;
  padding: 16px;
  font-weight: 500;
  border-radius: 6px;
  margin-bottom: 10px;
}

.col-order { width: 10%; }
.col-school { width: 30%; }
.col-major { width: 30%; }
.col-priority { width: 20%; }
.col-actions { width: 10%; text-align: center; }

.volunteer-item {
  display: flex;
  align-items: center;
  padding: 16px;
  margin-bottom: 10px;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
  background: white;
  transition: all 0.3s;
}

.volunteer-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
}

.order-number {
  font-weight: 500;
  color: #262626;
  text-align: center;
}

.school-input, .major-input {
  width: 100%;
}

.add-section {
  margin-top: 20px;
  text-align: center;
  padding: 16px;
  border-top: 1px solid #f0f0f0;
}

.add-button {
  height: 40px;
  padding: 0 20px;
}

.add-tips {
  color: #8c8c8c;
  font-size: 14px;
  margin-top: 8px;
}

.stats-content {
  padding: 0 10px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  margin: 12px 0;
  font-size: 14px;
}

.stat-value {
  color: #1890ff;
  font-weight: 500;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-btn {
  height: 45px;
}

.priority-list {
  padding: 0 10px;
}

.priority-item {
  display: flex;
  align-items: center;
  margin: 10px 0;
}

.priority-badge {
  display: inline-block;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  text-align: center;
  line-height: 24px;
  font-weight: bold;
  margin-right: 10px;
  color: white;
}

.priority-0 {
  background-color: #ff4d4f;
}

.priority-1 {
  background-color: #ff7a45;
}

.priority-2 {
  background-color: #ffa940;
}

.priority-3 {
  background-color: #52c41a;
}

.priority-label {
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .volunteer-content {
    grid-template-columns: 1fr 1fr;
  }

  .center-column {
    grid-column: span 2;
    order: 3;
  }
}

@media (max-width: 768px) {
  .volunteer-content {
    grid-template-columns: 1fr;
  }

  .center-column {
    grid-column: span 1;
  }

  .table-header {
    display: none;
  }

  .volunteer-item {
    flex-wrap: wrap;
    position: relative;
  }

  .col-order {
    position: absolute;
    left: 10px;
    top: 10px;
  }

  .col-school, .col-major {
    width: 100%;
    margin-bottom: 12px;
  }

  .col-priority {
    width: 60%;
  }

  .col-actions {
    position: absolute;
    right: 10px;
    top: 10px;
  }
}

.header-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
</style>
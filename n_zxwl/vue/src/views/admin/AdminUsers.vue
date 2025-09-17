<template>
  <div class="admin-users-container">
    <!-- 顶部导航栏 -->
    <div class="top-navbar">
      <div class="nav-content">
        <div class="brand-section">
          <img src="../../assets/zxwllogo.png" alt="Logo" class="page-logo">
          <h1 class="brand-title">智选未来·用户管理中心</h1>
        </div>
        <div class="action-buttons">
          <el-button
              type="primary"
              icon="Refresh"
              circle
              class="refresh-btn"
              @click="fetchUsers"
          />
          <el-button
              type="info"
              icon="ArrowLeft"
              class="back-btn"
              @click="handleGoBack"
              circle
          >
          </el-button>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="filter-section">
      <el-card shadow="never" class="filter-card">
        <div class="filter-content">
          <el-input
              v-model="searchForm.email"
              placeholder="搜索邮箱"
              prefix-icon="Search"
              class="search-input"
              @keyup.enter="handleSearch"
          />
          <el-input
              v-model="searchForm.keyword"
              placeholder="搜索用户名"
              prefix-icon="Search"
              class="search-input"
              @keyup.enter="handleSearch"
          />
          <div class="filter-group">
            <el-select
                v-model="searchForm.gender"
                placeholder="性别筛选"
                clearable
                class="filter-select"
            >
              <el-option label="男" value="1" />
              <el-option label="女" value="2" />
              <el-option label="未知" value="3" />
            </el-select>

            <el-date-picker
                v-model="searchForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                class="date-picker"
            />

            <el-button
                type="primary"
                @click="handleSearch"
                class="search-btn"
            >
              搜索
            </el-button>
            <el-button @click="resetSearch">
              重置
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 用户表格 -->
    <div class="table-section">
      <el-card shadow="never">
        <template #header>
          <div class="table-header">
            <span class="table-title">用户列表</span>
            <span class="table-total">共 {{ pagination.total }} 个用户</span>
          </div>
        </template>

        <el-table
            :data="userList"
            v-loading="loading"
            style="width: 100%"
            :row-class-name="tableRowClassName"
            @sort-change="handleSortChange"
            current-row-key="''"
            :lazy="false"
            :virtual-scroll="false"
        >
          <!-- 用户信息列 -->
          <el-table-column label="用户信息" min-width="250">
            <template #default="{ row }">
              <div class="user-info">
                <el-avatar
                    :size="50"
                    :src="row.avatarUrl || defaultAvatar(row.username)"
                    class="user-avatar"
                />
                <div class="user-details">
                  <div class="user-name">
                    <span class="name">{{ row.username }}</span>
                    <el-tag
                        v-if="row.isOnline"
                        size="small"
                        type="success"
                        effect="light"
                    >
                      在线
                    </el-tag>
                  </div>
                  <div class="user-email">
                    <el-icon><Message /></el-icon>
                    <span>{{ row.email }}</span>
                  </div>
                  <div class="user-display-name">
                    {{ row.displayName || '未设置昵称' }}
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>

          <!-- 基本信息列 -->
          <el-table-column label="基本信息" width="180">
            <template #default="{ row }">
              <div class="basic-info">
                <div class="info-item">
                  <el-icon><User /></el-icon>
                  <span>{{ formatGender(row.gender) }}</span>
                </div>
                <div class="info-item">
                  <el-icon><Calendar /></el-icon>
                  <span>{{ row.birthYear || '未设置' }}</span>
                </div>
                <div class="info-item">
                  <el-icon><Location /></el-icon>
                  <span>{{ row.location || '未设置位置' }}</span>
                </div>
              </div>
            </template>
          </el-table-column>

          <!-- 状态信息列 -->
          <el-table-column label="状态信息" width="150">
            <template #default="{ row }">
              <div class="status-info">
                <el-tag
                    :type="row.isOnline ? 'success' : 'info'"
                    effect="plain"
                    class="status-tag"
                >
                  {{ row.isOnline ? '在线' : '离线' }}
                </el-tag>
                <div class="last-online">
                  {{ formatLastOnline(row.lastOnlineTime) }}
                </div>
              </div>
            </template>
          </el-table-column>

          <!-- 注册时间 -->
          <el-table-column prop="createdAt" label="注册时间" width="180" sortable>
            <template #default="{ row }">
              <div class="time-info">
                <el-icon><Clock /></el-icon>
                <span>{{ formatDateTime(row.createdAt) }}</span>
              </div>
            </template>
          </el-table-column>

          <!-- 操作列 -->
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button
                    size="medium"
                    type="primary"
                    icon="View"
                    circle
                    @click="handleView(row)"
                >
                </el-button>
                <el-button
                    size="medium"
                    type="warning"
                    icon="Edit"
                    circle
                    @click="handleEdit(row)"
                >
                </el-button>
                <el-button
                    size="medium"
                    type="danger"
                    icon="Delete"
                    circle
                    @click="handleDelete(row)"
                >
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination">
          <el-pagination
              v-model:current-page="pagination.current"
              v-model:page-size="pagination.size"
              :page-sizes="[10, 20, 50, 100]"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 用户详情对话框 -->
    <user-detail-dialog
        v-model="detailVisible"
        :user="currentUser"
    />

    <!-- 编辑用户对话框 -->
    <user-edit-dialog
        v-model="editVisible"
        :user="currentUser"
        @success="handleEditSuccess"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, ElTable, ElTableColumn } from 'element-plus'
import {
  Message,
  User,
  Calendar,
  Location,
  Clock,
  View,
  Edit,
  Delete,
  Refresh,
  ArrowLeft,
  Search
} from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import UserDetailDialog from '@/components/admin/UserDetailDialog.vue'
import UserEditDialog from '@/components/admin/UserEditDialog.vue'

const router = useRouter()
const loading = ref(false)
const detailVisible = ref(false)
const editVisible = ref(false)
const currentUser = ref(null)

const userList = ref([])
const pagination = reactive({
  current: 1,
  size: 10,
  total: 0
})

const searchForm = reactive({
  email: '',
  keyword: '',
  gender: '',
  dateRange: []
})

// 加载用户列表
const fetchUsers = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.current,
      size: pagination.size,
      username: searchForm.keyword,
      email: searchForm.email,
      gender: searchForm.gender,
      startTime: searchForm.dateRange?.[0] ? formatDate(searchForm.dateRange[0]) : '',
      endTime: searchForm.dateRange?.[1] ? formatDate(searchForm.dateRange[1]) : ''
    }

    const response = await axios.get('/gapi/users', { params })

    // 处理嵌套的 data 结构
    const responseData = response.data.data

    if (responseData && Array.isArray(responseData.data)) {
      // 格式: {data: {data: [], total: 0}}
      userList.value = responseData.data || []
      pagination.total = responseData.total || 0
    } else if (Array.isArray(responseData)) {
      // 格式: {data: []}
      userList.value = responseData || []
      pagination.total = response.data.total || 0
    } else {
      userList.value = []
      pagination.total = 0
      console.warn('Unexpected response format:', response.data)
    }
  } catch (error) {
    console.error('Fetch users error:', error)
    ElMessage.error('获取用户数据失败')
    userList.value = []
  } finally {
    loading.value = false
  }
}

// 搜索用户
const handleSearch = () => {
  pagination.current = 1
  fetchUsers()
}

// 重置搜索
const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.gender = ''
  searchForm.dateRange = []
  handleSearch()
}

// 查看用户详情
const handleView = (user) => {
  currentUser.value = user
  detailVisible.value = true
}

// 编辑用户
const handleEdit = (user) => {
  currentUser.value = { ...user }
  editVisible.value = true
}

// 删除用户
const handleDelete = async (user) => {
  try {
    await ElMessageBox.confirm(
        `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
          center: true,
        }
    )

    // 直接使用axios发送请求
    await axios.delete(`/gapi/user/${user.id}`)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 返回首页
const handleGoBack = () => router.back()

// 格式化性别
const formatGender = (gender) => {
  const genderMap = { 0: '未知', 1: '男', 2: '女' }
  return genderMap[gender] || '未知'
}

// 格式化日期时间
const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 格式化最后在线时间
const formatLastOnline = (time) => {
  if (!time) return '从未在线'
  return new Date(time).toLocaleDateString('zh-CN')
}

// 生成默认头像
const defaultAvatar = (username) => {
  const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFA726']
  const char = username?.charAt(0)?.toUpperCase() || 'U'
  const color = colors[char.charCodeAt(0) % colors.length]

  return `data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">
    <rect width="100" height="100" rx="50" fill="${color}"/>
    <text x="50%" y="52%" dominant-baseline="middle" text-anchor="middle" font-family="Arial" font-size="40" fill="#fff">${char}</text>
  </svg>`
}

// 表格行样式
const tableRowClassName = ({ rowIndex }) => rowIndex % 2 === 1 ? 'even-row' : ''

// 分页大小改变
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.current = 1
  fetchUsers()
}

// 当前页改变
const handleCurrentChange = (page) => {
  pagination.current = page
  fetchUsers()
}

// 排序改变
const handleSortChange = ({ prop, order }) => {
  // 这里可以实现排序逻辑
  console.log('排序:', prop, order)
}

// 格式化日期
const formatDate = (date) => date.toISOString().split('T')[0]

// 编辑成功回调
const handleEditSuccess = () => {
  editVisible.value = false
  fetchUsers()
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped lang="scss">
.admin-users-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #f8fafc 100%);
  padding: 20px;
}

.top-navbar {
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  border-radius: 12px;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
}

.brand-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-logo {
  width: 60px;
  height: 60px;
  top: 0;
  border-radius: 8px;
}

.brand-title {
  margin: 0;
  color: white;
  font-size: 24px;
  font-weight: 600;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;

  &:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: rotate(360deg);
    transition: transform 0.5s ease;
  }
}

.back-btn {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.3);
  color: white;

  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
}

.filter-section {
  margin-bottom: 20px;
}

.filter-card {
  border: none;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.filter-content {
  display: flex;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 100px;
}

.filter-group {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.filter-select,
.date-picker {
  width: 180px;
}

.table-section {
  margin-bottom: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-title {
  font-size: 18px;
  font-weight: 600;
  color: #2d3748;
}

.table-total {
  color: #718096;
  font-size: 14px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  flex-shrink: 0;
}

.user-details {
  flex: 1;
}

.user-name {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.name {
  font-weight: 600;
  color: #2d3748;
}

.user-email {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #718096;
  font-size: 14px;
  margin-bottom: 2px;
}

.user-display-name {
  color: #a0aec0;
  font-size: 13px;
}

.basic-info {
  .info-item {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 6px;
    color: #4a5568;
    font-size: 14px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.status-tag {
  width: fit-content;
}

.last-online {
  color: #a0aec0;
  font-size: 12px;
}

.time-info {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #4a5568;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

// 响应式设计
@media (max-width: 768px) {
  .nav-content {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .filter-content {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input {
    min-width: auto;
  }

  .filter-group {
    justify-content: center;
  }

  .filter-select,
  .date-picker {
    width: 100%;
  }
}

// 表格行交替背景色
::v-deep(.even-row) {
  background-color: #fafafa;
}

::v-deep(.el-table__row:hover) {
  background-color: #f0f9ff !important;
}
</style>
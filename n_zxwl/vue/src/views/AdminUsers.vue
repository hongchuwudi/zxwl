<template>
  <div class="user-management-container">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <!-- 顶部导航栏 -->
    <div class="top-navbar">
      <div class="nav-content">
        <div class="brand-section brand-section-center">
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
          >返回</el-button>
        </div>
      </div>
    </div>

    <!-- 用户表格 -->
    <div class="table-container">
      <el-table
          v-loading="loading"
          :data="currentUsers"
          style="width: 100%"
          class="user-table"
          :header-cell-style="headerStyle"
          :row-style="rowStyle"
      >
        <el-table-column label="头像" width="100" align="center">
          <template #default="{ row }">
            <div class="avatar-wrapper">
              <el-avatar
                  :size="50"
                  :src="row.picture || defaultAvatar(row.name)"
                  class="hover-effect"
              />
              <div class="avatar-glow"></div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="name" label="姓名" width="120" sortable>
          <template #default="{ row }">
            <span class="name-highlight">{{ row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="email" label="邮箱" min-width="200">
          <template #default="{ row }">
            <div class="email-cell">
              <el-icon class="email-icon"><Message /></el-icon>
              <span class="email-text">{{ row.email }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="性别" width="100">
          <template #default="{ row }">
            <el-tag
                :type="row.sex === 1 ? 'primary' : row.sex === 2 ? 'danger' : 'info'"
                class="gender-tag"
                effect="dark"
            >
              {{ row.sex === 1 ? '男' : row.sex === 2 ? '女' : '未知' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="graduate" label="毕业年份" width="120" sortable>
          <template #default="{ row }">
            <div class="graduate-cell">
              <el-icon class="calendar-icon"><Calendar /></el-icon>
              <span class="graduate-text">{{ row.graduate || '未设置' }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="地址" min-width="200">
          <template #default="{ row }">
            <div class="address-cell">
              <el-icon class="location-icon"><Location /></el-icon>
              <span class="address-text">{{ formatAddress(row.address) || '地址未填写' }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
                type="danger"
                icon="Delete"
                circle
                class="delete-btn"
                @click="handleDelete(row.email)"
            />
          </template>
        </el-table-column>
      </el-table>
      <!-- 分页组件 -->
      <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[7, 10, 20]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="users.length"
          class="fixed-pagination"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Message,
  Calendar,
  Location,
  Delete,
  Refresh,
  ArrowLeft
} from '@element-plus/icons-vue'
import { provinceAndCityData } from 'element-china-area-data'
import axios from 'axios'
import { useRouter } from 'vue-router' // 导入 useRouter 函数

const users = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(7)
const router = useRouter() // 获取路由实例

// 计算当前页要显示的用户
const currentUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return users.value.slice(start, end)
})

// 地址数据映射
const addressMap = new Map()
provinceAndCityData.forEach(province => {
  addressMap.set(province.value, province.label)
  province.children?.forEach(city => {
    addressMap.set(city.value, city.label)
  })
})

// 格式化地址
const formatAddress = (code) => {
  if (!code) return ''
  try {
    const [provinceCode, cityCode] = code.split('/')
    return `${addressMap.get(provinceCode) || '未知省份'} / ${addressMap.get(cityCode) || '未知城市'}`
  } catch {
    return '地址格式异常'
  }
}

// 生成默认头像
const defaultAvatar = (name) => {
  const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4']
  const char = name?.charAt(0) || '?'
  return `data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">
    <rect width="100" height="100" rx="50" fill="${colors[Math.floor(Math.random()*colors.length)]}"/>
    <text x="50%" y="52%" dominant-baseline="middle" text-anchor="middle" font-family="Arial" font-size="60" fill="#fff">${char}</text>
  </svg>`
}

// 获取用户数据
const fetchUsers = async () => {
  try {
    loading.value = true
    const { data } = await axios.get('http://127.0.0.1:8792/profile/list')
    users.value = data.data.map(user => ({
      ...user,
      sex: user.sex ?? 0,
      graduate: user.graduate ?? null,
      address: user.address ?? ''
    }))
  } catch (error) {
    ElMessage.error('获取用户数据失败')
  } finally {
    loading.value = false
  }
}

// 删除用户
const handleDelete = async (email) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
      center: true,
    })

    // 发送POST请求，参数放在请求体中
    const { data } = await axios.post('http://127.0.0.1:8792/profile/delete', {
      email: email
    })

    if (data.code === 0) {
      const logData = {
        "email": localStorage.getItem('userEmail'),
        "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
        "operation": "删除用户"
      };
      const logResponse = await axios.post("gapi/log", logData, {
        headers: {
          "Content-Type": "application/json"
        }
      });
      ElMessage.success('删除成功')
      // 刷新列表时保持当前分页状态
      await fetchUsers()
      // 如果当前页无数据且不是第一页，自动返回前一页
      if (currentUsers.value.length === 0 && currentPage.value > 1) {
        currentPage.value--
      }
    } else if (data.code === 1002) {
      ElMessage.warning('用户不存在')
    } else {
      ElMessage.error(`删除失败，错误码：${data.code}`)
    }
  } catch (error) {
    // 处理请求错误
    if (error.response) {
      // 服务器返回4xx/5xx错误
      ElMessage.error(`服务器错误：${error.response.status}`)
    } else if (error.request) {
      // 请求无响应
      ElMessage.error('网络连接异常，请检查网络')
    } else {
      // 用户取消删除时不提示错误
      if (error !== 'cancel') {
        ElMessage.error(`操作失败：${error.message}`)
      }
    }
  }
}

// 返回处理
const handleGoBack = () => {
  router.push('/zxwl')
}

// 表格样式配置
const headerStyle = () => ({
  backgroundColor: '#f8f9fa',
  color: '#2c3e50',
  fontWeight: 600,
  fontSize: '14px',
  borderBottom: '2px solid #e9ecef'
})

const rowStyle = () => ({
  cursor: 'pointer',
  transition: 'all 0.3s ease',
  backgroundColor: '#ffffff'
})

// 每页数量改变时的处理函数
const handleSizeChange = (newSize) => {
  pageSize.value = newSize
}

// 当前页码改变时的处理函数
const handleCurrentChange = (newPage) => {
  currentPage.value = newPage
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-management-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  position: relative;
  padding-bottom: 60px; /* 为分页组件留出空间 */
}

.top-navbar {
  background: linear-gradient(135deg, #1a73e8 0%, #0d47a1 100%);
  padding: 1rem 2rem;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
}

.nav-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 页面 logo 样式 */
.page-logo {
  position: absolute;
  top: -0.55rem;
  left: 1rem;
  width: 100px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
}
.brand-section {
  flex: 1;
  text-align: center;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
}

.brand-title {
  margin: 0;
  color: white;
  font-size: 1.8rem;
  letter-spacing: 0.5px;
  text-shadow: 1px 1px 2px rgba(0,0,0,0.1);
  white-space: nowrap;
}

.action-buttons {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-left: auto;
}

.refresh-btn {
  background: rgba(255,255,255,0.1);
  border-color: rgba(255,255,255,0.2);
  transition: all 0.3s ease;
}
.brand-section-center {
  justify-content: center;
}
.refresh-btn:hover {
  transform: rotate(360deg);
  background: rgba(255,255,255,0.2);
}

.back-btn {
  background: rgba(255,255,255,0.1);
  border-color: rgba(255,255,255,0.2);
  color: white;
  padding: 0.75rem 1.5rem;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255,255,255,0.2);
  transform: translateX(-5px);
}

.table-container {
  max-width: 1400px;
  margin: 2rem auto;
  padding: 0 2rem;
}

.user-table {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(0,0,0,0.05);
  background: white;
}

.avatar-wrapper {
  position: relative;
  display: inline-block;
}

.avatar-glow {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  box-shadow: 0 0 12px rgba(64,158,255,0.3);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.hover-effect:hover .avatar-glow {
  opacity: 1;
}

.name-highlight {
  color: #2c3e50;
  font-weight: 600;
  position: relative;
}

.name-highlight::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: #1a73e8;
  transition: width 0.3s ease;
}

.el-table__row:hover .name-highlight::after {
  width: 100%;
}

.email-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.email-icon {
  color: #1a73e8;
  font-size: 18px;
}

.graduate-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.calendar-icon {
  color: #4CAF50;
  font-size: 18px;
}

.address-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.location-icon {
  color: #FF9800;
  font-size: 18px;
}

.gender-tag {
  font-weight: 500;
  letter-spacing: 0.5px;
}

.delete-btn {
  transition: all 0.3s ease;
}

.delete-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(255,58,58,0.3);
}

.el-table__row {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.el-table__row:hover {
  transform: translateX(8px);
  background: #f8f9fa !important;
}

.fixed-pagination {
  left: 50%;
  transform: translateX(-50%);
  position: fixed;        /* 固定定位 */
  bottom: 30px;           /* 距离底部20px */
}

@media (max-width: 768px) {
  .nav-content {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }

  .brand-title {
    font-size: 1.4rem;
  }

  .table-container {
    padding: 0 1rem;
  }

  .user-table {
    border-radius: 8px;
  }
}
</style>
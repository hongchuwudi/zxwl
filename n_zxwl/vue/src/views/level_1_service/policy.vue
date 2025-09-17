<template>
  <!-- 头部导航栏 -->
  <nav class="header-nav">
    <img src="../../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <div class="nav-content">
      <h1 class="logo">智选未来·最新高考政策解读</h1>
      <div class="button-group">
        <el-button class="back-prev-btn" @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回上一页
        </el-button>
      </div>
    </div>
  </nav>

  <!-- 政策容器 -->
  <div class="policy-container">
    <!-- 政策列表 -->
    <div class="policy-list">
      <el-card
          v-for="policy in paginatedData"
          :key="policy.id"
          class="policy-card"
          shadow="hover"
          @click="showPolicyDetail(policy.title)"
      >
        <template #header>
          <div class="card-header">
            <h3 class="title">{{ policy.title }}</h3>
            <div class="hover-indicator"></div>
          </div>
        </template>
        <p class="preview-content">{{ truncateContent(policy.content) }}</p>
        <div class="card-footer">
          <el-button type="primary" text>查看详情</el-button>
        </div>
      </el-card>
    </div>

    <!-- 分页控件 -->
    <div class="pagination">
      <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          layout="prev, pager, next, jumper"
          :total="policyList.length"
          :hide-on-single-page="true"
      />
    </div>

    <!-- 政策详情弹窗 -->
    <el-dialog
        v-model="showDetail"
        :title="selectedPolicy?.title"
        width="70%"
        top="10vh"
        destroy-on-close
        class="policy-detail-dialog"
    >
      <div class="modal-body" v-html="formatContent(selectedPolicy?.content)"></div>
      <template #footer>
        <span class="dialog-footer">
          <el-button v-if="isAdmin" type="danger" @click="confirmDeletePolicy">删除政策</el-button>
          <el-button @click="showDetail = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 删除确认对话框 -->
    <el-dialog
        v-model="showDeleteConfirm"
        title="确认删除"
        width="400px"
        class="delete-confirm-dialog"
    >
      <span>确定要删除「{{ selectedPolicy?.title }}」政策吗？此操作不可恢复。</span>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDeleteConfirm = false">取消</el-button>
          <el-button type="danger" @click="deletePolicy">确认删除</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 加载状态 -->
    <el-loading
        v-if="loading"
        :fullscreen="true"
        text="正在加载政策信息..."
        background-color="rgba(0, 0, 0, 0.4)"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useUserStore } from '@/utils/auth.js' // 请根据实际路径调整

const router = useRouter()
const userStore = useUserStore()
const policyList = ref([])
const currentPage = ref(1)
const pageSize = ref(6)
const showDeleteConfirm = ref(false)

// 计算属性
const totalPages = computed(() => Math.ceil(policyList.value.length / pageSize.value))
const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return policyList.value.slice(start, end)
})

// 检查是否为管理员
const isAdmin = computed(() => {
  const userEmail = userStore.userEmail.value || localStorage.getItem('userEmail')
  const userName = userStore.userName.value || localStorage.getItem('userName')
  console.log('userEmail:', userEmail)
  return userEmail === 'root@root.com' || userName === 'root' || userName === 'rootroot'
})

const goBack = () => router.back()

const selectedPolicy = ref(null)
const showDetail = ref(false)
const loading = ref(false)

// 获取政策列表
const fetchPolicies = async () => {
  try {
    loading.value = true
    const response = await axios.get('gapi/policy')
    if (response.data.code === 0) {
      policyList.value = response.data.data
    } else {
      ElMessage.error('获取政策列表失败')
    }
  } catch (error) {
    console.error('获取政策列表失败:', error)
    ElMessage.error('获取政策列表失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

// 获取政策详情
const showPolicyDetail = async (title) => {
  try {
    loading.value = true
    const response = await axios.post(
        'gapi/policy/search',
        { title },
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
    )

    if (response.data.code === 0 && response.data.data.length > 0) {
      selectedPolicy.value = response.data.data[0]
      showDetail.value = true
    } else {
      ElMessage.error('获取政策详情失败')
    }
  } catch (error) {
    console.error('获取政策详情失败:', error)
    ElMessage.error('获取政策详情失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 确认删除政策
const confirmDeletePolicy = () => {
  showDeleteConfirm.value = true
}

// 删除政策
const deletePolicy = async () => {
  try {
    loading.value = true
    const response = await axios.delete(`gapi/policy/${selectedPolicy.value.id}`)

    if (response.data.code === 0) {
      ElMessage.success('政策删除成功')
      // 从列表中移除已删除的政策
      policyList.value = policyList.value.filter(p => p.id !== selectedPolicy.value.id)
      showDetail.value = false
      showDeleteConfirm.value = false
    } else {
      ElMessage.error('政策删除失败')
    }
  } catch (error) {
    console.error('政策删除失败:', error)
    ElMessage.error('政策删除失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 格式化内容
const formatContent = content => content ? content.replace(/\r\n/g, '<br>').replace(/\n/g, '<br>') : ''
// 内容截断
const truncateContent = (content, length = 100) => content && content.length > length ? content.substring(0, length) + '...' : content

onMounted(async () => {
  const logData = {
    "email": localStorage.getItem('userEmail'),
    "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
    "operation": "用户查看高考政策"
  }

  await fetchPolicies()
})
</script>

<style scoped>
.header-nav {
  background: linear-gradient(135deg, #2c3e50, #3498db);
  box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  min-height: 71px;
}

.nav-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-logo {
  position: absolute;
  top: 0;
  left: 1rem;
  width: 100px;
  height: auto;
  z-index: 3;
}

.logo {
  color: white;
  font-size: 1.8rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  margin: 0;
  flex-grow: 1;
  text-align: center;
}

.button-group {
  display: flex;
  gap: 0.8rem;
}

/* 调整政策容器位置 */
.policy-container {
  position: fixed;
  left: 0;
  right: 0;
  top: 80px;
  bottom: 0;
  max-width: 1200px;
  margin: 0 auto;
  overflow-y: auto;
  padding: 20px;
  background-color: #f5f7fa;
}

.policy-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.policy-card {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  overflow: hidden;
}

.policy-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

::v-deep(.el-card__header) {
  padding: 1rem 1.2rem 0.5rem;
  border-bottom: 1px solid #eee;
  background: linear-gradient(to right, #f8fafc, #f1f5f9);
}

.card-header {
  position: relative;
}

.title {
  color: #2c3e50;
  font-size: 1.2rem;
  margin: 0;
  transition: color 0.3s ease;
  font-weight: 600;
}

.preview-content {
  color: #666;
  line-height: 1.6;
  font-size: 0.95rem;
  flex-grow: 1;
  padding: 0 1.2rem;
}

.hover-indicator {
  position: absolute;
  bottom: -0.5rem;
  left: 0;
  width: 0;
  height: 2px;
  background: #3498db;
  transition: width 0.3s ease;
}

.policy-card:hover .hover-indicator {
  width: 100%;
}

.card-footer {
  margin-top: 1rem;
  text-align: right;
  padding: 0 1.2rem 1rem;
}

/* 分页样式 */
.pagination {
  position: sticky;
  bottom: 0;
  background: white;
  padding: 1rem 0;
  display: flex;
  justify-content: center;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  border-radius: 12px;
  margin-top: 20px;
}

::v-deep(.policy-detail-dialog .el-dialog__headerbtn) {
  position: absolute;
  top: 20px;
  right: 25px; /* 增加右边距 */
  margin-top: 0;
}

/* 弹窗内容样式 */
.modal-body {
  line-height: 1.8;
  color: #444;
  max-height: 60vh;
  overflow-y: auto;
  padding: 1rem;
  font-size: 15px;
}

.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

/* 提高弹窗的z-index，使其显示在顶部栏上面 */
::v-deep(.el-overlay) {
  z-index: 1001 !important; /* 比导航栏的1000高 */
}

::v-deep(.policy-detail-dialog .el-dialog) {
  z-index: 1002 !important; /* 确保弹窗在最前面 */
  margin-top: 20px !important; /* 添加一点顶部间距 */
  border-radius: 12px;
  overflow: hidden;
}

::v-deep(.policy-detail-dialog .el-dialog__header) {
  background: linear-gradient(135deg, #2c3e50, #3498db);
  margin-right: 0;
  padding: 15px 20px;
}

::v-deep(.policy-detail-dialog .el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

::v-deep(.policy-detail-dialog .el-dialog__headerbtn) {
  top: 15px;
}

::v-deep(.policy-detail-dialog .el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 20px;
}

::v-deep(.policy-detail-dialog .el-dialog__body) {
  padding: 0;
}

::v-deep(.policy-detail-dialog .el-dialog__footer) {
  padding: 15px 20px;
  border-top: 1px solid #e8e8e8;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 删除确认对话框样式 */
::v-deep(.delete-confirm-dialog .el-dialog__header) {
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
}

::v-deep(.delete-confirm-dialog .el-dialog__title) {
  color: #f56c6c;
}

@media (max-width: 768px) {
  .nav-content {
    padding: 1rem;
    flex-direction: column;
    gap: 1rem;
  }

  .logo {
    font-size: 1.4rem;
    order: -1;
  }

  .button-group {
    width: 100%;
    justify-content: center;
  }

  .policy-container {
    top: 140px;
    padding: 15px;
  }

  .policy-list {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .page-logo {
    position: relative;
    margin: 0 auto;
    left: 0;
    top: 0;
  }

  /* 移动端弹窗调整 */
  ::v-deep(.policy-detail-dialog .el-dialog) {
    width: 90% !important;
    margin-top: 80px !important;
  }

  ::v-deep(.policy-detail-dialog .el-dialog__footer) {
    flex-direction: column-reverse;
  }

  ::v-deep(.policy-detail-dialog .el-button) {
    width: 100%;
    margin-left: 0 !important;
  }
}
</style>
<template>
  <!-- 头部导航栏 -->
  <nav class="header-nav">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <div class="nav-content">
      <h1 class="logo">智选未来·最新高考政策解读</h1>
      <button class="back-btn" @click="goBack">
        <span class="arrow"></span>
        返回首页
      </button>
    </div>
  </nav>

  <!-- 政策容器 -->
  <div class="policy-container">
    <!-- 政策列表 -->
    <div class="policy-list">
      <div
          v-for="policy in paginatedData"
          :key="policy.id"
          class="policy-card"
          @click="showPolicyDetail(policy.title)"
      >
        <div class="card-header">
          <h3 class="title">{{ policy.title }}</h3>
          <div class="hover-indicator"></div>
        </div>
        <p class="preview-content">{{ truncateContent(policy.content) }}</p>
      </div>
    </div>

    <!-- 分页控件 -->
    <div class="pagination">
      <button
          :disabled="currentPage === 1"
          @click="currentPage--"
          class="page-btn"
      >
        &lt;
      </button>
      <button
          v-for="page in totalPages"
          :key="page"
          @click="currentPage = page"
          :class="{ active: currentPage === page }"
          class="page-btn"
      >
        {{ page }}
      </button>
      <button
          :disabled="currentPage === totalPages"
          @click="currentPage++"
          class="page-btn"
      >
        &gt;
      </button>
    </div>

    <!-- 政策详情弹窗 -->
    <teleport to="body">
      <div v-if="showDetail" class="detail-modal">
        <div class="modal-content">
          <button class="close-btn" @click="showDetail = false">&times;</button>
          <h2 class="modal-title">{{ selectedPolicy?.title }}</h2>
          <div class="modal-body" v-html="formatContent(selectedPolicy?.content)"></div>
        </div>
      </div>
    </teleport>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <div class="loader"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const policyList = ref([])
const currentPage = ref(1)
const pageSize = ref(6)

// 计算属性
const totalPages = computed(() => Math.ceil(policyList.value.length / pageSize.value))
const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return policyList.value.slice(start, end)
})

// 导航方法
const goBack = () => {
  router.push('/zxwl')
}

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
    }
  } catch (error) {
    console.error('获取政策列表失败:', error)
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
    }
  } catch (error) {
    console.error('获取政策详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化内容
const formatContent = (content) => {
  return content.replace(/\r\n/g, '<br>').replace(/\n/g, '<br>')
}

// 内容截断
const truncateContent = (content, length = 100) => {
  return content.length > length ? content.substring(0, length) + '...' : content
}
onMounted(async () => {
  const logData = {
    "email": localStorage.getItem('userEmail'),
    "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
    "operation": "用户进入聊天室聊天"
  };
  const logResponse = await axios.post("gapi/log", logData, {
    headers: {
      "Content-Type": "application/json"
    }
  });
})
onMounted(() => {
  fetchPolicies()
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
  z-index: 1000;
}

.nav-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
}
.page-logo {
  position: absolute;
  top: -0rem;
  left: 1rem;
  width: 100px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
}
.logo {
  color: white;
  font-size: 1.8rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  text-align: center;
  left:180px;
  margin: 0;
}

.back-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid white;
  color: white;
  padding: 0.8rem 1.5rem;
  border-radius: 30px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  margin-left: auto
}

.back-btn:hover {
  background: white;
  color: #2c3e50;
  transform: translateY(-2px);
}

.arrow {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-left: 2px solid currentColor;
  border-bottom: 2px solid currentColor;
  transform: rotate(45deg);
  margin-right: 5px;
}

/* 调整政策容器位置 */
.policy-container {
  position: fixed;        /* 固定定位 */
  left: 0;
  right: 0;
  top: 50px;
  max-width: 1200px;
  justify-content: center;
  margin: 6rem auto 2rem;
  padding: 0 20px;
}

/* 新增分页样式 */
.pagination {
  position: fixed;        /* 固定定位 */
  bottom: -5px;           /* 距离底部20px */
  left: 0;
  right: 0;
  margin: 3rem 0;
  display: flex;
  justify-content: center;
  gap: 0.5rem;
}

.page-btn {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 8px;
  background: #f8f9fa;
  color: #2c3e50;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.page-btn:hover:not(:disabled) {
  background: #3498db;
  color: white;
  transform: translateY(-2px);
}

.page-btn.active {
  background: #2c3e50;
  color: white;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.policy-container {
  max-width: 1200px;
  margin: 2rem auto;
  padding: 0 20px;
}

.policy-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 2rem;
}

.policy-card {
  background: linear-gradient(145deg, #ffffff, #f8f9fa);
  border-radius: 15px;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.policy-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.2);
  background: linear-gradient(145deg, #f8f9fa, #ffffff);
}

.card-header {
  position: relative;
  margin-bottom: 1rem;
}

.title {
  color: #2c3e50;
  font-size: 1.2rem;
  margin-bottom: 0.5rem;
  transition: color 0.3s ease;
}

.preview-content {
  color: #666;
  line-height: 1.6;
  font-size: 0.95rem;
}

.hover-indicator {
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 0;
  height: 2px;
  background: #3498db;
  transition: width 0.3s ease;
}

.policy-card:hover .hover-indicator {
  width: 100%;
}

/* 详情弹窗样式 */
.detail-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(5px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 999;
}

.modal-content {
  margin-top: 60px;
  background: white;
  border-radius: 15px;
  padding: 2rem;
  max-width: 800px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  position: relative;
  animation: modalEnter 0.3s ease;
}

@keyframes modalEnter {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-title {
  color: #2c3e50;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
}

.modal-body {
  line-height: 1.8;
  color: #444;
  white-space: pre-wrap;
}

.close-btn {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  transition: color 0.3s ease;
}

.close-btn:hover {
  color: #e74c3c;
}

/* 加载动画 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.loader {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #3498db;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .nav-content {
    padding: 1rem;
  }

  .logo {
    font-size: 1.4rem;
  }

  .back-btn {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
  }

  .policy-container {
    margin-top: 5rem;
  }

  .pagination {
    flex-wrap: wrap;
  }
}
</style>
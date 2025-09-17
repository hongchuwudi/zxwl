<template>
  <div class="news-container">
    <!-- 查询条件区域 -->
    <el-card class="filter-card">
      <div class="filter-header">
        <h2 class="section-title">资讯论坛</h2>
        <div class="header-actions">
          <el-button type="success" @click="showAddDrawer = true" round circle size="large">
            <el-icon><Plus /></el-icon>
          </el-button>
          <el-button class="exit-button" type="primary" @click="handleExit" round circle size="large">
            <el-icon><Close /></el-icon>
          </el-button>
        </div>
      </div>
      <el-form :model="queryParams" label-width="100px" class="filter-form">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="标题关键词">
              <el-input v-model="queryParams.title" placeholder="输入标题关键词" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="内容关键词">
              <el-input v-model="queryParams.keywords" placeholder="输入内容关键词" clearable />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="资讯来源">
              <el-select
                  v-model="queryParams.from_source"
                  placeholder="选择或输入资讯来源"
                  clearable
                  filterable
                  allow-create
              >
                <el-option label="掌上高考" value="掌上高考" />
                <el-option label="教育部" value="教育部" />
                <el-option label="学校官网" value="学校官网" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="排序方式">
              <el-select v-model="queryParams.order_by" placeholder="选择排序方式" clearable>
                <el-option label="最新发布" value="publish_time" />
                <el-option label="最早发布" value="publish_time_asc" />
                <el-option label="最多阅读" value="news_num" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24" style="text-align: center">
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-container">
      <el-empty description="加载失败，请重试" />
      <el-button @click="fetchNews" type="primary">重新加载</el-button>
    </div>

    <!-- 空状态 -->
    <div v-else-if="newsList.length === 0" class="empty-container">
      <el-empty description="暂无相关资讯" />
    </div>

    <!-- 结果区域 -->
    <div v-else>
      <!-- 统计信息 -->
      <el-card class="stats-card">
        <div class="stats-info">
          <span>共找到 <strong>{{ total }}</strong> 条符合条件的资讯</span>
          <span>第 <strong>{{ queryParams.page }}</strong> 页 / 共 <strong>{{ totalPages }}</strong> 页</span>
        </div>
      </el-card>

      <!-- 资讯列表 -->
      <div class="news-list">
        <el-card
            v-for="news in newsList"
            :key="news.id"
            class="news-card"
            :class="{ 'card-hover': hoverCardId === news.id }"
            @mouseenter="hoverCardId = news.id"
            @mouseleave="hoverCardId = null"
            @click="navigateToNewsDetail(news.id)"
        >
          <div class="news-content">
            <div class="news-image" v-if="news.style_url">
              <el-image :src="news.style_url" fit="cover" class="image">
                <template #error>
                  <div class="image-error">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
                <template #placeholder>
                  <div class="image-loading">
                    <el-icon><Loading /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <div class="news-info">
              <h3 class="news-title">{{ news.title }}</h3>
              <p class="news-description">{{ news.description }}</p>
              <div class="news-meta">
                <div class="meta-item">
                  <el-icon><User /></el-icon>
                  <span>{{ news.from_source }}</span>
                </div>
                <div class="meta-item">
                  <el-icon><Clock /></el-icon>
                  <span>{{ formatTime(news.publish_time) }}</span>
                </div>
                <div class="meta-item" v-if="news.news_num">
                  <el-icon><View /></el-icon>
                  <span>{{ news.news_num }}阅读</span>
                </div>
              </div>
              <div class="news-tags">
                <el-tag v-if="news.class_name" size="small">{{ news.class_name }}</el-tag>
                <el-tag v-if="news.is_top === 2" type="danger" size="small">置顶</el-tag>
                <el-tag v-if="news.is_push === 1" type="warning" size="small">推荐</el-tag>
              </div>
            </div>
          </div>

          <!-- 悬浮提示 - 只保留高亮背景 -->
          <div v-if="hoverCardId === news.id" class="hover-tooltip"></div>
        </el-card>
      </div>

      <!-- 分页控件 -->
      <div class="pagination-container">
        <el-pagination
            v-model:current-page="queryParams.page"
            v-model:page-size="queryParams.page_size"
            :page-sizes="[12, 24, 48, 96]"
            :small="true"
            :background="true"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 添加资讯抽屉组件 -->
    <AddNewsDrawer v-model="showAddDrawer" @refresh="fetchNews" @submit="handlePublishSubmit"/>
  </div>
</template>
<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { Picture, Loading, Close, User, Clock, View, Plus } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/date.js'
import AddNewsDrawer from '../../components/addNewsDrawer.vue'

const router = useRouter()
const hoverCardId = ref(null)
const loading = ref(false)
const error = ref(false)
const newsList = ref([])
const total = ref(0)
const totalPages = ref(0)
const showAddDrawer = ref(false)

const queryParams = reactive({
  page: 1,
  page_size: 12,
  title: '',
  keywords: '',
  from_source: '',
  order_by: 'publish_time'
})

// 获取资讯数据
const fetchNews = async () => {
  try {
    loading.value = true
    error.value = false

    const response = await axios.post('/gapi/news/list', queryParams)

    if (response.data.error === 0) {
      // 确保 list 不会是 null
      newsList.value = response.data.data.list || []  // 关键修改！
      total.value = response.data.data.total
      totalPages.value = response.data.data.pages
    } else {
      // 处理 API 返回的错误
      newsList.value = []
      error.value = true
    }
  } catch (err) {
    console.error('获取资讯失败:', err)
    error.value = true
    newsList.value = []  // 确保异常情况下也是空数组
  } finally {
    loading.value = false
  }
}

// 处理发布提交
const handlePublishSubmit = async (formData) => {
  try {
    // 调用后端API
    const response = await axios.post('/gapi/news/insert', formData)

    if (response.data.error === 0) {
      ElMessage.success('发布成功！')
      await fetchNews()
    }
  } catch (error) {
    ElMessage.error('发布失败：' + error.message)
  }
}

// 处理搜索
const handleSearch = () => {
  queryParams.page = 1
  fetchNews()
}

// 处理重置
const handleReset = () => {
  Object.assign(queryParams, {
    page: 1,
    page_size: 12,
    title: '',
    keywords: '',
    from_source: '',
    order_by: 'publish_time'
  })
  fetchNews()
}

// 处理每页数量变化
const handleSizeChange = newSize => {
  queryParams.page_size = newSize
  queryParams.page = 1
  fetchNews()
}

// 处理页码变化
const handleCurrentChange = newPage => {
  queryParams.page = newPage
  fetchNews()
}

// 跳转到资讯详情页
const navigateToNewsDetail = id => {
  router.push({
    path: '/newsDetail',
    query: { id: id }
  })
}

// 返回主页
const handleExit = () => router.push('/zxwl')

// 初始化加载
onMounted(fetchNews)
</script>

<style lang="scss" scoped>
.news-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-container,
.error-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.error-container,
.empty-container {
  text-align: center;

  .el-button {
    margin-top: 20px;
  }
}

.filter-card {
  margin-bottom: 24px;
  border-radius: 12px;
}

.filter-form {
  margin-top: 20px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.section-title {
  color: #409EFF;
  font-size: 27px;
  margin: 0;
}

.stats-card {
  margin-bottom: 20px;
  border-radius: 12px;

  .stats-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
  }
}

.news-list {
  display: grid;
  grid-template-columns: 1fr;
  gap: 20px;
  margin-bottom: 24px;
}

.news-card {
  border-radius: 12px;
  transition: all 0.3s ease;
  position: relative;
  cursor: pointer;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
  border: 1px solid #ebeef5;

  &.card-hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
    border-color: #409EFF;

    .news-title {
      color: #409EFF;
    }
  }
}

.news-content {
  display: flex;
  gap: 16px;
  position: relative;
  z-index: 2;
}

.news-image {
  flex-shrink: 0;
  width: 200px;
  height: 150px;

  .image {
    width: 100%;
    height: 100%;
    border-radius: 8px;
    transition: transform 0.3s ease;
  }

  .image-error,
  .image-loading {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f5f7fa;
    color: #c0c4cc;
    border-radius: 8px;
  }
}

.news-card.card-hover .news-image .image {
  transform: scale(1.05);
}

.news-info {
  flex: 1;
  min-width: 0;

  .news-title {
    font-size: 1.2rem;
    color: #2c3e50;
    margin: 0 0 12px 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    transition: color 0.3s ease;
  }

  .news-description {
    color: #606266;
    margin: 0 0 16px 0;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .news-meta {
    display: flex;
    gap: 16px;
    margin-bottom: 12px;
    flex-wrap: wrap;

    .meta-item {
      display: flex;
      align-items: center;
      gap: 4px;
      color: #909399;
      font-size: 0.9rem;

      .el-icon {
        font-size: 0.9rem;
      }
    }
  }

  .news-tags {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }
}

.hover-tooltip {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, rgba(64, 158, 255, 0.05) 0%, rgba(64, 158, 255, 0.1) 100%);
  z-index: 1;
  border-radius: 12px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .news-content {
    flex-direction: column;
  }

  .news-image {
    width: 100%;
    height: 200px;
  }

  .stats-info {
    flex-direction: column;
    gap: 10px;
  }

  .hover-tooltip {
    display: none;
  }

  .news-card:active {
    transform: scale(0.98);
  }

  .filter-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
<template>
  <div class="news-detail-container">
    <!-- 顶部导航栏 -->
    <div class="news-header">
      <div class="header-left">
        <el-button @click="handleBack" :icon="ArrowLeft" class="back-btn">返回</el-button>
        <h1 class="news-title">{{ newsData.title }}</h1>
      </div>
      <div class="header-actions" v-if="showEditActions">
        <el-button type="primary" @click="handleEdit" :icon="Edit">修改</el-button>
        <el-button type="danger" @click="handleDelete" :icon="Delete">删除</el-button>
      </div>
    </div>

    <!-- 文章元信息 -->
    <el-card class="news-meta-card" :body-style="{ padding: '20px' }">
      <div class="meta-info">
        <div class="meta-left">
          <span class="source">{{ newsData.from_source }}</span>
          <span class="publish-time">{{ formatTime(newsData.publish_time) }}</span>
          <span class="views">{{ newsData.news_num }}阅读</span>
        </div>
        <div class="meta-right">
          <span class="class-name">{{ newsData.class_name }}</span>
        </div>
      </div>

    </el-card>

    <!-- 文章内容 -->
    <el-card class="news-content-card" :body-style="{ padding: '30px' }">

      <!-- 文章内容 -->
      <div class="content-wrapper">
        <div
            v-if="isHtmlContent"
            class="html-content"
            v-html="newsData.content"
        ></div>
        <div v-else class="plain-content">
          <pre>{{ newsData.content }}</pre>
        </div>
      </div>

      <!-- 关键词 -->
      <div class="keywords" v-if="newsData.keywords">
        <el-tag
            v-for="(keyword, index) in keywordsArray"
            :key="index"
            size="small"
            type="info"
            class="keyword-tag"
        >
          {{ keyword }}
        </el-tag>
      </div>
    </el-card>

    <!-- 互动按钮区域 -->
    <el-card class="interaction-card" :body-style="{ padding: '20px' }">
      <div class="interaction-buttons">
        <el-button
            @click="toggleLike"
            :type="isLiked ? 'danger' : 'default'"
            class="interaction-btn"

        >
          <template #icon>
            <LikeOutlined v-if="!isLiked" />
            <LikeFilled v-else />
          </template>
          {{ likeCount }}
        </el-button>

        <el-button
            :icon="Share"
            :type="isShared ? 'success' : ''"
            @click="toggleShare"
            class="interaction-btn"
        >
          {{ shareCount }}
        </el-button>

        <el-button
            :icon="Star"
            :type="isFavorited ? 'warning' : ''"
            @click="toggleFavorite"
            class="interaction-btn"
        >
          {{ favoriteCount }}
        </el-button>

        <el-button
            :icon="ChatDotRound"
            @click="toggleComments"
            class="interaction-btn comment-btn"
        >
          {{ commentCount }}
        </el-button>
      </div>
    </el-card>

    <!-- 评论区 -->
    <div class="comments-section" v-if="showComments">
      <NewsComments
          :news-id="newsId"
          @comment-added="handleCommentAdded"
      />
    </div>

    <!-- 删除确认对话框 -->
    <el-dialog
        v-model="deleteDialogVisible"
        title="确认删除"
        width="400px"
    >
      <span>确定要删除这篇资讯吗？此操作不可撤销。</span>
      <template #footer>
        <el-button @click="deleteDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmDelete" :loading="deleting">
          确定删除
        </el-button>
      </template>
    </el-dialog>

    <!-- 在模板底部添加编辑抽屉 -->
    <EditNewsDrawer
        v-model="editDrawerVisible"
        :news-id="newsId"
        @submitted="handleEditSubmitted"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'
import { useUserStore } from '@/utils/auth'
import EditNewsDrawer from './components/editNews.vue'
import NewsComments from './components/newsComments.vue'
import { LikeOutlined, LikeFilled } from '@ant-design/icons-vue'

// 统一导入所有图标
import {
  Share,
  Star,
  ChatDotRound,
  ArrowLeft,
  Edit,
  Delete
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const { userEmail,userName,getUser, checkLoginStatus } = useUserStore()

// 响应式数据
const editDrawerVisible = ref(false)
const newsId = ref(route.params.id || route.query.id)
const newsData = reactive({
  title: '',
  content: '',
  from_source: '',
  publish_time: '',
  news_num: '',
  class_name: '',
  style_url: '',
  keywords: '',
  publisher_email: '',
})
const loading = ref(false)
const isLiked = ref(false)
const isShared = ref(false)
const isFavorited = ref(false)
const showComments = ref(false)
const likeCount = ref(0)
const shareCount = ref(0)
const favoriteCount = ref(0)
const commentCount = ref(0)
const deleteDialogVisible = ref(false)
const deleting = ref(false)

// 计算属性
const keywordsArray = computed(() => newsData.keywords ? newsData.keywords.split(/[,，]/).filter(k => k.trim()) : [])
const isHtmlContent = computed(() => newsData.content && /<[a-z][\s\S]*>/i.test(newsData.content))
const showEditActions = computed(() => getUser().email === newsData.publisher_email)

// 获取新闻详情
const fetchNewsDetail = async () => {
  try {
    loading.value = true
    const response = await axios.get(`gapi/news/detail?id=${newsId.value}`)

    if (response.data.error === 0 && response.data.data) {
      Object.assign(newsData, response.data.data)

      // // 初始化互动数据（这里应该从后端获取真实数据）(测试用)
      // likeCount.value = Math.floor(Math.random() * 100)
      // shareCount.value = Math.floor(Math.random() * 50)
      // favoriteCount.value = Math.floor(Math.random() * 30)
      // commentCount.value = Math.floor(Math.random() * 20)
      await fetchInteractionStatus()
    } else {
      ElMessage.error('获取资讯详情失败')
      router.back()
    }
  } catch (error) {
    console.error('获取资讯详情失败:', error)
    ElMessage.error('获取资讯详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp * 1000 || timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 互动按钮处理
const toggleComments = () => {
  showComments.value = !showComments.value
}
const toggleLike = async () => {
  try {
    isLiked.value = !isLiked.value
    const change = isLiked.value ? 1 : -1
    likeCount.value += change

    const response = await axios.post('/gapi/news/count/update', {
      news_id: Number(newsId.value),
      field: 'like_count',
      action: isLiked.value ? 'increment' : 'decrement'
    })

    if (response.data.error !== 0) {
      isLiked.value = !isLiked.value
      likeCount.value -= change
      ElMessage.error('点赞操作失败')
    }
  } catch (error) {
    console.error('点赞操作失败:', error)
    isLiked.value = !isLiked.value
    likeCount.value -= isLiked.value ? 1 : -1
    ElMessage.error('网络错误，请重试')
  }
}
const toggleShare = async () => {
  try {
    isShared.value = !isShared.value
    const change = isShared.value ? 1 : -1
    shareCount.value += change

    const response = await axios.post('/gapi/news/count/update', {
      news_id: Number(newsId.value),
      field: 'share_count',
      action: isShared.value ? 'increment' : 'decrement'
    })

    if (response.data.error !== 0) {
      isShared.value = !isShared.value
      shareCount.value -= change
      ElMessage.error('分享操作失败')
    } else {
      ElMessage.success('分享成功！')
    }
  } catch (error) {
    console.error('分享操作失败:', error)
    isShared.value = !isShared.value
    shareCount.value -= isShared.value ? 1 : -1
    ElMessage.error('网络错误，请重试')
  }
}
const toggleFavorite = async () => {
  try {
    isFavorited.value = !isFavorited.value
    const change = isFavorited.value ? 1 : -1
    favoriteCount.value += change

    const response = await axios.post('/gapi/news/count/update', {
      news_id: Number(newsId.value),
      field: 'favorite_count',
      action: isFavorited.value ? 'increment' : 'decrement'
    })

    if (response.data.error !== 0) {
      isFavorited.value = !isFavorited.value
      favoriteCount.value -= change
      ElMessage.error('收藏操作失败')
    } else {
      const message = isFavorited.value ? '已收藏！' : '已取消收藏'
      ElMessage.success(message)
    }
  } catch (error) {
    console.error('收藏操作失败:', error)
    isFavorited.value = !isFavorited.value
    favoriteCount.value -= isFavorited.value ? 1 : -1
    ElMessage.error('网络错误，请重试')
  }
}

// 返回处理
const handleBack = () => {
  router.back()
}

// 编辑处理
const handleEdit = () => {
  // TODO: 跳转到编辑页面或打开编辑抽屉
  editDrawerVisible.value = true
  // ElMessage.info('暂未提供编辑功能,删除重新创建')
}
// 编辑成功回调
const handleEditSubmitted = () => {
  // 重新加载资讯详情
  fetchNewsDetail()
  ElMessage.success('资讯更新成功')
}

// 删除处理
const handleDelete = () => {
  deleteDialogVisible.value = true
}
const confirmDelete = async () => {
  try {
    deleting.value = true
    const response = await axios.delete(`/gapi/news/delete?id=${newsId.value}`)

    if (response.data.error === 0) {
      ElMessage.success('删除成功')
      router.push('/news')
    } else {
      ElMessage.error('删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
    deleteDialogVisible.value = false
  }
}

// 评论新增处理
const handleCommentAdded = () => {
  commentCount.value += 1
}

// api: 从服务器获取最新状态
const fetchInteractionStatus = async () => {
  try {
    const response = await axios.post('/gapi/news/count/get', {
      news_id: Number(newsId.value)
    })

    if (response.data.error === 0) {
      const counts = response.data.data.counts
      likeCount.value = counts.like_count
      shareCount.value = counts.share_count
      favoriteCount.value = counts.favorite_count
      commentCount.value = counts.comment_count
    }
  } catch (error) {
    console.error('获取互动状态失败:', error)
  }
}

// 监听路由变化
watch(() => route.params.id, (newId) => {
  if (newId) {
    newsId.value = newId
    fetchNewsDetail()
    // fetchInteractionStatus()
  }
})

// 组件挂载
onMounted(() => {
  console.log('newsId:', newsId.value)
  if (newsId.value) {
    fetchNewsDetail()
    // fetchInteractionStatus()
  } else {
    ElMessage.error('无效的资讯ID')
    router.back()
  }

  // 组件挂载时候获取最新数据
  fetchInteractionStatus()

})
</script>

<style scoped>
.news-detail-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.news-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-btn {
  flex-shrink: 0;
}

.news-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #1f2f3d;
  line-height: 1.4;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.news-meta-card {
  margin-bottom: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: 1px solid #ebeef5;
}

.meta-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #909399;
  font-size: 14px;
}

.meta-left {
  display: flex;
  gap: 20px;
}

.meta-right {
  font-weight: 500;
}

.news-content-card {
  margin-bottom: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: 1px solid #ebeef5;
}

.cover-image {
  margin-bottom: 30px;
  text-align: center;
}

.cover-img {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
}

.image-error {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #c0c4cc;
  border-radius: 8px;
}

.content-wrapper {
  line-height: 1.8;
  font-size: 16px;
  color: #2c3e50;
}

.html-content ::v-deep(*) {
  max-width: 100%;
}

.html-content ::v-deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.html-content ::v-deep(pre) {
  background: #f6f8fa;
  padding: 16px;
  border-radius: 6px;
  overflow-x: auto;
}

.html-content ::v-deep(code) {
  background: #f6f8fa;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.plain-content pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: inherit;
  margin: 0;
}

.keywords {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.keyword-tag {
  margin-bottom: 4px;
}

.interaction-card {
  margin-bottom: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: 1px solid #ebeef5;
}

.interaction-buttons {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
}

.interaction-btn {
  min-width: 100px;
  height: 40px;
}

.comment-btn {
  background-color: #409EFF;
  border-color: #409EFF;
  color: white;
}

.comment-btn:hover {
  background-color: #79bbff;
  border-color: #79bbff;
}

.comments-section {
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .news-detail-container {
    padding: 15px;
  }

  .news-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .header-actions {
    justify-content: center;
  }

  .news-title {
    font-size: 20px;
    text-align: center;
  }

  .meta-info {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .interaction-buttons {
    gap: 12px;
  }

  .interaction-btn {
    min-width: 80px;
    flex: 1;
  }
}
</style>
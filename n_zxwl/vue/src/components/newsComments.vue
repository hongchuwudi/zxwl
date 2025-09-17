<template>
  <div class="news-comments">
    <!-- 评论标题和统计 -->
    <div class="comments-header">
      <h3 class="comments-title">
        <el-icon><ChatDotRound /></el-icon>
        评论
      </h3>
      <span class="comments-count">{{ totalComments }} 条评论</span>
    </div>

    <!-- 发表评论区域 -->
    <div class="comment-input-section">
      <div class="input-avatar">
        <el-avatar :size="40" :src="userAvatar" v-if="userAvatar" />
        <el-avatar :size="40" v-else>
          <el-icon><User /></el-icon>
        </el-avatar>
      </div>
      <div class="input-main">
        <el-input
            v-model="newComment"
            type="textarea"
            :rows="3"
            :placeholder="isReply ? `回复 @${replyToUser}：` : '发表你的评论...'"
            :maxlength="500"
            show-word-limit
            resize="none"
        />
        <div class="input-actions">
          <div class="action-left" v-if="isReply">
            <span class="reply-hint">回复：@{{ replyToUser }}</span>
            <el-button
                size="small"
                @click="cancelReply"
                :icon="Close"
                class="cancel-reply-btn"
            >
              取消
            </el-button>
          </div>
          <el-button
              type="primary"
              size="small"
              @click="submitComment"
              :loading="submitting"
              :disabled="!newComment.trim()"
          >
            发表评论
          </el-button>
        </div>
      </div>
    </div>

    <!-- 评论排序 -->
    <div class="comments-sort">
      <el-radio-group v-model="sortType" size="small">
        <el-radio-button label="time">最新</el-radio-button>
        <el-radio-button label="hot">最热</el-radio-button>
      </el-radio-group>
    </div>

    <!-- 评论列表 -->
    <div class="comments-list" v-loading="loading">
      <div v-if="comments.length === 0 && !loading" class="empty-comments">
        <el-empty description="暂无评论，快来发表第一条评论吧～" />
      </div>

      <div
          v-for="comment in comments"
          :key="comment.id"
          class="comment-item"
          :class="{ 'has-reply': comment.replies && comment.replies.length > 0 }"
      >
        <!-- 主评论 -->
        <div class="comment-main">
          <div class="comment-avatar">
            <el-avatar :size="36" :src="getAvatar(comment.commenter_email)" />
          </div>
          <div class="comment-content">
            <div class="comment-header">
              <span class="comment-author">{{ comment.commenter_name }}</span>
              <span class="comment-time">{{ formatTime(comment.comment_time) }}</span>
            </div>
            <div class="comment-text">{{ comment.comment_content }}</div>
            <div class="comment-actions">
              <el-button
                  size="small"
                  @click="handleReply(comment)"
                  class="action-btn"
              >
                <template #icon>
                  <ChatDotRound />
                </template>
                回复
              </el-button>
              <el-button
                  size="small"
                  :type="comment.liked ? 'danger' : ''"
                  @click="toggleLike(comment)"
                  class="action-btn"
              >
                <template #icon>
                  <LikeOutlined v-if="!comment.liked" />
                  <LikeFilled v-else />
                </template>
                {{ comment.like_count || 0 }}
              </el-button>

              <!-- 展开/折叠回复按钮 -->
              <el-button
                  v-if="comment.reply_count > 0"
                  size="small"
                  @click="toggleReplies(comment)"
                  class="action-btn"
              >
                <template #icon>
                  <el-icon v-if="comment.showReplies">
                    <ArrowUp />
                  </el-icon>
                  <el-icon v-else>
                    <ArrowDown />
                  </el-icon>
                </template>
                {{ comment.showReplies ? '折叠' : '展开' }}回复 ({{ comment.reply_count }})
              </el-button>
            </div>
          </div>
        </div>

        <!-- 回复列表 -->
        <div
            v-if="comment.replies && comment.replies.length > 0 && comment.showReplies"
            class="replies-list"
        >
          <div
              v-for="reply in comment.replies"
              :key="reply.id"
              class="reply-item"
          >
            <div class="reply-avatar">
              <el-avatar :size="32" :src="getAvatar(reply.commenter_email)" />
            </div>
            <div class="reply-content">
              <div class="reply-header">
                <span class="reply-author">{{ reply.commenter_name }}</span>
                <span class="reply-time">{{ formatTime(reply.comment_time) }}</span>
              </div>
              <div class="reply-text">
                <span v-if="reply.parent_id !== comment.id" class="reply-to">
                  @{{ getReplyToUserName(reply) }}
                </span>
                {{ reply.comment_content }}
              </div>
              <div class="reply-actions">
                <el-button
                    size="small"
                    @click="handleReply(reply)"
                    class="action-btn"
                >
                  <template #icon>
                    <ChatDotRound />
                  </template>
                  回复
                </el-button>
                <el-button
                    size="small"
                    :type="reply.liked ? 'danger' : ''"
                    @click="toggleLike(reply)"
                    class="action-btn"
                >
                  <template #icon>
                    <LikeOutlined v-if="!reply.liked" />
                    <LikeFilled v-else />
                  </template>
                  {{ reply.like_count || 0 }}
                </el-button>
              </div>
            </div>
          </div>

          <!-- 查看更多回复 -->
          <div
              v-if="comment.reply_count > (comment.replies?.length || 0)"
              class="view-more-replies"
              @click="loadMoreReplies(comment)"
          >
            <el-button link type="primary" size="small" >
              展开更多 {{ comment.reply_count - (comment.replies?.length || 0) }} 条回复
              <el-icon><ArrowDown /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { useUserStore } from '@/utils/auth.js'
import { formatTime } from '@/utils/date.js'
import {
  ChatDotRound,
  User,
  Close,
  ArrowDown,
  ArrowUp
} from '@element-plus/icons-vue'
// 引入 ant-design 图标
import { LikeOutlined, LikeFilled } from '@ant-design/icons-vue'

const props = defineProps({
  newsId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['comment-added'])

const { getUser } = useUserStore()
const userLocal = computed(() => getUser())

// 响应式数据
const comments = ref([])
const newComment = ref('')
const sortType = ref('time')
const currentPage = ref(1)
const pageSize = ref(10)
const totalComments = ref(0)
const loading = ref(false)
const submitting = ref(false)
const isReply = ref(false)
const replyToComment = ref(null)
const replyToUser = ref('')

// 计算属性
const userAvatar = computed(() => {
  const user = getUser()
  return user?.avatar || ''
})

// api: 获取评论列表
const fetchComments = async () => {
  try {
    loading.value = true
    const response = await axios.post('/gapi/news/comments/list', {
      news_id: Number(props.newsId)
    })

    if (response.data.error === 0) {
      comments.value = response.data.data.comments || []
      totalComments.value = response.data.data.total || 0

      // 初始化评论状态
      comments.value.forEach(comment => {
        comment.liked = false
        // 默认展开回复
        comment.showReplies = false

        // 如果有回复，也初始化回复的点赞状态
        if (comment.replies) {
          comment.replies.forEach(reply => {
            reply.liked = false
          })
        }
      })
    } else {
      ElMessage.error('获取评论失败')
    }
  } catch (error) {
    console.error('获取评论失败:', error)
    ElMessage.error('获取评论失败')
  } finally {
    loading.value = false
  }
}

// api: 提交评论
const submitComment = async () => {
  if (!newComment.value.trim()) return

  try {
    submitting.value = true
    const commentData = {
      news_id: Number(props.newsId),
      commenter_name: getUser().name,
      commenter_email: getUser().email,
      comment_content: newComment.value.trim()
    }

    if (isReply.value && replyToComment.value) {
      commentData.parent_id = replyToComment.value.id
    }

    const response = await axios.post('/gapi/news/comments/add', commentData)

    if (response.data.error === 0) {
      ElMessage.success('评论发表成功')
      newComment.value = ''
      cancelReply()
      emit('comment-added')
      await fetchComments() // 刷新评论列表
    } else {
      ElMessage.error('评论发表失败')
    }
  } catch (error) {
    console.error('评论发表失败:', error)
    ElMessage.error('评论发表失败')
  } finally {
    submitting.value = false
  }
}

// api: 点赞/取消点赞
const toggleLike = async (comment) => {
  try {
    const action = comment.liked ? 'unlike' : 'like'
    const response = await axios.post('/gapi/news/comments/like', {
      comment_id: comment.id,
      action: action
    })

    if (response.data.error === 0) {
      comment.liked = !comment.liked
      comment.like_count = comment.liked ? (comment.like_count || 0) + 1 : Math.max(0, (comment.like_count || 1) - 1)
    } else {
      ElMessage.error('操作失败')
    }
  } catch (error) {
    console.error('点赞操作失败:', error)
    ElMessage.error('操作失败')
  }
}

// 回复评论操作
const handleReply = (comment) => {
  isReply.value = true
  replyToComment.value = comment
  replyToUser.value = comment.commenter_name
  // 滚动到输入框
  setTimeout(() => {
    const inputSection = document.querySelector('.comment-input-section')
    inputSection?.scrollIntoView({ behavior: 'smooth' })
  }, 100)
}

// 取消回复
const cancelReply = () => {
  isReply.value = false
  replyToComment.value = null
  replyToUser.value = ''
}

// 展开/折叠回复
const toggleReplies = (comment) => {
  comment.showReplies = !comment.showReplies
}

// 加载更多回复
const loadMoreReplies = async (comment) => {
  // 这里需要实现加载更多回复的逻辑
  ElMessage.info('加载更多回复功能开发中')
}

// 获取回复的用户名
const getReplyToUserName = (reply) => {
  // 这里需要根据parent_id找到被回复的用户名
  return '用户'
}

// 获取用户头像（根据邮箱生成默认头像）
const getAvatar = (email) => {
  // 这里可以根据实际情况实现头像获取逻辑
  // return  'https://t.alcy.cc/lai'
  // return `https://api.dicebear.com/7.x/avataaars/svg?seed=${email}`
  // return `https://api.dicebear.com/7.x/avataaars/svg?seed=${email}&backgroundColor=b6e3f4,c0aede,d1d4f9`
  // 方案1: 使用Adventurer系列（非常可爱）
  return `https://api.dicebear.com/7.x/adventurer/svg?seed=${email}`

  // 方案2: 使用Bottts机器人系列（卡通可爱）
  // return `https://api.dicebear.com/7.x/bottts/svg?seed=${email}`
  //
  // 方案3: 使用Lorelei系列（多彩有趣）
  // return `https://api.dicebear.com/7.x/lorelei/svg?seed=${email}`

  // 方案4: 使用Micah系列（表情丰富）
  // return `https://api.dicebear.com/7.x/micah/svg?seed=${email}&earrings=variant04`

  // return `https://api.dicebear.com/7.x/adventurer/svg?seed=${email}&hairColor=ffb3c6&backgroundColor=FDECF2`;

}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchComments()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchComments()
}

// 监听排序方式变化
watch(sortType, () => {
  currentPage.value = 1
  fetchComments()
})

// 监听新闻ID变化
watch(() => props.newsId, (newId) => {
  if (newId) {
    currentPage.value = 1
    fetchComments()
  }
})

// 组件挂载
onMounted(() => {
  if (props.newsId) {
    fetchComments()
  }
})
</script>

<style scoped>
/* 样式保持不变，与之前相同 */
.news-comments {
  padding: 20px 0;
}

.comments-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e5e9ef;
}

.comments-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2f3d;
}

.comments-count {
  color: #909399;
  font-size: 14px;
}

.comment-input-section {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.input-avatar {
  flex-shrink: 0;
}

.input-main {
  flex: 1;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
}

.action-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reply-hint {
  color: #409eff;
  font-size: 14px;
}

.cancel-reply-btn {
  padding: 4px 8px;
}

.comments-sort {
  margin-bottom: 16px;
}

.comments-list {
  min-height: 200px;
}

.empty-comments {
  padding: 40px 0;
  text-align: center;
}

.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid #e5e9ef;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-main {
  display: flex;
  gap: 12px;
}

.comment-avatar {
  flex-shrink: 0;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 500;
  color: #1f2f3d;
}

.comment-time {
  color: #909399;
  font-size: 12px;
}

.comment-text {
  margin-bottom: 12px;
  line-height: 1.6;
  color: #2c3e50;
}

.comment-actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.action-btn {
  padding: 4px 8px;
  font-size: 12px;
}

.replies-list {
  margin-left: 48px;
  margin-top: 16px;
  padding-left: 16px;
  border-left: 2px solid #e5e9ef;
}

.reply-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
}

.reply-item:not(:last-child) {
  border-bottom: 1px solid #f0f2f5;
}

.reply-avatar {
  flex-shrink: 0;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.reply-author {
  font-weight: 500;
  color: #1f2f3d;
  font-size: 14px;
}

.reply-time {
  color: #909399;
  font-size: 12px;
}

.reply-text {
  line-height: 1.5;
  color: #2c3e50;
  font-size: 14px;
}

.reply-to {
  color: #409eff;
  margin-right: 4px;
}

.reply-actions {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}

.view-more-replies {
  text-align: center;
  padding: 8px 0;
}

.comments-pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .comment-input-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .input-avatar {
    margin-bottom: 12px;
  }

  .replies-list {
    margin-left: 24px;
  }

  .comment-actions {
    flex-wrap: wrap;
  }
}
</style>
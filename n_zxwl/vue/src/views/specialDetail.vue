<template>

  <div class="special-detail-container">
    <!-- 头部区域 -->
    <div class="special-header">
      <div class="header-left">
        <h1 class="special-title">{{ specialDetail.name }}</h1>
        <p class="special-code">专业代码: {{ specialDetail.code }}</p>
      </div>

      <el-button
          class="exit-button"
          type="primary"
          @click="handleExit"
          round
      >
        <el-icon><Close /></el-icon>
        <span>返回专业列表</span>
      </el-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-container">
      <el-empty description="加载失败，请重试" />
      <el-button @click="fetchSpecialDetail" type="primary">重新加载</el-button>
    </div>

    <!-- 内容区域 -->
    <div v-else>
      <!-- 基本信息卡片 -->
      <el-card class="info-card">
        <h2 class="section-title"><el-icon><InfoFilled /></el-icon> 基本信息</h2>
        <div class="basic-info-grid">
          <div class="info-item">
            <span class="info-label">学科门类:</span>
            <span class="info-value">{{ specialDetail.type }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">专业类:</span>
            <span class="info-value">{{ specialDetail.type_detail }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">授予学位:</span>
            <span class="info-value">{{ specialDetail.degree === '' ? '未知' : specialDetail.degree }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">修业年限:</span>
            <span class="info-value">{{ specialDetail.limit_year }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">选科要求:</span>
            <span class="info-value">{{ specialDetail.subject_requirements === '' ? '无要求' : specialDetail.subject_requirements }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">男女比例:</span>
            <span class="info-value">{{ formatGenderRatio(specialDetail.gender_ratio) }}</span>
          </div>
        </div>
      </el-card>

      <!-- 知名院校 -->
      <el-card class="content-card" v-if="famousSchools.length > 0">
        <h2 class="section-title"><el-icon><School /></el-icon> 知名院校</h2>
        <div class="schools-list">
          <el-tag
              v-for="school in famousSchools"
              :key="school.id"
              type="info"
              size="large"
              class="school-tag"
              @click="handleSchoolClick(school.school_name)"
              @mouseenter="hoverSchool = school.school_name"
              @mouseleave="hoverSchool = null"
              :class="{ 'school-hover': hoverSchool === school.school_name }"
              style="cursor: pointer;"
          >
            {{ school.school_name }}
          </el-tag>
        </div>
      </el-card>

      <!-- 专业介绍 -->
      <el-card class="content-card" >
        <h2 class="section-title"><el-icon><Notebook /></el-icon> 专业介绍</h2>
        <div v-html="getContentByType(1)" class="content-text"></div>
      </el-card>

      <!-- 薪资趋势 -->
      <el-card class="content-card" v-if="salaryData.length > 0">
        <h2 class="section-title"><el-icon><Money /></el-icon> 薪资趋势</h2>
        <div class="salary-trend">
          <div v-for="item in sortedSalaryData" :key="item.id" class="salary-item">
            <span class="salary-year">{{ item.salary_year }}年经验</span>
            <span class="salary-value">¥{{ formatNumber(item.salary_value) }}/月</span>
            <span class="salary-type">{{ item.salary_type === 1 ? '平均薪资' : '中位数薪资' }}</span>
          </div>
        </div>
      </el-card>

      <!-- 就业率趋势图 -->
      <el-card class="content-card" v-if="employmentRates.length > 0">
        <h2 class="section-title"><el-icon><TrendCharts /></el-icon> 就业率趋势</h2>
        <div class="employment-trend">
          <div v-for="rate in employmentRates" :key="rate.id" class="trend-item">
            <span class="trend-year">{{ rate.year }}年</span>
            <span class="trend-rate">{{ rate.rate }}</span>
          </div>
        </div>
      </el-card>

      <!-- 就业与薪资信息 -->
      <div class="stats-container" >
        <el-card class="stat-card" v-if="specialDetail.employment_rate !== 0">
          <h3 class="stat-title">就业率</h3>
          <div class="stat-value">{{ specialDetail.employment_rate }}%</div>
          <div class="stat-trend">
            <el-tag v-if="employmentTrend === 'up'" type="success" size="small">逐年上升</el-tag>
            <el-tag v-else-if="employmentTrend === 'down'" type="danger" size="small">逐年下降</el-tag>
            <el-tag v-else type="info" size="small">保持稳定</el-tag>
          </div>
        </el-card>

        <el-card class="stat-card">
          <h3 class="stat-title">平均薪资</h3>
          <div class="stat-value">¥{{ Math.round(specialDetail.avg_salary / 12)}}/月</div>
          <div class="stat-subtitle">年薪 {{ formatNumber(specialDetail.avg_salary) }} 元</div>
        </el-card>

        <el-card class="stat-card" v-if="specialDetail.top_area">
          <h3 class="stat-title">热门就业地区</h3>
          <div class="stat-value">{{ specialDetail.top_area }}</div>
          <div class="stat-subtitle">最多就业岗位: {{ specialDetail.top_position }}</div>
        </el-card>
      </div >

      <!-- 就业方向 -->
      <el-card class="content-card" v-if="getContentByType(2) === '暂无内容'">
        <h2 class="section-title"><el-icon><Briefcase /></el-icon> 就业方向</h2>
        <div v-html="getContentByType(2)" class="content-text"></div>
      </el-card>

      <!-- 主要课程 -->
      <el-card class="content-card" v-if="getContentByType(2) === '暂无内容'">
        <h2 class="section-title"><el-icon><Collection /></el-icon> 主要课程</h2>
        <div v-html="getContentByType(4)" class="content-text"></div>
      </el-card>

      <!-- 就业分布 -->
      <el-card class="content-card" v-if = "jobDistributions.length > 0">
        <h2 class="section-title"><el-icon><DataAnalysis /></el-icon> 就业分布</h2>
        <div class="distribution-charts">
          <div class="chart-container">
            <h3>行业分布</h3>
            <div v-for="item in industryDistribution" :key="'industry-'+item.id" class="distribution-item">
              <div class="dist-header">
                <span class="dist-name">{{ item.name }}</span>
                <span class="dist-rate">{{ item.rate }}%</span>
              </div>
              <el-progress :percentage="parseFloat(item.rate)" :color="getProgressColor(parseFloat(item.rate))"></el-progress>
            </div>
          </div>

          <div class="chart-container">
            <h3>地区分布</h3>
            <div v-for="item in areaDistribution" :key="'area-'+item.id" class="distribution-item">
              <div class="dist-header">
                <span class="dist-name">{{ item.name }}</span>
                <span class="dist-rate">{{ item.rate }}%</span>
              </div>
              <el-progress :percentage="parseFloat(item.rate)" :color="getProgressColor(parseFloat(item.rate))"></el-progress>
            </div>
          </div>

          <div class="chart-container">
            <h3>岗位分布</h3>
            <div v-for="item in positionDistribution" :key="'position-'+item.id" class="distribution-item">
              <div class="dist-header">
                <span class="dist-name">{{ item.position }}</span>
                <span class="dist-rate">{{ item.rate }}%</span>
              </div>
              <el-progress :percentage="parseFloat(item.rate)" :color="getProgressColor(parseFloat(item.rate))"></el-progress>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 印象标签 -->
      <el-card class="content-card" v-if="impressionTags.length > 0">
        <h2 class="section-title"><el-icon><PriceTag /></el-icon> 专业印象</h2>
        <div class="tags-container">
          <div v-for="tag in impressionTags" :key="tag.id" class="tag-item">
            <el-image
                :src="tag.image_url"
                fit="cover"
                class="tag-image"
            >
              <template #error>
                <div class="tag-image-error">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <span class="tag-keyword">{{ tag.keyword }}</span>
          </div>
        </div>
      </el-card>

      <!-- 专业视频 -->
      <el-card class="content-card" v-if="videos.length > 0">
        <h2 class="section-title"><el-icon><VideoPlay /></el-icon> 专业视频</h2>
        <div class="videos-container">
          <div v-for="video in videos" :key="video.id" class="video-item">
            <el-image
                :src="video.cover_image"
                fit="cover"
                class="video-thumbnail"
                @click="playVideo(video)"
            >
              <template #error>
                <div class="thumbnail-error">
                  <el-icon><VideoCamera /></el-icon>
                </div>
              </template>
              <template #placeholder>
                <div class="thumbnail-loading">
                  <el-icon><Loading /></el-icon>
                </div>
              </template>
            </el-image>
            <p class="video-title">{{ video.title }}</p>
          </div>
        </div>
      </el-card>

      <!-- 视频播放对话框 -->
      <el-dialog
          v-model="videoDialogVisible"
          :title="currentVideo.title"
          width="70%"
          top="5vh"
      >
        <!-- 关键修改：添加 referrerpolicy 属性 -->
        <video
            v-if="videoDialogVisible"
            :src="currentVideo.video_url"
            controls
            autoplay
            class="video-player"
            referrerpolicy="unsafe-url"
        ></video>
      </el-dialog>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import {
  Close,
  InfoFilled,
  Notebook,
  Briefcase,
  Collection,
  DataAnalysis,
  School,
  VideoPlay,
  VideoCamera,
  Loading,
  Picture,
  TrendCharts,
  Money,
  PriceTag
} from '@element-plus/icons-vue'
// 确保已经导入 ElMessage
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref(false)
const specialDetail = ref({})
const specialContents = ref([])
const jobDistributions = ref([])
const famousSchools = ref([])
const videos = ref([])
const employmentRates = ref([])
const salaryData = ref([])
const impressionTags = ref([])
const universitySpecialInfo = ref([])

const videoDialogVisible = ref(false)
const currentVideo = ref({})

// 知名院校高亮显示
const hoverSchool = ref(null)
const loadingSchool = ref(false)

// 计算属性函数: 获取行业分布
const industryDistribution = computed(() => {
  return jobDistributions.value
      .filter(item => item.distribution_type === 1)
      .sort((a, b) => parseFloat(b.rate) - parseFloat(a.rate))
})

// 计算属性函数: 获取地区分布
const areaDistribution = computed(() => {
  return jobDistributions.value
      .filter(item => item.distribution_type === 2)
      .sort((a, b) => parseFloat(b.rate) - parseFloat(a.rate))
})

// 计算属性函数: 获取岗位分布
const positionDistribution = computed(() => {
  return jobDistributions.value
      .filter(item => item.distribution_type === 3)
      .sort((a, b) => parseFloat(b.rate) - parseFloat(a.rate))
})

// 计算属性函数: 获取地区分布
const sortedSalaryData = computed(() => {
  return [...salaryData.value].sort((a, b) => a.salary_year - b.salary_year)
})

// 计算属性函数: 获取就业趋势
const employmentTrend = computed(() => {
  if (employmentRates.value.length < 2) return 'stable'

  const rates = employmentRates.value.map(rate => {
    const rateValue = parseFloat(rate.rate.split('-')[0]) // 取范围的最低值
    return rateValue
  })

  // 检查趋势
  let increasing = true
  let decreasing = true

  for (let i = 1; i < rates.length; i++) {
    if (rates[i] <= rates[i - 1]) increasing = false
    if (rates[i] >= rates[i - 1]) decreasing = false
  }

  if (increasing) return 'up'
  if (decreasing) return 'down'
  return 'stable'
})

// 工具: 格式化数字
const formatNumber = (num) => {
  if (!num) return '0'
  return Number(num).toLocaleString()
}

// 方法: 获取内容
const getContentByType = (type) => {
  const content = specialContents.value.find(item => item.content_type === type)
  return content ? content.content : '暂无内容'
}

//  工具: 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage > 15) return '#409EFF'
  if (percentage > 10) return '#67C23A'
  if (percentage > 5) return '#E6A23C'
  return '#909399'
}

// 事件处理: 播放视频
const playVideo = (video) => {
  currentVideo.value = video
  videoDialogVisible.value = true
}

// 事件处理: 退出到专业列表
const handleExit = () => router.back()

// api: 获取专业详情
const fetchSpecialDetail = async () => {
  try {
    loading.value = true
    error.value = false

    const url = `gapi/specials/profiles/${route.query.id}`
    const response = await axios.get(url)

    if (response.data.code === 0) {
      const data = response.data.data

      specialDetail.value = data.special_detail
      specialContents.value = data.special_contents
      jobDistributions.value = data.job_distributions
      famousSchools.value = data.famous_schools
      videos.value = data.videos
      employmentRates.value = data.employment_rates
      salaryData.value = data.salary_data
      impressionTags.value = data.impression_tags
      universitySpecialInfo.value = data.university_special_info
    } else {
      throw new Error('API返回错误')
    }
  } catch (err) {
    console.error('获取专业详情失败:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

// 工具: 格式化性别比例显示 "60.00:40.00" -> "♂ 60% : ♀ 40%"
function formatGenderRatio(ratioStr) {
  // 分割字符串获取男女比例
  const ratios = ratioStr.split(':');

  // 转换为数字并四舍五入取整
  const maleRatio = Math.round(parseFloat(ratios[0]));
  const femaleRatio = Math.round(parseFloat(ratios[1]));

  // 返回格式化后的字符串（使用性别符号）
  return `♂ ${maleRatio}% : ♀ ${femaleRatio}%`;
}

// 功能: 添加学校点击处理方法
const handleSchoolClick = async (schoolName) => {
  if (loadingSchool.value) return

  loadingSchool.value = true
  try {
    const response = await axios.get(`/gapi/school/name?name=${encodeURIComponent(schoolName)}`)

    if (response.data.code === 0 && response.data.data && response.data.data.id) {
      // 找到学校ID，跳转到学校详情
      router.push({
        path: '/schoolDetail',
        query: { id: response.data.data.id }
      })
    } else if (response.data.msg && response.data.msg.includes('未找到')) {
      // 学校未找到，跳转到学校列表页
      router.push('/allSchool')
    } else {
      ElMessage.error('获取学校信息失败，请稍后重试')
    }
  } catch (error) {
    console.error('请求学校信息失败:', error)
    if (error.response?.status === 404) {
      ElMessage.warning(`"${schoolName}" 学校信息暂未收录`)
    } else {
      ElMessage.error('网络请求失败，请检查网络连接')
    }
  } finally {
    loadingSchool.value = false
  }
}

// 钩子函数: vue挂载后执行
onMounted(() => {
  // 检查是否已存在referrer meta标签
  let metaTag = document.querySelector('meta[name="referrer"]');

  if (!metaTag) {
    // 创建新的meta标签
    console.log('创建新的meta标签');
    metaTag = document.createElement('meta');
    metaTag.name = 'referrer';
    metaTag.content = 'no-referrer';
    document.getElementsByTagName('head')[0].appendChild(metaTag);
  } else {
    // 更新现有的meta标签
    console.log('更新现有的meta标签');
    metaTag.content = 'no-referrer';
  }
  fetchSpecialDetail()
})
</script>

<style lang="scss" scoped>
.special-detail-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

.error-container {
  text-align: center;

  .el-button {
    margin-top: 20px;
  }
}

.special-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eaeaea;

  .header-left {
    .special-title {
      font-size: 2rem;
      color: #2c3e50;
      margin: 0 0 8px 0;
    }

    .special-code {
      color: #909399;
      font-size: 1rem;
      margin: 0;
    }
  }

  .exit-button {
    background: linear-gradient(135deg, #ff6b6b, #ff8787);
    border: none;
    color: white;

    &:hover {
      background: linear-gradient(135deg, #ff8787, #ff6b6b);
    }
  }
}

.info-card {
  margin-bottom: 24px;
  border-radius: 12px;

  .basic-info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;

    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 8px 0;
      border-bottom: 1px dashed #ebeef5;

      .info-label {
        color: #606266;
        font-weight: 500;
      }

      .info-value {
        color: #303133;
      }
    }
  }
}

.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;

  .stat-card {
    text-align: center;
    border-radius: 12px;

    .stat-title {
      color: #909399;
      font-size: 1rem;
      margin: 0 0 12px 0;
    }

    .stat-value {
      color: #409EFF;
      font-size: 1.8rem;
      font-weight: bold;
      margin-bottom: 8px;
    }

    .stat-subtitle {
      color: #606266;
      font-size: 0.9rem;
    }

    .stat-trend {
      margin-top: 8px;
    }
  }
}

.employment-trend, .salary-trend {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;

  .trend-item, .salary-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    background: #f8f9fa;
    border-radius: 8px;

    .trend-year, .salary-year {
      color: #606266;
      font-weight: 500;
    }

    .trend-rate, .salary-value {
      color: #409EFF;
      font-weight: bold;
    }

    .salary-type {
      color: #909399;
      font-size: 0.8rem;
    }
  }
}

.content-card {
  margin-bottom: 24px;
  border-radius: 12px;

  .content-text {
    line-height: 1.8;
    color: #606266;
  }
}

.distribution-charts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 30px;

  .chart-container {
    h3 {
      color: #409EFF;
      margin-bottom: 16px;
    }

    .distribution-item {
      margin-bottom: 16px;

      .dist-header {
        display: flex;
        justify-content: space-between;
        margin-bottom: 6px;

        .dist-name {
          color: #606266;
        }

        .dist-rate {
          color: #409EFF;
          font-weight: 500;
        }
      }
    }
  }
}

.schools-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;

  .school-tag {
    font-size: 1rem;
    padding: 8px 16px;
  }
}

/* 添加学校标签悬浮样式 */
.school-tag {
  transition: all 0.3s ease;
  cursor: pointer;
}

.school-hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(144, 147, 153, 0.3);
  background: linear-gradient(135deg, #f4f4f5, #e9e9eb) !important;
}

/* 加载状态 */
.school-loading {
  opacity: 0.7;
  cursor: not-allowed;
}

.tags-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;

  .tag-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;

    .tag-image {
      width: 80px;
      height: 80px;
      border-radius: 8px;
      margin-bottom: 8px;

      .tag-image-error {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #f5f7fa;
        color: #c0c4cc;
      }
    }

    .tag-keyword {
      font-weight: 500;
      color: #303133;
    }
  }
}

.videos-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;

  .video-item {
    cursor: pointer;
    transition: transform 0.3s;

    &:hover {
      transform: translateY(-5px);
    }

    .video-thumbnail {
      width: 100%;
      height: 160px;
      border-radius: 8px;
      overflow: hidden;
    }

    .video-title {
      margin: 10px 0 0 0;
      font-size: 0.9rem;
      color: #606266;
      line-height: 1.4;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
    }
  }
}

.video-player {
  width: 100%;
  border-radius: 8px;
}

.thumbnail-error,
.thumbnail-loading {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #c0c4cc;
}

.section-title {
  color: #409EFF;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;

  .el-icon {
    font-size: 1.2rem;
  }
}

@media (max-width: 768px) {
  .special-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;

    .exit-button {
      align-self: flex-end;
    }
  }

  .distribution-charts {
    grid-template-columns: 1fr;
  }

  .videos-container {
    grid-template-columns: 1fr;
  }

  .employment-trend, .salary-trend {
    grid-template-columns: 1fr;
  }

  .tags-container {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  }
}

.badge-style {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
}
</style>
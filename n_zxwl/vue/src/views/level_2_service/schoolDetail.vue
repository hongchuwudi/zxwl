<template>
  <div class="school-detail-container">
    <!-- 头部区域 -->
    <div class="school-header">
      <div class="header-left">
        <h1 class="school-title">{{ schoolDetail.name }}</h1>
        <p class="school-motto">{{ schoolDetail.motto }}</p>
      </div>

      <el-button
          class="exit-button"
          type="primary"
          @click="handleExit"
          round
      >
        <el-icon>
          <Close/>
        </el-icon>
        <span>返回学校列表</span>
      </el-button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated/>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-container">
      <el-empty description="加载失败，请重试"/>
      <el-button @click="fetchSchoolDetail" type="primary">重新加载</el-button>
    </div>

    <!-- 内容区域 -->
    <div v-else>
      <!-- 学校基本信息卡片 -->
      <el-card class="info-card">
        <div class="school-basic-info">
          <div class="logo-container">
            <el-image
                :src="schoolDetail.logo_url"
                fit="cover"
                class="school-logo"
            >
              <template #error>
                <div class="logo-error">
                  <el-icon>
                    <Picture/>
                  </el-icon>
                </div>
              </template>
            </el-image>
          </div>
          <div class="basic-info-grid">
            <div class="info-item">
              <span class="info-label">创办时间:</span>
              <span class="info-value">{{ schoolDetail.create_date }}年</span>
            </div>
            <div class="info-item">
              <span class="info-label">学校性质:</span>
              <span class="info-value">{{ schoolDetail.nature_name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">学校类型:</span>
              <span class="info-value">{{ schoolDetail.type_name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">主管部门:</span>
              <span class="info-value">{{ schoolDetail.belong }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">学校层次:</span>
              <span class="info-value">{{ schoolDetail.level_name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">双一流:</span>
              <span class="info-value">{{ schoolDetail.dual_class_name || '否' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">211:</span>
              <span class="info-value">{{ schoolDetail.f211 === 1 ? '是' : schoolDetail.f211 === 2 ? '否': '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">985:</span>
              <span class="info-value">{{ schoolDetail.f985 === 1 ? '是' : schoolDetail.f985 === 2 ? '否': '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">占地面积:</span>
              <span class="info-value">{{ schoolDetail.area }}亩</span>
            </div>
            <div class="info-item">
              <span class="info-label">QS世界排名:</span>
              <span class="info-value">{{ schoolDetail.qs_world || '暂无' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">软科排名:</span>
              <span class="info-value">{{ schoolDetail.ruanke_rank || '暂无' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">US排名:</span>
              <span class="info-value">{{ schoolDetail.us_rank || '暂无' }}</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 学校简介 -->
      <el-card class="content-card">
        <h2 class="section-title">
          <el-icon>
            <InfoFilled/>
          </el-icon>
          学校简介
        </h2>
<!--        :style="{ maxHeight: isCollapsed ? collapsedHeight + 'px' : contentHeight + 'px' }"-->
        <div
            ref="contentRef"
            class="content-text rich-text-content"
            :class="{ 'collapsed': isCollapsed && canCollapse }"
            v-html="formattedContent"
        ></div>
        <div v-if="canCollapse" class="toggle-btn">
          <el-button
              link
              type="primary"
              @click="toggleContent"
              class="toggle-button"
          >
            {{ isCollapsed ? '展开全部' : '收起' }}
            <el-icon>
              <ArrowDown v-if="isCollapsed"/>
              <ArrowUp v-else/>
            </el-icon>
          </el-button>
        </div>
      </el-card>

      <!-- 联系信息 -->
      <el-card class="content-card">
        <h2 class="section-title">
          <el-icon>
            <Connection/>
          </el-icon>
          联系信息
        </h2>
        <div class="contact-grid">
          <div class="contact-item">
            <span class="contact-label">地址:</span>
            <span
                class="contact-value address-highlight"
                @click="showMapDrawer"
                @mouseenter="hoverAddress = true"
                @mouseleave="hoverAddress = false"
                :class="{ 'address-hover': hoverAddress }"
            >{{ getFirstCampusAddress(schoolDetail.address) }}</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">邮编:</span>
            <span class="contact-value">{{ schoolDetail.postcode }}</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">电话:</span>
            <span class="contact-value">{{ schoolDetail.phone }}</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">邮箱:</span>
            <span class="contact-value">{{ schoolDetail.email }}</span>
          </div>
          <div class="contact-item">
            <span class="contact-label">招生网站:</span>
            <span class="contact-value">
              <a :href="schoolDetail.site" target="_blank" rel="noopener">{{ schoolDetail.site }}</a>
            </span>
          </div>
          <div class="contact-item">
            <span class="contact-label">学校官网:</span>
            <span class="contact-value">
              <a :href="schoolDetail.school_site" target="_blank" rel="noopener">{{ schoolDetail.school_site }}</a>
            </span>
          </div>
        </div>
      </el-card>

      <!-- 院系设置 -->
      <el-card class="content-card" v-if="collegesDepartments.length > 0">
        <h2 class="section-title">
          <el-icon>
            <OfficeBuilding/>
          </el-icon>
          院系设置
        </h2>
        <div class="departments-container">
          <div v-for="(campus, index) in groupedDepartments" :key="index" class="campus-section">
            <h3 class="campus-name">{{ campus.campusName }}</h3>
            <div class="departments-list">
              <el-tag
                  v-for="dept in campus.departments"
                  :key="dept.department_id"
                  type="info"
                  size="large"
                  class="department-tag"
              >
                {{ dept.department_name }}
              </el-tag>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 学科排名 -->
      <el-card class="content-card" v-if="disciplineRankings.length > 0">
        <h2 class="section-title">
          <el-icon>
            <Trophy/>
          </el-icon>
          学科排名
        </h2>
        <div class="rankings-container">
          <div v-for="rank in sortedRankings" :key="rank.rank_level" class="ranking-item">
            <span class="rank-level">{{ rank.rank_level }}级学科</span>
            <span class="rank-count">{{ rank.count }}个</span>
          </div>
        </div>
      </el-card>

      <!-- 双一流学科 -->
      <el-card class="content-card" v-if="dualClassSubjects.length > 0">
        <h2 class="section-title">
          <el-icon>
            <Star/>
          </el-icon>
          双一流学科
        </h2>
        <div class="subjects-list">
          <el-tag
              v-for="subject in dualClassSubjects"
              :key="subject.id"
              type="success"
              size="large"
              class="subject-tag"
          >
            {{ subject.subject_name }}
          </el-tag>
        </div>
      </el-card>

      <!-- 特色专业 -->
      <el-card class="content-card" v-if="specialPrograms.length > 0">
        <h2 class="section-title">
          <el-icon>
            <Collection/>
          </el-icon>
          特色专业
        </h2>
        <el-table :data="specialPrograms" style="width: 100%">
          <el-table-column label="专业名称" min-width="180">
            <template #default="scope">
        <span
            class="special-name-link"
            @click="handleSpecialClick(scope.row.special_name)"
            @mouseenter="hoverSpecial = scope.row.special_name"
            @mouseleave="hoverSpecial = null"
            :class="{ 'special-hover': hoverSpecial === scope.row.special_name }"
        >
          {{ scope.row.special_name }}
          <el-icon v-if="loadingSpecial" class="is-loading">
            <Loading />
          </el-icon>
        </span>
            </template>
          </el-table-column>
          <el-table-column prop="level_name" label="层次" width="100"/>
          <el-table-column prop="limit_year" label="学制" width="80"/>
          <el-table-column label="国家级特色" width="100">
            <template #default="scope">
              {{ scope.row.nation_feature === '1' ? '是' : '否' }}
            </template>
          </el-table-column>
          <el-table-column label="软科等级" width="100">
            <template #default="scope">
              {{ scope.row.ruanke_level || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="软科排名" width="100">
            <template #default="scope">
              {{ scope.row.ruanke_rank === 9999 ? '-' : scope.row.ruanke_rank }}
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <!-- 录取分数线(使用抽离的地图组件) -->
      <ChinaMap :admission-scores="admissionScores"/>

      <!-- 学校视频 -->
      <el-card class="content-card" v-if="videos.length > 0">
        <h2 class="section-title">
          <el-icon>
            <VideoPlay/>
          </el-icon>
          学校视频
        </h2>
        <div class="videos-container">
          <div v-for="video in videos" :key="video.id" class="video-item">
            <el-image
                :src="video.img_url"
                fit="cover"
                class="video-thumbnail"
                @click="playVideo(video)"
            >
              <template #error>
                <div class="thumbnail-error">
                  <el-icon>
                    <VideoCamera/>
                  </el-icon>
                </div>
              </template>
              <template #placeholder>
                <div class="thumbnail-loading">
                  <el-icon>
                    <Loading/>
                  </el-icon>
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
        <video
            v-if="videoDialogVisible"
            :src="currentVideo.url"
            controls
            autoplay
            class="video-player"
            referrerpolicy="unsafe-url"
        ></video>
      </el-dialog>

      <!-- 百度地图抽屉 -->
      <el-drawer
          v-model="mapDrawerVisible"
          :title="'学校主位置:  ' + getSchoolName()"
          direction="rtl"
          size="46%"
          :before-close="handleDrawerClose"
      >
        <BaiduMap
            :address="locationAddress()"
            :zoom="16"
            class="map-drawer-content"
        />
      </el-drawer>
    </div>
  </div>
</template>

<script setup>
import BaiduMap from "@/components/baiduMap.vue"
import ChinaMap from "@/components/chinaMap.vue"
// 添加新的引用
import { ElMessage } from 'element-plus'

// 这个函数的功能仅仅只是为了解决IDE的bug:仅在script中检查组件使用,如果格式化代码就会删掉组件从而产生更加严重的bug
const useComponent = () => {
  return {
    BaiduMap,
    ChinaMap
  }
}

import {computed, onMounted, ref,nextTick} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import { InfoFilled, ArrowDown, ArrowUp } from '@element-plus/icons-vue'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const error = ref(false)

const loading = ref(true)           // 加载函数

const collapsedHeight = ref(100) // 折叠时的高度
const isTransitioning = ref(false)

// 后端返回的数据
const schoolDetail = ref({})        // 学校详情
const admissionScores = ref([])
const collegesDepartments = ref([])
const disciplineRankings = ref([])
const dualClassSubjects = ref([])
const specialPrograms = ref([])
const videos = ref([])

const videoDialogVisible = ref(false)
const currentVideo = ref({})
const selectedYear = ref(2024)
// const selectedYear = ref(new Date().getFullYear() - 1) // 默认显示去年数据
const selectedProvince = ref('')

// 特设专业高亮度变量
const hoverSpecial = ref(null)
const loadingSpecial = ref(false)

// 简介高亮显示
const contentRef = ref(null)
const isCollapsed = ref(true)
const canCollapse = ref(true)
const contentHeight = ref(0)

// baidu地图
const hoverAddress = ref(false)
const mapDrawerVisible = ref(false)

// vue宏: 接受父组件传入的参数
const props = defineProps({
  admissionScores: {
    type: Array,
    default: () => []
  }
})

// 计算属性函数: 获取部门分组
const groupedDepartments = computed(() => {
  const campuses = {}
  collegesDepartments.value.forEach(item => {
    if (!campuses[item.campus_name]) {
      campuses[item.campus_name] = {
        campusName: item.campus_name,
        departments: []
      }
    }
    campuses[item.campus_name].departments.push({
      department_id: item.department_id,
      department_name: item.department_name
    })
  })
  return Object.values(campuses)
})

// 计算属性函数: 获取排序后的isco排名
const sortedRankings = computed(() => {
  const rankOrder = {'A+': 1, 'A': 2, 'A-': 3, 'B+': 4, 'B': 5, 'B-': 6, 'C+': 7, 'C': 8, 'C-': 9}
  return [...disciplineRankings.value].sort((a, b) => {
    return (rankOrder[a.rank_level] || 10) - (rankOrder[b.rank_level] || 10)
  })
})

// 计算属性函数: 可用年份计算
const availableYears = computed(() => {
  const years = new Set()
  props.admissionScores.forEach(score => {
    years.add(score.year)
  })
  return Array.from(years).sort((a, b) => b - a) // 降序排列
})

// 功能: 显示地图抽屉
const showMapDrawer = () => {
  mapDrawerVisible.value = true
}

// 功能: 处理抽屉关闭
const handleDrawerClose = (done) => {
  done()
}

// 工具: 播放视频
const playVideo = (video) => {
  currentVideo.value = video
  videoDialogVisible.value = true
}

// 工具: 退出
const handleExit = () => router.back()

// 添加点击处理方法
const handleSpecialClick = async (specialName) => {
  if (loadingSpecial.value) return

  loadingSpecial.value = true
  try {
    const response = await axios.get(`/gapi/special/name?name=${encodeURIComponent(specialName)}`)

    if (response.data.code === 0 && response.data.data && response.data.data.id) {
      // 找到专业ID，跳转到专业详情
      router.push({
        path: '/specialDetail',
        query: { id: response.data.data.id }
      })
    } else if (response.data.msg && response.data.msg.includes('未找到')) {
      router.push('/professional')
    } else {
      ElMessage.error('获取专业信息失败，请稍后重试')
    }
  } catch (error) {
    console.error('请求专业信息失败:', error)
    if (error.response?.status === 404) {
      ElMessage.warning(`"${specialName}" 专业信息暂未收录`)
    } else {
      ElMessage.error('网络请求失败，请检查网络连接')
    }
  } finally {
    loadingSpecial.value = false
  }
}

// 工具: 返回学校名称拼接
const getSchoolName = () => schoolDetail.value.address

// 工具: 返回学校地址
const locationAddress = () => {
  const address = getFirstCampusAddress(schoolDetail.value.address)
  // 检查地址是否存在且是字符串
  if (!address || typeof address !== 'string' || address === '地址信息暂无') {
    // 默认地址:油专
    return '西安石油大学鄠邑校区'
  }

  // 地址格式可能直接是 地址,也可能是 校区:地址
  console.log('学校地址:', address)
  return address.includes('：') ? address.split('：')[1] : address
}

// 工具: 提取第一个校区地址
const getFirstCampusAddress = (address) => {
  if (!address || typeof address !== 'string') {
    return '地址信息暂无'
  }
  const trimmedAddress = address.trim()
  if (!trimmedAddress) {
    return '地址信息暂无'
  }
  // 检查是否包含多个校区（用逗号分隔）
  if (trimmedAddress.includes(',')) {
    const firstCampus = trimmedAddress.split(',')[0].trim()
    // 检查第一个校区是否包含冒号分隔
    if (firstCampus.includes(':')) {
      return firstCampus.split(':')[1].trim()
    }
    return firstCampus
  }
  // 如果没有逗号，但有冒号（单个校区但包含校区名称）
  if (trimmedAddress.includes(':')) {
    return trimmedAddress.split(':')[1].trim()
  }
  // 直接返回地址
  return trimmedAddress
}

// api函数: 获取学校详情
const fetchSchoolDetail = async () => {
  try {
    loading.value = true
    error.value = false

    const url = `gapi/schools/profiles/${route.query.id}`
    const response = await axios.get(url)
    if (response.data.code === 0) {

      const data = response.data.data

      schoolDetail.value = data.universities_detail
      admissionScores.value = data.admission_scores
      collegesDepartments.value = data.colleges_departments
      disciplineRankings.value = data.discipline_rankings
      dualClassSubjects.value = data.dual_class_subjects
      specialPrograms.value = data.special_programs
      videos.value = data.videos

      // 设置默认选中的年份为最新年份
      if (availableYears.value.length > 0) {
        selectedYear.value = availableYears.value[0]
      }
    } else {
      throw new Error('API返回错误')
    }
  } catch (err) {
    console.error('获取学校详情失败:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

// 处理内容，添加高亮
const formattedContent = computed(() => {
  if (!schoolDetail.value.content) return '暂无学校简介'

  let content = schoolDetail.value.content

  // 高亮重点词汇
  const highlightWords = ['211', '985', '双一流', 'ESI', 'QS', '博士', '硕士', '国家重点', '世界一流']
  highlightWords.forEach(word => {
    const regex = new RegExp(word, 'g')
    content = content.replace(regex, `<span class="highlight">${word}</span>`)
  })

  return content
})

// 检查是否需要折叠
const toggleContent = () => {
  if (isTransitioning.value) return;

  isTransitioning.value = true;

  // 如果需要展开
  if (isCollapsed.value) {
    isCollapsed.value = false;
    // 展开时不需要特殊处理，让CSS过渡自然工作
    setTimeout(() => {
      isTransitioning.value = false;
    }, 300); // 与CSS过渡时间保持一致
  }
  // 如果需要收起
  else {
    // 先滚动到可见区域顶部
    const cardElement = contentRef.value.closest('.content-card');
    if (cardElement) {
      cardElement.scrollIntoView({
        behavior: 'smooth',
        block: 'start'
      });
    }

    // 等待滚动完成后再收起内容
    setTimeout(() => {
      isCollapsed.value = true;
      setTimeout(() => {
        isTransitioning.value = false;
      }, 300); // ← 修改这个数字来控制滑动速度
    }, 400); // 稍微长于滚动时间
  }
}

// 修改检查折叠的方法
const checkCollapse = async () => {
  await nextTick()
  if (contentRef.value) {
    contentHeight.value = contentRef.value.scrollHeight
    canCollapse.value = contentHeight.value > collapsedHeight.value
  }
}

// vue生命周期函数: 创建时执行
onMounted(() => {
  // 检查是否已存在referrer meta标签
  let metaTag = document.querySelector('meta[name="referrer"]');

  if (!metaTag) {
    // 创建新的meta标签
    console.log('创建新的meta标签 referrer:no-referrer');
    metaTag = document.createElement('meta');
    metaTag.name = 'referrer';
    metaTag.content = 'no-referrer';
    document.getElementsByTagName('head')[0].appendChild(metaTag);
  } else {
    // 更新现有的meta标签
    console.log('更新现有的meta标签 referrer:no-referrer ');
    metaTag.content = 'no-referrer';
  }

  fetchSchoolDetail()

  checkCollapse()
})
</script>

<style lang="scss" scoped>
.school-detail-container {
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

.school-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eaeaea;

  .header-left {
    .school-title {
      font-size: 2rem;
      color: #2c3e50;
      margin: 0 0 8px 0;
    }

    .school-motto {
      color: #909399;
      font-size: 1.2rem;
      font-style: italic;
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

  .school-basic-info {
    display: flex;
    gap: 30px;

    .logo-container {
      flex-shrink: 0;

      .school-logo {
        width: 150px;
        height: 150px;
        border-radius: 8px;
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
      }

      .logo-error {
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

    .basic-info-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
      gap: 16px;
      flex-grow: 1;

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
}

/* 添加专业名称链接样式 */
.special-name-link {
  color: #409EFF;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 4px 8px;
  border-radius: 4px;
  display: inline-block;
}
.special-hover {
  background-color: #ecf5ff;
  color: #67c23a;
  text-decoration: underline;
  transform: translateX(2px);
}
.special-name-link .el-icon {
  margin-left: 4px;
}

.rich-text-content p {
  margin-bottom: 16px;
  line-height: 1.8;
  font-size: 14px;
  color: #606266;
}

.rich-text-content {
  line-height: 1.8;
  font-size: 14px;
  color: #606266;
  text-align: justify;
}

.content-card {
  margin-bottom: 24px;
  border-radius: 12px;

  .content-text {
    line-height: 1.8;
    color: #606266;
    text-align: justify;
  }
}

.contact-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;

  .contact-item {
    display: flex;
    align-items: center;

    .contact-label {
      color: #606266;
      font-weight: 500;
      margin-right: 10px;
      min-width: 60px;
    }

    .contact-value {
      color: #303133;

      a {
        color: #409EFF;
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}

.departments-container {
  .campus-section {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }

    .campus-name {
      color: #409EFF;
      margin-bottom: 12px;
      font-size: 1.1rem;
    }

    .departments-list {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;

      .department-tag {
        font-size: 1rem;
        padding: 8px 16px;
      }
    }
  }
}

.rankings-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;

  .ranking-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    background: #f8f9fa;
    border-radius: 8px;

    .rank-level {
      color: #606266;
      font-weight: 500;
    }

    .rank-count {
      color: #409EFF;
      font-weight: bold;
    }
  }
}

.subjects-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;

  .subject-tag {
    font-size: 1rem;
    padding: 8px 16px;
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

}

@media (max-width: 768px) {
  .school-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;

    .exit-button {
      align-self: flex-end;
    }
  }

  .info-card .school-basic-info {
    flex-direction: column;
  }

  .contact-grid {
    grid-template-columns: 1fr;
  }


  .videos-container {
    grid-template-columns: 1fr;
  }
}

.address-highlight {
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 1px 1px;
  border-radius: 4px;
}

.address-hover {
  background-color: #ecf5ff;
  color: #409EFF;
  text-decoration: underline;
}

.map-drawer-content {
  height: 100%;
  width: 100%;
}

.content-card {
  margin: 12px 0;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-title {
  display: flex;
  align-items: center;
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2f3d;
}

.section-title .el-icon {
  margin-right: 8px;
  color: #409eff;
}

.toggle-btn {
  text-align: center;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed #e8e8e8;
}

.toggle-button {
  font-size: 13px;
  font-weight: 500;
}

.toggle-button .el-icon {
  margin-left: 4px;
  transition: transform 0.3s ease;
}

.rich-text-content {
  line-height: 1.7;
  font-size: 15px;
  color: #4a5568;
  text-align: justify;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;

  /* 添加平滑过渡效果 */
  transition: all 0.3s ease-in-out;
  overflow: hidden;

  /* 移除之前的mask渐变效果 */
  mask-image: none;
  -webkit-mask-image: none;
}

/* 移除之前的.collapsed类样式，改用内联样式控制高度 */
.rich-text-content p {
  margin: 8px 0;
  text-indent: 2em;
  transition: opacity 0.3s ease-in-out;
}


.rich-text-content.collapsed {
  max-height: 90px;
  overflow: hidden;
  mask-image: linear-gradient(to bottom, black 60%, transparent 100%);
  -webkit-mask-image: linear-gradient(to bottom, black 60%, transparent 100%);
}

.rich-text-content .highlight {
  background: linear-gradient(120deg, #e6f7ff 0%, #bae7ff 100%);
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
  color: #7cb3ea;
  border: 1px solid #a6d3f5;
  margin: 0 2px;
  transition: all 0.3s ease-in-out;
}

.rich-text-content strong {
  color: #1f2f3d;
  font-weight: 600;
}


.rich-text-content br {
  display: block;
  content: "";
  margin: 6px 0;
}

/* 添加一个微妙的阴影指示可折叠 */
.content-card {
  position: relative;
  transition: box-shadow 0.3s ease;
}

.content-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

/* 确保卡片内容有适当的内边距 */
.content-card :deep(.el-card__body) {
  padding: 20px;
}

::v-deep(.el-drawer__header){
  margin-bottom: 0 !important;
}
</style>
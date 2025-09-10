<template>
  <div class="college-container">


    <!-- 查询条件区域 -->
    <el-card class="filter-card">
      <div class="filter-header">
      <h2 class="section-title">大学查询</h2>
      <el-button
          class="exit-button"
          type="primary"
          @click="handleExit"
          round
      >
        <el-icon><Close /></el-icon>
        <span>返回主页</span>
      </el-button>
      </div>
      <el-form :model="queryParams" label-width="100px" class="filter-form">
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="学校名称">
              <el-input
                  v-model="queryParams.name"
                  placeholder="输入学校名称"
                  clearable
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="省份">
              <el-select
                  v-model="queryParams.province_id"
                  placeholder="选择省份"
                  clearable
              >
                <!-- 省份选项需要根据实际情况填充 -->
                <el-option
                    v-for="province in provinces"
                    :key="province.value"
                    :label="province.label"
                    :value="province.value"
                />
                <!-- 更多省份选项 -->
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="学校类型">
              <el-select
                  v-model="queryParams.type_name"
                  placeholder="选择学校类型"
                  clearable
              >
                <el-option label="综合类" value="综合类" />
                <el-option label="理工类" value="理工类" />
                <el-option label="师范类" value="师范类" />
                <el-option label="农林类" value="农林类" />
                <el-option label="医药类" value="医药类" />
                <el-option label="语言类" value="语言类" />
                <el-option label="财经类" value="财经类" />
                <el-option label="政法类" value="政法类" />
                <el-option label="艺术类" value="艺术类" />
                <el-option label="民族类" value="民族类" />
                <el-option label="体育类" value="体育类" />
                <el-option label="军事类" value="军事类" />
                <el-option label="其他" value="其他" />
                <el-option label="未知" value="" />

              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="办学层次">
              <el-select
                  v-model="queryParams.level_name"
                  placeholder="选择办学层次"
                  clearable
              >
                <el-option label="本科" value="本科" />
                <el-option label="专科（高职）" value="专科（高职）" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="211高校">
              <el-select
                  v-model="queryParams.is_211"
                  placeholder="是否211"
                  clearable
              >
                <el-option label="是" :value="1" />
                <el-option label="否" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="985高校">
              <el-select
                  v-model="queryParams.is_985"
                  placeholder="是否985"
                  clearable
              >
                <el-option label="是" :value="1" />
                <el-option label="否" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="双一流">
              <el-select
                  v-model="queryParams.dual_class"
                  placeholder="是否双一流"
                  clearable
              >
                <el-option label="是" value="双一流" />
                <el-option label="否" value="" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12" :md="8" :lg="6">
            <el-form-item label="排序方式">
              <el-select
                  v-model="queryParams.order_by"
                  placeholder="选择排序字段"
                  clearable
              >
                <el-option label="默认" value="" />
                <el-option label="学校名称" value="name" />
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
      <el-button @click="fetchColleges" type="primary">重新加载</el-button>
    </div>

    <!-- 结果区域 -->
    <div v-else>
      <!-- 统计信息 -->
      <el-card class="stats-card">
        <div class="stats-info">
          <span>共找到 <strong>{{ total }}</strong> 所符合条件的大学</span>
          <span>第 <strong>{{ queryParams.page }}</strong> 页 / 共 <strong>{{ totalPages }}</strong> 页</span>
        </div>
      </el-card>

      <!-- 大学列表 -->
      <div class="college-list">
        <el-card
            v-for="university in universities"
            :key="university.id"
            class="college-card"
            :class="{ 'card-hover': hoverCardId === university.id }"
            @mouseenter="handleMouseEnter(university.id)"
            @mouseleave="handleMouseLeave"
            @click="navigateToCollegeDetail(university)"
        >
          <div class="college-content">
            <div class="college-logo">
              <el-image
                  :src="university.logo_url"
                  fit="cover"
                  class="logo-image"
              >
                <template #error>
                  <div class="logo-error">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
                <template #placeholder>
                  <div class="logo-loading">
                    <el-icon><Loading /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <div class="college-info">
              <h3 class="college-name">{{ university.name }}</h3>
              <p class="college-motto" v-if="university.motto">{{ university.motto }}</p>
              <div class="college-details">
                <div class="detail-item">
                  <span class="detail-label">所在地：</span>
                  <span class="detail-value">{{ university.province_name }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">办学类型：</span>
                  <span class="detail-value">{{ university.type_name }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">办学层次：</span>
                  <span class="detail-value">{{ university.level_name }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">创办时间：</span>
                  <span class="detail-value">{{ university.create_date }}</span>
                </div>
              </div>
              <div class="college-badges">
                <el-tag v-if="university.f211 === 1" type="danger" size="small">211</el-tag>
                <el-tag v-if="university.f985 === 1" type="warning" size="small">985</el-tag>
                <el-tag v-if="university.dual_class_name === '双一流'" type="success" size="small">双一流</el-tag>
              </div>
            </div>
          </div>

          <!-- 悬浮提示 -->
          <div v-if="hoverCardId === university.id" class="hover-tooltip">
          </div>
        </el-card>
      </div>

      <!-- 分页控件 -->
      <div class="pagination-container">
        <el-pagination
            v-model:current-page="queryParams.page"
            v-model:page-size="queryParams.page_size"
            :page-sizes="[12, 21, 51, 99]"
            :small="true"
            :background="true"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import {Picture, Loading, Link, Close} from '@element-plus/icons-vue'

const router = useRouter()
const hoverCardId = ref(null)
const loading = ref(false)
const error = ref(false)
const universities = ref([])
const total = ref(0)
const totalPages = ref(0)
const isHovering = ref(false)

const queryParams = reactive({
  page: 1,
  page_size: 21,
  name: '',
  province_id: null,
  type_name: '',
  level_name: '',
  is_211: null,
  is_985: null,
  dual_class: '',
  order_by: '',
  order_desc: false
})
const provinces = [
  // 中国34个省级行政区
  {label: '北京', value: 11}, {label: '天津', value: 12}, {label: '河北', value: 13}, {label: '山西', value: 14}, {label: '内蒙古', value: 15},
  {label: '辽宁', value: 21}, {label: '吉林', value: 22}, {label: '黑龙江', value: 23}, {label: '上海', value: 31}, {label: '江苏', value: 32},
  {label: '浙江', value: 33}, {label: '安徽', value: 34}, {label: '福建', value: 35}, {label: '江西', value: 36}, {label: '山东', value: 37},
  {label: '河南', value: 41}, {label: '湖北', value: 42}, {label: '湖南', value: 43}, {label: '广东', value: 44}, {label: '广西', value: 45},
  {label: '海南', value: 46}, {label: '重庆', value: 50}, {label: '四川', value: 51}, {label: '贵州', value: 52}, {label: '云南', value: 53},
  {label: '西藏', value: 54}, {label: '陕西', value: 61}, {label: '甘肃', value: 62}, {label: '青海', value: 63}, {label: '宁夏', value: 64},
  {label: '新疆', value: 65}, {label: '台湾', value: 71}, {label: '香港', value: 81}, {label: '澳门', value: 82}
]
// 获取大学数据
const fetchColleges = async () => {
  try {
    loading.value = true
    error.value = false

    // 构建查询参数
    const params = new URLSearchParams()
    Object.keys(queryParams).forEach(key => {
      if (queryParams[key] !== null && queryParams[key] !== '') {
        params.append(key, queryParams[key])
      }
    })
    const response = await axios.get(`gapi/gAllSchools/?${params}`)
    console.log(response)

    if (response.data) {
      universities.value = response.data.universities
      total.value = response.data.total
      totalPages.value = response.data.total_pages
    }
  } catch (err) {
    console.error('获取大学数据失败:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

// 处理搜索
const handleSearch = () => {
  queryParams.page = 1 // 重置到第一页
  fetchColleges()
}

// 处理重置
const handleReset = () => {
  Object.assign(queryParams, {
    page: 1,
    page_size: 21,
    name: '',
    province_id: null,
    type_name: '',
    level_name: '',
    is_211: null,
    is_985: null,
    dual_class: '',
    order_by: '',
    order_desc: false
  })
  fetchColleges()
}

// 处理每页数量变化
const handleSizeChange = (newSize) => {
  queryParams.page_size = newSize
  queryParams.page = 1
  fetchColleges()
}

// 处理页码变化
const handleCurrentChange = (newPage) => {
  queryParams.page = newPage
  fetchColleges()
}

// 处理鼠标进入
const handleMouseEnter = (id) => {
  isHovering.value = true
  hoverCardId.value = id
}

// 处理鼠标离开
const handleMouseLeave = () => {
  isHovering.value = false
  hoverCardId.value = null
}

// 跳转到学校详情页
const navigateToCollegeDetail = (university) => {
  router.push({
    path: '/schoolDetail',
    query: { id: university.id }
  })
}

// 工具: 退出
const handleExit = () => {
  router.push('/zxwl')
}


// 初始化加载
onMounted(() => {
  fetchColleges()
})
</script>

<style lang="scss" scoped>
.college-container {
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

.filter-card {
  margin-bottom: 24px;
  border-radius: 12px;
}

.filter-form {
  margin-top: 20px;
}

.section-title {
  color: #409EFF;
  font-size: 27px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
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

.college-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.college-card {
  border-radius: 12px;
  transition: transform 0.3s, box-shadow 0.3s;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  }
}

.college-content {
  display: flex;
  gap: 16px;
}

.college-logo {
  flex-shrink: 0;
  width: 80px;
  height: 80px;

  .logo-image {
    width: 100%;
    height: 100%;
    border-radius: 8px;
  }

  .logo-error,
  .logo-loading {
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
.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.college-info {
  flex: 1;
  min-width: 0;

  .college-name {
    font-size: 1.2rem;
    color: #2c3e50;
    margin: 0 0 8px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .college-motto {
    color: #909399;
    font-style: italic;
    margin: 0 0 12px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .college-details {
    margin-bottom: 12px;

    .detail-item {
      display: flex;
      margin-bottom: 4px;

      .detail-label {
        color: #606266;
        font-weight: 500;
        min-width: 70px;
      }

      .detail-value {
        color: #303133;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }

  .college-badges {
    display: flex;
    gap: 8px;
  }
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

@media (max-width: 768px) {
  .college-list {
    grid-template-columns: 1fr;
  }

  .college-content {
    flex-direction: column;
    text-align: center;
  }

  .college-logo {
    align-self: center;
  }

  .stats-info {
    flex-direction: column;
    gap: 10px;
  }
}

.college-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.college-card {
  border-radius: 12px;
  transition: all 0.3s ease;
  position: relative;
  cursor: pointer;
  overflow: hidden;

  // 默认状态
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
  border: 1px solid #ebeef5;

  // 悬浮状态
  &.card-hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
    border-color: #409EFF;

    .college-name {
      color: #409EFF;
    }
  }
}

.college-content {
  display: flex;
  gap: 16px;
  position: relative;
  z-index: 2;
}

.college-logo {
  flex-shrink: 0;
  width: 80px;
  height: 80px;

  .logo-image {
    width: 100%;
    height: 100%;
    border-radius: 8px;
    transition: transform 0.3s ease;
  }

  .logo-error,
  .logo-loading {
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

.college-card.card-hover .college-logo .logo-image {
  transform: scale(1.05);
}

.college-info {
  flex: 1;
  min-width: 0;

  .college-name {
    font-size: 1.2rem;
    color: #2c3e50;
    margin: 0 0 8px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: color 0.3s ease;
  }

  .college-motto {
    color: #909399;
    font-style: italic;
    margin: 0 0 12px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .college-details {
    margin-bottom: 12px;

    .detail-item {
      display: flex;
      margin-bottom: 4px;

      .detail-label {
        color: #606266;
        font-weight: 500;
        min-width: 70px;
      }

      .detail-value {
        color: #303133;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }

  .college-badges {
    display: flex;
    gap: 8px;
  }
}

// 悬浮提示
.hover-tooltip {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
      to bottom,
      rgba(64, 158, 255, 0.05) 0%,
      rgba(64, 158, 255, 0.1) 100%
  );
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  animation: fadeIn 0.3s ease forwards;
  z-index: 1;
  border-radius: 12px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.tooltip-content {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  color: #409EFF;
  font-weight: 500;

  .el-icon {
    font-size: 1.1rem;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .college-list {
    grid-template-columns: 1fr;
  }

  .college-content {
    flex-direction: column;
    text-align: center;
  }

  .college-logo {
    align-self: center;
  }

  .hover-tooltip {
    display: none; // 在移动端隐藏悬浮提示
  }

  .college-card {
    // 移动端点击效果
    &:active {
      transform: scale(0.98);
    }
  }
}
</style>
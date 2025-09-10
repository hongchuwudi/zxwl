<template>
  <div class="search-results-container">
    <div class="search-bar">
      <span class="search-title">搜索结果</span>
      <el-button
          class="exit-button"
          type="primary"
          @click="handleExit"
          round
      >
        <el-icon>
          <Close/>
        </el-icon>
        <span>返回主页</span>
      </el-button>
    </div>
    <div>
      <!-- 搜索结果统计 -->
      <el-card class="stats-card" v-if="showResults">
        <div class="stats-info">
          <span>共找到 <strong>{{ totalResults }}</strong> 条相关结果</span>
          <span v-if="hasSchools">包含 <strong>{{ schoolCount }}</strong> 所学校</span>
          <span v-if="hasProfessionals">包含 <strong>{{ professionalCount }}</strong> 个专业</span>
        </div>
      </el-card>

      <!-- 学校搜索结果 -->
      <div v-if="hasSchools" class="results-section">
        <h3 class="section-title">学校搜索结果</h3>
        <div class="schools-list">
          <el-card
              v-for="school in resData.school_profile_res"
              :key="'school-' + school.id"
              class="result-card school-card"
              @click="navigateToSchoolDetail(school)"
          >
            <div class="result-content">
              <div class="result-logo">
                <el-image
                    :src="school.logo_url"
                    fit="cover"
                    class="logo-image"
                >
                  <template #error>
                    <div class="logo-error">
                      <el-icon>
                        <Picture/>
                      </el-icon>
                    </div>
                  </template>
                  <template #placeholder>
                    <div class="logo-loading">
                      <el-icon>
                        <Loading/>
                      </el-icon>
                    </div>
                  </template>
                </el-image>
              </div>
              <div class="result-info">
                <h4 class="result-name">{{ school.name }}</h4>
                <p class="result-motto" v-if="school.motto">{{ school.motto }}</p>
                <div class="result-details">
                  <div class="detail-item">
                    <span class="detail-label">所在地：</span>
                    <span class="detail-value">{{ school.province_name }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">办学类型：</span>
                    <span class="detail-value">{{ school.type_name }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">办学层次：</span>
                    <span class="detail-value">{{ school.level_name }}</span>
                  </div>
                </div>
                <div class="result-badges">
                  <el-tag v-if="school.f211 === 1" type="danger" size="small">211</el-tag>
                  <el-tag v-if="school.f985 === 1" type="warning" size="small">985</el-tag>
                  <el-tag v-if="school.dual_class_name === '双一流'" type="success" size="small">双一流</el-tag>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 专业搜索结果 -->
      <div v-if="hasProfessionals" class="results-section">
        <h3 class="section-title">专业搜索结果</h3>
        <div class="professionals-list">
          <el-card
              v-for="professional in resData.professional_item_res_arr"
              :key="'professional-' + professional.id"
              class="result-card professional-card"
              @click="navigateToProfessionalDetail(professional)"
          >
            <div class="result-content">
              <div class="professional-icon">
                <el-icon class="icon">
                  <Reading/>
                </el-icon>
              </div>
              <div class="result-info">
                <h4 class="result-name">{{ professional.name }}</h4>
                <div class="result-details">
                  <div class="detail-item">
                    <span class="detail-label">专业代码：</span>
                    <span class="detail-value">{{ professional.code }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">学位：</span>
                    <span class="detail-value">{{ professional.degree }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">学制：</span>
                    <span class="detail-value">{{ professional.limit_year }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">学科门类：</span>
                    <span class="detail-value">{{ professional.level1_name }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">专业类：</span>
                    <span class="detail-value">{{ professional.level2_name }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">平均薪资：</span>
                    <span class="detail-value salary">{{
                        professional.salaryavg ? `${professional.salaryavg}元` : '暂无数据'
                      }}</span>
                  </div>
                </div>
                <div class="result-badges">
                  <el-tag v-if="professional.employment_rate" type="info" size="small">
                    就业率: {{ (professional.employment_rate * 100).toFixed(1) }}%
                  </el-tag>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!showResults" class="empty-state">
        <el-empty description="暂无搜索结果"/>
        <p class="empty-tip">尝试使用不同的关键词进行搜索</p>
      </div>
    </div>
  </div>

</template>

<script setup>
import {computed, onMounted, ref, watch} from 'vue' // 添加watch导入
import {useRoute, useRouter} from 'vue-router'

import axios from "axios";

const router = useRouter()
const route = useRoute()
const resData = ref({})

// 计算属性
const showResults = computed(() => hasSchools.value || hasProfessionals.value)
const hasSchools = computed(() => resData.value.school_profile_res?.length > 0)
const hasProfessionals = computed(() => resData.value.professional_item_res_arr?.length > 0)
const schoolCount = computed(() => resData.value.school_profile_res?.length || 0)
const professionalCount = computed(() => resData.value.professional_item_res_arr?.length || 0)
const totalResults = computed(() => schoolCount.value + professionalCount.value)

// 导航方法
const navigateToSchoolDetail = (school) => router.push({path: '/schoolDetail', query: {id: school.id}})
const navigateToProfessionalDetail = (professional) => router.push({path: '/specialDetail', query: {id: professional.id}})

const handleExit = () => router.push('/zxwl')
const fetchUAS = async (keyword) => {  // 添加keyword参数
  console.log('搜索关键词:', keyword)
  if (!keyword) return;  // 添加空值检查
  try {
    const response = await axios.get(`gapi/searchSchAndSpe/?keyword=${route.query.keyword}`)
    resData.value = response.data.data
    console.log('搜索结果:', resData.value)
  } catch (error) {
    console.error('搜索请求失败:', error)
  }
}

onMounted(() => fetchUAS(route.query.keyword)) // 组件挂载时立即执行一次搜索
watch(() => route.query.keyword, newKeyword => fetchUAS(newKeyword)) // 添加对keyword的监听
</script>

<style lang="scss" scoped>
.search-results-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.stats-card {
  margin-bottom: 24px;
  border-radius: 12px;

  .stats-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    flex-wrap: wrap;
    gap: 10px;

    span {
      display: flex;
      align-items: center;
      gap: 5px;
    }
  }
}

.results-section {
  margin-bottom: 32px;

  .section-title {
    color: #409EFF;
    font-size: 1.5rem;
    margin-bottom: 20px;
    padding-left: 10px;
    border-left: 4px solid #409EFF;
  }
}

.schools-list, .professionals-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.result-card {
  border-radius: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
  border: 1px solid #ebeef5;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
    border-color: #409EFF;

    .result-name {
      color: #409EFF;
    }
  }
}

.result-content {
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.result-logo, .professional-icon {
  flex-shrink: 0;
  width: 60px;
  height: 60px;

  .logo-image {
    width: 100%;
    height: 100%;
    border-radius: 8px;
    transition: transform 0.3s ease;
  }

  .logo-error, .logo-loading {
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

.professional-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #409EFF, #79BBFF);
  border-radius: 8px;

  .icon {
    font-size: 2rem;
    color: white;
  }
}

.result-card:hover .result-logo .logo-image {
  transform: scale(1.05);
}

.result-info {
  flex: 1;
  min-width: 0;

  .result-name {
    font-size: 1.1rem;
    color: #2c3e50;
    margin: 0 0 8px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: color 0.3s ease;
  }

  .result-motto {
    color: #909399;
    font-style: italic;
    margin: 0 0 12px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 0.9rem;
  }

  .result-details {
    margin-bottom: 12px;

    .detail-item {
      display: flex;
      margin-bottom: 4px;
      font-size: 0.9rem;

      .detail-label {
        color: #606266;
        font-weight: 500;
        min-width: 70px;
        flex-shrink: 0;
      }

      .detail-value {
        color: #303133;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &.salary {
          color: #e6a23c;
          font-weight: 600;
        }
      }
    }
  }

  .result-badges {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }
}

.empty-state {
  text-align: center;
  padding: 60px 0;

  .empty-tip {
    color: #909399;
    margin-top: 10px;
    font-size: 0.9rem;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .search-results-container {
    padding: 10px;
  }

  .schools-list, .professionals-list {
    grid-template-columns: 1fr;
  }

  .stats-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .result-content {
    flex-direction: column;
    text-align: center;
  }

  .result-logo, .professional-icon {
    align-self: center;
  }

  .result-details .detail-item {
    justify-content: center;
  }
}
@media (max-width: 480px) {
  .result-card {
    margin: 0 -10px;
    border-radius: 0;
    border-left: none;
    border-right: none;
  }
}

// 搜索栏样式
.search-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  border-radius: 12px;
  margin-bottom: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border: 1px solid #dcdfe6;
  .search-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #409EFF;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    letter-spacing: 0.5px;
  }
  .exit-button {
    background: linear-gradient(135deg, #409EFF 0%, #79BBFF 100%);
    border: none;
    padding: 10px 20px;
    transition: all 0.3s ease;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
      background: linear-gradient(135deg, #79BBFF 0%, #409EFF 100%);
    }

    &:active {
      transform: translateY(0);
      box-shadow: 0 2px 6px rgba(64, 158, 255, 0.3);
    }

    .el-icon {
      margin-right: 6px;
      font-size: 1.1em;
    }

    span {
      font-weight: 500;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .search-bar {
    flex-direction: column;
    gap: 16px;
    padding: 16px;

    .search-title {
      font-size: 1.3rem;
    }

    .exit-button {
      width: 100%;
      justify-content: center;
    }
  }
}
@media (max-width: 480px) {
  .search-bar {
    border-radius: 8px;
    padding: 12px;

    .search-title {
      font-size: 1.2rem;
    }
  }
}
</style>
<template>
  <div class="score-container">
    <a-page-header
        title="高考分数分析"
        @back="() => $router.go(-1)"
        class="page-header"
    >
      <template #extra>
        <div class="header-buttons">
          <a-tooltip placement="bottom" title="刷新数据">
            <a-button type="primary" @click="loadScoreData" size="large" shape="circle">
              <template #icon><ReloadOutlined /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip placement="bottom" title="导出数据">
            <a-button type="default" @click="exportData" size="large" shape="circle">
              <template #icon><DownloadOutlined /></template>
            </a-button>
          </a-tooltip>
        </div>
      </template>
    </a-page-header>

    <a-spin :spinning="loading">
      <div class="score-content">
        <!-- 筛选条件 -->
        <a-card :bordered="false" class="filter-card">
          <a-space :size="20">
            <!--省份-->
            <div class="filter-item">
              <label>省份：</label>
              <a-select
                  v-model:value="filterParams.provinceId"
                  style="width: 200px"
                  placeholder="请选择省份"
                  @change="handleProvinceChange"
              >
                <a-select-option v-for="province in provinces" :key="province.value" :value="province.value" >
                  {{ province.label }}
                </a-select-option>
              </a-select>
            </div>
            <!--类型-->
            <div class="filter-item">
              <label>高考类型：</label>
              <a-select
                  v-model:value="filterParams.typeId"
                  style="width: 200px"
                  placeholder="请选择高考类型"
                  @change="handleTypeChange"
              >
                <a-select-option v-for="type in examTypeOptions" :key="type.id" :value="type.id">
                  {{ type.name }}
                </a-select-option>
              </a-select>
            </div>
            <!--年份-->
            <div class="filter-item">
              <label>年份：</label>
              <a-select
                  v-model:value="filterParams.year"
                  style="width: 120px"
                  placeholder="请选择年份"
                  @change="handleYearChange"
              >
                <a-select-option v-for="year in availableYears" :key="year" :value="year">
                  {{ year }}年
                </a-select-option>
              </a-select>
            </div>
            <!--批-->
            <div class="filter-item" v-if="filterParams.provinceId === 11 || filterParams.provinceId === 31 || filterParams.provinceId === 12">
              <label>批次：</label>
              <a-select
                  v-model:value="filterParams.batch"
                  :label="availableBatches[filterParams.batch]"
                  style="width: 120px"
                  placeholder="请选择批次"
              >
                <a-select-option v-for="batch in availableBatches" :key="batch.value" :value="batch.value">
                  {{batch.label}}
                </a-select-option>
              </a-select>
            </div>


            <a-button type="primary" @click="loadScoreData" :loading="loading">
              <template #icon><SearchOutlined /></template>
              查询
            </a-button>
          </a-space>
        </a-card>

        <!-- 图表展示区域 -->
        <a-card :bordered="false" class="chart-card">
          <template #title>
            <div class="chart-title">
              <span>分数分布图表</span>
              <a-radio-group v-model:value="chartType" button-style="solid">
                <a-radio-button value="line">折线图</a-radio-button>
                <a-radio-button value="bar">柱状图</a-radio-button>
<!--                <a-radio-button value="pie">饼图</a-radio-button>-->
              </a-radio-group>
            </div>
          </template>

          <div class="chart-container">
            <div v-if="scoreData.sections.length > 0" ref="chartRef" style="height: 400px; width: 100%;"></div>
            <a-empty v-else description="暂无数据" />
          </div>
        </a-card>

        <!-- 统计信息 -->
        <a-row :gutter="20" class="stats-row">
          <a-col :span="6">
            <a-statistic title="平均分" :value="scoreData.scoreStat?.average_score || 0" :precision="1" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="最高分" :value="scoreData.scoreStat?.max_score || 0" :precision="1" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="最低分" :value="scoreData.scoreStat?.min_score || 0" :precision="1" />
          </a-col>
          <a-col :span="6">
            <a-statistic title="中位数" :value="scoreData.scoreStat?.median_score || 0" :precision="1" />
          </a-col>
        </a-row>

        <!-- 分数段表格 -->
        <a-card :bordered="false" class="table-card">
          <template #title>
            <span>一分一段表</span>
            <span class="total-students">总人数：{{ scoreData.examInfo?.total_num || 0 }}</span>
          </template>

          <a-table
              :columns="columns"
              :data-source="scoreData.sections"
              :pagination="false"
              :scroll="{ y: 400 }"
              size="middle"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'percentage'">
                {{ record.percentage.toFixed(3) }}%
              </template>
            </template>
          </a-table>
        </a-card>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import {
  SearchOutlined,
  ReloadOutlined,
  DownloadOutlined
} from '@ant-design/icons-vue'
import axios from "axios"

const router = useRouter()

// 筛选参数
const filterParams = reactive({
  provinceId: 11, // 默认北京
  typeId: 3,     // 默认理科
  year: new Date().getFullYear(), // 默认去年
  batch: 1
})

// 状态变量
const loading = ref(false)
const chartType = ref('line')
const chartRef = ref(null)
let chartInstance = null

// 数据
const scoreData = reactive({
  examInfo: null,
  sections: [],
  scoreStat: null
})

const availableYears = ref([2025,2024,2023, 2022, 2021,2020])
const availableBatches = ref([
  {label: '本科', value: 1},
  {label: '专科', value: 2},
])
// 省份选项
const provinces = [
  {label: '北京', value: 11}, {label: '天津', value: 12}, {label: '河北', value: 13}, {label: '山西', value: 14}, {label: '内蒙古', value: 15},
  {label: '辽宁', value: 21}, {label: '吉林', value: 22}, {label: '黑龙江', value: 23}, {label: '上海', value: 31}, {label: '江苏', value: 32},
  {label: '浙江', value: 33}, {label: '安徽', value: 34}, {label: '福建', value: 35}, {label: '江西', value: 36}, {label: '山东', value: 37},
  {label: '河南', value: 41}, {label: '湖北', value: 42}, {label: '湖南', value: 43}, {label: '广东', value: 44}, {label: '广西', value: 45},
  {label: '海南', value: 46}, {label: '重庆', value: 50}, {label: '四川', value: 51}, {label: '贵州', value: 52}, {label: '云南', value: 53},
  {label: '西藏', value: 54}, {label: '陕西', value: 61}, {label: '甘肃', value: 62}, {label: '青海', value: 63}, {label: '宁夏', value: 64},
  {label: '新疆', value: 65}, {label: '台湾', value: 71}, {label: '香港', value: 81}, {label: '澳门', value: 82}
]
// 高考类型选项
const examTypeOptions = [
  { id: 2073, name: '新高考物理类' }, { id: 2074, name: '新高考历史类' },
  { id: 1, name: '旧高考理科' }, { id: 2, name: '旧高考文科' },
  { id: 3, name: '综合' },
  { id: 5, name: '体育类' }, { id: 23, name: '体育文' }, { id: 24, name: '体育理' },
  { id: 25, name: '艺术文' }, { id: 26, name: '艺术理' },   { id: 4, name: '艺术类' }, { id: 2292, name: '艺术类(历史)' }, { id: 2293, name: '艺术类(物理)' },
  { id: 31, name: '蒙授体育' }, { id: 32, name: '蒙授文科' }, { id: 33, name: '蒙授理科' },
  { id: 2294, name: '体育类(历史)' }, { id: 2295, name: '体育类(物理)' }
]
// 表格列定义
const columns = [
  {
    title: '分数段',
    dataIndex: 'score_range',
    key: 'score_range',
    width: 100
  },
  {
    title: '最低分',
    dataIndex: 'min_score',
    key: 'min_score',
    width: 80
  },
  {
    title: '最高分',
    dataIndex: 'max_score',
    key: 'max_score',
    width: 80
  },
  {
    title: '排名范围',
    dataIndex: 'rank_range',
    key: 'rank_range',
    width: 120
  },
  {
    title: '人数',
    dataIndex: 'total_students',
    key: 'total_students',
    width: 80
  },
  {
    title: '占比',
    dataIndex: 'percentage',
    key: 'percentage',
    width: 80
  }
]

const loadAvailableYears = async () => {
  try {
    const response = await axios.get('/gapi/scores/years', {
      params: {
        province_id: filterParams.provinceId,
        type_id: filterParams.typeId
      }
    })
    if (response.data.code === 0) {
      availableYears.value = response.data.data
    }
  } catch (error) {
    console.error('获取可用年份失败:', error)
  }
}

// 加载分数数据
const loadScoreData = async () => {
  loading.value = true
  try {
    const response = await axios.get('/gapi/scores/data', {
      params: {
        province_id: filterParams.provinceId,
        type_id: filterParams.typeId,
        year: filterParams.year,
        batch: filterParams.batch
      }
    })

    if (response.data.code === 0) {
      scoreData.examInfo = response.data.data.exam_info
      scoreData.sections = response.data.data.sections
      scoreData.scoreStat = response.data.data.score_stat
      console.log('scoreData:', scoreData)
      await nextTick(() => {
        renderChart()
      })
    } else {
      console.log('获取分数数据失败:', response.data.message)
      if(response.data.message === '获取分数数据失败: record not found') message.error('暂无相关数据')
      else message.error('获取分数数据失败')
    }
  } catch (error) {
    message.error('暂未相关数据')
  } finally {
    loading.value = false
  }
}


const renderChart = () => {
  if (!chartRef.value || scoreData.sections.length === 0) return

  if (chartInstance) {
    chartInstance.dispose()
  }

  chartInstance = echarts.init(chartRef.value)

  const chartData = scoreData.sections.map(section => ({
    name: section.score_range,
    value: section.total_students,
    minScore: section.min_score,
    maxScore: section.max_score
  }))

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: (params) => {
        const data = params[0].data
        return `分数段: ${data.name}<br/>人数: ${data.value}<br/>分数: ${data.minScore}-${data.maxScore}`
      }
    },
    xAxis: {
      type: 'category',
      data: chartData.map(item => item.name),
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: {
      type: 'value',
      name: '人数'
    },
    series: [{
      type: chartType.value,
      data: chartData,
      ...(chartType.value === 'pie' ? {
        radius: '50%',
        label: {
          formatter: '{b}: {c} ({d}%)'
        }
      } : {})
    }]
  }

  chartInstance.setOption(option)
}
const exportData = () => {
  if (scoreData.sections.length === 0) {
    message.warning('没有可导出的数据')
    return
  }

  // 导出逻辑
  message.success('数据导出功能开发中')
}
const handleProvinceChange = () => loadAvailableYears()
const handleTypeChange = () => loadAvailableYears()
const handleYearChange = () => loadScoreData()

// 生命周期
onMounted(async () => {
  await loadAvailableYears()
  await loadScoreData()
})

// 监听图表类型变化
watch(chartType, () => {
  if (scoreData.sections.length > 0) {
    renderChart()
  }
})

// 监听窗口大小变化，重新渲染图表
window.addEventListener('resize', () => {
  if (chartInstance) {
    chartInstance.resize()
  }
})
</script>

<style scoped>
.score-container {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.page-header {
  background-color: #fff;
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.score-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.filter-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.filter-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-item label {
  font-weight: 500;
  color: #262626;
  min-width: 60px;
}

.chart-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.chart-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  margin-top: 16px;
}

.stats-row {
  margin-bottom: 20px;
  ::v-deep(.ant-statistic) {
    text-align: center;
  }
}


.table-card {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.total-students {
  margin-left: 16px;
  color: #8c8c8c;
  font-size: 14px;
}

.header-buttons {
  display: flex;
  gap: 8px;
}

::v-deep(.ant-table-thead > tr > th) {
  background-color: #fafafa;
  font-weight: 500;
}

::v-deep(.ant-table-tbody > tr:hover > td) {
  background-color: #f5f5f5;
}
</style>
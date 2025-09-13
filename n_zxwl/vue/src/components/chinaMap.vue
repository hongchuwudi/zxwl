<template>
  <!-- 录取分数线 -->
  <el-card class="content-card" v-if="admissionScores.length > 0">
    <h2 class="section-title"><el-icon><DataLine /></el-icon> 录取分数线</h2>

    <!-- 筛选控件 -->
    <div class="scores-filter">
      <el-select v-model="selectedYear" placeholder="选择年份" style="width: 120px; margin-right: 10px;">
        <el-option
            v-for="year in availableYears"
            :key="year"
            :label="`${year}年`"
            :value="year"
        />
      </el-select>
      <el-button
          v-if="currentMapLevel !== 'china'"
          @click="backToChinaMap"
          icon="Back"
          size="small"
      >
        返回全国
      </el-button>
    </div>

    <!-- ECharts地图容器 -->
    <div ref="mapChart" class="map-container"></div>

    <!-- 图例说明 -->
    <div class="map-legend">
      <span>分数越低</span>
      <div class="legend-gradient"></div>
      <span>分数越高</span>
    </div>
  </el-card>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick ,computed} from 'vue'
import * as echarts from 'echarts'
import chinaJson from '@/assets/中华人民共和国省.json' // 需要先准备中国地图GeoJSON数据

const props = defineProps({
  admissionScores: {
    type: Array,
    default: () => []
  }
})

const selectedYear = ref(new Date().getFullYear() - 1)
const mapChart = ref(null)
const chartInstance = ref(null)
const currentMapLevel = ref('china')
const currentProvince = ref('')

// 可用年份计算
const availableYears = computed(() => {
  const years = new Set()
  props.admissionScores.forEach(score => {
    years.add(score.year)
  })
  return Array.from(years).sort((a, b) => b - a)
})

// 筛选后的分数数据
const filteredScores = computed(() => {
  return props.admissionScores.filter(score =>
      score.year === selectedYear.value
  )
})

// 获取省份名称
const getProvinceName = (provinceId) => {
  const proMap = {
    '11': '北京', '12': '天津', '13': '河北', '14': '山西', '15': '内蒙古',
    '21': '辽宁', '22': '吉林', '23': '黑龙江', '31': '上海', '32': '江苏',
    '33': '浙江', '34': '安徽', '35': '福建', '36': '江西', '37': '山东',
    '41': '河南', '42': '湖北', '43': '湖南', '44': '广东', '45': '广西',
    '46': '海南', '50': '重庆', '51': '四川', '52': '贵州', '53': '云南',
    '54': '西藏', '61': '陕西', '62': '甘肃', '63': '青海', '64': '宁夏',
    '65': '新疆', '71': '台湾', '81': '香港', '82': '澳门'
  }
  return proMap[provinceId] || `省份${provinceId}`
}

// 获取分数类型名称
const getScoreTypeName = (scoreType) => {
  const typeMap = {
    '1': '旧高考理科', '2': '旧高考文科', '3': '综合', '4': '艺术类', '5': '体育类',
    '23': '体育文', '24': '体育理', '25': '艺术文', '26': '艺术理',
    '31': '蒙授体育', '32': '蒙授文科', '33': '蒙授理科',
    '2073': '新高考物理类', '2074': '新高考历史类', '2292': '艺术类(历史)',
    '2293': '艺术类(物理)', '2294': '体育类(历史)', '2295': '体育类(物理)',
  }
  return typeMap[scoreType] || `类型${scoreType}`
}

// 初始化图表
const initChart = () => {
  if (!mapChart.value) return

  // 先检查是否已有实例，避免重复创建
  if (chartInstance.value) echarts.dispose(chartInstance.value)

  chartInstance.value = echarts.init(mapChart.value)
  echarts.registerMap('china', chinaJson)

  const option = getMapOption()
  chartInstance.value.setOption(option)
  chartInstance.value.on('click', handleMapClick)

  // 确保只添加一次resize监听
  window.removeEventListener('resize', handleResize)
  window.addEventListener('resize', handleResize)
}

// 获取地图配置
const getMapOption = () => {
// 按省份分组数据
  const provinceData = {}
  filteredScores.value.forEach(score => {
    const provinceName = getProvinceName(score.province_id)

    if (!provinceData[provinceName]) {
      provinceData[provinceName] = []
    }
    provinceData[provinceName].push({
      min_score: score.min_score,
      score_type: score.score_type,
      year: score.year,
      provinceId: score.province_id
    })
  })

  // 创建地图数据
  const mapData = Object.entries(provinceData).map(([provinceName, scores]) => {
    const minScore = Math.min(...scores.map(s => s.min_score))
    return {
      name: provinceName,
      value: minScore,
      scores: scores,
      provinceId: scores[0].provinceId
    }
  })

  // 找出最小和最大分数用于视觉映射
  const values = mapData.map(item => item.value)
  const minScore = Math.min(...values)
  const maxScore = Math.max(...values)

  return {
    tooltip: {
      trigger: 'item',
      formatter: function(params) {
        const { name, data } = params
        if (!data || !data.scores) return `${name}<br/>暂无数据`

        // 按分数排序
        const sortedScores = data.scores.sort((a, b) => a.min_score - b.min_score)

        let html = `
          <div style="font-weight:bold;margin-bottom:8px;font-size:14px;">${name}</div>
          <div style="margin-bottom:8px;color:#666;">${selectedYear.value}年录取分数线</div>
        `

        sortedScores.forEach(score => {
          html += `
            <div style="display:flex;justify-content:space-between;margin-bottom:4px;">
              <span style="color:#606266;">${getScoreTypeName(score.score_type)}:</span>
              <span style="font-weight:bold;color:#409EFF;">${score.min_score}分</span>
            </div>
          `
        })

        return html
      }
    },
    visualMap: {
      type: 'continuous',
      min: minScore,
      max: maxScore,
      text: ['低分', '高分'],
      realtime: false,
      calculable: true,
      inRange: {
        color: ['#C6EBFF', '#5B9BD5', '#204D74']
      },
      textStyle: {
        color: '#606266'
      }
    },
    series: [{
      name: '录取分数线',     // 系列名称
      type: 'map',          // 地图类型
      zoom: 1.6,            // 初始缩放比例（默认 1），根据需要调整
      center: [105, 35],   // 初始中心经纬度 [lng, lat]，可微调以把关注区域放中间
      map: currentMapLevel.value === 'china' ? 'china' : currentProvince.value, // 切换省份显示
      roam: false,        // 禁止鼠标滚轮缩放地图
      emphasis: {
        label: {
          show: true,
          color: '#ffff1f',
          fontSize: 16
        },
        itemStyle: {
          areaColor: '#409EFF',
          borderWidth: 2,
          borderColor: '#fff'
        }
      },
      data: mapData,
      nameMap: {
        '北京市': '北京', '天津市': '天津', '河北省': '河北', '山西省': '山西', '内蒙古自治区': '内蒙古',
        '辽宁省': '辽宁', '吉林省': '吉林', '黑龙江省': '黑龙江', '上海市': '上海', '江苏省': '江苏',
        '浙江省': '浙江', '安徽省': '安徽', '福建省': '福建', '江西省': '江西', '山东省': '山东', '河南省': '河南',
        '湖北省': '湖北', '湖南省': '湖南', '广东省': '广东', '广西壮族自治区': '广西', '海南省': '海南',
        '重庆市': '重庆', '四川省': '四川', '贵州省': '贵州', '云南省': '云南', '西藏自治区': '西藏',
        '陕西省': '陕西', '甘肃省': '甘肃', '青海省': '青海', '宁夏回族自治区': '宁夏', '新疆维吾尔自治区': '新疆',
        '台湾省': '台湾', '香港特别行政区': '香港', '澳门特别行政区': '澳门'
      }
    }]
  }
}

// 处理地图点击事件
const handleMapClick = async (params) => {
  if (currentMapLevel.value === 'china') {
    // 全国地图点击，下钻到省份
    const provinceName = params.name
    const provinceId = Object.entries(proMap).find(
        ([id, name]) => name === provinceName
    )?.[0]

    if (provinceId) {
      try {
        // 动态加载省份地图数据
        const provinceJson = await import(`@/assets/geo/province/${provinceId}.json`)
        echarts.registerMap(provinceName, provinceJson)

        currentMapLevel.value = 'province'
        currentProvince.value = provinceName

        // 更新图表显示省份地图
        const option = getMapOption()
        chartInstance.value.setOption(option)
      } catch (error) {
        console.error('加载省份地图失败:', error)
        // 可以在这里添加用户提示
      }
    }
  }
}

// 返回全国地图
const backToChinaMap = () => {
  currentMapLevel.value = 'china'
  currentProvince.value = ''

  const option = getMapOption()
  chartInstance.value.setOption(option)
}

// 处理窗口大小变化
const handleResize = () => {
  if (chartInstance.value) {
    chartInstance.value.resize()
  }
}

// 只更新数据而不重新初始化
const updateChartData = () => {
  if (!chartInstance.value) {
    // 如果实例不存在，才初始化
    initChart()
    return
  }

  const option = getMapOption()
  chartInstance.value.setOption(option)
}


// 监听年份变化
watch(selectedYear, () => {
  if (chartInstance.value) {
    const option = getMapOption()
    chartInstance.value.setOption(option)
  }
})

// 监听分数变化
watch(() => props.admissionScores, (newScores) => {
  // 不调用 initChart()，而是更新数据
  if (newScores && newScores.length > 0) nextTick(() => updateChartData())
}, { deep: true })

onMounted(() => {
  // 确保DOM已经完全渲染
  if (props.admissionScores.length > 0)
    nextTick(() => setTimeout(initChart, 100))
})

onUnmounted(() => {
  if (chartInstance.value) echarts.dispose(chartInstance.value)
  window.removeEventListener('resize', handleResize)
})

</script>

<style lang="scss" scoped>
.map-container {
  width: 100%;
  height: 400px;
  margin-top: 16px;
}

.map-legend {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 16px;
  font-size: 12px;
  color: #606266;

  .legend-gradient {
    width: 200px;
    height: 16px;
    margin: 0 10px;
    background: linear-gradient(to right, #C6EBFF, #5B9BD5, #204D74);
    border-radius: 3px;
  }
}

.scores-filter {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

@media (max-width: 768px) {
  .map-container {
    height: 300px;
  }

  .map-legend {
    flex-direction: column;
    gap: 8px;

    .legend-gradient {
      width: 150px;
    }
  }
}

</style>
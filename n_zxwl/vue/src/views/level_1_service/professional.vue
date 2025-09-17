<template>
  <div class="container">
    <!-- 头部导航 -->
    <header class="header">
      <div class="header-content">
        <img src="../../assets/zxwllogo.png" alt="Logo" class="page-logo">

        <!-- 导航分类区域 -->
        <div class="nav-container">
          <div class="nav-scroll-wrapper">
            <nav class="nav">
              <div
                  v-for="category in categories"
                  :key="category"
                  class="nav-item"
                  :class="{ active: currentCategory === category }"
                  @click="currentCategory = category"
              >
                {{ category }}
              </div>
            </nav>
          </div>
        </div>

        <!-- 右上角返回按钮 -->
        <div class="back-button" @click="handleBack">
          <el-icon :size="28" class="back-icon">
            <Back />
          </el-icon>
        </div>
      </div>
    </header>

    <div class="main-content">
      <!-- 侧边栏 -->
      <aside class="sidebar">
        <div
            v-for="sub in subCategories[currentCategory]"
            :key="sub"
            class="sidebar-item"
            :class="{ active: currentSubCategory === sub }"
            @click="selectSubCategory(sub)"
        >
          {{ sub }}
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="content">
        <div class="chart-data-container">
          <div class="pie-chart" ref="pieChart"></div>
          <div class="stats-box">
            <div class="stat-item">
              <label>平均年薪</label>
              <div class="value">{{ formatNumber(currentData[0]?.salaryavg) }}元</div>
            </div>
            <div class="stat-item">
              <label>平均五年工作经验月薪</label>
              <div class="value">{{ formatNumber(currentData[0]?.fivesalaryavg) }}元</div>
            </div>
            <div class="chart-legend" v-if="currentData.length > 0">
              <div class="legend-item">
                <span class="legend-color" style="background-color: #5470c6;"></span>
                <span>男生比例: {{ currentData[0].boy_rate }}%</span>
              </div>
              <div class="legend-item">
                <span class="legend-color" style="background-color: #ee6666;"></span>
                <span>女生比例: {{ currentData[0].girl_rate }}%</span>
              </div>
            </div>
          </div>
        </div>

        <div class="table-container">
          <table>
            <thead>
            <tr>
              <th>专业名称</th>
              <th>学位层次</th>
              <th>学制</th>
              <th>平均年薪</th>
              <th>五年月薪</th>
              <th class="action-column"></th>
            </tr>
            </thead>
            <tbody>
            <tr
                v-for="item in currentData"
                :key="item.name"
                @click="selectItem(item.id)"
                @mouseenter="hoverItemId = item.id"
                @mouseleave="hoverItemId = null"
                :class="{ 'hover-row': hoverItemId === item.id }"
            >
              <td>{{ item.name }}</td>
              <td>{{ item.level1_name }}</td>
              <td>{{ item.limit_year }}</td>
              <td>{{ formatNumber(item.salaryavg) }}元</td>
              <td>{{ formatNumber(item.fivesalaryavg) }}元</td>
              <td class="action-cell">
                <span class="arrow-icon" v-if="hoverItemId === item.id">↗</span>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </main>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted, watch } from 'vue';
import * as echarts from 'echarts';
import axios from 'axios';
import { Back } from '@element-plus/icons-vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// TODO subCategories的获取方法
/*
SELECT
    CONCAT('{',
           GROUP_CONCAT(
               CONCAT('''', level2_name, ''': ',
                      '[', GROUP_CONCAT(DISTINCT CONCAT('''', level3_name, '''') ORDER BY level3_name SEPARATOR ', '), ']'
               )
               ORDER BY level2_name SEPARATOR ', '
           ),
           '}'
    ) AS js_object
FROM special_detail
WHERE level2_name IN ('哲学', '经济学', '法学', '教育学', '文学', '历史学', '理学', '工学', '农学', '医学', '管理学', '艺术学', '农林牧渔大类', '资源环境与安全大类', '能源动力与材料大类', '土木建筑大类', '水利大类', '装备制造大类', '生物与化工大类', '轻工纺织大类', '文化艺术大类', '食品药品与粮食大类', '交通运输大类', '财经商贸大类', '电子与信息大类', '医药卫生大类', '旅游大类', '新闻传播大类', '教育与体育大类', '公安与司法大类', '公共管理与服务大类')
GROUP BY level2_name;
*/
const categories = [
  '医学', '工学', '经济学', '电子与信息大类', '法学', '教育学', '文学', '历史学',
  '理学',  '农学',  '管理学', '艺术学', '农林牧渔大类', '资源环境与安全大类',
  '能源动力与材料大类', '土木建筑大类', '水利大类', '哲学', '装备制造大类',
  '生物与化工大类', '轻工纺织大类', '文化艺术大类', '食品药品与粮食大类',
  '交通运输大类', '财经商贸大类',  '医药卫生大类',
  '旅游大类', '新闻传播大类', '教育与体育大类', '公安与司法大类', '公共管理与服务大类',];
const subCategories = {
  医学:['中医学类', '中药学类', '中西医结合类', '临床医学类', '公共卫生与预防医学类', '医学技术类', '口腔医学类', '基础医学类', '护理学类', '法医学类', '药学类'],
  交通运输大类:['城市轨道交通类', '水上运输类', '管道运输类', '航空运输类', '道路运输类', '邮政类', '铁道运输类'],
  公共管理与服务大类:['公共事业类', '公共服务类', '公共管理类', '文秘类'],
  公安与司法大类:['侦查类', '公安技术类', '公安管理类', '司法技术类', '安全防范类', '法律实务类', '法律执行类'],
  农学:['动物医学类', '动物生产类', '林学类', '植物生产类', '水产类', '自然保护与环境生态类', '草学类'],
  农林牧渔大类:['农业类', '林业类', '渔业类', '畜牧业类'],
  医药卫生大类:['中医药类', '临床医学类', '健康管理与促进类', '公共卫生与卫生管理类', '医学技术类', '康复治疗类', '护理类', '眼视光类', '药学类'],
  历史学:['历史学类'],
  哲学:['哲学类'],
  土木建筑大类:['土建施工类', '城乡规划与管理类', '市政工程类', '建筑设备类', '建筑设计类', '建设工程管理类', '房地产类'],
  工学:['计算机类','交叉工程类', '交通运输类', '仪器类', '公安技术类', '兵器类', '农业工程类', '力学类', '化工与制药类', '土木类', '地质类', '安全科学与工程类', '建筑类', '机械类', '材料类', '林业工程类', '核工程类', '水利类', '测绘类', '海洋工程类', '环境科学与工程类', '生物医学工程类', '生物工程类', '电子信息类', '电气类', '矿业类', '纺织类', '能源动力类', '自动化类', '航空航天类',  '轻工类', '食品科学与工程类'],
  教育与体育大类:['体育类', '教育类', '语言类'],
  教育学:['体育学类', '教育学类'],
  文化艺术大类:['文化服务类', '民族文化艺术类', '艺术设计类', '表演艺术类'],
  文学:['中国语言文学类', '外国语言文学类', '新闻传播学类'],
  新闻传播大类:['广播影视类', '新闻出版类'],
  资源环境与安全大类:['地质类', '安全类', '气象类', '测绘地理信息类', '煤炭类', '环境保护类', '石油与天然气类', '资源勘查类', '金属与非金属矿类'],
  旅游大类:['旅游类', '餐饮类'],
  水利大类:['水利工程与管理类', '水利水电设备类', '水土保持与水环境类', '水文水资源类'],
  法学:['公安学类', '政治学类', '民族学类', '法学类', '社会学类', '马克思主义理论类'],
  理学:['化学类', '地球物理学类', '地理科学类', '地质学类', '大气科学类', '天文学类', '心理学类', '数学类', '海洋科学类', '物理学类', '生物科学类', '统计学类'],
  生物与化工大类:['化工技术类', '生物技术类'],
  电子与信息大类:['电子信息类', '计算机类', '通信类', '集成电路类'],
  管理学:['公共管理类', '农业经济管理类', '图书情报与档案管理类', '工业工程类', '工商管理类', '旅游管理类', '物流管理与工程类', '电子商务类', '管理科学与工程类'],
  经济学:['经济与贸易类', '经济学类', '财政学类', '金融学类'],
  能源动力与材料大类:['建筑材料类', '新能源发电工程类', '有色金属材料类', '热能与发电工程类', '电力技术类', '非金属材料类', '黑色金属材料类'],
  艺术学:['戏剧与影视学类', '美术学类', '艺术学理论类', '设计学类', '音乐与舞蹈学类'],
  装备制造大类:['机械设计制造类', '机电设备类', '汽车制造类', '自动化类', '航空装备类', '船舶与海洋工程装备类', '轨道装备类'],
  财经商贸大类:['工商管理类', '物流类', '电子商务类', '经济贸易类', '统计类', '财务会计类', '财政税务类', '金融类'],
  轻工纺织大类:['包装类', '印刷类', '纺织服装类', '轻化工类'],
  食品药品与粮食大类:['粮食类', '药品与医疗器械类', '食品类'],
};

const currentCategory = ref('工学');
const currentSubCategory = ref('计算机类');
const currentData = ref([]);
const pieChart = ref(null);

const hoverItemId = ref(null); // 添加hoverItemId来跟踪当前悬停的行
let chartInstance = null;

// 格式化数字显示
const formatNumber = (num) => Number(num).toLocaleString()

// 获取数据
const fetchData = async (category) => {
  try {
    const response = await axios.post('gapi/professional/querys', {
      level3_name: category,
    });
    currentData.value = response.data.data;
    updateChart();
  } catch (error) {
    console.error('数据获取失败:', error);
  }
};

// 初始化图表
const initChart = () => {
  chartInstance = echarts.init(pieChart.value);
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c}%',
    },
    series: [
      {
        name: '性别比例',
        type: 'pie',
        radius: '70%',
        color: ['#5470c6', '#ee6666'],
        label: {
          formatter: '{b}: {d}%',
        },
        data: [],
      },
    ],
  };
  chartInstance.setOption(option);
};

// 更新图表数据
const updateChart = () => {
  if (currentData.value.length > 0) {
    const item = currentData.value[0];
    const option = {
      series: [
        {
          data: [
            { value: item.boy_rate, name: '男生比例' },
            { value: item.girl_rate, name: '女生比例' },
          ],
        },
      ],
    };
    chartInstance.setOption(option);
  }
};

const selectSubCategory = (sub) => {
  currentSubCategory.value = sub;
  fetchData(sub);
};

const selectItem = (id) => {
  router.push({
    path: '/specialDetail',
    query: { id: id }
  });
};

// 处理返回按钮点击事件
const handleBack = () => router.back()

onMounted(async () => {
  const logData = {
    "email": localStorage.getItem('userEmail'),
    "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
    "operation": "用户浏览专业"
  };
  const logResponse = await axios.post("gapi/log", logData, {
    headers: {
      "Content-Type": "application/json"
    }
  });
})
onMounted(() => {
  initChart();
  fetchData('计算机类');
})
</script>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f0f2f5;
}

/* 头部样式优化 */
.header {
  background: linear-gradient(135deg, #2c3e50, #3498db);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  position: relative;
}

.page-logo {
  width: 100px;
  height: auto;
  flex-shrink: 0;
}

.nav-container {
  flex: 1;
  margin: 0 2rem;
  overflow: hidden;
}

.nav-scroll-wrapper {
  overflow-x: auto;
  padding-bottom: 8px; /* 为滚动条留出空间 */
}

.nav-scroll-wrapper::-webkit-scrollbar {
  height: 7px;
}

.nav-scroll-wrapper::-webkit-scrollbar-thumb {
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
}

.nav-scroll-wrapper::-webkit-scrollbar-track {
  background-color: rgba(255, 255, 255, 0.1);
}

.nav {
  display: flex;
  gap: 1rem;
  padding: 0 0.5rem;
  min-width: max-content;
}

.nav-item {
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
  flex-shrink: 0;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.1);
}

.nav-item.active {
  background: rgba(255, 255, 255, 0.2);
}

.back-button {
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  padding: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.back-icon {
  color: #606266;
  transition: color 0.3s ease;
}

.back-button:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.back-button:hover .back-icon {
  color: #6a11cb;
}

/* 主内容区域 */
.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 侧边栏宽度调整 */
.sidebar {
  width: 200px; /* 缩小侧边栏宽度 */
  background: white;
  padding: 1rem;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  overflow-y: auto;
  overflow-x: hidden;
  flex-shrink: 0;
}

.sidebar-item {
  padding: 0.8rem 1rem;
  margin: 0.5rem 0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar-item:hover {
  background: #f5f7fa;
}

.sidebar-item.active {
  background: #3498db;
  color: white;
}

.content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

/* 图表和数据容器 */
.chart-data-container {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 2rem;
  margin-bottom: 2rem;
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.pie-chart {
  height: 300px;
}

.stats-box {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-item {
  margin: 1rem 0;
}

.stat-item label {
  display: block;
  color: #666;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

.stat-item .value {
  font-size: 1.5rem;
  color: #2c3e50;
  font-weight: bold;
}

.chart-legend {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #eee;
}

.legend-item {
  display: flex;
  align-items: center;
  margin: 0.5rem 0;
}

.legend-color {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  margin-right: 8px;
  display: inline-block;
}

/* 表格样式优化 */
.table-container {
  background: white;
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

th {
  background: #f8f9fa;
  color: #666;
  font-weight: 600;
}

.action-column {
  width: 50px;
}

.action-cell {
  text-align: center;
  padding: 1rem 0.5rem;
}

.arrow-icon {
  color: #2196f3;
  font-weight: bold;
  font-size: 1.2rem;
}

tr {
  transition: all 0.2s ease;
}

tr:hover {
  background: #f8f9fa;
  cursor: pointer;
}

/* 高亮行样式 */
.hover-row {
  background: linear-gradient(135deg, #e3f2fd, #bbdefb) !important;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .chart-data-container {
    grid-template-columns: 1fr;
  }

  .pie-chart {
    height: 250px;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
  }

  .nav-container {
    margin: 0;
    width: 100%;
  }

  .main-content {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
    max-height: 200px;
    overflow-x: auto;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    padding: 1rem;
  }

  .sidebar-item {
    margin: 0;
    flex: 1;
    min-width: 120px;
    text-align: center;
  }

  .content {
    padding: 1rem;
  }

  .chart-data-container {
    padding: 1rem;
  }
}
</style>
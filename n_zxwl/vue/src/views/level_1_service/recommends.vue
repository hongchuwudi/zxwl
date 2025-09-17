<template>
  <div :style="containerStyle" class="container">
    <img src="../../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <div class="back-button" @click="handleBack">
      <el-icon :size="28" class="back-icon">
        <Back/>
      </el-icon>
    </div>

    <!-- 输入表单 -->
    <transition name="form-fade">
      <div class="form-card" v-if="!isLoading && !resultData && !apiError">
        <h2 class="form-title">🎓 智选未来高校推荐服务</h2>

        <!-- 紧凑的表单行 -->
        <div class="form-row">
          <div class="input-group">
            <label>📊 高考分数</label>
            <input type="number" v-model="form.score" placeholder="输入分数" class="compact-input">
          </div>

          <div class="input-group">
            <label>📍 生源地</label>
            <el-select v-model="form.province_name" placeholder="选择省份" class="full-width-input">
              <el-option
                  v-for="province in provinceOptions"
                  :key="province.value"
                  :label="province.label"
                  :value="province.value"
              />
            </el-select>
          </div>

          <div class="input-group">
            <label>📝 高考类型</label>
            <el-select v-model="form.type_name" placeholder="选择类型" class="full-width-input">
              <el-option label="物理类" value="物理类"/>
              <el-option label="历史类" value="历史类"/>
              <el-option label="综合类" value="综合"/>
            </el-select>
          </div>

        </div>
        <div class="form-row">

          <div class="input-group" v-if="isSpecialProvince">
            <label>🎯 批次类型</label>
            <el-select v-model="form.batch_name" placeholder="选择批次" class="full-width-input">
              <el-option label="本科" value="本科"/>
              <el-option label="专科" value="专科"/>
            </el-select>
          </div>

          <div class="input-group">
            <label>💰 毕业最低年薪(选填)</label>
            <input type="text" v-model="form.salary" placeholder="期望年薪" class="compact-input">
          </div>

          <div class="input-group">
            <label>🏙️ 目标地区(选填)</label>
            <el-select v-model="form.goal_province_name" placeholder="选择目标省份" class="full-width-input">
              <el-option
                  v-for="province in provinceOptions"
                  :key="province.value"
                  :label="province.label"
                  :value="province.value"
              />
            </el-select>
          </div>
        </div>

        <div class="input-group">
          <label>❤️ 兴趣专业或个人特点(选填)</label>
          <textarea v-model="form.interest" placeholder="请输入你感兴趣的专业或个人特点" rows="2"
                    class="form-textarea"></textarea>
        </div>

        <div class="input-group">
          <label>❤️ 家庭意见(选填)</label>
          <textarea v-model="form.family_pref" placeholder="请输入家庭意见" rows="2" class="form-textarea"></textarea>
        </div>

        <button
            class="generate-btn"
            @click="submitForm"
            :disabled="!isFormComplete"
            :class="{ 'btn-active': isFormComplete }"
        >
          ✨ 生成智能推荐方案
        </button>
      </div>
    </transition>

    <!-- 加载动画 - Ant Design Vue组件 -->
    <transition name="fade">
      <div v-if="isLoading" class="loading-container">
        <div class="antd-loading-wrapper">
          <a-spin size="large" />
          <h3 class="loading-text">正在扫描全国高校数据库，可能需要10s+，请稍候...</h3>
          <p class="loading-subtext">AI正在分析数据并生成个性化推荐</p>
        </div>
      </div>
    </transition>
    <!-- 错误提示 -->
    <transition name="error-slide">
      <div v-if="apiError" class="error-alert">
        ⚠️ {{ apiError }}
        <button @click="apiError = null" class="close-btn">×</button>
      </div>
    </transition>

    <!-- 结果展示 -->
    <transition name="result-scale">
      <div v-if="isNullResult === '非空'" class="result-container">
        <!-- 录取分析 -->
        <div class="analysis-section" v-if="resultData.admission_analysis !== '用户没有家庭偏好,暂时不使用AI分析'">
          <h3 class="section-title">📊 录取分析</h3>
          <div class="analysis-content" v-html="renderMarkdown(resultData.admission_analysis)"></div>
        </div>

        <!-- 兴趣分析 -->
        <div class="analysis-section" v-if="resultData.interest_analysis !== '用户未提供兴趣爱好信息'">
          <h3 class="section-title">🎯 兴趣分析</h3>
          <p class="analysis-content">{{ resultData.interest_analysis }}</p>
        </div>

        <!-- 学校推荐 - 冲稳保展示 -->
        <div class="recommendation-section">
          <h3 class="section-title">🏫 院校推荐 - 冲稳保分析</h3>

          <!-- 冲刺院校  -->
          <div class="category-section chong">
            <div class="category-header">
              <span class="category-icon">🚀</span>
              <h4 class="category-title">冲刺院校</h4>
              <span class="category-desc">录取概率40%以下，有挑战性的理想院校</span>
            </div>
            <div class="university-grid">
              <div
                  v-for="(uni, index) in resultData.university_recommendations.filter(item => item.admission_probability< 0.4)"
                  :key="uni.school_name + index"
                  class="uni-card rush"
              >
                <div class="card-header">
                  <h4 class="uni-name">{{ uni.school_name }}</h4>
                  <div class="probability-badge">
                    {{ (uni.admission_probability * 100).toFixed(0) }}%
                  </div>
                </div>
                <div class="card-body">
                  <div class="uni-info">
                    <span class="info-item">📍 {{ uni.school_address }}</span>
                    <span class="info-item">📅{{ form.goal_year }}年最低分: {{ uni.min_score }}</span>
                    <span v-if="uni.has_rk_rank === 1" class="info-item">📈 软科排名: {{ uni.ruanke_rank }}</span>
                    <span v-if="uni.has_xyh_rank === 1" class="info-item">🎖️ 校友会排名: {{ uni.xyh_rank }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 稳妥院校  -->
          <div class="category-section wen">
            <div class="category-header">
              <span class="category-icon">✅</span>
              <h4 class="category-title">稳妥院校</h4>
              <span class="category-desc">录取概率50%-90%，与成绩匹配的合适院校</span>
            </div>
            <div class="university-grid">
              <div
                  v-for="(uni, index) in resultData.university_recommendations.filter(item => item.admission_probability< 0.9 && item.admission_probability>= 0.5)"
                  :key="uni.school_name + index"
                  class="uni-card stable"
              >
                <div class="card-header">
                  <h4 class="uni-name">{{ uni.school_name }}</h4>
                  <div class="probability-badge">
                    {{ (uni.admission_probability * 100).toFixed(0) }}%
                  </div>
                </div>
                <div class="card-body">
                  <div class="uni-info">
                    <span class="info-item">📍 {{ uni.school_address }}</span>
                    <span class="info-item">📅{{ form.goal_year }}年最低分: {{ uni.min_score }}</span>
                    <span v-if="uni.has_rk_rank === 1" class="info-item">📈 软科排名: {{ uni.ruanke_rank }}</span>
                    <span v-if="uni.has_xyh_rank === 1" class="info-item">🎖️ 校友会排名: {{ uni.xyh_rank }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 保底院校  -->
          <div class="category-section bao">
            <div class="category-header">
              <span class="category-icon">🛡️</span>
              <h4 class="category-title">保底院校</h4>
              <span class="category-desc">录取概率90%以上，确保能被录取的院校</span>
            </div>
            <div class="university-grid">
              <div
                  v-for="(uni, index) in resultData.university_recommendations.filter(item => item.admission_probability >= 0.9)"
                  :key="uni.school_name + index"
                  class="uni-card safe"
              >
                <div class="card-header">
                  <h4 class="uni-name">{{ uni.school_name }}</h4>
                  <div class="probability-badge">
                    {{ (uni.admission_probability * 100).toFixed(0) }}%
                  </div>
                </div>
                <div class="card-body">
                  <div class="uni-info">
                    <span class="info-item">📍 {{ uni.school_address }}</span>
                    <span class="info-item">📅{{ form.goal_year }}年最低分: {{ uni.min_score }}</span>
                    <span v-if="uni.has_rk_rank === 1" class="info-item">📈 软科排名: {{ uni.ruanke_rank }}</span>
                    <span v-if="uni.has_xyh_rank === 1" class="info-item">🎖️ 校友会排名: {{ uni.xyh_rank }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 专业推荐 -->
        <div class="recommendation-section"
             v-if="resultData.special_recommendations && resultData.special_recommendations.length">
          <h3 class="section-title">📚 专业推荐</h3>
          <div class="special-grid">
            <div
                v-for="(special, index) in resultData.special_recommendations"
                :key="special.special_name + index"
                class="special-card"
            >
              <div class="card-header">
                <h4 class="special-name">{{ special.special_name }}</h4>
                <div class="probability-badge">
                  {{ (special.admission_probability * 100).toFixed(0) }}%
                </div>
              </div>
              <div class="card-body">
                <div class="special-info">
                  <span class="info-item">🏫 {{ special.school_name }}</span>
                  <span class="info-item">📍 {{ special.school_address }}</span>
                  <span class="info-item">📖 {{ special.special_level1_name }} - {{
                      special.special_level2_name
                    }} -{{ special.special_level3_name }}</span>
                  <span v-if="special.special_avg_salary"
                        class="info-item">💰 平均毕业年资: ¥{{ special.special_avg_salary }}</span>
                  <span v-if="special.special_subject_requirements"
                        class="info-item">🎯 选科要求: {{ special.special_subject_requirements }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else-if="isNullResult === '空'" class="no-result">
        <div class="no-result-content">
          <div class="no-result-icon">🎓</div>
          <h3 class="no-result-title">暂无推荐结果</h3>
          <p class="no-result-desc">
            根据您提供的信息，暂时没有找到合适的院校和专业推荐。<br>
            建议您调整查询条件或尝试不同的分数和地区组合。
          </p>
          <div class="no-result-actions">
            <button class="retry-btn" @click="resetForm">
              🔄 重新填写
            </button>
            <button class="back-home-btn" @click="handleBack">
              🏠 返回首页
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import {computed, onMounted, reactive, ref} from 'vue'
import axios from 'axios'
import {gsap} from 'gsap'
import {useRouter} from 'vue-router'
import {console} from "vuedraggable/src/util/console.js";
import {ElMessage} from "element-plus";
import {useUserStore} from "@/utils/auth.js";
const { userName, userEmail, getUser, checkLoginStatus } = useUserStore();
import { LoadingOutlined } from '@ant-design/icons-vue';
import { Spin } from 'ant-design-vue';
const router = useRouter()

const form = reactive({
  province_name: '',
  goal_province_name: '',         // 默认为空
  score: null,
  type_name: '',
  batch_name: '',
  interest: '',
  family_pref: '',                   // 用户输入的家庭偏好
  salary: 0,                       // 用户输入的希望工资
  year: new Date().getFullYear(),  // 默认当前年份
  goal_year: 2024                  // 默认2024

})

const isLoading = ref(false)
const resultData = ref(null)
const apiError = ref(null)
const isNullResult = ref("默认")
const containerStyle = computed(() => {
  if (isLoading.value || !resultData.value) {
    return {minHeight: '100vh'}
  } else {
    return {minHeight: 'auto'}
  }
})
const isFormComplete = computed(() => form.province_name && form.score && form.type_name)
// const containerStyle = computed(() => (isLoading.value || !resultData.value)  ? { minHeight: '100vh' } : { minHeight: 'auto' })
// const isFormComplete = computed(() => form.province_name && form.score && form.type_name)
// 简写

// 判断是否为“北上津”（名称或代码）
const specialProvinces = ['北京', '上海', '天津', '11', '31', '12']
const isSpecialProvince = computed(() => specialProvinces.includes(form.province_name))

// 年份选项（近5年）
const yearOptions = ref(Array.from({length: 5}, (_, i) => new Date().getFullYear() - i))

// 简单Markdown转换
const renderMarkdown = (text) => text
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>') // 加粗
    .replace(/\*(.*?)\*/g, '<em>$1</em>') // 斜体
    .replace(/^- (.*?)(?=\n|$)/gm, '<li>$1</li>') // 无序列表
    .replace(/^# (.*?)(?=\n|$)/gm, '<h2>$1</h2>') // 一级标题
    .replace(/^## (.*?)(?=\n|$)/gm, '<h3>$1</h3>') // 二级标题
    .replace(/^### (.*?)(?=\n|$)/gm, '<h4>$1</h4>') // 3级标题
    .replace(/^#### (.*?)(?=\n|$)/gm, '<h5>$1</h5>') // 4级标题
    .replace(/\n/g, '<br>') // 换行

// 重置表单
const resetForm = () => {
  form.province_name = ''
  form.goal_province_name = ''
  form.score = null
  form.type_name = ''
  form.batch_name = ''
  form.interest = ''
  form.family_pref = ''
  form.salary = 0
  resultData.value = null
  apiError.value = null

  // 添加一个简单的动画效果
  gsap.from('.form-card', {
    duration: 0.8,
    opacity: 0,
    y: 50,
    ease: 'back.out(1.7)'
  })

  isNullResult.value = "默认"
}

// 省份选项
const provinceOptions = ref([
  {label: '北京市', value: '北京'},
  {label: '天津市', value: '天津'},
  {label: '河北省', value: '河北省'},
  {label: '山西省', value: '山西省'},
  {label: '内蒙古自治区', value: '内蒙古自治区'},
  {label: '辽宁省', value: '辽宁省'},
  {label: '吉林省', value: '吉林省'},
  {label: '黑龙江省', value: '黑龙江省'},
  {label: '上海市', value: '上海'},
  {label: '江苏省', value: '江苏省'},
  {label: '浙江省', value: '浙江省'},
  {label: '安徽省', value: '安徽省'},
  {label: '福建省', value: '福建省'},
  {label: '江西省', value: '江西省'},
  {label: '山东省', value: '山东省'},
  {label: '河南省', value: '河南省'},
  {label: '湖北省', value: '湖北省'},
  {label: '湖南省', value: '湖南省'},
  {label: '广东省', value: '广东省'},
  {label: '广西壮族自治区', value: '广西壮族自治区'},
  {label: '海南省', value: '海南省'},
  {label: '重庆市', value: '重庆市'},
  {label: '四川省', value: '四川省'},
  {label: '贵州省', value: '贵州省'},
  {label: '云南省', value: '云南省'},
  {label: '陕西省', value: '陕西省'},
  {label: '甘肃省', value: '甘肃省'},
  {label: '青海省', value: '青海省'},
  {label: '宁夏回族自治区', value: '宁夏回族自治区'},
  {label: '新疆维吾尔自治区', value: '新疆维吾尔自治区'}
])

const submitForm = async () => {
  try {
    isLoading.value = true
    apiError.value = null
    resultData.value = null

    gsap.to('.form-card', {
      duration: 0.8,
      opacity: 0,
      y: -50,
      ease: 'power3.inOut'
    })

    const response = await axios.post('gapi/recommends', form)

    if (response.data.code !== 0) {
      throw new Error(response.data.message || '服务器返回未知错误')
    }

    resultData.value = response.data.data

    if((resultData.value.special_recommendations ===  null || resultData.value.special_recommendations.length === 0)
        && (resultData.value.university_recommendations === null || resultData.value.university_recommendations.length === 0)
    ) {
      console.log("空",resultData.value)
      isNullResult.value = "空"
    } else {
      console.log("f空",resultData.value)
      isNullResult.value = "非空"
    }
    playResultAnimation()

  } catch (error) {
    handleError(error)
    resetFormAnimation()
  } finally {
    isLoading.value = false
  }
}

const playResultAnimation = () => {
  gsap.from('.uni-card, .special-card, .analysis-section', {
    duration: 0.8,
    opacity: 0,
    y: 50,
    stagger: 0.1,
    ease: 'back.out(1.7)'
  })
}

const handleError = (error) => {
  if (error.response) {
    apiError.value = error.response.data.message || `请求错误：${error.response.status}`
  } else if (error.request) {
    apiError.value = '网络连接异常，请检查网络后重试'
  } else {
    apiError.value = error.message || '系统异常，请稍后重试'
  }
}

const resetFormAnimation = () => {
  gsap.to('.form-card', {
    duration: 0.8,
    opacity: 1,
    y: 0,
    ease: 'power3.out'
  })
}

const handleBack = () => {
  router.back()
}
onMounted(async () => {
  getUser()
  if (!checkLoginStatus()) {
    ElMessage.error('请先登录！')
    router.push('/login')
  }

  const logData = {
    "email": localStorage.getItem('userEmail'),
    "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
    "operation": "用户执行智能推荐"
  };
  const logResponse = await axios.post("gapi/log", logData, {
    headers: {
      "Content-Type": "application/json"
    }
  });
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  height: 100%;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  font-family: 'Inter', 'Microsoft YaHei', sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-sizing: border-box; /* 确保内边距不影响总高度 */
}

.page-logo {
  position: absolute;
  top: 1rem;
  left: 1rem;
  width: 100px;
  height: auto;
  z-index: 3;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.back-button {
  position: absolute;
  top: 2.5rem;
  right: 4rem;
  cursor: pointer;
  z-index: 1000;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 50%;
  padding: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
  background: rgba(52, 152, 219, 0.1);
}

.back-icon {
  color: #606266;
  transition: color 0.3s ease;
}

.back-button:hover .back-icon {
  color: #3498db;
}

.form-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 2.5rem;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  width: 100%;
  max-width: 800px;
  margin: 2rem 0;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.form-title {
  text-align: center;
  margin-bottom: 2rem;
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.8rem;
  background: linear-gradient(135deg, #3498db 0%, #2c3e50 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 10px;
  align-items: start;
}

.input-group {
  margin-bottom: 1.5rem;
  display: flex;
  flex-direction: column;
}

label {
  display: block;
  margin-bottom: 0.8rem;
  color: #34495e;
  font-weight: 500;
  font-size: 0.95rem;
}

/* 统一所有输入框的样式 */
.compact-input {
  width: 100%;
  padding: 0.8rem 1rem;
  border: 2px solid #e1e8ed;
  border-radius: 12px;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  height: 44px;
  min-height: 44px;
  box-sizing: border-box;
  font-family: inherit;
}

.compact-input:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.15);
  outline: none;
  transform: translateY(-1px);
}

.full-width-input {
  width: 100%;
  border: 2px solid #e1e8ed;
  border-radius: 12px;
  font-size: 0.95rem;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  height: 100%;
  min-height: 44px;
  box-sizing: border-box;
  font-family: inherit;

}

/* 针对 Vue 3的样式穿透 */
.full-width-input ::v-deep(.el-select__wrapper) {
  box-shadow: none !important;
  box-sizing: content-box !important;
  border-radius: 12px !important; /* 圆角 */
}

/* textarea 样式 */
.form-textarea {
  height: auto;
  min-height: 100px;
  width: 100%;
  padding: 0.9rem 1.2rem;
  border: 2px solid #e1e8ed;
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.3s ease;
  resize: vertical;
  background: rgba(255, 255, 255, 0.9);
  box-sizing: border-box;
}

.form-textarea:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.15);
  outline: none;
  transform: translateY(-1px);
}

.generate-btn {
  width: 100%;
  padding: 1.1rem;
  margin-top: 2rem;
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1.05rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  opacity: 0.7;
  letter-spacing: 0.5px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.generate-btn.btn-active {
  opacity: 1;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(52, 152, 219, 0.35);
}

.generate-btn:disabled {
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.generate-btn:not(:disabled):hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 30px rgba(52, 152, 219, 0.4);
}

/* 更新加载样式 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  width: 100%;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.antd-loading-wrapper {
  text-align: center;
  padding: 2rem;
}

/* 自定义旋转动画 */
::v-deep(.ant-spin-dot-item) {
  background-color: #1890ff; /* Ant Design 主蓝色 */
}

::v-deep(.ant-spin-lg .ant-spin-dot) {
  font-size: 40px;
}

@keyframes quantum-pulse {
  0%, 100% {
    transform: scale(0.8);
    opacity: 0.7;
    box-shadow: 0 2px 6px rgba(52, 152, 219, 0.2);
  }
  50% {
    transform: scale(1.3);
    opacity: 1;
    box-shadow: 0 6px 15px rgba(52, 152, 219, 0.4);
  }
}

.loading-text {
  color: #2c3e50;
  margin: 20px 0 10px;
  font-size: 1.2rem;
  font-weight: 600;
}

.loading-subtext {
  color: #7f8c8d;
  font-size: 0.95rem;
  max-width: 300px;
  margin: 0 auto;
  line-height: 1.5;
}

/* 错误提示 */
.error-alert {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: linear-gradient(135deg, #e74c3c, #c0392b);
  color: white;
  padding: 1rem 2rem;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 6px 20px rgba(231, 76, 60, 0.3);
  z-index: 1000;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0 0.5rem;
  transition: transform 0.2s ease;
}

.close-btn:hover {
  transform: scale(1.2);
}

/* 结果容器 */
.result-container {
  width: 95%;
  max-width: 1200px;
  margin: 2rem auto;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.section-title {
  color: #2c3e50;
  margin-bottom: 1.5rem;
  padding-bottom: 0.8rem;
  border-bottom: 2px solid #ecf0f1;
  font-weight: 600;
  font-size: 1.4rem;
  background: linear-gradient(135deg, #2c3e50, #3498db);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.analysis-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  padding: 1.8rem;
  margin-bottom: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(52, 152, 219, 0.1);
}

.analysis-content {
  color: #34495e;
  line-height: 1.7;
  margin: 0;
  font-size: 1rem;
}

.recommendation-section {
  margin-bottom: 3rem;
}

.university-grid, .special-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.8rem;
  margin-top: 1.5rem;
}

.uni-card, .special-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 1.8rem;
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border-left: 5px solid #3498db;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.uni-card.rush {
  border-left-color: #e74c3c;
}

.uni-card.stable {
  border-left-color: #f39c12;
}

.uni-card.safe {
  border-left-color: #27ae60;
}

.uni-card:hover, .special-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.2rem;
  gap: 1rem;
}

.uni-name, .special-name {
  margin: 0;
  color: #2c3e50;
  font-weight: 600;
  flex: 1;
  font-size: 1.1rem;
  line-height: 1.4;
}

.probability-badge {
  background: linear-gradient(135deg, #3498db, #2980b9);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-weight: 600;
  font-size: 0.9rem;
  min-width: 55px;
  text-align: center;
  box-shadow: 0 3px 8px rgba(52, 152, 219, 0.3);
}

.uni-info, .special-info {
  display: flex;
  flex-direction: column;
  gap: 0.7rem;
}

.info-item {
  color: #5d6d7e;
  font-size: 0.92rem;
  line-height: 1.5;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.info-item::before {
  content: "•";
  color: #3498db;
  font-weight: bold;
  font-size: 1.2rem;
}

/* 分类区域 */
.category-section {
  margin-bottom: 3rem;
  padding: 1.5rem;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.category-header {
  display: flex;
  align-items: center;
  margin-bottom: 1.5rem;
  padding: 1rem 1.5rem;
  border-radius: 12px;
  color: white;
  gap: 1rem;
}

.chong .category-header {
  background: linear-gradient(135deg, #ff6b6b, #e74c3c);
}

.wen .category-header {
  background: linear-gradient(135deg, #42b983, #27ae60);
}

.bao .category-header {
  background: linear-gradient(135deg, #3498db, #2980b9);
}

.category-icon {
  font-size: 1.8rem;
}

.category-title {
  font-size: 1.3rem;
  font-weight: 600;
  margin-right: auto;
}

.category-desc {
  font-size: 0.9rem;
  opacity: 0.95;
  text-align: right;
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 968px) {
  .form-row {
    grid-template-columns: repeat(1, 1fr);
  }

  .university-grid, .special-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (max-width: 768px) {
  .container {
    padding: 1rem;
  }

  .page-logo {
    left: 0.5rem;
    width: 80px;
  }

  .back-button {
    top: 1.5rem;
    right: 1rem;
    padding: 8px;
  }

  .form-card {
    padding: 1.8rem;
    margin: 1.5rem 0;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 15px;
  }

  .form-title {
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
  }

  .loading-container {
    height: 100vh;
    padding: 1.5rem;
  }

  .quantum-spinner {
    width: 80px;
    height: 80px;
  }

  .quantum-dot {
    width: 16px;
    height: 16px;
  }

  .loading-text {
    font-size: 1.1rem;
  }

  .result-container {
    padding: 1.5rem;
    margin: 1.5rem auto;
  }

  .university-grid, .special-grid {
    grid-template-columns: 1fr;
    gap: 1.2rem;
  }

  .category-header {
    flex-direction: column;
    text-align: center;
    gap: 0.8rem;
  }

  .category-desc {
    text-align: center;
  }
}

@media (max-width: 480px) {
  .form-card {
    padding: 1.2rem;
  }

  .form-title {
    font-size: 1.3rem;
  }

  .compact-input, .form-textarea {
    padding: 0.7rem 0.9rem;
    font-size: 0.9rem;
  }

  .generate-btn {
    padding: 0.9rem;
    font-size: 1rem;
  }

  .error-alert {
    padding: 0.8rem 1.2rem;
    font-size: 0.9rem;
    max-width: 90%;
  }
}

.no-result {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
  width: 100%;
  padding: 2rem;
}

.no-result-content {
  text-align: center;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 3rem 2rem;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  max-width: 500px;
  width: 100%;
}

.no-result-icon {
  font-size: 4rem;
  margin-bottom: 1.5rem;
  animation: bounce 2s infinite;
}

/* 保持原有的响应式设计 */
@media (max-width: 768px) {
  .loading-container {
    padding: 1.5rem;
  }

  ::v-deep(.ant-spin-lg .ant-spin-dot) {
    font-size: 32px;
  }

  .loading-text {
    font-size: 1.1rem;
  }
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-20px);
  }
  60% {
    transform: translateY(-10px);
  }
}

.no-result-title {
  color: #2c3e50;
  font-size: 1.8rem;
  margin-bottom: 1rem;
  font-weight: 600;
  background: linear-gradient(135deg, #2c3e50, #3498db);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.no-result-desc {
  color: #7f8c8d;
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 2rem;
}

.no-result-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.retry-btn, .back-home-btn {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-btn {
  background: linear-gradient(135deg, #3498db, #2980b9);
  color: white;
}

.retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(52, 152, 219, 0.35);
}

.back-home-btn {
  background: rgba(255, 255, 255, 0.9);
  color: #2c3e50;
  border: 2px solid #e1e8ed;
}

.back-home-btn:hover {
  background: #f8f9fa;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .no-result {
    min-height: 50vh;
    padding: 1rem;
  }

  .no-result-content {
    padding: 2rem 1.5rem;
  }

  .no-result-icon {
    font-size: 3rem;
  }

  .no-result-title {
    font-size: 1.5rem;
  }

  .no-result-desc {
    font-size: 1rem;
  }

  .no-result-actions {
    flex-direction: column;
    gap: 0.8rem;
  }

  .retry-btn, .back-home-btn {
    width: 100%;
  }
}

</style>


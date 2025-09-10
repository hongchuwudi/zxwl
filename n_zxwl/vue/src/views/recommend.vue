<template>
  <div class="container">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <div class="back-button" @click="handleBack">
      <el-icon :size="28" class="back-icon">
        <Back />
      </el-icon>
    </div>

    <!-- 输入表单 -->
    <transition name="form-fade">
      <div class="form-card" v-if="!isLoading &&!resultData &&!apiError">
        <h2 class="form-title">🎓 智选未来高校推荐服务</h2>
        <div class="input-group">
          <label>📍 生源地</label>
          <input type="text" v-model="form.origin" placeholder="请输入你的生源地">
        </div>
        <div class="input-group">
          <label>🏙️ 目标地区</label>
          <input type="text" v-model="form.target_location" placeholder="请输入心仪的地区">
        </div>
        <div class="input-group">
          <label>📊 高考分数</label>
          <input type="number" v-model="form.score" placeholder="请输入你的高考分数">
        </div>
        <div class="input-group">
          <label>🏅 全省排名</label>
          <input type="number" v-model="form.rank" placeholder="请输入你的全省排名">
        </div>
        <div class="input-group">
          <label>❤️ 兴趣专业或个人特点</label>
          <input type="text" v-model="form.preferred_subject" placeholder="请输入你感兴趣的专业或个人特点">
        </div>

        <button
            class="generate-btn"
            @click="submitForm"
            :disabled="!isFormComplete"
            :class="{ 'btn-active': isFormComplete }"
        >
          ✨ 生成推荐方案
        </button>
      </div>
    </transition>

    <!-- 加载动画 -->
    <transition name="fade">
      <div v-if="isLoading" class="loading-container">
        <div class="quantum-spinner">
          <div class="particle"></div>
          <div class="particle"></div>
          <div class="particle"></div>
        </div>
        <h3 class="loading-text">正在扫描全国高校数据库...</h3>
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
      <div v-if="resultData" class="result-container">
        <div
            v-for="(category, type) in categories"
            :key="type"
            class="category-block"
            :class="type"
        >
          <div class="category-header">
            <h3 class="category-title">
              {{ category.label }}
              <span class="tag-bubble">{{ category.tag }}</span>
            </h3>
          </div>
          <div class="university-list">
            <div
                v-for="(uni, index) in resultData[type]"
                :key="uni.name + index"
                class="uni-card"
            >
              <div class="probability-bar" :style="getProbabilityStyle(uni.probability, type)"></div>
              <div class="card-content">
                <h4 class="uni-name">{{ uni.name }}</h4>
                <div class="probability-text">
                  <span class="percent">{{ (uni.probability * 100).toFixed(0) }}%</span>
                  <span class="label">录取概率</span>
                </div>
                <div class="major-container">
                  <div
                      v-for="(major, mIndex) in uni.recommended_majors"
                      :key="mIndex"
                      class="major-bubble"
                  >
                    {{ major }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="result-container">
        <div class="no-result-text">
          暂无数据，请重新填写信息
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import {ref, reactive, computed, onMounted} from 'vue'
import axios from 'axios'
import { gsap } from 'gsap'
import { Back } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  origin: '',
  target_location: '',
  score: null,
  rank: null,
  preferred_subject: ''
})

const isLoading = ref(false)
const resultData = ref(null)
const apiError = ref(null)

const categories = {
  rush: { label: '冲刺院校', tag: '冲', color: '#FF6B6B' },
  stable: { label: '稳妥院校', tag: '稳', color: '#4ECDC4' },
  safe: { label: '保底院校', tag: '保', color: '#45B7D1' }
}

const isFormComplete = computed(() => {
  return Object.values(form).every(v => {
    if (typeof v === 'number') return v !== null && v >= 0
    return v.trim() !== ''
  })
})

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

    const response = await axios.post('gapi/recommend', {
      origin: form.origin,
      target_location: form.target_location,
      score: Number(form.score),
      rank: Number(form.rank),
      preferred_subject: form.preferred_subject
    })

    if (response.data.code !== 0) {
      throw new Error(response.data.message || '服务器返回未知错误')
    }

    resultData.value = response.data.data
    playResultAnimation()

  } catch (error) {
    handleError(error)
    resetFormAnimation()
  } finally {
    isLoading.value = false
  }
}

const playResultAnimation = () => {
  gsap.from('.uni-card', {
    duration: 0.8,
    opacity: 0,
    y: 100,
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

const getProbabilityStyle = (probability, type) => ({
  width: `${probability * 100}%`,
  background: `linear-gradient(90deg, ${categories[type].color} 0%, ${categories[type].color}80 100%)`
})

const handleBack = () => {
  router.push('/zxwl')
}
onMounted(async () => {
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
  padding: 2rem;
  min-height: 100vh;
  background: url(../assets/colleges.jpg) no-repeat center center fixed;
  background-size: cover;
  font-family: 'Microsoft YaHei', sans-serif;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.page-logo {
  position: absolute;
  top: 1rem;
  left: 1rem;
  width: 100px;
  height: auto;
  z-index: 3;
}

.back-button {
  position: absolute;
  top: 2.5rem;
  right: 4rem;
  cursor: pointer;
  z-index: 1000;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  padding: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.form-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 20px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(31, 38, 135, 0.1);
  backdrop-filter: blur(12px);
  width: 100%;
  max-width: 500px;
  margin: 2rem 0;
}

/* 结果展示样式 */
.result-container {
  display: flex;
  gap: 2rem;
  width: 95%;
  max-width: 1200px;
  padding: 2rem 0;
  margin: 0 auto;
  margin-top: 4rem;
}

.category-block {
  flex: 1;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 15px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-width: 320px;
  max-height: 70vh;
}

.category-header {
  padding: 1.5rem;
  background: linear-gradient(135deg, var(--header-color), #ffffff);
  border-bottom: 2px solid rgba(0, 0, 0, 0.05);
}

.rush .category-header { --header-color: #FF6B6B; }
.stable .category-header { --header-color: #4ECDC4; }
.safe .category-header { --header-color: #45B7D1; }

.university-list {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: var(--scroll-color) transparent;
}

.rush .university-list { --scroll-color: #FF6B6B; }
.stable .university-list { --scroll-color: #4ECDC4; }
.safe .university-list { --scroll-color: #45B7D1; }

.university-list::-webkit-scrollbar {
  width: 6px;
}

.university-list::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
}

.university-list::-webkit-scrollbar-thumb {
  background-color: var(--scroll-color);
  border-radius: 4px;
}

.uni-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 16px;
  padding: 1.5rem;
  margin: 1rem 0;
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  width: 100%;
}

@media (max-width: 1024px) {
  .result-container {
    flex-direction: column;
    width: 100%;
  }

  .category-block {
    max-width: 100%;
    max-height: 50vh;
  }
}

.input-group {
  margin: 1.2rem 0;
}

label {
  display: block;
  margin-bottom: 0.6rem;
  color: #4a5568;
  font-size: 0.95rem;
}

input {
  width: 90%;
  padding: 0.8rem 1.2rem;
  border: 2px solid #c3e4ff;
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

input:focus {
  border-color: #4ECDC4;
  box-shadow: 0 0 15px rgba(78, 205, 196, 0.2);
  outline: none;
}

.generate-btn {
  width: 100%;
  padding: 1rem;
  margin-top: 1.5rem;
  background: linear-gradient(135deg, #4ECDC4 0%, #45B7D1 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  opacity: 0.8;
}

.generate-btn.btn-active {
  opacity: 1;
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(78, 205, 196, 0.3);
}

.loading-container {
  position: fixed; /* 改为 fixed 定位 */
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999; /* 确保在最上层 */
  background: rgba(255, 255, 255, 0.9); /* 添加半透明背景 */
  margin-top: 0; /* 移除 margin-top */
}
.quantum-spinner {
  position: relative;
  width: 80px;
  height: 80px;
}

.particle {
  position: absolute;
  width: 12px;
  height: 12px;
  background: #4ECDC4;
  border-radius: 50%;
  animation: quantum-orb 2.4s infinite cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes quantum-orb {
  0% {
    transform: rotate(0deg) translateX(40px) rotate(0deg);
    opacity: 1;
  }
  100% {
    transform: rotate(360deg) translateX(40px) rotate(-360deg);
    opacity: 0.3;
  }
}

.loading-text {
  color: white;
  margin-top: 1rem;
  font-size: 1.2rem;
  text-align: center;
}

.error-alert {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: #ff6b6b;
  color: white;
  padding: 0.8rem 2rem;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 4px 15px rgba(255, 107, 107, 0.3);
  z-index: 1000;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0 0.5rem;
}

.probability-bar {
  position: absolute;
  top: 0;
  left: 0;
  height: 4px;
  transition: width 1s ease;
}

.uni-name {
  color: #2d3748;
  margin-bottom: 0.8rem;
  font-size: 1.1rem;
}

.probability-text {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.percent {
  color: #4ECDC4;
  font-size: 1.4rem;
  font-weight: 700;
}

.label {
  color: #718096;
  font-size: 0.9rem;
}

.major-container {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
}

.major-bubble {
  background: rgba(78, 205, 196, 0.1);
  color: #2a827b;
  padding: 0.4rem 0.8rem;
  border-radius: 8px;
  font-size: 0.9rem;
}

.category-title {
  margin: 0;
  font-size: 1.2rem;
  color: #4a5568;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.tag-bubble {
  padding: 0.3rem 1rem;
  border-radius: 15px;
  font-size: 0.9rem;
  background: var(--tag-color);
  color: white;
}

.form-fade-enter-active,
.form-fade-leave-active {
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}
.form-fade-enter-from,
.form-fade-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

.error-slide-enter-active,
.error-slide-leave-active {
  transition: all 0.4s ease;
}
.error-slide-enter-from,
.error-slide-leave-to {
  opacity: 0;
  transform: translate(-50%, -20px);
}

.result-scale-enter-active {
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}
.result-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
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

.particle:nth-child(1) { animation-delay: 0s; }
.particle:nth-child(2) { animation-delay: 0.4s; }
.particle:nth-child(3) { animation-delay: 0.8s; }
</style>
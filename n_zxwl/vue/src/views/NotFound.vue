<template>
  <div class="not-found-container">
    <!-- 星空背景 -->
    <div class="stars">
      <div v-for="n in 150" :key="n" class="star" :style="starStyle"></div>
    </div>

    <!-- 漂浮的大学名字粒子 -->
    <div class="university-particles">
      <div
          v-for="(name, index) in universityNames"
          :key="index"
          class="uni-particle"
          :style="particleStyle"
      >
        {{ name }}
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="content">
      <div class="logo-container">
        <img src="../assets/zxwllogo.png" alt="智选未来" class="logo">
      </div>

      <div class="main-content">
        <h1 class="error-code">404</h1>
        <h2 class="error-title">页面迷失在星辰大海</h2>
        <p class="error-desc">
          您寻找的页面可能正在宇宙中遨游，<br>
          让我们带您回到正确的轨道
        </p>

        <div class="action-buttons">
          <button @click="goHome" class="action-btn primary-btn">
            🌟 返回首页
          </button>
          <button @click="goBack" class="action-btn secondary-btn">
            🔙 返回上一页
          </button>
        </div>

        <div class="search-box">
          <input
              v-model="searchQuery"
              @keypress.enter="search"
              placeholder="搜索您需要的内容..."
              class="search-input"
          >
          <button @click="search" class="search-btn">
            🔍
          </button>
        </div>
      </div>
    </div>

    <!-- 装饰性元素 -->
    <div class="decoration">
      <div class="satellite"></div>
      <div class="comet"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { gsap } from 'gsap'

const router = useRouter()
const searchQuery = ref('')

// 大学名字列表
const universityNames = ref([
  '清华大学', '北京大学', '复旦大学', '上海交大', '浙江大学',
  '南京大学', '中国科大', '武汉大学', '华中科大', '中山大学',
  '四川大学', '南开大学', '西安交大', '哈工大', '同济大学'
])

// 星星样式
const starStyle = {
  top: `${Math.random() * 100}%`,
  left: `${Math.random() * 100}%`,
  animationDelay: `${Math.random() * 3}s`,
  animationDuration: `${2 + Math.random() * 2}s`
}

// 粒子样式
const particleStyle = {
  left: `${Math.random() * 100}%`,
  top: `${Math.random() * 100}%`,
  animationDelay: `${Math.random() * 2}s`,
  fontSize: `${12 + Math.random() * 8}px`,
  opacity: `${0.6 + Math.random() * 0.4}`
}

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.go(-1)
}

const search = () => {
  if (searchQuery.value.trim()) {
    window.open(`https://www.baidu.com/s?wd=${searchQuery.value}`)
  }
}

onMounted(() => {
  // 添加入场动画
  gsap.from('.main-content', {
    duration: 1.5,
    y: 50,
    opacity: 0,
    ease: 'back.out(1.7)'
  })

  gsap.from('.uni-particle', {
    duration: 2,
    scale: 0,
    opacity: 0,
    stagger: 0.1,
    ease: 'elastic.out(1, 0.3)'
  })
})
</script>

<style scoped>
.not-found-container {
  height: 100vh;
  background: linear-gradient(135deg,
  #667eea 0%,
  #764ba2 25%,
  #f093fb 50%,
  #f5576c 75%,
  #4facfe 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  font-family: 'Inter', 'Microsoft YaHei', sans-serif;
  color: white;
}

.stars {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.star {
  position: absolute;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: twinkle 3s infinite ease-in-out;
}

.star:nth-child(3n) {
  width: 2px;
  height: 2px;
  opacity: 0.6;
}

.star:nth-child(3n+1) {
  width: 3px;
  height: 3px;
  opacity: 0.8;
}

.star:nth-child(3n+2) {
  width: 1px;
  height: 1px;
  opacity: 0.4;
}

.university-particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.uni-particle {
  position: absolute;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
  white-space: nowrap;
  animation: float 6s infinite ease-in-out;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
  transform: translate(-50%, -50%);
}

.content {
  text-align: center;
  z-index: 2;
  position: relative;
  padding: 2rem;
}

.logo-container {
  margin-bottom: 2rem;
}

.logo {
  width: 120px;
  height: auto;
  filter: drop-shadow(0 0 20px rgba(255, 255, 255, 0.3));
}

.main-content {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 30px;
  padding: 3rem 2rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow:
      0 20px 40px rgba(0, 0, 0, 0.1),
      inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.error-code {
  font-size: 8rem;
  font-weight: 800;
  margin: 0;
  background: linear-gradient(135deg, #fff, #e0e7ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  line-height: 1;
}

.error-title {
  font-size: 2.5rem;
  margin: 1rem 0;
  color: #fff;
  font-weight: 600;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.error-desc {
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.9);
  line-height: 1.6;
  margin-bottom: 2rem;
  text-shadow: 0 1px 5px rgba(0, 0, 0, 0.2);
}

.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.action-btn {
  padding: 1rem 2rem;
  border: none;
  border-radius: 50px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  display: inline-block;
}

.primary-btn {
  background: linear-gradient(135deg, #ff6b6b, #ee5a52);
  color: white;
  box-shadow: 0 8px 25px rgba(255, 107, 107, 0.4);
}

.primary-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 35px rgba(255, 107, 107, 0.6);
}

.secondary-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.secondary-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-3px);
}

.search-box {
  display: flex;
  max-width: 400px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50px;
  padding: 0.5rem;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.search-input {
  flex: 1;
  padding: 1rem 1.5rem;
  border: none;
  background: transparent;
  color: white;
  font-size: 1rem;
  outline: none;
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.7);
}

.search-btn {
  padding: 1rem 1.5rem;
  border: none;
  background: linear-gradient(135deg, #4ecdc4, #44a08d);
  color: white;
  border-radius: 50px;
  cursor: pointer;
  font-size: 1.2rem;
  transition: all 0.3s ease;
}

.search-btn:hover {
  background: linear-gradient(135deg, #44a08d, #4ecdc4);
}

.decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.satellite {
  position: absolute;
  top: 20%;
  right: 10%;
  width: 40px;
  height: 20px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 10px;
  animation: orbit 20s linear infinite;
}

.satellite::before {
  content: '';
  position: absolute;
  top: 50%;
  left: -10px;
  width: 20px;
  height: 4px;
  background: rgba(255, 255, 255, 0.6);
  transform: translateY(-50%);
}

.comet {
  position: absolute;
  top: 10%;
  left: -50px;
  width: 80px;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.8));
  animation: comet 15s linear infinite;
}

@keyframes twinkle {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

@keyframes float {
  0%, 100% { transform: translate(-50%, -50%) translateY(0px); }
  50% { transform: translate(-50%, -50%) translateY(-20px); }
}

@keyframes orbit {
  0% { transform: rotate(0deg) translateX(200px) rotate(0deg); }
  100% { transform: rotate(360deg) translateX(200px) rotate(-360deg); }
}

@keyframes comet {
  0% {
    transform: translateX(0) translateY(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  100% {
    transform: translateX(100vw) translateY(100vh);
    opacity: 0;
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content {
    padding: 1rem;
  }

  .logo {
    width: 80px;
  }

  .error-code {
    font-size: 6rem;
  }

  .error-title {
    font-size: 2rem;
  }

  .error-desc {
    font-size: 1rem;
  }

  .action-buttons {
    flex-direction: column;
    align-items: center;
  }

  .action-btn {
    width: 100%;
    max-width: 250px;
  }

  .search-box {
    flex-direction: column;
    gap: 0.5rem;
    border-radius: 25px;
  }

  .search-input {
    border-radius: 25px;
  }

  .search-btn {
    border-radius: 25px;
    padding: 0.8rem;
  }

  .main-content {
    padding: 2rem 1.5rem;
  }
}

@media (max-width: 480px) {
  .error-code {
    font-size: 4rem;
  }

  .error-title {
    font-size: 1.5rem;
  }

  .action-btn {
    padding: 0.8rem 1.5rem;
    font-size: 1rem;
  }
}
</style>
<template>
  <div class="home-container">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <!-- 动态粒子背景 -->
    <div class="particle-background">
      <canvas ref="canvasRef" class="particle-canvas"></canvas>
    </div>

    <!-- 主内容 -->
    <div class="main-content">
      <!-- 动态导航 -->
      <nav class="floating-nav">
        <!-- 左侧的导航栏 -->
        <div class="nav-inner">
          <router-link
              v-for="(item, index) in navItems"
              :key="item.path"
              :to="item.path"
              class="nav-item"
              @mouseenter="animateNavHover(index)"
              @mouseleave="resetNavHover(index)"
          >
            <!-- 遍历navItems数组，为每个元素创建一个导航项 -->
            <component :is="item.icon" class="nav-icon" />
            <span class="nav-text">{{ item.name }}</span>
            <div class="nav-underline" />
          </router-link>
        </div>

        <!-- 右侧用户功能区 -->
        <div class="nav-right">
          <div class="user-section" ref="userSectionRef">
            <router-link v-if="!isLoggedIn" to="/login" class="login-btn">登录/注册</router-link>
            <div v-else class="user-icon" @click.stop="toggleUserMenu">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
              </svg>
            </div>
            <!-- 用户下拉菜单 -->
            <transition name="slide-fade">
              <div v-show="isUserMenuVisible" class="user-menu glass-effect">
                <router-link to="/profile" class="menu-item" @click="isUserMenuVisible = false">
                  <svg class="menu-icon" viewBox="0 0 24 24">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
                  </svg>
                  个人中心
                </router-link>

                <template v-if="userData.email === 'root@root.com'">
                  <div class="menu-divider"></div>

                  <router-link to="/admin/users" class="menu-item" @click="isUserMenuVisible = false">
                    <svg class="menu-icon" viewBox="0 0 24 24">
                      <path d="M12 5.5c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2M6 22c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2m6 0c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2m6-11c0-3.3-2.7-6-6-6s-6 2.7-6 6c0 1.2.4 2.4 1 3.3v1.2c0 .3.2.5.5.5h9c.3 0 .5-.2.5-.5v-1.2c.6-.9 1-2.1 1-3.3"/>
                    </svg>
                    用户管理
                  </router-link>

                  <router-link to="/admin/policies/create" class="menu-item" @click="isUserMenuVisible = false">
                    <svg class="menu-icon" viewBox="0 0 24 24">
                      <path d="M19 3h-4.18C14.4 1.84 13.3 1 12 1s-2.4.84-2.82 2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2m-7 0c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1M8 17H6v-2h2v2m0-4H6v-2h2v2m0-4H6V7h2v2m6 8h-4v-2h4v2m4-4h-8v-2h8v2m0-4h-8V7h8v2z"/>
                    </svg>
                    添加政策
                  </router-link>

                  <router-link to="/admin/logs" class="menu-item" @click="isUserMenuVisible = false">
                    <svg class="menu-icon" viewBox="0 0 24 24">
                      <path d="M14 12v-2h-4v2H9v4h6v-4h-1zm6-6v12c0 1.1-.9 2-2 2H6c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2h12c1.1 0 2 .9 2 2zm-2 0H6v12h12V6z"/>
                    </svg>
                    系统日志
                  </router-link>
                </template>

                <div class="menu-divider"></div>

                <div class="menu-item" @click="handleLogout">
                  <svg class="menu-icon" viewBox="0 0 24 24">
                    <path d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"/>
                  </svg>
                  退出登录
                </div>
              </div>
            </transition>
          </div>
        </div>
      </nav>

      <!-- 3D轮播图 -->
      <div class="carousel-container">
        <div class="carousel-track">
          <div
              v-for="(item, index) in carouselItems"
              :key="index"
              class="carousel-item"
              :class="{ active: currentIndex === index }"
              @mouseenter="pauseAutoPlay"
              @mouseleave="resumeAutoPlay"
              ref="carouselItemsRef"
          >
            <div class="card-content">
              <div
                  class="card-image"
                  :style="{ backgroundImage: `url(${item.image})` }"
                  @load="handleImageLoad"
              ></div>
              <div class="card-text">
                <h3>{{ item.title }}</h3>
                <p>{{ item.description }}</p>
                <router-link :to="item.path" class="card-link">了解更多</router-link>

              </div>
            </div>
          </div>
        </div>
        <div class="carousel-dots">
          <button
              v-for="(_, index) in carouselItems"
              :key="index"
              :class="{ active: currentIndex === index }"
              @click="goToSlide(index)"
          />
        </div>
      </div>

      <!-- 动态搜索框 -->
      <div class="search-wrapper">
        <div class="search-bar" :class="{ focused: isSearchFocused }">
          <input
              v-model="searchKeyword"
              @focus="isSearchFocused = true"
              @blur="isSearchFocused = false"
              placeholder="搜索院校/专业..."
              @keyup.enter="handleSearch"
          />
          <button class="search-button" @mousedown.prevent="handleSearch">
            <Magnify class="search-icon" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { gsap } from 'gsap';
import { useRouter } from 'vue-router';
import { ParticleEngine } from '@/utils/particleEngine';
import { useUserStore } from  '@/utils/auth.js'
import SearchResults from '../components/searchSchAndSpe.vue'
import Magnify from 'vue-material-design-icons/Magnify.vue';
import {ElMessage} from "element-plus";
import schoolsImage from '../assets/chun.jpg';
import majorsImage from '../assets/xia.jpg';
import policyImage from '../assets/qiu.jpg';
import simulateImage from '../assets/dong.jpg';
import materialsImage from '../assets/bbjj.jpg';
import axios from "axios";

const router = useRouter();
const userData = ref({
  name: '',
  email: ''
})

const {deleteUser,getUser } = useUserStore()
const isLoggedIn = ref(true);
const isUserMenuVisible = ref(false);
const userSectionRef = ref(null);
const userEmail = localStorage.getItem('userEmail')
// 粒子引擎实例
const canvasRef = ref(null);
let particleEngine = null;

// 导航项数据
const navItems = [
  { name: '智选对话', path: '/aismartsel'},
  { name: '智能推荐', path: '/recommends' },
  { name: '一分一段', path: '/score-section' },
  { name: '志愿模拟', path: '/simulate' },
  { name: '政策阅读', path: '/policy'},
  {name: '家庭共享',  path: '/user-friend'}
];

const carouselItems = ref([
  // {
  //   title: '院校名录',
  //   description: '涵盖广泛，各类院校信息一应俱全',
  //   path: '/colleges',
  //   image: schoolsImage,
  // },
  {
    title: '院校名录',
    description: '涵盖广泛，各类院校信息一应俱全',
    path: '/allSchool',
    image: schoolsImage,
  },
  {
    title: '专业诠衡',
    description: '深度剖析专业前景与就业方向',
    path: '/professional',
    image: majorsImage,
  },
  {
    title: '政策解读',
    description: '最新高考政策权威解读',
    path: '/policy',
    image: policyImage,
  },
  {
    title: '志愿指南',
    description: '科学填报志愿方法与技巧',
    path: '/simulate',
    image: simulateImage,
  },
  {
    title: '咨询论坛',
    description: '和大牛们一起聊高考吧',
    path: '/news',
    image: materialsImage,
  },
]);

// 搜索功能
const searchKeyword = ref('');
const isSearchFocused = ref(false);
const loading = ref(false)
const UASData = ref()

// 在搜索成功后设置数据
const handleSearch = async () => {
  if (searchKeyword.value.trim() === '') return;
  loading.value = true
  // 将数据传递给搜索结果页面
  router.push({
    path: '/searchUAS',
    query: {
      keyword: searchKeyword.value,
    }
  })
  loading.value = false
}

// 轮播图逻辑
const currentIndex = ref(0);
const carouselItemsRef = ref([]);
let autoPlayTimer = null;

const animateSlide = (newIndex) => {
  const currentItem = carouselItemsRef.value[currentIndex.value];
  const nextItem = carouselItemsRef.value[newIndex];

  gsap.to(currentItem, {
    duration: 1.2,
    rotationY: -180,
    opacity: 0,
    ease: 'power4.inOut',
    onComplete: () => {
      currentIndex.value = newIndex;
    },
  });

  gsap.fromTo(
      nextItem,
      { rotationY: 180, opacity: 0 },
      {
        duration: 1.2,
        rotationY: 0,
        opacity: 1,
        ease: 'power4.inOut',
      }
  );
};

// 用户菜单功能
const toggleUserMenu = () => {
  isUserMenuVisible.value = !isUserMenuVisible.value;
};
const handleClickOutside = (event) => {
  if (userSectionRef.value && !userSectionRef.value.contains(event.target)) {
    isUserMenuVisible.value = false;
  }
};
const handleImageLoad = (event) => {
  event.target.parentElement.classList.add('loaded');
};
const handleLogout = () => {
  isLoggedIn.value = false;
  isUserMenuVisible.value = false;
  deleteUser()
  router.push('/');
};

// 自动播放控制
const startAutoPlay = () => {
  autoPlayTimer = setInterval(() => {
    const newIndex = (currentIndex.value + 1) % carouselItems.value.length;
    animateSlide(newIndex);
  }, 5000);
};
const pauseAutoPlay = () => clearInterval(autoPlayTimer);
const resumeAutoPlay = () => startAutoPlay();
const goToSlide = (index) => {
  if (index === currentIndex.value) return;
  animateSlide(index);
};

// 导航栏动画 - 显示
const animateNavHover = (index) => {
  // 先杀死该元素的所有进行中的动画
  gsap.killTweensOf(`.nav-item:nth-child(${index + 1}) .nav-underline`);

  gsap.to(`.nav-item:nth-child(${index + 1}) .nav-underline`, {
    width: '80%',
    duration: 0.3, // 缩短动画时间
    ease: 'power2.out',
  });
};

// 导航栏动画 - 重置
const resetNavHover = (index) => {
  // 先杀死该元素的所有进行中的动画
  gsap.killTweensOf(`.nav-item:nth-child(${index + 1}) .nav-underline`);

  gsap.to(`.nav-item:nth-child(${index + 1}) .nav-underline`, {
    width: '0%',
    duration: 0.2, // 缩短动画时间
    ease: 'power2.out',
  });
};

// 生命周期
onMounted(() => {

  particleEngine = new ParticleEngine(canvasRef.value);
  particleEngine.init();

  gsap.set(carouselItemsRef.value, {
    transformOrigin: '50% 50% 0',
    perspective: 1000,
    transformStyle: 'preserve-3d',
    position: 'absolute',
    backfaceVisibility: 'hidden',
  });

  gsap.set(carouselItemsRef.value, { opacity: 0, rotationY: 180 });
  gsap.set(carouselItemsRef.value[0], { opacity: 1, rotationY: 0 });

  startAutoPlay();
  window.addEventListener('click', handleClickOutside);

  const user = getUser();
  if (user)
    userData.value = {
      name: user.name || '',
      email: user.email || ''
    }
});

onUnmounted(() => {
  particleEngine?.destroy();
  pauseAutoPlay();
  window.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.home-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: url(../assets/colleges.jpg) no-repeat center center fixed;
  background-size: cover;
}

.particle-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
}

.particle-canvas {
  position: fixed;
  top: 0;
  left: 0;
  pointer-events: none;
}

.floating-nav {
  position: fixed;
  top: 1.5rem;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(255, 255, 255, 0.30);
  backdrop-filter: blur(12px);
  border-radius: 1.5rem;
  padding: 0.6rem 2rem;
  z-index: 100;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  display: flex;
  justify-content: space-between;
  width:67%;

}

.nav-inner {
  display: flex;
  gap: 2rem;
  position: relative;
  left: 0%;
}

.page-logo {
  position: absolute;
  top: 1rem;
  left: 1rem;
  width: 110px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
}
.nav-item {
  position: relative;
  cursor: pointer;
  padding: 0.4rem 0.8rem;
  transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
  perspective: 500px;
}

.nav-text {
  display: inline-block;
  font-size: 1.1rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  color: rgb(0, 0, 0);
  transform: translateZ(10px);
  text-shadow:
      0.05em 0.05em 0 rgba(255,255,255,0.4),
      -0.025em -0.025em 0 rgba(0,0,0,0.3),
      0.1em 0.15em 0.05em rgba(0,0,0,0.15),
      0 0 0.3em rgba(76,161,255,0.5);
  transition:
      transform 0.3s ease,
      text-shadow 0.3s ease;
}

.nav-item:hover .nav-text {
  color: #fff;
  transform:
      translateZ(20px)
      rotateX(-10deg)
      rotateY(5deg);
  text-shadow:
      0.1em 0.1em 0 rgba(255,255,255,0.4),
      -0.05em -0.05em 0 rgba(0,0,0,0.3),
      0.2em 0.3em 0.1em rgba(0,0,0,0.2),
      0 0 0.5em rgba(76,161,255,0.8),
      0 0 1em rgba(123,44,191,0.5);
}

.nav-icon {
  width: 24px;
  height: 24px;
  margin-right: 0.5rem;
  vertical-align: middle;
  transition: transform 0.3s ease;
  filter: drop-shadow(1px 1px 1px rgba(0,0,0,0.2));
}

.nav-item:hover .nav-icon {
  transform: scale(1.15) translateY(-1px);
}

.nav-underline {
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, #4CA1FF, #7B2CBF);
}

.nav-right {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
}

.user-section {
  position: relative;
  display: flex;
  align-items: center;
}

.login-btn {
  color: #ffffff;
  padding: 8px 16px;
  border-radius: 20px;
  background: linear-gradient(45deg, #4CA1FF, #7B2CBF);
  text-decoration: none;
  transition: transform 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
}

.user-icon {
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
}

.user-icon svg {
  fill: #ffffff;
  transition: fill 0.3s;
}

.user-icon:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: rotate(15deg);
}

.user-menu {
  position: absolute;
  right: -60px;
  top: 120%;
  min-width: 160px;
  border-radius: 12px;
  padding: 8px 0;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  z-index: 1001;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  color: #333;
  text-decoration: none;
  transition: all 0.2s;
  font-size: 14px;
}

.menu-item:hover {
  background: rgba(76, 161, 255, 0.1);
  color: #4CA1FF;
}

.menu-icon {
  width: 18px;
  height: 18px;
  margin-right: 12px;
  fill: currentColor;
}

.menu-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.1);
  margin: 6px 16px;
}

.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

.carousel-container {
  position: relative;
  width: 1200px;
  height: 450px;
  margin: 6.5rem auto;
  perspective: 2000px;
}

.carousel-track {
  position: relative;
  width: 100%;
  height: 100%;
}

.carousel-item {
  width: 100%;
  height: 100%;
  padding: 1rem;
  backface-visibility: hidden;
}

.card-content {
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(12px);
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.3);
  transform-style: preserve-3d;
}

.card-image {
  height: 450px;
  background-size: cover;
  background-position: center;
  clip-path: none;
  background-repeat: no-repeat;
  background-color: rgba(255, 255, 255, 0.1);
  transition: opacity 0.5s ease;
  position: relative;
}

.card-image.loaded {
  opacity: 1;
}

.card-text {
  background-color: rgba(255, 255, 255, 0.4);
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: auto;
  padding: 2rem;
  color: #000000;
  transform: translateZ(50px);
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  padding: 1rem;
  border-radius: 0.5rem;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.card-text h3 {
  font-size: 2rem;
  margin-bottom: 1rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.card-text p {
  font-size: 1.1rem;
  line-height: 1.6;
  opacity: 0.9;
  margin-bottom: 2rem;
}

.card-link {
  display: inline-block;
  padding: 0.8rem 2rem;
  background: linear-gradient(45deg, #4CA1FF, #7B2CBF);
  border-radius: 30px;
  color: white;
  text-decoration: none;
  transition: transform 0.3s ease;
}

.card-link:hover {
  transform: translateY(-2px);
}

.carousel-dots {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
}

.carousel-dots button {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: none;
  background: rgba(255, 255, 255, 0.4);
  cursor: pointer;
  transition: all 0.3s ease;
}

.carousel-dots button.active {
  background: #4CA1FF;
  transform: scale(1.3);
}

.search-wrapper {
  position: fixed;
  bottom: 10vh;
  left: 50%;
  transform: translateX(-50%);
}
.search-bar {
  width: 500px;
  height: 56px;
  background: rgb(255, 255, 255);
  border-radius: 2rem;
  display: flex;
  align-items: center;
  padding: 0 1rem;
  transition: all 0.3s ease;
}

.search-bar.focused {
  width: 600px;
  background: rgb(250, 250, 250);
}

.search-bar input {
  flex: 1;
  background: transparent;
  border: none;
  color: #333;
  font-size: 1.1rem;
  padding: 0 1rem;
  outline: none;
}
.search-button {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(45deg, #4CA1FF, #7B2CBF);
  border: none;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.search-button:hover {
  transform: scale(1.1);
}

@media (max-width: 768px) {
  .floating-nav {
    top: 1rem;
    padding: 0.4rem 0.8rem;
    width: 100%;
  }

  .nav-inner {
    gap: 1rem;
    left: 0;
  }

  .nav-text {
    font-size: 0.9rem;
  }

  .nav-icon {
    width: 20px;
    height: 20px;
  }

  .nav-right {
    right: 1rem;
  }

  .carousel-container {
    width: 90%;
    height: auto;
    aspect-ratio: 1200 / 450;
    margin: 3rem auto;
  }

  .card-image {
    height: auto;
    aspect-ratio: 1200 / 300;
    clip-path: none;
  }

  .search-bar {
    width: 90vw;
  }

  .search-bar.focused {
    width: 90vw;
  }

  .card-text h3 {
    font-size: 1.5rem;
  }

  .card-text p {
    font-size: 1rem;
  }
  .page-logo {
    width: 60px; /* 小屏幕下调整 logo 大小 */
  }
}
</style>
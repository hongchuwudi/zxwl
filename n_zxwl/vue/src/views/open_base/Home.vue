<template>
  <div class="container">
    <!-- logo 图片 -->
    <div class="logo-container">
      <img src="../../assets/zxwllogo.png" alt="智选未来logo">
    </div>
    <!-- 星空背景 -->
    <div class="star-field">
      <div
          v-for="(star, index) in stars"
          :key="'star'+index"
          class="star"
          :style="star.style"
      ></div>
    </div>

    <!-- 大学粒子 -->
    <div class="particles-container">
      <div
          v-for="(name, index) in universityNames"
          :key="index"
          class="university-particle"
          :ref="el => { if(el) particlesRefs[index] = el }"
      >
        {{ name }}
      </div>
    </div>

    <!-- 主内容 -->
    <div class="content-container">
      <div class="particle-title">
        <span class="title-text">智选未来</span>
        <span class="subtitle-text">高考志愿服务平台</span>
        <div
            v-for="(point, index) in titleParticles"
            :key="'title-'+index"
            class="title-particle"
            :style="point.style"
        ></div>
      </div>
      <div class="button-container">

        <button
            class="auth-button"
            @click="handleLoginClick"
            :disabled="isAnimating"
            ref="loginButton"
        >
          登录 / 注册
        </button>
        <button
            class="auth-button"
            @click="handleGuestClick"
            :disabled="isAnimating"
            ref="guestButton"
        >
          游客进入
        </button>
        <!-- 按钮粒子容器 -->
        <div
            v-if="showButtonParticles"
            class="button-particles-container"
            ref="buttonParticlesContainer"
        >
          <div
              v-for="(p, index) in buttonParticles"
              :key="'btn-particle-'+index"
              class="button-particle"
              :style="p.style"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onUnmounted } from 'vue'
import { gsap } from 'gsap'
import { useRouter } from 'vue-router'

const router = useRouter()

// 星空配置
const stars = ref(Array.from({ length: 150 }, (_, i) => ({
  style: {
    left: gsap.utils.random(0, 100) + '%',
    top: gsap.utils.random(0, 100) + '%',
    animationDelay: gsap.utils.random(0, 3) + 's'
  }
})))

// 大学列表
const universityNames = ref([
  '清华大学', '北京大学', '复旦大学', '上海交通大学',
  '浙江大学', '南京大学', '中国科大', '武汉大学',
  '华中科大', '中山大学', '四川大学', '南开大学',
  '西安交大', '哈尔滨工业大学', '同济大学', '东南大学',
  '中国科学技术大学', '北京航空航天大学', '北京理工大学', '中国人民大学',
  '北京师范大学', '天津大学', '山东大学', '西北工业大学',
  '厦门大学', '中南大学', '吉林大学', '中国农业大学',
  '大连理工大学', '华东师范大学', '华南理工大学', '电子科技大学',
  '天津师范大学', '西安石油大学','北京工业大学', '江南大学'
])

// 标题粒子配置
const titleParticles = ref(Array.from({ length: 800 }, (_, i) => ({
  style: {
    x: gsap.utils.random(-50, 650),
    y: gsap.utils.random(-50, 170),
    scale: gsap.utils.random(0.5, 1.5),
    opacity: 0,
    rotation: gsap.utils.random(0, 360)
  }
})))

// 按钮粒子相关
const isAnimating = ref(false)
const showButtonParticles = ref(false)
const buttonParticles = ref([])
const loginButton = ref(null)
const guestButton = ref(null)
const buttonParticlesContainer = ref(null)

const particlesRefs = ref([])
let ctx

// 粒子初始化方法
const initParticlePosition = (el) => {
  gsap.set(el, {
    left: gsap.utils.random(5, 95) + '%',
    top: gsap.utils.random(5, 95) + '%',
    rotation: gsap.utils.random(-15, 15),
    scale: 0,
    opacity: 0,
    willChange: 'transform, opacity'
  })
}

const createButtonParticles = (buttonEl) => {
  if(!buttonEl) return

  const rect = buttonEl.getBoundingClientRect()
  const particles = []
  const particleCount = 30

  for(let i = 0; i < particleCount; i++) {
    particles.push({
      style: {
        x: rect.left + gsap.utils.random(0, rect.width),
        y: rect.top + gsap.utils.random(0, rect.height),
        opacity: 1,
        scale: 1
      }
    })
  }

  buttonParticles.value = particles
  showButtonParticles.value = true
}

const animateButtonParticles = (buttonEl) => {
  return new Promise(resolve => {
    if(!buttonParticlesContainer.value) return resolve()

    const tl = gsap.timeline({
      onComplete: () => resolve()
    })

    tl.to('.button-particle', {
      duration: 1.2,
      x: () => gsap.utils.random(-200, 200),
      y: () => gsap.utils.random(-150, 150),
      scale: 0,
      opacity: 0,
      ease: 'power4.out',
      stagger: 0.02
    }, 0)

    tl.to(buttonEl, {
      duration: 0.6,
      scale: 0.8,
      opacity: 0,
      ease: 'back.in(2)'
    }, 0)
  })
}

const handleAction = async (buttonEl, routePath) => {
  if(isAnimating.value) return
  isAnimating.value = true

  createButtonParticles(buttonEl)
  await animateButtonParticles(buttonEl)

  router.push(routePath)
  isAnimating.value = false
  showButtonParticles.value = false
}

const handleLoginClick = async () => {
  await handleAction(loginButton.value, '/login')
}

const handleGuestClick = async () => {
  window.localStorage.removeItem("userEmail");
  await handleAction(guestButton.value, '/zxwl')
}

onMounted(async () => {
  await nextTick()

  ctx = gsap.context(() => {
    gsap.set('.content-container', { autoAlpha: 0, scale: 0.8 })
    particlesRefs.value.forEach(el => el && initParticlePosition(el))

    const particleTimelines = particlesRefs.value.map((el, index) => {
      if (!el) return Promise.resolve()

      return new Promise(resolve => {
        const tl = gsap.timeline({
          delay: index * 0.12,
          onComplete: resolve
        })

        tl.to(el, {
          duration: 0.8,
          opacity: 1,
          scale: 1,
          rotation: 0,
          ease: 'back.out(1.8)'
        })
            .to(el, {
              duration: 1.2,
              x: '+=random(-15%,15%)',
              y: '+=random(-10%,10%)',
              ease: 'sine.inOut',
            }, 0)
            .to(el, {
              duration: 0.6,
              opacity: 0,
              scale: 0.5,
              rotation: gsap.utils.random(-25, 25),
              ease: 'power3.in',
              delay: gsap.utils.random(0.3, 0.7),
              onComplete: () => gsap.set(el, { display: 'none' })
            }, '+=0.3')
      })
    })

    gsap.to('.title-particle', {
      duration: 8,
      x: '+=random(-800,800)',
      y: '+=random(-600,600)',
      rotation: '+=random(-180,180)',
      scale: 'random(0.4,1.6)',
      opacity: 'random(0.7,1)',
      repeat: -1,
      yoyo: true,
      yoyoEase: "sine.inOut",
      ease: 'sine.inOut',
      stagger: {
        each: 0.2,
        from: 'random'
      },
      delay: 0.5
    })

    Promise.all(particleTimelines)
        .then(() => {
          gsap.to('.content-container', {
            duration: 1.5,
            autoAlpha: 1,
            scale: 1,
            ease: 'elastic.out(1, 0.3)'
          })
          gsap.to('.button-container', {
            duration: 0.8,
            y: 0,
            opacity: 1,
            ease: 'power4.out'
          }, '-=0.5')
        })
        .catch(() => {
          gsap.set('.content-container', { autoAlpha: 1 })
          gsap.set('.button-container', { opacity: 1 })
        })
  })
})

onUnmounted(() => {
  ctx && ctx.revert()
})
</script>

<style scoped>
.container {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: url(../../assets/loginloge.png) no-repeat center center fixed;
  background-size: cover;
}

.logo-container {
  position: absolute;
  top: 10px;
  left: 10px;
  z-index: 4;
}

.logo-container img {
  max-width: 100px;
  height: auto;
}

.star-field {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.particles-container {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 2;
}

.content-container {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  z-index: 3;
  pointer-events: none;
}

.star {
  position: absolute;
  width: 2px;
  height: 2px;
  background: rgb(255, 255, 255);
  border-radius: 50%;
  animation: twinkle 3s infinite ease-in-out;
  filter: blur(0.5px);
}

@keyframes twinkle {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

.university-particle {
  position: absolute;
  color: rgb(255, 255, 255);
  font-size: 1.1rem;
  font-weight: 600;
  white-space: nowrap;
  transform: translate(-50%, -50%);
  filter:
      drop-shadow(0 0 8px rgba(255,255,255,0.6))
      drop-shadow(0 0 15px rgba(255,255,255,0.3));
}

.particle-title {
  position: relative;
  width: 600px;
  height: 160px;
  margin: 0 auto 2rem;
}

.title-text {
  position: absolute;
  top: -30%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 72px;
  color: rgb(255, 255, 255);
  font-family: system-ui, -apple-system, sans-serif;
  text-shadow: 0 0 20px rgba(255,255,255,0.8);
}

.title-particle {
  position: absolute;
  width: 5px;
  height: 5px;
  background: rgb(255, 255, 255);
  border-radius: 50%;
  filter: blur(1px);
  box-shadow: 0 0 15px rgba(255,255,255,0.8);
  transform: translate(-50%, -50%);
  will-change: transform;
}

.subtitle-text {
  position: absolute;
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 60px;
  color: rgba(255,255,255,0.9);
  font-family: system-ui, -apple-system, sans-serif;
  text-shadow: 0 0 20px rgba(255,255,255,0.8);
  white-space: nowrap;
  width: 100%;
  text-align: center;
}

.button-container {
  opacity: 0;
  transform: translateY(20px);
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  align-items: center;
}

.auth-button {
  padding: 14px 40px;
  font-size: 1.2rem;
  background: transparent;
  border: 2px solid rgba(255,255,255,0.8);
  color: #fff;
  border-radius: 30px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow:
      0 0 15px rgba(255,255,255,0.3),
      inset 0 0 10px rgba(255,255,255,0.2);
  pointer-events: auto;
}

.auth-button:hover {
  background: rgba(255,255,255,0.05);
  transform: translateY(-3px);
  box-shadow:
      0 0 25px rgba(255,255,255,0.5),
      inset 0 0 15px rgba(255,255,255,0.3);
}

.button-particles-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.button-particle {
  position: absolute;
  width: 8px;
  height: 8px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  filter: blur(2px);
  box-shadow:
      0 0 15px rgba(255,255,255,0.8),
      0 0 30px rgba(255,255,255,0.6);
  transform: translate(-50%, -50%);
}

@media (max-width: 768px) {
  .particle-title {
    height: 100px;
  }
  .title-text {
    top: 30%;
    font-size: 36px;
  }
  .subtitle-text {
    top: 60%;
    font-size: 18px;
  }
  .auth-button {
    padding: 10px 30px;
    font-size: 1rem;
  }
  .button-particle {
    width: 6px;
    height: 6px;
  }
  .button-container {
    gap: 15px;
  }
}
</style>
<template>
  <div class="login-container">
    <!-- 页面左上角添加 logo -->
    <img src="../../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <!-- 星空背景 -->
    <div class="star-field">
      <div
          v-for="(star, index) in stars"
          :key="'star'+index"
          class="star"
          :style="star.style"
      ></div>
    </div>

    <!-- 登录框 -->
    <div class="login-box">
      <div class="form-wrapper">
        <el-form ref="formRef" :rules="data.rules" :model="data.form">
          <div class="form-header">
            <h1 class="title-text">智选未来</h1>
            <h2 class="subtitle-text">高考志愿服务平台</h2>
          </div>

          <el-form-item prop="email">
            <el-input
                size="large"
                v-model="data.form.email"
                placeholder="邮箱地址"
                prefix-icon="Message"
                class="custom-input"
            />
          </el-form-item>

          <el-form-item prop="passwd">
            <el-input
                show-password
                size="large"
                v-model="data.form.passwd"
                placeholder="密码"
                prefix-icon="Lock"
                class="custom-input"
            />
          </el-form-item>

          <el-button
              @click="login"
              size="large"
              class="login-btn"
              :loading="isLogging"
          >
            <span class="btn-text">立即登录</span>
            <div class="btn-glow"></div>
          </el-button>

          <div class="action-links">
            <router-link to="/register" class="link-text">立即注册</router-link>
            <router-link to="/Reset" class="link-text">忘记密码？</router-link>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { Message, Lock } from "@element-plus/icons-vue";
import axios from "axios";
import { useUserStore } from '@/utils/auth.js';
const { saveUser } = useUserStore();
import DeviceUtils from '@/utils/deviceInfo.js'
import { initWebSocket } from '@/utils/wsUtil.js'; // 导入 WebSocket
import { useRouter } from 'vue-router';
// 星空配置
const stars = Array.from({ length: 150 }, (_, i) => ({
  style: {
    left: Math.random() * 100 + "%",
    top: Math.random() * 100 + "%",
    animationDelay: Math.random() * 3 + "s",
  },
}));
const router = useRouter();
const data = reactive({
  form: {
    email: "",
    passwd: "",
  },
  rules: {
    email: [
      { required: true, message: "请输入邮箱地址或用户名", trigger: "blur" },
    ],
    passwd: [
      { required: true, message: "请输入密码", trigger: "blur" },
      { min: 6, max: 18, message: "长度在6到18个字符", trigger: "blur" },
    ],
  },
});

const formRef = ref();
const isLogging = ref(false);

const login = () => {
  formRef.value.validate(async (valid) => {
    if (!valid) return;

    try {
      isLogging.value = true;
      const response = await axios.post("gapi/user/login", {
        login: data.form.email,
        password: data.form.passwd,
        device_info: DeviceUtils.getDeviceInfoJSON()
      }, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.data.error === 0 && response.data.code === 0) {
        ElMessage.success("登录成功,欢迎来到智选未来！");

        // 保存用户信息
        const userInfo = {
          name: response.data.data.user.username,
          id: response.data.data.user.id,
          email: data.form.email,
          displayName: response.data.data.user.displayName,
          avatarUrl: response.data.data.user.avatarUrl,
          gender: response.data.data.user.gender,
          birthYear: response.data.data.user.birthYear,
          location: response.data.data.user.location,
          bio: response.data.data.user.bio,
          isOnline: response.data.data.user.isOnline,
          token: response.data.data.token // 保存 token
        };
        saveUser(userInfo);

        // 开启 WebSocket 连接
        try {
          initWebSocket(userInfo.id, userInfo.email);
          console.log('WebSocket连接初始化成功');
        } catch (error) {
          console.error('WebSocket连接失败:', error);
          ElMessage.error('WebSocket连接失败');
        }

        // 添加日志
        // const logData = {
        //   "email": data.form.email,
        //   "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
        //   "operation": "用户登录"
        // };
        // const logResponse = await axios.post("gapi/log", logData, {
        //   headers: {
        //     "Content-Type": "application/json"
        //   }
        // });

        setTimeout(() => {
          router.push('/Zxwl');
        }, 750);
      } else {
        ElMessage.error(response.data.message || "登录失败");
      }
    } catch (error) {
      ElMessage.error(`请求失败：${error.message}`);
    } finally {
      isLogging.value = false;
    }
  });
};
</script>

<style scoped>
.login-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: radial-gradient(ellipse at center,
  rgb(171, 193, 220) 0%,
  rgb(150, 159, 225) 100%);
}

.star-field {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.star {
  position: absolute;
  width: 2px;
  height: 2px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: twinkle 3s infinite ease-in-out;
  filter: blur(0.5px);
}

@keyframes twinkle {
  0%,
  100% {
    opacity: 0.3;
  }
  50% {
    opacity: 1;
  }
}

.login-box {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  max-width: 480px;
  padding: 2rem;
  z-index: 2;
}

/* 页面 logo 样式 */
.page-logo {
  position: absolute;
  top: 1rem;
  left: 1rem;
  width: 100px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
}

.form-wrapper {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  padding: 2.5rem 2rem;
  box-shadow: 0 0 40px rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.form-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.title-text {
  font-size: 2.9rem;
  color: #8f9fff;
  margin-bottom: 0.5rem;
  font-weight: 600;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subtitle-text {
  font-size: 1.5rem;
  color: #889ffd;
  font-weight: 400;
}

.custom-input {
  margin-bottom: 1.8rem;
  transition: all 0.3s ease;
}

.custom-input:hover {
  transform: translateY(-2px);
}

.login-btn {
  width: 100%;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  padding: 1rem;
  border-radius: 12px;
  transition: all 0.3s ease;
  margin: 1.5rem 0;
}

.login-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
}

.btn-text {
  position: relative;
  z-index: 2;
  letter-spacing: 1px;
}

.btn-glow {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg,
  transparent,
  rgba(255, 255, 255, 0.3),
  transparent);
  transition: all 0.6s ease;
}

.login-btn:hover .btn-glow {
  left: 100%;
}

.action-links {
  display: flex;
  justify-content: space-between;
  margin-top: 1.5rem;
}

.link-text {
  color: #667eea;
  font-weight: 500;
  text-decoration: none;
  position: relative;
  transition: all 0.3s ease;
}

.link-text:hover {
  color: #764ba2;
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .login-box {
    padding: 1rem;
    width: 90%;
  }

  .form-wrapper {
    padding: 2rem 1.5rem;
  }

  .title-text {
    font-size: 1.8rem;
  }

  .subtitle-text {
    font-size: 1rem;
  }

  .page-logo {
    width: 60px; /* 小屏幕下调整 logo 大小 */
  }
}
</style>
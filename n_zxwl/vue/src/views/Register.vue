<template>
  <div class="login-container">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <div class="login-box">
      <div class="form-wrapper">
        <div class="decorative-line"></div>
        <el-form ref="formRef" :rules="rules" :model="form" label-position="top">
          <div class="form-header">
            <span class="title-text">Welcome to Register</span>
            <div class="title-underline">
              <div class="underline"></div>
              <div class="circle"></div>
            </div>
          </div>

          <el-form-item prop="username">
            <el-input
                v-model="form.username"
                placeholder="用户名"
                prefix-icon="User"
                size="large"
                class="custom-input"
            />
          </el-form-item>

          <el-form-item prop="email">
            <el-input
                v-model="form.email"
                placeholder="邮箱地址"
                prefix-icon="Message"
                size="large"
                class="custom-input"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
                v-model="form.password"
                type="password"
                show-password
                placeholder="密码"
                prefix-icon="Lock"
                size="large"
                class="custom-input"
            />
          </el-form-item>

          <el-form-item prop="confirmPassword">
            <el-input
                v-model="form.confirmPassword"
                type="password"
                show-password
                placeholder="确认密码"
                prefix-icon="Lock"
                size="large"
                class="custom-input"
            />
          </el-form-item>

          <div class="verification-container">
            <el-form-item prop="verificationCode" class="verification-input">
              <el-input
                  v-model="form.verificationCode"
                  placeholder="验证码"
                  prefix-icon="Key"
                  size="large"
                  class="custom-input"
              />
            </el-form-item>
            <el-button
                @click="getVerificationCode"
                size="large"
                class="verification-btn"
                :disabled="isCodeSending"
            >
              {{ codeButtonText }}
            </el-button>
          </div>

          <el-button
              @click="handleRegister"
              size="large"
              class="register-btn"
              :loading="isRegistering"
          >
            <span class="btn-text">立即注册</span>
            <div class="fill-container"></div>
          </el-button>

          <div class="login-link">
            已有账号？
            <router-link to="/login" class="link-text">马上登录</router-link>
          </div>
        </el-form>
      </div>
    </div>
    <div class="decorative-bubbles">
      <div v-for="i in 10" :key="i" class="bubble"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Lock, Message, Key } from '@element-plus/icons-vue'
import axios from 'axios'

const formRef = ref()
const isCodeSending = ref(false)
const codeButtonText = ref('获取验证码')
const isRegistering = ref(false)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  verificationCode: ''
})

const rules = reactive({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 12, message: '长度在3到12个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: ['blur', 'change'] }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 18, message: '长度在6到18个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== form.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  verificationCode: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
})

const getVerificationCode = async () => {
  try {
    isCodeSending.value = true
    const email = form.email;
    const response = await axios.post(`gapi/get_varifycode`, { email });//返回的响应是 {"error":0,"msg":"验证码已发送"}

    if (response.data.code === 0) {
      ElMessage.success('验证码已发送')
      let countdown = 60
      const timer = setInterval(() => {
        codeButtonText.value = `${countdown}秒后重发`
        if (countdown-- <= 0) {
          clearInterval(timer)
          codeButtonText.value = '获取验证码'
          isCodeSending.value = false
        }
      }, 1000)
    }
  } catch (error) {
    ElMessage.error(error.response.data.msg || '发送失败')
    isCodeSending.value = false
  }
}

const handleRegister = async () => {
  try {
    await formRef.value.validate()
    isRegistering.value = true

    const registerData = {
      user: form.username,
      email: form.email,
      passwd: form.password,
      confirm: form.confirmPassword,
      varifycode: form.verificationCode
    }

    const response = await axios.post('gapi/user_register', registerData)

    if (response.data.error === 0) {
      ElMessage.success('注册成功')
      setTimeout(() => {
        window.location.href = '/login'
      }, 1500)
    }else if(response.data.error === 1003){
      ElMessage.error('验证码过期')
    }else if(response.data.error === 1004){
      ElMessage.error('验证码错误')
    }else if(response.data.error === 1005){
      ElMessage.error('用户或电子邮件已存在')
    }
  }
    catch (error) {
    ElMessage.error(error.response.data.msg || '注册失败')
  } finally {
    isRegistering.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}
.page-logo {
  position: absolute;
  top: 1rem;
  left: 1rem;
  width: 100px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
}

.login-container::before {
  content: '';
  position: absolute;
  width: 200%;
  height: 200%;
  background: linear-gradient(45deg, rgba(255,255,255,0.05) 25%,
  transparent 25%, transparent 50%,
  rgba(255,255,255,0.05) 50%,
  rgba(255,255,255,0.05) 75%,
  transparent 75%);
  background-size: 4px 4px;
  animation: animateBg 20s linear infinite;
  z-index: 0;
}

@keyframes animateBg {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.login-box {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 480px;
  padding: 2rem;
}

.form-wrapper {
  background: rgba(255, 255, 255, 0.97);
  border-radius: 20px;
  padding: 2.5rem 2rem;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.15);
  position: relative;
  backdrop-filter: blur(10px);
}

.decorative-line {
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
  background: linear-gradient(to bottom, #667eea, #764ba2);
}

.form-header {
  text-align: center;
  margin-bottom: 2rem;
}

.title-text {
  font-size: 1.8rem;
  font-weight: 600;
  color: #2d3748;
  letter-spacing: 1px;
  position: relative;
  display: inline-block;
}

.title-underline {
  position: relative;
  width: 120px;
  margin: 0.5rem auto;
}

.underline {
  height: 3px;
  background: linear-gradient(90deg, transparent, #667eea, transparent);
}

.circle {
  position: absolute;
  top: -3px;
  left: 50%;
  transform: translateX(-50%);
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #667eea;
  box-shadow: 0 0 8px rgba(102, 126, 234, 0.3);
}

.custom-input {
  margin-bottom: 1.5rem;
  transition: all 0.3s ease;
}

.custom-input:hover {
  transform: translateY(-2px);
}

.custom-input::v-deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.verification-container {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.verification-input {
  flex: 1;
}

.verification-btn {
  flex-shrink: 0;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.verification-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.register-btn {
  width: 100%;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  padding: 1rem;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.btn-text {
  position: relative;
  z-index: 2;
  letter-spacing: 1px;
}

.fill-container {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.2),
      transparent
  );
  transition: all 0.6s ease;
}

.register-btn:hover .fill-container {
  left: 100%;
}

.login-link {
  text-align: center;
  margin-top: 1.5rem;
  color: #718096;
  font-size: 0.9rem;
}

.link-text {
  color: #667eea;
  font-weight: 500;
  text-decoration: none;
  position: relative;
  transition: color 0.3s ease;
}

.link-text:hover {
  color: #764ba2;
}

.link-text::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: #667eea;
  transition: width 0.3s ease;
}

.link-text:hover::after {
  width: 100%;
}

.decorative-bubbles .bubble {
  position: absolute;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  animation: float 15s infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-100px) rotate(180deg); }
}

@media (max-width: 768px) {
  .login-box {
    padding: 1rem;
  }

  .form-wrapper {
    padding: 2rem 1.5rem;
  }

  .title-text {
    font-size: 1.5rem;
  }
}
</style>
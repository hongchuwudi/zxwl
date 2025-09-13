<template>
  <div class="login-container">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <div class="login-box">
      <div class="form-wrapper">
        <div class="decorative-line"></div>
        <el-form ref="formRef" :rules="data.rules" :model="data.form">
          <div class="form-header">
            <span class="title-text">密码重置</span>
            <div class="title-underline">
              <div class="underline"></div>
              <div class="circle"></div>
            </div>
          </div>

          <el-form-item prop="username">
            <el-input
                size="large"
                v-model="data.form.username"
                placeholder="用户名"
                prefix-icon="User"
                class="custom-input"
            ></el-input>
          </el-form-item>

          <el-form-item prop="email">
            <el-input
                size="large"
                v-model="data.form.email"
                placeholder="邮箱地址"
                prefix-icon="Message"
                class="custom-input"
            ></el-input>
          </el-form-item>

          <el-form-item prop="password">
            <el-input
                show-password
                size="large"
                v-model="data.form.password"
                placeholder="新密码"
                prefix-icon="Lock"
                class="custom-input"
            ></el-input>
          </el-form-item>

          <div class="verification-container">
            <el-form-item prop="verificationCode" class="verification-input">
              <el-input
                  size="large"
                  v-model="data.form.verificationCode"
                  placeholder="验证码"
                  prefix-icon="Eleme"
                  class="custom-input"
              ></el-input>
            </el-form-item>
            <el-button
                @click="getvar"
                size="large"
                class="verification-btn"
                :disabled="codeSent"
            >
              {{ codeButtonText }}
            </el-button>
          </div>

          <el-button
              @click="register"
              size="large"
              class="register-btn"
              :loading="isSubmitting"
          >
            <span class="btn-text">重置密码</span>
            <div class="fill-container"></div>
          </el-button>

          <div class="login-link">
            想起密码？请
            <router-link to="/Login" class="link-text">立即登录</router-link>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from "vue"
import { ElMessage } from "element-plus"
import { User, Lock, Message, Eleme } from "@element-plus/icons-vue"
import axios from 'axios'

const data = reactive({
  form: {},
  rules: {
    email: [{ required: true, message: '请输入邮箱地址', trigger: 'blur' }],
    password: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
    verificationCode: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
  }
})

const formRef = ref()
const codeSent = ref(false)
const codeButtonText = ref('获取验证码')
const isSubmitting = ref(false)


const getvar = async () => {
  try {
    const email = data.form.email

    // 添加邮箱验证
    if (!email) {
      ElMessage.warning('请输入邮箱地址')
      return
    }

    codeSent.value = true
    const response = await axios.post('gapi/get_varifycode', { email })

    if (response.data.code === 0 || response.data.meg === '验证码已发送') {
      ElMessage.success('验证码已发送')

      // 添加倒计时逻辑
      let countdown = 60
      codeButtonText.value = `${countdown}秒后重发`

      const timer = setInterval(() => {
        countdown--
        codeButtonText.value = `${countdown}秒后重发`

        if (countdown <= 0) {
          clearInterval(timer)
          codeButtonText.value = '获取验证码'
          codeSent.value = false
        }
      }, 1000)

    } else {
      ElMessage.error(response.data.msg || '发送失败')
      codeSent.value = false
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '发送失败')
    codeSent.value = false
  }
}

const register = async () => {
  try {
    await formRef.value.validate()
    isSubmitting.value = true

    const postData = {
      username: data.form.username,
      email: data.form.email,
      password: data.form.password,
      verify_code: data.form.verificationCode
    }

    const response = await axios.post('gapi/user/change-password', postData)

    if (response.data.error === 0) {
      ElMessage.success('密码重置成功')
      setTimeout(() => {
        window.location.href = '/Login'
      }, 1500)
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    isSubmitting.value = false
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
.page-logo {
  position: absolute;
  top: 1rem;
  left: 1rem;
  width: 100px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
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
}

.fill-container {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg,
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
}

.link-text {
  color: #667eea;
  font-weight: 500;
  text-decoration: none;
  position: relative;
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
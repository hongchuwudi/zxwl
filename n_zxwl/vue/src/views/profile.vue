<template>
  <div class="personal-center-container">

    <div class="back-button" @click="handleBack">
      <el-icon :size="28" class="back-icon">
        <Back />
      </el-icon>
    </div>

    <div class="profile-card glassmorphism">
      <div class="header">
        <h2 class="gradient-title">智选未来个人中心</h2>
      </div>

      <div class="avatar-section">
        <el-upload
            class="avatar-uploader"
            action="#"
            :show-file-list="false"
            :on-change="handleAvatarChange"
        >
          <div class="avatar-wrapper">
            <div class="hover-mask">
              <i class="el-icon-upload" />
            </div>
            <transition name="zoom">
              <img v-if="form.picture" :src="form.picture" class="avatar" />
            </transition>
            <div v-if="!form.picture" class="empty-avatar">
              <i class="el-icon-user" />
            </div>
          </div>
        </el-upload>
      </div>

      <el-form :model="form" label-width="100px" class="profile-form">
        <el-form-item label="用户名">
          <el-input
              v-model="form.name"
              placeholder="请输入用户名"
              class="elegant-input input-width"
          />
        </el-form-item>

        <el-form-item label="性 别">
          <el-radio-group v-model="form.sex">
            <el-radio :label="1" class="gender-radio">男</el-radio>
            <el-radio :label="2" class="gender-radio">女</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="毕业年份">
          <el-date-picker
              v-model="form.graduate"
              type="year"
              placeholder="选择年份"
              value-format="YYYY"
              class="year-picker input-width"
          />
        </el-form-item>

        <el-form-item label="所在地区">
          <el-cascader
              v-model="form.address"
              :options="cityOptions"
              placeholder="请选择省/市"
              class="region-selector"
          />
        </el-form-item>

        <el-form-item label="电子邮箱">
          <el-input
              v-model="form.email"
              disabled
              class="disabled-input input-width"
          />
        </el-form-item>
      </el-form>

      <el-button type="primary" @click="saveProfile" round class="floating-button">
        <span class="button-content">
          <i class="el-icon-check" />
          <span>save&nbsp&nbsp&nbsp </span>
        </span>
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { provinceAndCityData } from 'element-china-area-data'
import { Back } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

const form = ref({
  name: '',
  email: '',
  sex: 1,
  graduate: '',
  address: [],
  picture: ''
})

const cityOptions = provinceAndCityData

const router = useRouter()

const handleBack = () => {
  router.back()
}

// 获取个人资料
const fetchProfile = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail')
    if (!userEmail) {
      ElMessage.error('请先登录')
      return
    }

    const response = await axios.post('gapi/profile', {
      email: userEmail
    })

    if (response.data.code === 0) {
      form.value = {
        ...response.data.data,
        graduate: response.data.data.graduate ?
            String(response.data.data.graduate) : // 转换为字符串
            '',
        address: response.data.data.address ?
            response.data.data.address.split('/') : []
      }
    }
  } catch (error) {
    ElMessage.error('获取资料失败')
    console.error('Error fetching profile:', error)
  }
}

// 保存修改
const saveProfile = async () => {
  try {
    const userEmail = localStorage.getItem('userEmail')
    if (!userEmail) {
      ElMessage.error('请先登录')
      return
    }

    const payload = {
      ...form.value,
      email: userEmail,
      address: form.value.address.join('/'),
      graduate: form.value.graduate ? parseInt(form.value.graduate) : null

    }

    const response = await axios.post('gapi/profile/update', payload)

    if (response.data.code === 0) {
      const logData = {
        "email": localStorage.getItem('userEmail'),
        "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
        "operation": "用户更新资料"
      };
      const logResponse = await axios.post("gapi/log", logData, {
        headers: {
          "Content-Type": "application/json"
        }
      });
      ElMessage.success('资料更新成功')
      if (payload.email !== userEmail) {
        localStorage.setItem('userEmail', payload.email)
      }
      router.push('/Zxwl')
    }
    else {
      ElMessage.error('更新失败,请勿空修改')
    }
  } catch (error) {
    ElMessage.error('更新失败')
    console.error('Error updating profile:', error)
  }
}

// 头像修改处理
const handleAvatarChange = (file) => {
  const reader = new FileReader()
  reader.readAsDataURL(file.raw)
  reader.onload = () => {
    form.value.picture = reader.result
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped>

.back-button {
  position: absolute;
  top: 2rem;
  left: 2rem;
  cursor: pointer;
  z-index: 1000;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  padding: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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


.personal-center-container {
  min-height: 100vh;
  background: url(../assets/colleges.jpg) no-repeat center center fixed;
  background-size: cover;
  padding: 1rem 2rem;
  display: flex;
  justify-content: center;
  align-items: center;
}

.profile-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
  padding: 2.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
  transition: transform 0.3s ease;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.glassmorphism {
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.profile-card:hover {
  transform: translateY(-5px);
}

.profile-card::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: linear-gradient(45deg, transparent, rgba(255,255,255,0.1), transparent);
  transform: rotate(45deg);
  pointer-events: none;
}

.header {
  text-align: center;
  margin-bottom: 0rem;
}

.gradient-title {
  background: linear-gradient(45deg, #000000, #000000);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 2.2rem;
  letter-spacing: 1px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
}

.floating-button {
  background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
  border: none;
  padding: 1rem 2rem;
  font-weight: 600;
  letter-spacing: 1px;
  box-shadow: 0 4px 12px rgba(106, 17, 203, 0.3);
  transition: all 0.3s ease;
  align-self: center;
  margin-top: 0rem;
  display: flex; /* 添加这一行，使按钮内部使用flex布局 */
  justify-content: center; /* 添加这一行，使内部元素水平居中 */
  align-items: center; /* 添加这一行，使内部元素垂直居中 */
}


.floating-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(106, 17, 203, 0.4);
}

.button-content {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.avatar-section {
  text-align: center;
  margin-bottom: 2rem;
  position: relative;
}

.avatar-wrapper {
  display: inline-block;
  border-radius: 50%;
  position: relative;
  transition: transform 0.3s ease;
  cursor: pointer;
}

.avatar-wrapper:hover {
  transform: scale(1.05);
}

.hover-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.4);
  border-radius: 50%;
  opacity: 0;
  transition: opacity 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.5rem;
}

.avatar-wrapper:hover .hover-mask {
  opacity: 1;
}

.avatar {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 8px 24px rgba(0,0,0,0.15);
  border: 3px solid white;
}

.empty-avatar {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background: linear-gradient(45deg, #f3e7e9, #e3eeff);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  color: #6a11cb;
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}

.profile-form {
  padding: 0.8rem;
  padding-left: 4.5rem;
}

/* 深度样式覆盖 */
.el-form-item__label {
  font-weight: 600;
  color: #606266;
  letter-spacing: 0.5px;
  font-size: 1rem;
}

.el-input__inner {
  border-radius: 12px;
  transition: all 0.3s ease;
  border: 2px solid #e0e0e0;
  height: 48px;
  font-size: 1rem;
}

.el-input__inner:focus {
  border-color: #6a11cb;
  box-shadow: 0 0 8px rgba(106, 17, 203, 0.2);
}

.el-radio__inner {
  border-color: #6a11cb;
  width: 20px;
  height: 20px;
}

.el-radio__input.is-checked .el-radio__inner {
  background: #6a11cb;
  border-color: #6a11cb;
}

.el-radio__label {
  font-size: 1rem;
  color: #606266;
}

.el-date-editor {
  width: 100%;
}

.el-cascader {
  width: 100%;
}

/* 新增样式，设置统一宽度 */
.input-width {
  width: 69.5%;
}

/* 动画效果 */
.zoom-enter-active,
.zoom-leave-active {
  transition: transform 0.3s cubic-bezier(0.68, -0.55, 0.27, 1.55);
}

.zoom-enter-from,
.zoom-leave-to {
  transform: scale(0);
}

@media (max-width: 768px) {
  .profile-card {
    width: 95%;
    padding: 1.5rem;
  }

  .gradient-title {
    font-size: 1.8rem;
  }

  .header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }

  .avatar {
    width: 120px;
    height: 120px;
  }

  .empty-avatar {
    width: 120px;
    height: 120px;
    font-size: 2.5rem;
  }
}
</style>
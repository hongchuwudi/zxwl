<template>
  <div class="admin-policy-container">
    <!-- 头部导航栏 -->
    <nav class="header-nav">
      <img src="../../assets/zxwllogo.png" alt="Logo" class="page-logo">
      <div class="nav-content">
        <h1 class="logo">智选未来·政策信息发布系统</h1>
        <el-button class="back-btn" @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回上一页
        </el-button>
      </div>
    </nav>

    <div class="policy-content">
      <el-card class="policy-card">
        <el-form
            :model="policyForm"
            :rules="formRules"
            ref="formRef"
            label-position="top"
            class="policy-form"
        >
          <!-- 政策标题 -->
          <el-form-item label="政策标题" prop="title">
            <el-input
                v-model="policyForm.title"
                placeholder="请输入政策标题（5-100字符）"
                clearable
                class="input-glow"
            />
          </el-form-item>

          <!-- 政策前言 -->
          <el-form-item label="政策摘要" prop="foreword">
            <el-input
                v-model="policyForm.foreword"
                type="textarea"
                :rows="3"
                placeholder="请输入政策摘要（200字以内）"
                show-word-limit
                maxlength="200"
                class="input-gradient"
            />
          </el-form-item>

          <!-- 政策内容 -->
          <el-form-item label="政策正文" prop="content">
            <el-input
                v-model="policyForm.content"
                type="textarea"
                :rows="8"
                placeholder="请输入详细政策内容（至少50字符）"
                show-word-limit
                maxlength="2000"
                class="input-scroll"
            />
          </el-form-item>

          <!-- 提交按钮 -->
          <div class="action-container">
            <el-button
                type="primary"
                @click="handleSubmit"
                :loading="submitting"
                class="submit-btn"
            >
              <span v-if="!submitting">立即发布</span>
              <span v-else class="loading-text">提交中...</span>
            </el-button>

            <el-button
                @click="goBack"
                class="cancel-btn"
            >
              取消
            </el-button>
          </div>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const formRef = ref()
const submitting = ref(false)

const policyForm = reactive({
  title: '',
  content: '',
  foreword: ''
})

const formRules = reactive({
  title: [
    { required: true, message: '政策标题不能为空', trigger: 'blur' },
    { min: 5, max: 100, message: '标题长度需在5-100字符之间', trigger: 'blur' }
  ],
  foreword: [
    { required: true, message: '政策摘要不能为空', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '政策正文不能为空', trigger: 'blur' },
    { min: 50, message: '正文内容至少需要50字符', trigger: 'blur' }
  ]
})

// 返回上一页
const goBack = () => {
  router.go(-1)
}

const handleSubmit = async () => {
  try {
    // 表单验证
    await formRef.value.validate()

    submitting.value = true

    // API请求配置
    const response = await axios.post(
        '/gapi/policy/instert',
        {
          title: policyForm.title,
          content: policyForm.content,
          foreword: policyForm.foreword
        },
        {
          headers: {
            'Content-Type': 'application/json'
          },
          timeout: 10000 // 10秒超时
        }
    )

    // 处理响应
    if (response.data.code === 0) {
      const logData = {
        "email": localStorage.getItem('userEmail'),
        "date": new Date().toISOString().slice(0, 19).replace('T', ' '),
        "operation": "添加政策资料"
      };
      const logResponse = await axios.post("/gapi/log", logData, {
        headers: {
          "Content-Type": "application/json"
        }
      });
      ElMessage.success({
        message: `创建成功，政策ID：${response.data.id}`,
        duration: 3000,
        customClass: 'success-message'
      })
      formRef.value.resetFields()
    } else {
      ElMessage.error(response.data.msg || '提交失败，请检查数据格式')
    }
  } catch (error) {
    // 错误处理
    let errorMessage = '操作失败，请稍后重试'
    if (error.response) {
      errorMessage = `服务器错误：${error.response.data.msg || '未知错误'}`
    } else if (error.request) {
      errorMessage = '网络连接异常，请检查网络状态'
    } else if (error.message) {
      errorMessage = error.message
    }
    ElMessage.error({ message: errorMessage, duration: 2000 })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.admin-policy-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 0;
}

.header-nav {
  background: linear-gradient(135deg, #2c3e50, #3498db);
  box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
  min-height: 71px;
}

.nav-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-logo {
  position: absolute;
  top: 0;
  left: 1rem;
  width: 100px;
  height: auto;
  z-index: 3;
}

.logo {
  color: white;
  font-size: 1.5rem;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  margin: 0;
  flex-grow: 1;
  text-align: center;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border: 1px solid rgba(255, 255, 255, 0.5);
  color: white;
  background: rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: white;
  color: #2c3e50;
  transform: translateY(-2px);
}

.policy-content {
  max-width: 900px;
  margin: 2rem auto;
  padding: 0 1rem;
}

.policy-card {
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  border: none;
}

.policy-form {
  padding: 1.5rem;
}

::v-deep(.el-form-item__label) {
  font-weight: 600;
  color: #4a5568;
  font-size: 0.95rem;
  margin-bottom: 0.8rem;
}

.input-glow {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.input-glow:hover {
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.input-glow:focus-within {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.input-gradient {
  background: linear-gradient(145deg, #ffffff, #f8faff);
  border-radius: 8px;
}

.submit-btn {
  width: 180px;
  height: 44px;
  font-size: 1.1rem;
  background: linear-gradient(45deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  transition: all 0.3s ease;
  margin-right: 1rem;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(118, 75, 162, 0.25);
}

.cancel-btn {
  width: 120px;
  height: 44px;
  font-size: 1rem;
  border-radius: 8px;
}

.loading-text {
  letter-spacing: 1px;
  color: rgba(255, 255, 255, 0.9);
}

.action-container {
  text-align: center;
  margin-top: 2rem;
  display: flex;
  justify-content: center;
  gap: 1rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .nav-content {
    padding: 0 1rem;
    flex-direction: column;
    gap: 1rem;
  }

  .logo {
    font-size: 1.2rem;
    order: -1;
  }

  .policy-content {
    margin: 1rem auto;
    padding: 0 0.5rem;
  }

  .policy-form {
    padding: 1rem;
  }

  .action-container {
    flex-direction: column;
    align-items: center;
  }

  .submit-btn, .cancel-btn {
    width: 100%;
    margin-right: 0;
  }

  .page-logo {
    position: relative;
    margin: 0 auto;
    left: 0;
    top: 0;
    margin-bottom: 0.5rem;
  }
}
</style>
<template>
  <div class="policy-container">
    <h2 class="title">智选未来·政策信息发布系统</h2>
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
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

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
    // 错误处理（包含网页6的推荐方案）
    let errorMessage = '操作失败，请稍后重试'
    if (error.response) {
      errorMessage = `服务器错误：${error.response.data.msg || '未知错误'}`
    } else if (error.request) {
      errorMessage = '网络连接异常，请检查网络状态'
    }
    ElMessage.error({ message: errorMessage, duration: 2000 })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.policy-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2.5rem;
  background: linear-gradient(145deg, #f8faff 0%, #ffffff 100%);
  border-radius: 16px;
  box-shadow: 0 12px 24px -6px rgba(99, 102, 241, 0.1),
  0 0 15px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(160, 174, 192, 0.1);
}

.title {
  text-align: center;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 2rem;
  margin-bottom: 2rem;
  letter-spacing: 1px;
}

.policy-form {
  padding: 1.5rem;
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #4a5568;
  font-size: 0.95rem;
  margin-bottom: 0.8rem;
}

.input-glow {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

  &:hover {
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
  }

  &:focus-within {
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
  }
}

.input-gradient {
  background: linear-gradient(145deg, #ffffff, #f8faff);
  border-radius: 8px;
}

.input-scroll {
  &::after {
    content: '';
    position: absolute;
    right: 0;
    bottom: 0;
    width: 100%;
    height: 30px;
    background: linear-gradient(transparent, #f8faff);
    pointer-events: none;
  }
}

.submit-btn {
  width: 180px;
  height: 44px;
  font-size: 1.1rem;
  background: linear-gradient(45deg, #667eea, #764ba2);
  border: none;
  border-radius: 8px;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(118, 75, 162, 0.25);
  }
}

.loading-text {
  letter-spacing: 1px;
  color: rgba(255, 255, 255, 0.9);
}

.action-container {
  text-align: center;
  margin-top: 2rem;

}

/* 响应式设计 */
@media (max-width: 768px) {
  .policy-container {
    margin: 1rem;
    padding: 1.5rem;
  }

  .title {
    font-size: 1.5rem;
  }

  .submit-btn {
    width: 100%;
  }
}
</style>
<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const props = defineProps({
  modelValue: Boolean,
  user: Object
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const loading = ref(false)
const form = ref({
  username: '',
  email: '',
  displayName: '',
  gender: 0,
  birthYear: '',
  location: '',
  bio: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

watch(() => props.user, (newUser) => {
  if (newUser) {
    form.value = {
      username: newUser.username || '',
      email: newUser.email || '',
      displayName: newUser.displayName || '',
      gender: newUser.gender || 0,
      birthYear: newUser.birthYear || '',
      location: newUser.location || '',
      bio: newUser.bio || ''
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  try {
    loading.value = true

    // 转换字段格式以匹配后端期望
    const requestData = {
      username: form.value.username,
      email: form.value.email,
      displayName: form.value.displayName,  // 改为蛇形命名
      gender: Number(form.value.gender),     // 确保是数字
      birthYear: form.value.birthYear ? Number(form.value.birthYear) : 0, // 转为数字
      location: form.value.location,
      bio: form.value.bio
    }

    console.log('转换后的请求数据:', requestData)

    const response = await axios.put(`/gapi/user/update/${props.user.id}`, requestData, {
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (response.data && response.data.error === 0) {
      ElMessage.success(response.data.message || '用户信息更新成功')
      dialogVisible.value = false
      emit('success')
    } else {
      ElMessage.error(response.data?.message || '更新失败')
    }
  } catch (error) {
    // 详细的错误信息
    if (error.response) {
      // 服务器返回了错误状态码
      console.error('服务器错误:', error.response.status, error.response.data)
      ElMessage.error(`服务器错误: ${error.response.data?.message || error.response.status}`)
    } else if (error.request) {
      // 请求已发出但没有收到响应
      console.error('网络错误:', error.request)
      ElMessage.error('网络连接错误，请检查网络设置')
    } else {
      // 其他错误
      console.error('请求配置错误:', error.message)
      ElMessage.error(`请求错误: ${error.message}`)
    }
  } finally {
    loading.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
}
</script>

<template>
  <el-dialog
      v-model="dialogVisible"
      :title="`编辑用户 - ${user?.username || '未知用户'}`"
      width="500px"
      destroy-on-close
      @closed="handleClose"
  >
    <el-form
        v-loading="loading"
        :model="form"
        :rules="rules"
        label-width="80px"
        label-position="left"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            disabled=""
        />
      </el-form-item>

      <el-form-item label="邮箱" prop="email">
        <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            type="email"
            disabled=""
        />
      </el-form-item>

      <el-form-item label="显示名称">
        <el-input
            v-model="form.displayName"
            placeholder="请输入显示名称"
        />
      </el-form-item>

      <el-form-item label="性别">
        <el-radio-group v-model="form.gender">
          <el-radio :label="3" disabled="">未知</el-radio>
          <el-radio :label="1">男</el-radio>
          <el-radio :label="2">女</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="出生年份">
        <el-input
            v-model="form.birthYear"
            placeholder="请输入出生年份"
            type="number"
            min="1900"
            :max="new Date().getFullYear()"
        />
      </el-form-item>

      <el-form-item label="位置">
        <el-input
            v-model="form.location"
            placeholder="请输入位置信息"
        />
      </el-form-item>

      <el-form-item label="个人简介">
        <el-input
            v-model="form.bio"
            type="textarea"
            :rows="3"
            placeholder="请输入个人简介"
            maxlength="200"
            show-word-limit
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
::v-deep(.el-dialog__body) {
  padding: 20px;
}

::V-deep(.el-form-item) {
  margin-bottom: 16px;
}
</style>
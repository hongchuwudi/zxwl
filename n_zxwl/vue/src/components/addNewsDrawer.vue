<template>
  <el-drawer
      v-model="visible"
      title="发布资讯文章"
      direction="rtl"
      size="74%"
      class="news-drawer"
      :before-close="handleClose"
  >
    <div class="drawer-content">
      <el-form
          :model="formData"
          :rules="rules"
          ref="formRef"
          label-width="100px"
          class="news-form"
      >
        <!-- 文章标题 -->
        <el-form-item label="文章标题" prop="title">
          <el-input
              v-model="formData.title"
              placeholder="请输入文章标题"
              clearable
              maxlength="100"
              show-word-limit
          />
        </el-form-item>

        <!-- 关键词 -->
        <el-form-item label="关键词" prop="keywords">
          <el-input
              v-model="formData.keywords"
              placeholder="请输入关键词，多个关键词用逗号分隔"
              clearable
          />
        </el-form-item>

        <!-- 省份选择 -->
        <el-form-item label="所属省份" prop="province_id">
          <el-select
              v-model="formData.province_id"
              placeholder="请选择省份"
              clearable
              filterable
              style="width: 100%"
          >
            <el-option
                v-for="province in provinces"
                :key="province.value"
                :label="province.label"
                :value="province.value"
            />
          </el-select>
        </el-form-item>

        <!-- 封面图片 -->
        <el-form-item label="封面图片" prop="style_url">
          <el-upload
              class="cover-uploader"
              action="#"
              :auto-upload="false"
              :show-file-list="false"
              :on-change="handleCoverChange"
              accept="image/*"
          >
            <img v-if="formData.style_url" :src="formData.style_url" class="cover-image" />
            <el-icon v-else class="cover-uploader-icon"><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">点击上传封面图片，建议尺寸 800×400px</div>
        </el-form-item>

        <!-- 文章描述 -->
        <el-form-item label="文章描述" prop="description">
          <el-input
              v-model="formData.description"
              type="textarea"
              :rows="3"
              placeholder="请输入文章简要描述"
              maxlength="500"
              show-word-limit
          />
        </el-form-item>

        <!-- 文章内容 - 富文本编辑器 -->
        <el-form-item label="文章内容" prop="content">
          <div class="editor-container">
            <el-alert
                title="请编写完整的文章内容（支持HTML格式:前往https://bi.cool/project编辑）"
                type="info"
                :closable="false"
                style="margin-bottom: 10px"
            >
              <template #description>
                <span>如需使用高级HTML编辑功能，请</span>
                <el-link
                    type="primary"
                    href="https://bi.cool/project"
                    target="_blank"
                    style="margin: 0 4px"
                >
                  前往编辑平台
                </el-link>
                <span>进行编辑</span>
              </template>
            </el-alert>
            <div class="editor-wrapper">
      <textarea
          v-model="formData.content"
          class="rich-text-editor"
          placeholder="请输入文章详细内容..."
          rows="15"
      ></textarea>
            </div>
          </div>
        </el-form-item>
      </el-form>

      <el-divider />

      <div class="drawer-actions">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          发布文章
        </el-button>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, reactive, watch, computed ,onMounted} from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {useUserStore} from "@/utils/auth.js";
import { Plus } from '@element-plus/icons-vue'
import {useRoute, useRouter} from 'vue-router'
const router = useRouter()
const route = useRoute()

const { userName, userEmail, getUser ,checkLoginStatus} = useUserStore();

// 定义组件属性
const props = defineProps({
  modelValue: Boolean
})
// 定义事件
const emit = defineEmits(['update:modelValue', 'submit'])

// 响应式数据
const visible = computed({
  get: () => props.modelValue,
  set: value => emit('update:modelValue', value)
})

const formRef = ref(null)
const submitting = ref(false)

// 表单数据
// 表单数据
const formData = reactive({
  title: '',
  keywords: '',
  province_id: '',
  style_url: '',
  description: '',
  content: '',
  video_detail: '',
  video_type: '',
  video_img: '',
  news_num: '',
  is_push: 2,
  is_top: 1,
  style_type: '',
  card_school_id: '',
  card_live_id: '',
  class_name: ''
})



// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' },
    { min: 5, message: '标题长度不能少于5个字符', trigger: 'blur' }
  ],
  keywords: [
    { required: true, message: '请输入关键词', trigger: 'blur' }
  ],
  province_id: [
    { required: true, message: '请选择所属省份', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入文章描述', trigger: 'blur' },
    { min: 10, message: '描述长度不能少于10个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' },
    { min: 50, message: '内容长度不能少于50个字符', trigger: 'blur' }
  ]
}

// 省份数据
const provinces = [
  { label: '北京', value: '11' }, { label: '天津', value: '12' }, { label: '河北', value: '13' },
  { label: '山西', value: '14' }, { label: '内蒙古', value: '15' }, { label: '辽宁', value: '21' },
  { label: '吉林', value: '22' }, { label: '黑龙江', value: '23' }, { label: '上海', value: '31' },
  { label: '江苏', value: '32' }, { label: '浙江', value: '33' }, { label: '安徽', value: '34' },
  { label: '福建', value: '35' }, { label: '江西', value: '36' }, { label: '山东', value: '37' },
  { label: '河南', value: '41' }, { label: '湖北', value: '42' }, { label: '湖南', value: '43' },
  { label: '广东', value: '44' }, { label: '广西', value: '45' }, { label: '海南', value: '46' },
  { label: '重庆', value: '50' }, { label: '四川', value: '51' }, { label: '贵州', value: '52' },
  { label: '云南', value: '53' }, { label: '西藏', value: '54' }, { label: '陕西', value: '61' },
  { label: '甘肃', value: '62' }, { label: '青海', value: '63' }, { label: '宁夏', value: '64' },
  { label: '新疆', value: '65' }, { label: '台湾', value: '71' }, { label: '香港', value: '81' },
  { label: '澳门', value: '82' }
]

// 封面图片上传处理
const handleCoverChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    formData.style_url = e.target.result
  }
  reader.readAsDataURL(file.raw)
}

// 关闭抽屉
const handleClose = () => {
  if (formData.title || formData.content) {
    ElMessageBox.confirm('确定要取消发布吗？已填写的内容将不会保存', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      resetForm()
      visible.value = false
    }).catch(() => {})
  } else {
    visible.value = false
  }
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(formData, {
    title: '',
    keywords: '',
    province_id: '',
    style_url: '',
    description: '',
    content: ''
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    const valid = await formRef.value.validate()
    if (!valid) return

    submitting.value = true

    // 准备提交数据
    const submitData = {
      ...formData,
      publisher_email: userEmail.value , // 发布者邮箱
      from_source: userName.value, // 来源字段先空着
      // 其他后端需要的字段...
    }

    // 触发提交事件
    emit('submit', submitData)

    // 重置表单并关闭抽屉
    resetForm()
    visible.value = false

    // ElMessage.success('文章发布成功！')

  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    submitting.value = false
  }
}

// 监听抽屉显示状态
watch(visible, (newVal) => {
  if (!newVal) {
    resetForm()
  }
})

onMounted(() => {
  // 渲染前加载用户身份
 getUser()
  console.log('用户身份：', userName.value)
  console.log('用户邮箱：', userEmail.value)
  if (!checkLoginStatus()) {
    // 可以在这里执行重定向或其他操作
    ElMessage.error('请先登录！')
    router.push('/login')
  }
})
</script>

<style scoped>
.news-drawer {
  z-index: 2000;
}

.drawer-content {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.news-form {
  flex: 1;
  overflow-y: auto;
  padding-right: 10px;
}

.cover-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 200px;
  height: 120px;
}

.cover-uploader:hover {
  border-color: #409EFF;
}

.cover-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 200px;
  height: 120px;
  line-height: 120px;
  text-align: center;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.editor-container {
  width: 100%;
}

.editor-wrapper {
  border: 1px solid #DCDFE6;
  border-radius: 4px;
  padding: 5px;
}

.rich-text-editor {
  width: 100%;
  border: none;
  outline: none;
  resize: vertical;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.5;
  padding: 8px;
  box-sizing: border-box;
}

.rich-text-editor:focus {
  border-color: #409EFF;
}

.drawer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #EBEEF5;
}

/* 滚动条样式 */
.news-form::-webkit-scrollbar {
  width: 6px;
}

.news-form::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.news-form::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 3px;
}

.news-form::-webkit-scrollbar-thumb:hover {
  background: #909399;
}
</style>
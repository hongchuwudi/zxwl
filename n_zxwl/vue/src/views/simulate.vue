<template>
  <div class="container">
    <img src="../assets/zxwllogo.png" alt="Logo" class="page-logo">
    <!-- 新增返回按钮 -->
    <div class="back-button" @click="handleBack">
      <el-icon :size="28" class="back-icon">
        <Back />
      </el-icon>
    </div>
    <header class="header">
      <div class="header-content">
        <h1 class="title">智选未来 · 高考志愿模拟填报</h1>
      </div>
    </header>

    <main class="main-content">
      <!-- 查看状态提示 -->
      <div v-if="isViewingFamily" class="viewing-notice">
        <span class="notice-part">正在查看</span>
        <span class="notice-part">{{ viewedEmail }}</span>
        <span class="notice-part">的志愿</span>
        <span class="notice-part">(只读模式)</span>
        <el-button @click="returnToSelf" type="primary" size="small">
          返回我的填报
        </el-button>
      </div>

      <!-- 左侧批次选择 -->
      <div class="batch-section">
        <div class="section-title">填报批次</div>
        <el-scrollbar height="calc(100vh - 280px)">
          <div
              v-for="(batch, index) in batches"
              :key="index"
              class="batch-item"
              :class="{ 'active-batch': activeBatch === index }"
              @click="switchBatch(index)"
          >
            <div class="batch-order">{{ index + 1 }}</div>
            <div class="batch-info">
              <h3>{{ batch.name }}</h3>
              <p>最多填报 {{ batch.max }} 个志愿</p>
            </div>
          </div>
        </el-scrollbar>
      </div>

      <!-- 中间志愿填报区 -->
      <div class="volunteer-section">
        <div class="volunteer-table">
          <div class="table-header">
            <div class="col-order">序号</div>
            <div class="col-school">院校信息</div>
            <div class="col-majors">专业志愿（1 - 6）</div>
            <div class="col-actions">操作</div>
          </div>

          <div class="volunteer-list-container">
            <draggable
                v-model="currentVolunteers"
                item-key="id"
                handle=".drag-handle"
                ghost-class="ghost-item"
                @end="onDragEnd"
                :disabled="isViewingFamily"
            >
              <template #item="{ element, index }">
                <div class="volunteer-item" :class="{ 'is-disabled': isViewingFamily }">
                  <div class="col-order">
                    <el-icon class="drag-handle"><Rank /></el-icon>
                    <span class="order-number">{{ index + 1 }}</span>
                  </div>

                  <div class="col-school">
                    <el-input
                        v-model="element.schoolName"
                        :disabled="isViewingFamily"
                        placeholder="请输入院校名称"
                        class="school-input"
                        clearable
                    />
                  </div>

                  <div class="col-majors">
                    <div class="major-grid">
                      <div
                          v-for="(major, idx) in element.majors"
                          :key="idx"
                          class="major-item"
                      >
                        <div class="major-label">专业{{ idx + 1 }}</div>
                        <el-input
                            v-model="element.majors[idx]"
                            :disabled="isViewingFamily"
                            :placeholder="`请输入第 ${idx + 1} 专业志愿`"
                            class="major-input"
                            clearable
                        />
                      </div>
                    </div>
                  </div>

                  <div class="col-actions">
                    <el-button
                        type="danger"
                        :disabled="isViewingFamily"
                        :icon="Delete"
                        circle
                        plain
                        @click="removeVolunteer(index)"
                    />
                  </div>
                </div>
              </template>
            </draggable>
          </div>

          <div class="add-section">
            <el-button
                type="primary"
                :icon="Plus"
                class="add-button"
                :disabled="!canAddMore || isViewingFamily"
                @click="addVolunteer"
            >
              添加志愿（{{ remainingSlots }}个剩余）
            </el-button>
            <div class="add-tips">
              当前批次最多可填报 {{ currentBatch.max }} 个志愿
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧信息栏 -->
      <div class="info-sidebar">
        <!-- 统计卡片 -->
        <div class="stats-card">
          <h3 class="stats-title">
            <el-icon><DataAnalysis /></el-icon>
            填报统计
          </h3>
          <div class="stats-content">
            <div class="stat-item">
              <label>已选院校：</label>
              <span class="stat-value">{{ selectedSchoolCount }}</span>
            </div>
            <div class="stat-item">
              <label>已填专业：</label>
              <span class="stat-value">{{ selectedMajorCount }}</span>
            </div>
            <div class="stat-item">
              <label>剩余名额：</label>
              <span class="stat-value">{{ remainingSlots }}</span>
            </div>
          </div>
        </div>

        <!-- 家庭共享卡片 -->
        <div class="stats-card family-sharing">
          <h3 class="stats-title">
            <el-icon><User /></el-icon>
            家庭共享
          </h3>
          <div class="stats-content">
            <div class="member-list">
              <div v-for="(email, index) in familyEmails" :key="index" class="member-item">
                <div class="member-info">
                  <span class="member-email" @click="viewFamilyVolunteer(email)">{{ email }}</span>
                  <span class="member-status">已绑定</span>
                </div>
                <el-button
                    type="danger"
                    size="small"
                    :icon="Delete"
                    circle
                    plain
                    @click="removeMember(email)"
                />
              </div>
            </div>

            <el-button
                type="primary"
                class="add-member-btn"
                @click="showAddDialog = true"
            >
              <el-icon><Plus /></el-icon>
              添加家庭成员
            </el-button>
          </div>
        </div>
      </div>
    </main>

    <!-- 底部操作按钮 -->
    <footer class="footer">
      <div class="footer-actions">
        <el-button
            type="primary"
            size="large"
            class="save-btn"
            :disabled="isViewingFamily"
            @click="saveDraft"
        >
          <el-icon><Folder /></el-icon>
          保存草稿
        </el-button>
        <el-button type="success" size="large" class="submit-btn" @click="submitForm">
          <el-icon><Promotion /></el-icon>
          导出文件
        </el-button>
        <el-button
            type="info"
            size="large"
            :disabled="isViewingFamily"
            @click="resetForm"
        >
          <el-icon><Refresh /></el-icon>
          重置填报
        </el-button>
      </div>
    </footer>

    <!-- 添加成员对话框 -->
    <el-dialog v-model="showAddDialog" title="添加家庭成员" width="400px">
      <el-form :model="newMember" :rules="rules" ref="memberForm">
        <el-form-item label="邮箱地址" prop="email">
          <el-input v-model="newMember.email" placeholder="请输入成员邮箱" />
        </el-form-item>
        <el-form-item label="共享密码" prop="passwd">
          <el-input
              v-model="newMember.passwd"
              type="password"
              placeholder="请输入共享密码"
              show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addMember">确认添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import draggable from 'vuedraggable'
import { ElMessage, ElLoading } from 'element-plus'
import axios from 'axios'
import {
  User,
  Rank,
  Delete,
  Plus,
  DataAnalysis,
  Folder,
  Promotion,
  Refresh,
  Back
} from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import * as XLSX from 'xlsx'




const router = useRouter()

// 状态变量
const batches = reactive([
  { name: '本科提前批', max: 3, type: 0 },
  { name: '本科一批', max: 6, type: 1 },
  { name: '本科二批', max: 6, type: 2 }
])
const activeBatch = ref(0)
const currentVolunteers = ref([])
const userEmail = ref('')
const isLoading = ref(false)
const familyEmails = ref([])
const showAddDialog = ref(false)
const newMember = reactive({ email: '', passwd: '' })
const isViewingFamily = ref(false)
const viewedEmail = ref('')

// 验证规则
const rules = reactive({
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: ['blur', 'change'] }
  ],
  passwd: [
    { required: true, message: '请输入共享密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ]
})

// 计算属性
const currentBatch = computed(() => batches[activeBatch.value])
const remainingSlots = computed(() => currentBatch.value.max - currentVolunteers.value.length)
const canAddMore = computed(() => remainingSlots.value > 0 && !isViewingFamily.value)
const selectedSchoolCount = computed(() => currentVolunteers.value.filter(v => v.schoolName.trim()).length)
const selectedMajorCount = computed(() =>
    currentVolunteers.value.reduce((sum, v) => sum + v.majors.filter(m => m.trim()).length, 0)
)

// 生命周期
onMounted(async () => {
  const email = localStorage.getItem('userEmail')
  if (email) {
    userEmail.value = email
    await fetchVolunteers(0)
    await loadFamilyMembers()
  } else {
    ElMessage.warning('未检测到用户登录信息')
    window.location.href = '/login'
  }
})

// 方法
const loadFamilyMembers = async () => {
  try {
    const response = await axios.post('gapi/family/find', { myemail: userEmail.value })
    familyEmails.value = response.data.code === 0 ? response.data.emails : []
  } catch (error) {
    ElMessage.error('获取家庭成员失败')
  }
}

const fetchVolunteers = async (batchType, email = userEmail.value) => {
  const loading = ElLoading.service({ lock: true, text: '加载数据中...', background: 'rgba(0, 0, 0, 0.7)' })
  try {
    const response = await axios.post('gapi/volunteer/fetch', {
      user_email: email,
      batch_type: batchType
    }, { timeout: 10000 })

    if (response.data.status === 0) {
      currentVolunteers.value = response.data.records
          .sort((a, b) => a.sequence - b.sequence)
          .map(record => ({
            id: `${Date.now()}_${record.sequence}`,
            schoolName: record.institution,
            majors: Array(6).fill().map((_, i) => record[`${['first', 'second', 'third', 'fourth', 'fifth', 'sixth'][i]}_major`] || '')
          }))
    } else {
      currentVolunteers.value = []
    }
  } catch (error) {
    ElMessage.error('数据加载失败')
  } finally {
    loading.close()
  }
}

const switchBatch = async (index) => {
  if (index === activeBatch.value) return
  activeBatch.value = index
  const email = isViewingFamily.value ? viewedEmail.value : userEmail.value
  await fetchVolunteers(batches[index].type, email)
}

const onDragEnd = () => {
  currentVolunteers.value = currentVolunteers.value.map((item, index) => ({ ...item, order: index + 1 }))
}

const addVolunteer = () => {
  if (canAddMore.value) {
    currentVolunteers.value.push({ id: `new_${Date.now()}`, schoolName: '', majors: Array(6).fill('') })
  }
}

const removeVolunteer = (index) => currentVolunteers.value.splice(index, 1)
const resetForm = () => currentVolunteers.value = []

const saveDraft = async () => {
  if (!validateForm()) return
  const loading = ElLoading.service({ lock: true })
  try {
    const payload = currentVolunteers.value
        .filter(item => item.schoolName.trim())
        .map((item, index) => ({
          college: item.schoolName.trim(),
          benke: currentBatch.value.type,
          yitianzhuanye: index,
          ...Object.fromEntries(item.majors.map((m, i) => [`major${i + 1}`, m.trim()])),
          email: userEmail.value
        }))

    const response = await axios.post('gapi/professional/upsert', payload)
    response.data.code === 0 ? ElMessage.success('保存成功') : ElMessage.error(response.data.message)
    await fetchVolunteers(currentBatch.value.type)
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    loading.close()
  }
}

const validateForm = () => {
  if (!currentVolunteers.value.length) {
    ElMessage.warning('请添加至少一个志愿')
    return false
  }
  return !currentVolunteers.value.some(v => !v.schoolName.trim())
}


const addMember = async () => {
    const authResponse = await axios.post('gapi/auth', {Email: newMember.email,
        Passwd: newMember.passwd})
    if (authResponse.data.code !== 0) {
      showAddDialog.value = false
      ElMessage.error('共享密码错误')
    }else{
    const res = await axios.post('gapi/family/add', {
      myemail: userEmail.value,
      familyemail: newMember.email
    })
    if (res.data.code === 0) {
      familyEmails.value.push(newMember.email)
      showAddDialog.value = false
      Object.assign(newMember, { email: '', password: '' })
      ElMessage.success('添加成功')
    }}
}

const removeMember = async (email) => {
  try {
    await axios.post('gapi/family/remove', { myemail: userEmail.value, familyemail: email })
    familyEmails.value = familyEmails.value.filter(e => e !== email)
    ElMessage.success('移除成功')
  } catch (error) {
    ElMessage.error('移除失败')
  }
}

const viewFamilyVolunteer = async (email) => {
  viewedEmail.value = email
  isViewingFamily.value = true
  await fetchVolunteers(currentBatch.value.type, email)
}

const returnToSelf = async () => {
  isViewingFamily.value = false
  viewedEmail.value = ''
  await fetchVolunteers(currentBatch.value.type)
}
// 在现有代码中添加导出方法
const exportToFile = () => {
  if (currentVolunteers.value.length === 0) {
    ElMessage.warning('没有可导出的数据')
    return
  }

  // 准备数据
  const data = currentVolunteers.value.map((vol, index) => {
    const majors = vol.majors.reduce((acc, major, idx) => {
      acc[`专业${idx + 1}`] = major
      return acc
    }, {})

    return {
      序号: index + 1,
      批次: currentBatch.value.name,
      院校名称: vol.schoolName,
      ...majors
    }
  })

  // 创建工作簿
  const worksheet = XLSX.utils.json_to_sheet(data)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, '志愿填报')

  // 生成文件名
  const fileName = `志愿填报_${currentBatch.value.name}_${
      isViewingFamily.value ? viewedEmail.value + '_' : ''
  }${new Date().toLocaleDateString().replace(/\//g, '-')}.xlsx`

  // 导出文件
  XLSX.writeFile(workbook, fileName)
  ElMessage.success('文件导出成功')
}
// 修改提交方法
const submitForm = async () => {
  if (validateForm()) {
    await saveDraft()
    exportToFile()
  }
}
// 处理返回按钮点击事件
const handleBack = () => {
  router.push('/zxwl')
}

</script>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", Arial, sans-serif;
}

.header {
  background: linear-gradient(135deg, #1a73e8, #0d47a1);
  color: white;
  padding: 0 2rem;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 1rem 0;
}

.title {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 500;
  letter-spacing: 1px;
  text-align: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: rgba(255,255,255,0.1);
  padding: 0.8rem 1.2rem;
  border-radius: 6px;
}

.user-icon {
  font-size: 1.8rem;
}

.user-details p {
  margin: 0;
  line-height: 1.4;
}

.exam-number {
  font-size: 0.9rem;
  opacity: 0.8;
}

.main-content {
  flex: 1;
  display: flex;
  gap: 1.5rem;
  padding: 1.5rem;
  overflow: hidden;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.batch-section {
  width: 260px;
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.section-title {
  font-size: 1.2rem;
  font-weight: 500;
  margin-bottom: 1rem;
  color: #2c3e50;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #eee;
}

.batch-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  margin-bottom: 0.5rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid transparent;
}

.batch-item:hover {
  background: #f8f9fa;
}

.active-batch {
  background: #e8f4ff !important;
  border-color: #1a73e8 !important;
}

.batch-order {
  width: 32px;
  height: 32px;
  background: #1a73e8;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  flex-shrink: 0;
}

.batch-info h3 {
  margin: 0;
  font-size: 1rem;
  color: #303133;
}

.batch-info p {
  margin: 4px 0 0;
  font-size: 0.9rem;
  color: #909399;
}

.volunteer-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.volunteer-table {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  overflow: hidden;
}

.table-header {
  display: flex;
  background: #f8f9fa;
  padding: 1rem;
  font-weight: 500;
  border-radius: 6px;
  margin-bottom: 0.5rem;
}

.col-order { width: 8%; }
.col-school { width: 25%; }
.col-majors { width: 55%; }
.col-actions { width: 12%; text-align: right; }

.volunteer-list-container {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.volunteer-item {
  display: flex;
  align-items: center;
  padding: 1rem;
  margin-bottom: 8px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  transition: all 0.3s;
  background: white;
}

.volunteer-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

.drag-handle {
  cursor: move;
  color: #909399;
  margin-right: 0.5rem;
}

.order-number {
  font-weight: 500;
  color: #606266;
}

.school-input {
  width: 100%;
}

.major-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.8rem;
}

.major-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.major-label {
  min-width: 50px;
  color: #606266;
  font-size: 0.9rem;
}

.major-input {
  flex: 1;
}

.add-section {
  margin-top: 1rem;
  text-align: center;
  padding: 1rem;
  border-top: 1px solid #eee;
}

.add-button {
  padding: 1rem 2rem;
  font-size: 1rem;
  width: 100%;
}

.add-tips {
  color: #909399;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.info-sidebar {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.stats-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stats-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.stats-content {
  padding: 0 1rem;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  margin: 0.8rem 0;
  font-size: 0.95rem;
}

.stat-value {
  color: #1a73e8;
  font-weight: 500;
}

.family-sharing {
  margin-top: 20px;
  right: 10px;
  top: 200px;
  width: 272px;
  height: 500px;
  overflow-y: auto;
}

.member-list {
  margin-bottom: 15px;
}

.member-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  margin: 8px 0;
  background: #f8f9fa;
  border-radius: 4px;
}

.member-info {
  flex: 1;
  margin-right: 15px;
}

.page-logo {
  position: absolute;
  top: -0.5rem;
  left: 1rem;
  width: 100px; /* 可按需调整 logo 大小 */
  height: auto;
  z-index: 3;
}

.member-email {
  display: block;
  color: #606266;
  font-size: 0.9rem;
}

.member-status {
  font-size: 0.8rem;
  color: #67C23A;
}

.add-member-btn {
  width: 100%;
  margin-top: 10px;
}

.footer {
  background: white;
  padding: 1rem 2rem;
  box-shadow: 0 -2px 12px rgba(0,0,0,0.08);
}

.footer-actions {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.save-btn,
.submit-btn {
  padding: 1rem 2.5rem;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .main-content {
    flex-wrap: wrap;
    height: auto;
  }

  .batch-section {
    width: 100%;
    margin-bottom: 1rem;
  }

  .volunteer-section {
    order: 2;
    width: 100%;
  }

  .info-sidebar {
    order: 3;
    width: 100%;
    flex-direction: row;
  }
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: flex-start;
  }

  .table-header {
    display: none;
  }

  .volunteer-item {
    flex-wrap: wrap;
    position: relative;
    padding: 1.5rem;
  }

  .col-order {
    position: absolute;
    left: 10px;
    top: 10px;
  }

  .col-school {
    width: 100%;
    margin-bottom: 1rem;
  }

  .col-majors {
    width: 100%;
  }

  .col-actions {
    position: absolute;
    right: 10px;
    top: 10px;
  }

  .major-grid {
    grid-template-columns: 1fr;
  }
}

/* 动画 */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.volunteer-item {
  animation: slideIn 0.3s ease;
}

.ghost-item {
  opacity: 0.5;
  background: #f8f9fa;
}

.sortable-chosen {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.viewing-notice {
  background: #f0f9eb;
  color: #67c23a;
  padding: 8px;
  margin: -8px 0 16px;
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 1000;
  font-size: 14px;
  text-align: center;
}

.notice-part {
  width: 100%;
  word-wrap: break-word;
}

.is-disabled .el-input__inner {
  background-color: #f5f7fa !important;
  color: #909399 !important;
  cursor: not-allowed;
}

.member-email {
  cursor: pointer;
  text-decoration: underline;
  color: #1a73e8 !important;
  transition: color 0.3s;
}

.member-email:hover {
  color: #0d47a1 !important;
}

.volunteer-item.is-disabled {
  background: #f8f9fa;
  opacity: 0.8;
}

.back-button {
  position: absolute;
  top: 0.75rem;
  right: 4rem;
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
</style>
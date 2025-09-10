<template>
  <div class="log-container">
    <div class="header">
      <h1 class="title">智选未来 - 操作日志</h1>
      <el-button
          type="primary"
          icon="el-icon-back"
          @click="$router.push('/zxwl')"
          class="back-btn"
      >
        返回首页
      </el-button>
    </div>

    <el-card class="box-card">
      <el-table
          :data="tableData"
          style="width: 100%"
          :max-height="600"
          v-loading="loading"
      >
        <el-table-column prop="email" label="用户邮箱" width="220" />
        <el-table-column prop="date" label="操作时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.date) }}
          </template>
        </el-table-column>
        <el-table-column prop="operation" label="操作内容" />
      </el-table>

      <el-pagination
          background
          layout="prev, pager, next"
          :total="total"
          :page-size="pageSize"
          v-model:current-page="currentPage"
          class="pagination"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import dayjs from 'dayjs'

const logs = ref([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const tableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return logs.value.slice(start, end)
})

const formatDate = (dateString) => {
  return dayjs(dateString).format('YYYY-MM-DD HH:mm:ss')
}

const fetchData = async () => {
  try {
    loading.value = true
    const response = await axios.get('http://127.0.0.1:8792/logs')
    if (response.data.code === 0) {
      logs.value = response.data.data
      total.value = response.data.total
    }
  } catch (error) {
    console.error('Error fetching logs:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.log-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
}

.header {
  display: flex;
  justify-content: center;
  align-items: center;
  margin-bottom: 20px;
  border-bottom: 2px solid #e0e0e0;
  padding-bottom: 10px;
  position: relative;
}

.title {
  color: #2c3e50;
  font-size: 32px;
  font-weight: 700;
  letter-spacing: 1px;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.1);
}

.box-card {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  padding: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.back-btn {
  background: #409eff;
  border-color: #409eff;
  border-radius: 6px;
  transition: all 0.3s ease;
  position: absolute;
  right: 0;
}

.back-btn:hover {
  background: #3a8ee6;
  border-color: #3a8ee6;
  transform: scale(1.05);
}

.el-table {
  margin-top: 10px;
  font-size: 14px;
}

.el-table__header th {
  background-color: #f5f7fa;
  color: #606266;
  font-weight: 600;
}
</style>

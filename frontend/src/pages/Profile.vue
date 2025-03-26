<template>
    <div class="profile-page">
      <el-button type="primary" @click="$router.push('/problems')" style="margin-bottom: 20px;">返回题目列表</el-button>
      <el-button type="danger" @click="logout" style="float: right; margin-bottom: 20px;">登出</el-button>
      <h2>我的提交记录</h2>
      <el-table :data="submissions" style="width: 100%" v-if="submissions.length">
        <el-table-column prop="ID" label="编号" width="100" />
        <el-table-column prop="problem_id" label="题目 ID" width="100" />
        <el-table-column prop="language" label="语言" width="100" />
        <el-table-column prop="result" label="评测结果" width="150" />
        <el-table-column prop="submit_time" label="提交时间" width="200">
          <template #default="{ row }">
            {{ formatTime(row.submit_time) }}
          </template>
        </el-table-column>
        <el-table-column label="通过 / 总数">
          <template #default="{ row }">
            {{ row.passed_count }} / {{ row.total_count }}
          </template>
        </el-table-column>
        <el-table-column label="详情" width="100">
          <template #default="{ row }">
            <el-button type="text" @click="$router.push(`/submissions/${row.ID}/results`)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-else>暂无提交记录</div>
      <el-pagination
        v-if="total > 0"
        :current-page="page"
        :page-size="10"
        :total="total"
        @current-change="fetchSubmissions"
        layout="total, prev, pager, next"
      />
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRouter } from 'vue-router'
  import axios from '../api/axios'
  
  const router = useRouter()
  const submissions = ref([])
  const total = ref(0)
  const page = ref(1)
  
  const fetchSubmissions = async (currentPage) => {
    page.value = currentPage
    try {
      const res = await axios.get(`/auth/submissions?page=${page.value}`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      const { submissions: data, total: totalCount } = res.data
      submissions.value = data
      total.value = totalCount
      console.log('Submissions fetched:', { total: total.value, page: page.value, submissions: submissions.value })
    } catch (err) {
      console.error('获取提交记录失败', err)
    }
  }
  
  const logout = () => {
    localStorage.removeItem('token')
    router.push('/login')
  }
  
  onMounted(() => {
    fetchSubmissions(page.value)
  })
  
  function formatTime(timeStr) {
    const date = new Date(timeStr)
    return date.toLocaleString()
  }
  </script>
  
  <style scoped>
  .profile-page {
    max-width: 900px;
    margin: 40px auto;
    padding: 20px;
    background-color: #ffffff;
    border-radius: 8px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  }
  </style>
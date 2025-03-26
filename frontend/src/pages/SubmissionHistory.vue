<template>
    <div class="submission-history">
        <div class="button-row" style="display: flex; justify-content: space-between;">
      <el-button @click="$router.push('/problems')">返回题目列表</el-button>
      <el-button @click="$router.push('/profile')">我的主页</el-button>
    </div>
      <h2>
        <a @click.prevent="goToProblem" style="cursor: pointer; text-decoration: underline;">
          {{ problemTitle }}
        </a> 提交记录
      </h2>
      
      <table>
        <thead>
          <tr>
            <th>提交编号</th>
            <th>用户</th>
            <th>语言</th>
            <th>结果</th>
            <th>提交时间</th>
            <th>通过 / 总数</th>
            <th>详情</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="submission in submissionHistory.data" :key="submission.ID">
            <td>{{ submission.ID }}</td>
            <td>{{ submission.user_id }}</td>
            <td>{{ submission.language }}</td>
            <td><span :class="getStatusClass(submission.result)">{{ submission.result }}</span></td>
            <td>{{ new Date(submission.submit_time).toLocaleString() }}</td>
            <td>{{ submission.passed_count }} / {{ submission.total_count }}</td>
            <td>
              <el-button size="small" type="primary" @click="goToDetail(submission.ID)">查看</el-button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- 分页控件 -->
      <el-pagination
        v-if="submissionHistory.total > submissionHistory.limit"
        :current-page="submissionHistory.page"
        :page-size="submissionHistory.limit"
        :total="submissionHistory.total"
        layout="prev, pager, next, jumper"
        @current-change="handlePageChange">
      </el-pagination>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import axios from '../api/axios'
  
  const route = useRoute()
  const router = useRouter()
  const problemTitle = ref('')
  const submissionHistory = ref({
    data: [],
    total: 0,
    page: 1,
    limit: 10
  })
  
  const fetchSubmissionHistory = async (page = 1) => {
  const problemID = route.params.id
  try {
    const response = await axios.get(`/auth/problems/${problemID}/submissions`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      },
      params: {
        page: page,
        limit: submissionHistory.value.limit
      }
    })

    // ✅ 微调：确保 data、total 正确赋值
    submissionHistory.value.data = response.data.submissions
    submissionHistory.value.page = page
    submissionHistory.value.total = response.data.total || response.data.submissions.length // 如果后端没返回 total
  } catch (error) {
    console.error('获取提交记录失败：', error)
  }
}
  
  const handlePageChange = (page) => {
    submissionHistory.value.page = page
    fetchSubmissionHistory(page)
  }
  
  onMounted(async () => {
    const problemID = route.params.id
    try {
      const res = await axios.get(`/problems/${problemID}`)
      problemTitle.value = res.data.title
      fetchSubmissionHistory()
    } catch (error) {
      console.error('获取题目失败：', error)
    }
  })
  
  const getStatusClass = (status) => {
    return status.toLowerCase()
  }

  const goToDetail = (submissionId) => {
    router.push(`/submissions/${submissionId}/results`)
  }

  const goToProblemList = () => {
    router.push('/problems')
  }

  const goToProfile = () => {
    router.push('/profile')
  }

  const goToProblem = () => {
    const problemID = route.params.id
    router.push(`/problems/${problemID}`)
  }
  </script>
  
  <style scoped>
  .submission-history {
    width: 80%;
    margin: 20px auto;
  }
  
  h2 {
    font-size: 28px;
    margin-bottom: 20px;
  }
  
  table {
    width: 100%;
    border-collapse: collapse;
  }
  
  th, td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }
  
  th {
    background-color: #f2f2f2;
  }
  
  .accepted {
    color: green;
  }
  
  .wrong {
    color: red;
  }
  
  .tle {
    color: orange;
  }
  
  .mle {
    color: purple;
  }
  
  .re {
    color: blue;
  }

  .el-button {
    padding: 4px 8px;
    font-size: 14px;
  }

  .button-row {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
  }
  </style>
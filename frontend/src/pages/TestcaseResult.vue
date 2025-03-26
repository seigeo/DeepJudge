<template>
    <div class="result-page">
      <div class="button-container">
        <el-button @click="$router.push('/problems')">返回题目列表</el-button>
        <el-button @click="$router.push('/profile')">我的主页</el-button>
      </div>
      <h2 class="title">测试点结果</h2>
      <div class="problem-info">
        题目编号：{{ problemId }}
        <el-button @click="$router.push(`/problems/${problemId}`)">查看题目</el-button>
      </div>
      <div class="submission-card">
        <div><strong>提交时间：</strong>{{ submission.submit_time && new Date(submission.submit_time).toLocaleString() }}</div>
        <div><strong>用户ID：</strong>{{ submission.user_id }}</div>
        <div><strong>提交语言：</strong>{{ submission.language }}</div>
        <div><strong>通过数量：</strong><strong>{{ submission.passed_count }} / <strong>总数量：</strong>{{ submission.total_count }}</strong></div>
      </div>
      <div>
        <span>源代码：</span>
        <el-button @click="copyToClipboard">复制</el-button>
      </div>
      <pre class="code-area">{{ submission.code }}</pre>
      <el-table :data="results" style="width: 100%; margin-top: 20px;" v-if="results.length" class="result-table">
        <el-table-column prop="case_id" label="测试点编号" width="100" />
        <el-table-column label="评测状态" width="160">
          <template v-slot="{ row }">
            <span :style="{ color: row.status === 'Accepted' ? 'green' : row.status === 'Wrong Answer' ? 'red' : row.status === 'Runtime Error' ? 'orange' : row.status === 'Time Limit Exceeded' ? 'blue' : row.status === 'MLE' ? 'purple' : 'black' }">
              {{ row.status }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="runtime_ms" label="运行时间 (ms)" width="160" />
        <el-table-column label="实际输出">
          <template v-slot="{ row }">
            <div>
              <span>
                {{ row.output && row.output.length > 50 ? row.output.slice(0, 1000) : row.output }}
              </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="期望输出">
          <template v-slot="{ row }">
            <div>
              <span>
                {{ row.expected && row.expected.length > 50 ? row.expected.slice(0, 1000) : row.expected }}
              </span>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="empty-message">暂无测试点信息</div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const results = ref([])
const problemId = ref('')
const submission = ref({})

import { useRoute } from 'vue-router'
import axios from '../api/axios'

const route = useRoute()

const copyToClipboard = () => {
  navigator.clipboard.writeText(submission.value.code).then(() => {
    ElMessage.success('代码已复制到剪贴板！')
  }).catch(err => {
    console.error('复制失败', err)
  });
}

onMounted(async () => {
  const submissionId = route.params.id
  try {
      const res = await axios.get(`/auth/submissions/${submissionId}/results`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    results.value = res.data.results
    submission.value = res.data.submission
    problemId.value = submission.value.problem_id
  } catch (err) {
    console.error('获取测试点结果失败', err)
  }
})
</script>

<style scoped>
.result-page {
  max-width: 1000px;
  margin: 40px auto;
  padding: 20px;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}
.button-container {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}
.title {
  text-align: center;
  margin-bottom: 20px;
}
.problem-info {
  text-align: center;
  margin-bottom: 20px;
}
.submission-card {
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #f9f9f9;
}
.code-area {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #e0e0e0;
  background-color: #f0f0f0;
  font-family: monospace;
  padding: 10px;
  border-radius: 8px;
}
.result-table {
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.empty-message {
  text-align: center;
}
</style>
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
      <el-button type="primary" @click="analyzeSubmission">开始分析</el-button>
      <div class="analysis-records" v-if="analysisRecords.length">
        <h3>分析记录</h3>
        <el-table :data="analysisRecords" style="width: 100%; margin-top: 20px;">
          <el-table-column prop="ID" label="分析ID" width="120" />
          <el-table-column prop="status" label="状态" width="160">
            <template v-slot="{ row }">
              <span :style="{ color: row.status === 'completed' ? 'green' : row.status === 'failed' ? 'red' : row.status === 'pending' ? 'gray' : row.status === 'analyzing' ? 'blue' : 'black' }">
                {{ row.status }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="结果" width="700">
            <template v-slot="{ row }">
              <div v-if="row.content">
                <div
                  v-html="renderMarkdown(row.expanded ? row.content : (row.content.length > 50 ? row.content.slice(0, 50) + '...' : row.content))"
                ></div>
                <el-button
                  type="text"
                  size="small"
                  @click="row.expanded = !row.expanded"
                  style="margin-top: 4px;"
                >
                  {{ row.expanded ? '收起' : '展开' }}
                </el-button>
              </div>
              <span v-else style="color: gray;">无摘要</span>
            </template>
          </el-table-column>
        </el-table>
      </div>
      <div v-else class="empty-message">暂无分析记录</div>
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
import { useRoute } from 'vue-router'
import axios from '../api/axios'
import MarkdownIt from 'markdown-it'
import markdownItKatex from 'markdown-it-katex'
const md = new MarkdownIt().use(markdownItKatex)
const renderMarkdown = (text) => md.render(text)

const results = ref([])
const problemId = ref('')
const submission = ref({})
const analysisRecords = ref([])

console.log('analysisRecords', analysisRecords.value)
const route = useRoute()

const copyToClipboard = () => {
  navigator.clipboard.writeText(submission.value.code).then(() => {
    ElMessage.success('代码已复制到剪贴板！')
  }).catch(err => {
    console.error('复制失败', err)
  });
}

const analyzeSubmission = async () => {
  const submissionId = route.params.id
  try {
    await axios.post(`/auth/submissions/${submissionId}/analyze`, {}, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    ElMessage.success('分析请求已发送！')
  } catch (err) {
    console.error('分析请求失败', err)
  }
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

    const analysisRes = await axios.get(`/auth/submissions/${submissionId}/analysis`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    analysisRecords.value = analysisRes.data.analyses.reverse().map(record => ({ ...record, expanded: false }))
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
.analysis-records {
  margin-top: 20px;
}
.result-table {
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
.empty-message {
  text-align: center;
}
</style>
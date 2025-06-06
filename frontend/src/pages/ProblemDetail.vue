<template>
  <div class="problem-detail" v-if="problem">
    <div class="button-row">
      <el-button @click="$router.push('/problems')" class="back-button">返回题目列表</el-button>
      <el-button @click="$router.push('/profile')">我的主页</el-button>
    </div>
    
    <h2>{{ problem.title }}</h2>

    <section class="difficulty-section">
      <p><strong>难度：</strong><span :class="'difficulty-' + problem.difficulty.toLowerCase()">{{ problem.difficulty }}</span></p>
      <el-button @click="goToSubmissionHistory" class="view-submissions-button">查看提交记录</el-button>
    </section>

    <section class="card">
      <h3>题目描述</h3>
      <div class="markdown-body" v-html="descriptionHtml"></div>
    </section>

    <section class="card">
      <h3>输入描述</h3>
      <div class="markdown-body" v-html="inputHtml"></div>
    </section>

    <section class="card">
      <h3>输出描述</h3>
      <div class="markdown-body" v-html="outputHtml"></div>
    </section>

    <section class="card">
      <h3>样例输入</h3>
      <pre>{{ problem.sample_input }}</pre>
    </section>

    <section class="card">
      <h3>样例输出</h3>
      <pre>{{ problem.sample_output }}</pre>
    </section>

    <section class="card">
      <h3>提交代码</h3>
      <el-select v-model="selectedLanguage" placeholder="选择语言" class="code-select">
        <el-option label="Python" value="python"></el-option>
        <el-option label="Java" value="java"></el-option>
        <el-option label="C++" value="cpp"></el-option>
      </el-select>
      <el-input type="textarea" v-model="code" placeholder="输入你的代码" rows="10" class="code-input"></el-input>
      <el-button type="primary" @click="submitCode" class="submit-button">提交</el-button>
    </section>

    <section v-if="result" class="card">
      <h3>评测结果</h3>
      <p><strong>结果信息：</strong>{{ result }}</p>
      <p><strong>通过测试点数：</strong>{{ passed }}</p>
      <p><strong>总测试点数：</strong>{{ total }}</p>
    </section>

    <section v-if="caseResults.length" class="card">
      <h3>测试点评测详情</h3>
      <table>
        <thead>
          <tr>
            <th>CaseID</th>
            <th>状态</th>
            <th>耗时</th>
            <th>输出</th>
            <th>期望</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="caseResult in caseResults" :key="caseResult.caseId">
            <td>{{ caseResult.caseId }}</td>
            <td>
              <span :class="getStatusClass(caseResult.status)">{{ caseResult.status }}</span>
            </td>
            <td>{{ caseResult.time }}</td>
            <td>{{ caseResult.output }}</td>
            <td>{{ caseResult.expected }}</td>
          </tr>
        </tbody>
      </table>
    </section>
  </div>

  <div v-else>
    <el-spinner />
    <p>加载中...</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from '../api/axios'
import { useRoute, useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'
import mathjax3 from 'markdown-it-mathjax3'
import '../styles/markdown.css'

const md = new MarkdownIt({
  html: true,
  breaks: true,
  linkify: true
})
md.use(mathjax3)

const route = useRoute()
const problem = ref(null)
const descriptionHtml = ref('')
const inputHtml = ref('')
const outputHtml = ref('')
const selectedLanguage = ref('')
const code = ref('')
const result = ref('')
const passed = ref(0)
const total = ref(0)
const caseResults = ref([])
const router = useRouter()

const renderMarkdown = (text) => {
  if (!text) return ''
  return md.render(text)
}

const submitCode = async () => {
  const id = route.params.id
  if (!code.value || !selectedLanguage.value) {
    alert("请输入代码并选择语言")
    return
  }

  try {
    const res = await axios.post(`/auth/problems/${id}/submit`, {
      code: code.value,
      language: selectedLanguage.value
    })

    result.value = "Pending"
    passed.value = 0
    total.value = 0
    caseResults.value = []
    alert("提交成功，已加入评测队列，稍后请查看提交记录")
  } catch (err) {
    console.error("提交失败：", err)
    alert("提交失败，请稍后重试")
  }
}

const goToSubmissionHistory = () => {
  const problemID = route.params.id
  router.push(`/submissions/${problemID}`)
}

onMounted(async () => {
  const id = route.params.id
  try {
    const res = await axios.get(`/problems/${id}`)
    problem.value = res.data
    descriptionHtml.value = renderMarkdown(problem.value.description)
    inputHtml.value = renderMarkdown(problem.value.input)
    outputHtml.value = renderMarkdown(problem.value.output)
  } catch (error) {
    console.error('加载题目失败：', error)
  }
})

const getStatusClass = (status) => {
  return {
    'Accepted': 'status-success',
    'Wrong Answer': 'status-error',
    'Runtime Error': 'status-error',
    'Time Limit Exceeded': 'status-warning',
    'Compilation Error': 'status-error',
    'Pending': 'status-pending'
  }[status] || 'status-default'
}
</script>

<style scoped>
.problem-detail {
  width: 80%;
  margin: 20px auto;
}

.button-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.problem-detail h2 {
  font-size: 28px;
  margin-bottom: 10px;
}

.problem-detail h3 {
  font-size: 20px;
  color: #333;
  border-left: 4px solid #409EFF;
  padding-left: 10px;
  margin-bottom: 10px;
}

.card {
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

section {
  margin-bottom: 20px;
}

.code-select {
  width: 120px;
  margin-bottom: 10px;
}

.code-input {
  margin-bottom: 10px;
  font-family: monospace;
}

.submit-button {
  width: 100%;
}

.difficulty-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.difficulty-easy {
  color: #67C23A;
}

.difficulty-medium {
  color: #E6A23C;
}

.difficulty-hard {
  color: #F56C6C;
}

.status-success {
  color: #67C23A;
}

.status-error {
  color: #F56C6C;
}

.status-warning {
  color: #E6A23C;
}

.status-pending {
  color: #909399;
}

.status-default {
  color: #606266;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}

th, td {
  border: 1px solid #EBEEF5;
  padding: 12px;
  text-align: left;
}

th {
  background-color: #F5F7FA;
  color: #606266;
}

tr:hover {
  background-color: #F5F7FA;
}

/* 添加 MathJax 相关样式 */
:deep(.mjx-chtml) {
  margin: 1em 0;
  font-size: 1.1em !important;
}

:deep(.mjx-chtml.MJXc-display) {
  margin: 1em 0;
  overflow-x: auto;
  overflow-y: hidden;
  padding: 1em 0;
}
</style>
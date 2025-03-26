<template>
    <div class="problem-detail" v-if="problem">
      <div class="button-row">
        <el-button @click="$router.push('/problems')" class="back-button">返回题目列表</el-button>
        <el-button @click="$router.push('/profile')">我的主页</el-button>
      </div>
      <h2>{{ problem.title }}</h2>
      <p><strong>难度：</strong><span :class="'difficulty-' + problem.difficulty.toLowerCase()">{{ problem.difficulty }}</span></p>
  
      <section class="card">
        <h3>题目描述</h3>
        <p v-html="descriptionHtml"></p>
      </section>
  
      <section class="card">
        <h3>输入描述</h3>
        <p v-html="inputHtml"></p>
      </section>
  
      <section class="card">
        <h3>输出描述</h3>
        <p v-html="outputHtml"></p>
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
      <section v-if="submissionHistory.length" class="card">
        <h3>历史提交记录</h3>
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
            <tr v-for="submission in submissionHistory" :key="submission.id">
              <td>{{ submission.ID }}</td>
              <td>{{ submission.user_id }}</td>
              <td>{{ submission.language }}</td>
              <td><span :class="getStatusClass(submission.result)">{{ submission.result }}</span></td>
              <td>{{ new Date(submission.submit_time).toLocaleString() }}</td>
              <td>{{ submission.passed_count }} / {{ submission.total_count }}</td>
              <td><el-button @click="$router.push(`/submissions/${submission.ID}/results`)" size="small">详情</el-button></td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import axios from '../api/axios'
  import { marked } from 'marked'
  import katexExtension from 'marked-katex-extension'
  import 'katex/dist/katex.min.css'

  marked.use(katexExtension())
  import { ElSelect, ElOption, ElInput, ElButton } from 'element-plus'

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
  const caseResults = ref([]) // 新增
  const submissionHistory = ref([]) // 新增

  const getStatusClass = (status) => {
    return status.toLowerCase();
  }

  const submitCode = async () => {
    const id = route.params.id
    console.log("Submit request body:", { code: code.value, language: selectedLanguage.value })
    try {
      const response = await axios.post(`/auth/problems/${id}/submit`, {
        code: code.value,
        language: selectedLanguage.value // ✅ 加上 .value
      }, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      const data = response.data
      result.value = data.result
      passed.value = data.passed
      total.value = data.total
      caseResults.value = data.caseResults // 更新测试点评测详情
      alert('代码提交成功！')
    } catch (error) {
      console.error('提交代码失败：', error)
      alert('提交代码失败，请重试！')
    }
  }
  
  onMounted(async () => {
    const id = route.params.id
    try {
      const res = await axios.get(`/problems/${id}`)
      problem.value = res.data
      descriptionHtml.value = marked(problem.value.description)
      inputHtml.value = marked(problem.value.input)
      outputHtml.value = marked(problem.value.output)

      // 获取历史提交记录
      const submissionsResponse = await axios.get(`/auth/problems/${id}/submissions`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      submissionHistory.value = submissionsResponse.data // 更新历史提交记录
    } catch (err) {
      console.error('获取题目失败：', err)
    }
  })
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
  }

  .katex {
  font-size: 1.05em;
}
  
  section {
    margin-bottom: 20px;
  }

  .code-input {
    margin-top: 10px;
    font-family: monospace;
  }

  .code-select {
    width: 200px;
    margin-bottom: 10px;
  }

  .submit-button {
    margin-top: 10px;
  }

  .back-button {
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
  .difficulty-easy {
    color: green;
  }
  .difficulty-medium {
    color: orange;
  }
  .difficulty-hard {
    color: red;
  }
  </style>
<template>
    <div v-if="!token" style="text-align: center; margin-bottom: 20px">
      <el-input v-model="loginForm.username" placeholder="用户名" style="width: 200px; margin-right: 10px" />
      <el-input v-model="loginForm.password" type="password" placeholder="密码" style="width: 200px; margin-right: 10px" />
      <el-button type="primary" @click="login">登录</el-button>
    </div>
    <div class="admin-dashboard">
      <div class="top-bar">
        <h2>题目管理后台</h2>
        <el-button type="primary" @click="showCreateDialog = true">新建题目</el-button>
      </div>
  
      <el-table :data="problems" style="width: 100%; margin-top: 20px">
        <el-table-column prop="id" label="题目ID" width="80" />
        <el-table-column prop="title" label="题目标题" />
        <el-table-column label="操作" width="300">
          <template #default="scope">
            <el-button type="primary" size="small" @click="openEdit(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="deleteProblem(scope.row.id)">删除</el-button>
            <el-upload
              :action="`http://localhost:8080/auth/edit/${scope.row.id}/upload`"
              :headers="{ Authorization: `Bearer ${token}` }"
              :show-file-list="false"
              :before-upload="beforeUpload"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
            >
              <el-button type="success" size="small">上传测试用例</el-button>
            </el-upload>
            <el-popover
              placement="bottom"
              title="上传说明"
              :width="300"
              trigger="hover"
            >
              <template #reference>
                <el-button type="info" size="small" circle><i class="el-icon-question"></i></el-button>
              </template>
              <div>
                <p>支持两种上传方式：</p>
                <ol>
                  <li>单个测试用例文件（.in或.out）</li>
                  <li>包含多个测试用例的zip压缩包</li>
                </ol>
                <p>文件命名规则：</p>
                <ul>
                  <li>单个文件：数字.in 或 数字.out（例如：1.in, 1.out）</li>
                  <li>zip包中的文件同样需要遵循上述命名规则</li>
                </ul>
              </div>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>
  
      <!-- 创建/编辑弹窗 -->
      <el-dialog :title="editMode ? '编辑题目' : '新建题目'" v-model="showCreateDialog">
        <el-form :model="form">
          <el-form-item label="标题">
            <el-input v-model="form.title" />
          </el-form-item>
          <el-form-item label="难度">
            <el-select v-model="form.difficulty" placeholder="选择难度">
              <el-option label="简单" value="easy" />
              <el-option label="中等" value="medium" />
              <el-option label="困难" value="hard" />
            </el-select>
          </el-form-item>
          <el-form-item label="描述">
            <el-input type="textarea" v-model="form.description" rows="4" />
          </el-form-item>
          <el-form-item label="输入描述">
            <el-input type="textarea" v-model="form.input" rows="2" />
          </el-form-item>
          <el-form-item label="输出描述">
            <el-input type="textarea" v-model="form.output" rows="2" />
          </el-form-item>
          <el-form-item label="样例输入">
            <el-input v-model="form.sample_input" />
          </el-form-item>
          <el-form-item label="样例输出">
            <el-input v-model="form.sample_output" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="submitForm">{{ editMode ? '保存' : '创建' }}</el-button>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script setup>
  import { ref, reactive, onMounted } from 'vue'
  import axios from '../api/axios'
  import { ElMessage } from 'element-plus'
  
  const problems = ref([])
  const showCreateDialog = ref(false)
  const editMode = ref(false)
  const form = ref({
    title: '',
    difficulty: '',
    description: '',
    input: '',
    output: '',
    sample_input: '',
    sample_output: ''
  })
  let editingId = null
  
  const loginForm = reactive({
    username: '',
    password: ''
  })

  const token = ref(localStorage.getItem('token') || '')

  const login = async () => {
    try {
      const res = await axios.post('/login', loginForm)
      localStorage.setItem('token', res.data.token)
      token.value = res.data.token
      ElMessage.success('登录成功')
      loadProblems()
    } catch {
      ElMessage.error('登录失败')
    }
  }

  const loadProblems = async () => {
    if (!token.value) return
    const res = await axios.get('/problems', {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    problems.value = res.data
  }
  
  const submitForm = async () => {
    try {
      if (editMode.value) {
        await axios.put(`/auth/edit/${editingId}`, form.value, {
          headers: { Authorization: `Bearer ${token.value}` }
        })
        ElMessage.success('编辑成功')
      } else {
        await axios.post('/auth/edit/add', form.value, {
          headers: { Authorization: `Bearer ${token.value}` }
        })
        ElMessage.success('创建成功')
      }
      showCreateDialog.value = false
      loadProblems()
    } catch (err) {
      ElMessage.error('操作失败')
    }
  }
  
  const deleteProblem = async (id) => {
    try {
      await axios.delete(`/auth/edit/${id}`, {
        headers: { Authorization: `Bearer ${token.value}` }
      })
      ElMessage.success('删除成功')
      loadProblems()
    } catch {
      ElMessage.error('删除失败')
    }
  }
  
  const openEdit = (problem) => {
    editingId = problem.id
    editMode.value = true
    showCreateDialog.value = true
    form.value = { ...problem }
  }
  
  const beforeUpload = (file) => {
    const isZIP = file.type === 'application/zip' || file.name.endsWith('.zip')
    const isTestCase = file.name.match(/^\d+\.(in|out)$/)
    
    if (!isZIP && !isTestCase) {
      ElMessage.error('只能上传zip压缩包或者.in/.out格式的测试用例文件！')
      return false
    }
    
    return true
  }

  const handleUploadSuccess = (response) => {
    if (response.pairs) {
      ElMessage.success(`成功上传${response.pairs.length}组测试用例`)
    } else {
      ElMessage.success('文件上传成功')
    }
    loadProblems()
  }

  const handleUploadError = () => {
    ElMessage.error('上传失败，请检查文件格式是否正确')
  }
  
  onMounted(loadProblems)
  </script>
  
  <style scoped>
  .admin-dashboard {
    width: 90%;
    margin: 20px auto;
  }
  .top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  </style>
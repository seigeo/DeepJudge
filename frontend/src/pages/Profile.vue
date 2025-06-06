<template>
    <div class="profile-page">
      <el-button type="primary" @click="$router.push('/problems')" style="margin-bottom: 20px;">返回题目列表</el-button>
      <el-button type="danger" @click="logout" style="float: right; margin-bottom: 20px;">登出</el-button>

      <!-- 用户信息卡片 -->
      <el-card class="user-info-card">
        <template #header>
          <div class="card-header">
            <span>个人信息</span>
            <el-button type="primary" link @click="editSignature" v-if="!isEditing">编辑签名</el-button>
            <div v-else>
              <el-button type="success" link @click="saveSignature">保存</el-button>
              <el-button type="danger" link @click="cancelEdit">取消</el-button>
            </div>
          </div>
        </template>
        <div class="user-info">
          <p><strong>用户 ID：</strong>{{ userInfo.user_id }}</p>
          <p><strong>用户名：</strong>{{ userInfo.username }}</p>
          <p><strong>注册时间：</strong>{{ formatTime(userInfo.created_at) }}</p>
          <p><strong>解题数量：</strong>{{ userInfo.solved_count }}</p>
          <div class="signature-section">
            <strong>个性签名：</strong>
            <div v-if="!isEditing" class="signature-text">
              {{ userInfo.signature || '这个人很懒，还没有写签名～' }}
            </div>
            <el-input
              v-else
              v-model="editingSignature"
              type="textarea"
              :rows="3"
              placeholder="写点什么吧..."
              maxlength="200"
              show-word-limit
            />
          </div>
        </div>
      </el-card>

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
  import { ElMessage } from 'element-plus'
  
  const router = useRouter()
  const submissions = ref([])
  const total = ref(0)
  const page = ref(1)
  const userInfo = ref({})
  const isEditing = ref(false)
  const editingSignature = ref('')
  
  const fetchUserInfo = async () => {
    try {
      const res = await axios.get('/auth/profile', {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`
        }
      })
      userInfo.value = res.data
    } catch (err) {
      console.error('获取用户信息失败', err)
      ElMessage.error('获取用户信息失败')
    }
  }

  const editSignature = () => {
    isEditing.value = true
    editingSignature.value = userInfo.value.signature || ''
  }

  const cancelEdit = () => {
    isEditing.value = false
    editingSignature.value = ''
  }

  const saveSignature = async () => {
    try {
      await axios.put('/auth/profile', 
        { signature: editingSignature.value },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        }
      )
      userInfo.value.signature = editingSignature.value
      isEditing.value = false
      ElMessage.success('签名更新成功')
    } catch (err) {
      console.error('更新签名失败', err)
      ElMessage.error('更新签名失败')
    }
  }

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
    } catch (err) {
      console.error('获取提交记录失败', err)
      ElMessage.error('获取提交记录失败')
    }
  }
  
  const logout = () => {
    localStorage.removeItem('token')
    router.push('/login')
  }
  
  onMounted(() => {
    fetchUserInfo()
    fetchSubmissions(page.value)
  })
  
  function formatTime(timeStr) {
    if (!timeStr) return ''
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

  .user-info-card {
    margin-bottom: 30px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .user-info {
    padding: 10px;
  }

  .user-info p {
    margin: 10px 0;
    font-size: 14px;
  }

  .signature-section {
    margin-top: 20px;
  }

  .signature-text {
    margin-top: 10px;
    padding: 10px;
    background-color: #f5f7fa;
    border-radius: 4px;
    min-height: 60px;
    white-space: pre-wrap;
  }
  </style>
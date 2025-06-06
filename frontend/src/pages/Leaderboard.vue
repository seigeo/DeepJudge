<template>
  <div class="leaderboard">
    <div class="header-row">
      <h2>排行榜</h2>
      <div>
        <el-button @click="$router.push('/problems')">题目列表</el-button>
        <el-button @click="$router.push('/profile')">我的主页</el-button>
      </div>
    </div>

    <div v-if="userRank" class="user-rank">
      <h3>我的排名</h3>
      <div class="rank-info">
        <p>当前排名：<span class="highlight">{{ userRank.rank }}</span></p>
        <p>解题数量：<span class="highlight">{{ userRank.solved_count }}</span></p>
      </div>
    </div>

    <div class="rankings">
      <h3>TOP 50</h3>
      <el-table v-loading="loading" :data="rankings" style="width: 100%">
        <el-table-column prop="rank" label="排名" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="signature" label="个性签名">
          <template #default="{ row }">
            <span class="signature-text">{{ row.signature || '这个人很懒，还没有写签名～' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="solved_count" label="解题数" width="100" />
      </el-table>
      <div v-if="!loading && rankings.length === 0" class="empty-message">
        暂无排名数据
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from '../api/axios'
import { ElMessage } from 'element-plus'

const rankings = ref([])
const userRank = ref(null)
const loading = ref(true)

const fetchLeaderboard = async () => {
  try {
    loading.value = true
    const res = await axios.get('/leaderboard')
    rankings.value = res.data.rankings || []
  } catch (err) {
    console.error('获取排行榜失败：', err)
    ElMessage.error('获取排行榜失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const fetchUserRank = async () => {
  try {
    const res = await axios.get('/auth/rank', {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    userRank.value = res.data
  } catch (err) {
    console.error('获取用户排名失败：', err)
    if (err.response?.status !== 401) { // 忽略未登录错误
      ElMessage.error('获取用户排名失败，请稍后重试')
    }
  }
}

onMounted(() => {
  fetchLeaderboard()
  fetchUserRank()
})
</script>

<style scoped>
.leaderboard {
  width: 80%;
  margin: 20px auto;
  padding: 20px;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.user-rank {
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.rank-info {
  display: flex;
  gap: 40px;
}

.highlight {
  color: #409eff;
  font-weight: bold;
  font-size: 1.2em;
}

h2 {
  margin: 0;
  color: #303133;
}

h3 {
  margin: 0 0 15px 0;
  color: #606266;
}

.empty-message {
  text-align: center;
  color: #909399;
  padding: 20px;
  font-size: 14px;
}

.signature-text {
  color: #606266;
  font-size: 14px;
  white-space: pre-wrap;
  word-break: break-all;
}
</style> 
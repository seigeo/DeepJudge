<template>
  <div class="problem-list">
    <div class="header-row">
      <h2>题目列表</h2>
      <div>
        <el-button @click="$router.push('/leaderboard')">排行榜</el-button>
        <el-button @click="$router.push('/profile')">我的主页</el-button>
      </div>
    </div>
    <el-table :data="problems" style="width: 100%" @row-click="handleRowClick">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="题目名称" />
      <el-table-column label="通过率" width="150">
        <template #default="scope">
          <div>
            {{ scope.row.accepted_count }}/{{ scope.row.submission_count }}
            ({{ formatPassRate(scope.row.pass_rate) }})
          </div>
        </template>
      </el-table-column>
      <el-table-column label="难度" width="100">
        <template #default="scope">
          <span :class="scope.row.difficulty">{{ scope.row.difficulty }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from '../api/axios'

const problems = ref([])
const router = useRouter()

const handleRowClick = (row) => {
  router.push({ name: 'ProblemDetail', params: { id: row.id } })
}

const formatPassRate = (rate) => {
  return rate ? `${(rate * 100).toFixed(1)}%` : '0.0%'
}

onMounted(async () => {
  try {
    const res = await axios.get('/problems')
    problems.value = res.data
  } catch (err) {
    console.error('获取题目失败：', err)
  }
})
</script>

<style scoped>
.problem-list {
  width: 80%;
  margin: 20px auto;
  text-align: center;
}
.problem-table {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
.easy {
  color: green;
}
.medium {
  color: orange;
}
.hard {
  color: red;
}
.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
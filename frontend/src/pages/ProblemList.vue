<template>
  <div class="problem-list">
    <div class="header-row">
      <h2>题目列表</h2>
      <el-button @click="$router.push('/profile')">我的主页</el-button>
    </div>
    <el-table :data="problems" style="width: 100%" @row-click="handleRowClick">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="题目名称" />
      <el-table-column prop="accepted_count" label="通过数" width="100" />
      <el-table-column prop="submission_count" label="提交数" width="100" />
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

onMounted(async () => {
  try {
    const res = await axios.get('/problems')
    problems.value = res.data // 根据你的后端返回结构，可能是 res.data 或 res.data.data
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
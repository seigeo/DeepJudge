<template>
  <div class="register">
    <el-card class="box-card">
      <h2>注册</h2>
      <el-form :model="form">
        <el-form-item label="用户名">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleRegister">注册</el-button>
          <el-button type="text" @click="goToLogin">已有账号？去登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const form = reactive({
  username: '',
  password: ''
})

const router = useRouter()

const handleRegister = async () => {
  try {
    await axios.post('http://localhost:8080/register', form)
    alert('注册成功，请登录')
    router.push('/login')
  } catch (err) {
    alert('注册失败：' + err.response.data.message)
  }
}

const goToLogin = () => {
  router.push('/login')
}
</script>

<style scoped>
.register {
  width: 400px;
  margin: 100px auto;
}
</style>
<template>
  <div class="login">
    <el-card class="box-card">
      <h2>登录</h2>
      <el-form :model="form">
        <el-form-item label="用户名">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" />
        </el-form-item>
        <el-form-item>
           <el-button type="primary" @click="handleLogin">登录</el-button>
           <el-button type="text" @click="goToRegister">没有账号？去注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()

const goToRegister = () => {
  router.push('/register')
}

const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  try {
    const res = await axios.post('http://localhost:8080/login', form)
    localStorage.setItem('token', res.data.token)
    alert('登录成功')
    // 登录成功后跳转
    router.push('/problems')
    // TODO: 跳转到题目列表
  } catch (err) {
    alert('登录失败')
  }
}
</script>

<style scoped>
.login {
  width: 400px;
  margin: 100px auto;
}
</style>
import { createRouter, createWebHistory } from 'vue-router'
import Login from '../pages/Login.vue'
import Register from '../pages/Register.vue'
import ProblemList from '../pages/ProblemList.vue'
import ProblemDetail from '../pages/ProblemDetail.vue'
import Profile from '../pages/Profile.vue'
import TestcaseResult from '../pages/TestcaseResult.vue'
import SubmissionHistory from '../pages/SubmissionHistory.vue'
import admin from '../pages/AdminDashboard.vue'
import Leaderboard from '../pages/Leaderboard.vue'

const routes = [
  { path: '/', redirect: '/login' },
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/problems', component: ProblemList },
  { path: '/problems/:id', name: 'ProblemDetail', component: ProblemDetail },
  { path: '/profile', component: Profile },
  { path: '/submissions/:id/results', component: TestcaseResult },
  { path: '/submissions/:id', component: SubmissionHistory },
  { path: '/admin', component: admin },
  { path: '/leaderboard', component: Leaderboard }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
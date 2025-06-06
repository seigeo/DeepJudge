<template>
  <footer class="footer">
    <div class="footer-content">
      <div class="main-row">
        <div class="logo-section">
          <img src="/logo.png" alt="DeepJudge Logo" class="footer-logo" />
        </div>
        <div class="info-section">
          <p class="copyright">© 2025 Deep Judge. All rights reserved.</p>
          <p class="author">作者：seigeo</p>
          <div class="links">
            <a href="https://github.com/seigeo/DeepJudge" target="_blank" class="social-link">GitHub</a>
            <span class="divider">|</span>
            <a href="mailto:tengjialiang@qq.com" class="social-link">Email</a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  padding: 20px 0;
  background-color: #f5f7fa;
  border-top: 1px solid #e4e7ed;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.main-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo-section {
  display: flex;
  align-items: center;
}

.footer-logo {
  height: 60px;
  width: auto;
  margin-right: 12px;
}

.logo-text {
  font-size: 28px;
  font-weight: bold;
  color: #409EFF;
  letter-spacing: 0.5px;
}

.info-section {
  text-align: right;
  color: #606266;
  font-size: 14px;
}

.info-section p {
  margin: 3px 0;
}

.copyright, .author {
  color: #909399;
}

.links {
  margin-top: 5px;
}

.social-link {
  color: #409EFF;
  text-decoration: none;
  transition: color 0.3s;
  font-size: 15px;
}

.social-link:hover {
  color: #66b1ff;
}

.divider {
  color: #dcdfe6;
  margin: 0 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .main-row {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .info-section {
    text-align: center;
    margin-top: 15px;
  }

  .logo-section {
    margin-bottom: 10px;
  }
}
</style> 
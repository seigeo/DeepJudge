# 🚀 DeepJudge

一个基于 **Go 语言** 构建的高并发智能 Online Judge 系统，集成 DeepSeek API 代码分析能力，具备完善的用户系统、题目管理、异步评测与 Docker 沙箱运行，使用 Redis 消息队列实现任务调度，致力于打造教学 & 训练场景下的轻量级评测平台。

---

## 🧠 项目特色

**DeepJudge** 注重性能与工程实践，具备以下核心特性：

- 👤 用户注册 / 登录 / JWT 鉴权
- 📋 题目增删改查 + 测试点上传（支持 .in/.out）
- ⚙️ 支持代码提交异步化 + 多 Worker 并发评测
- 🐳 Docker 容器沙箱运行代码，隔离安全
- 📬 使用 Redis 队列保障评测任务可控、可持久
- 📊 自动记录每组测试点评测状态（AC/WA/TLE）
- 🤖 预留 DeepSeek 智能代码分析接口

---

## 🔧 技术栈与架构亮点

| 模块         | 技术选型                              |
| ------------ | ------------------------------------- |
| 开发语言     | Go 1.20+                              |
| Web 框架     | Gin                                   |
| 用户鉴权     | JWT                                   |
| 数据存储     | SQLite（可平滑切换 MySQL/PostgreSQL） |
| 消息队列     | Redis Stream + Goroutine Worker Pool  |
| 排行榜系统   | Redis Sorted Set                      |
| 缓存层       | Redis（提交记录、用户信息）           |
| 沙箱运行     | Docker 容器（运行用户代码）           |
| 智能分析     | DeepSeek API                          |
| 前端         | Vue 3 + Element Plus         |

---

## 💡 后端技术亮点

- 🚀 基于 Redis Stream 的消息队列，保证评测任务的可靠性和持久性
- 🔄 多 Worker 并发评测，Worker 自动负载均衡
- 📊 使用 Redis Sorted Set 实现实时排行榜，支持多维度排序
- 🔒 Docker 沙箱隔离执行用户代码，严格限制系统资源
- 🎯 精确的评测结果（AC/WA/TLE/MLE/RE）和详细运行时信息
- 💾 Redis 多级缓存设计，提升系统响应速度
- 🛡️ 基于 JWT + 中间件的权限校验体系
- 🔍 完善的日志记录和错误追踪机制

---

## ✅ 当前进展

- ✅ 用户注册 / 登录 / 权限校验
- ✅ 题目管理（创建 / 编辑 / 删除 / 上传测试点）
- ✅ 支持代码提交 + Docker 编译/运行
- ✅ Redis 消息队列驱动并发评测任务
- ✅ 多 Worker 调度机制，提升评测吞吐量
- ✅ 每次提交记录所有测试点状态
- ✅ 提交记录分页接口
- ✅ DeepSeek 智能分析结果展示
- ✅ 安全限制：资源限额 / 防止刷题攻击
- ✅ 用户评测排行 / 提交统计

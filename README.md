# deepjudge

基于 Go 实现的轻量级 Online Judge 平台，集成 DeepSeek 的代码分析功能，支持用户在线注册、登录和题目管理，后续将扩展为完整的评测系统。

## 🧠 项目简介

**deepjudge** 是一个面向教学与开发实践的 OJ 平台，致力于打造一个集编程练习、代码提交、智能分析于一体的在线评测平台。目前平台初步实现：

- 用户账号系统（注册 / 登录）
- 题目管理功能（创建 / 编辑 / 删除）

未来将支持在线代码提交、编译运行、安全沙箱评测、智能反馈等功能。

## 🚧 当前进度

- ✅ 用户注册与登录
- ✅ 题目增删改查功能
- ⏳ 待开发：代码提交、评测队列、DeepSeek AI代码分析接口集成、判题后反馈等功能模块

## 🛠 技术栈

- **后端框架**：Go + Gin
- **数据库**：SQLite （开发中）
- **AI 分析**：DeepSeek API 接入（规划中）
- **部署支持**：Docker（可选）


## **📘 DeepJudge API 文档**

> 所有带 **/auth** 的接口均需要携带 JWT 令牌进行身份验证

> 请求头需包含：Authorization: Bearer <your_token>

---

**🧑 用户认证接口**

| **方法** | **路径** | **描述**          |
| -------------- | -------------- | ----------------------- |
| POST           | /register      | 用户注册                |
| POST           | /login         | 用户登录，返回 JWT 令牌 |

---

**📚 公共题目接口（无需登录）**

| **方法** | **路径** | **描述**   |
| -------------- | -------------- | ---------------- |
| GET            | /problems      | 获取题目列表     |
| GET            | /problems/:id  | 获取指定题目详情 |

---

**🔐 题目管理接口（需登录）**

| **方法** | **路径**     | **描述** |
| -------------- | ------------------ | -------------- |
| POST           | /auth/problems     | 创建新题目     |
| PUT            | /auth/problems/:id | 更新题目信息   |
| DELETE         | /auth/problems/:id | 删除题目       |

---

**💻 代码提交与评测接口**

| **方法** | **路径**            | **描述**             |
| -------------- | ------------------------- | -------------------------- |
| POST           | /auth/problems/:id/submit | 向指定题目提交代码进行评测 |

请求体示例（JSON）：

```
{
  "language": "cpp",
  "code": "#include<iostream>\nusing namespace std;\nint main() { int a, b; cin >> a >> b; cout << a + b; return 0; }"
}
```

---

**📂 测试点上传接口**

| **方法** | **路径**            | **描述**                                             |
| -------------- | ------------------------- | ---------------------------------------------------------- |
| POST           | /auth/problems/:id/upload | 上传测试点**.in/.out**文件到指定题目目录（multipart 格式） |

curl 示例：

```
curl -X POST http://localhost:8080/auth/problems/1/upload \
  -H "Authorization: Bearer <token>" \
  -F "files=@testcases/1/1.in" \
  -F "files=@testcases/1/1.out"
```

---

**📈 提交记录与评测结果接口**

| **方法** | **路径**                              | **描述**                             |
| -------------- | ------------------------------------------- | ------------------------------------------ |
| GET            | /auth/submissions                           | 获取当前用户所有题目的提交记录（支持分页） |
| GET            | /auth/problems/:id/submissions/:sid/results | 获取某次提交的所有测试点评测结果           |

分页参数示例：

```
/auth/submissions?page=1&limit=10
```

---

**✅ 返回字段参考**

**/auth/submissions 响应：**

```
{
  "total": 23,
  "page": 1,
  "limit": 10,
  "submissions": [
    {
      "id": 42,
      "problem_id": 3,
      "language": "cpp",
      "result": "Accepted",
      "passed_count": 3,
      "total_count": 3,
      "submit_time": "2025-03-25T15:42:00.123Z"
    }
  ]
}
```

---

**/auth/problems/:id/submissions/:sid/results 响应：**

```
[
  {
    "case_id": "1",
    "status": "Accepted",
    "output": "3",
    "expected": "3",
    "runtime_ms": 15
  },
  {
    "case_id": "2",
    "status": "Wrong Answer",
    "output": "4",
    "expected": "5",
    "runtime_ms": 18
  }
]
```

## 🏃 快速开始

1. 克隆项目：
   ```bash
   git clone https://github.com/your-username/deepjudge.git
   cd deepjudge
   ```

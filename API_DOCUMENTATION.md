# API 文档

本文档详细描述了 REST API 的所有端点、请求格式和响应示例。

## 📋 目录

- [基础信息](#基础信息)
- [认证系统](#认证系统)
- [学生管理 API](#学生管理-api)
- [教师管理 API](#教师管理-api)
- [管理员 API](#管理员-api)
- [错误处理](#错误处理)
- [查询参数](#查询参数)

## 🌐 基础信息

### Base URL
```
https://localhost:8443
```

### 通用响应格式

**成功响应**
```json
{
  "status": "success",
  "count": 10,
  "data": [...]
}
```

**分页响应**
```json
{
  "status": "success",
  "count": 100,
  "page": 1,
  "page_size": 10,
  "data": [...]
}
```

**错误响应**
```json
{
  "error": "错误描述信息"
}
```

### HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 204 | 无内容（更新/删除成功） |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 🔐 认证系统

### 1. 用户登录

**端点**: `POST /execs/login`

**请求体**:
```json
{
  "username": "admin",
  "password": "password123"
}
```

**响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

**cURL 示例**:
```bash
curl -k -X POST https://localhost:8443/execs/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

---

### 2. 用户登出

**端点**: `POST /execs/logout`

**请求头**:
```
Authorization: Bearer <token>
```

**响应**:
```json
{
  "status": "success",
  "message": "Logged out successfully"
}
```

---

### 3. 修改密码

**端点**: `POST /execs/{id}/updatepassword`

**请求头**:
```
Authorization: Bearer <token>
```

**请求体**:
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**响应**:
```json
{
  "token": "new_jwt_token...",
  "password_updated": true
}
```

---

### 4. 忘记密码

**端点**: `POST /execs/forgotpassword`

**请求体**:
```json
{
  "email": "user@example.com"
}
```

**响应**:
```json
{
  "status": "success",
  "message": "Password reset email sent"
}
```

---

### 5. 重置密码

**端点**: `POST /execs/resetpassword/reset/{resetcode}`

**路径参数**:
- `resetcode`: 重置密码的验证码（通过邮件发送）

**请求体**:
```json
{
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

**响应**:
```json
{
  "status": "success",
  "message": "Password reset successfully"
}
```

---

## 👨‍🎓 学生管理 API

### 数据模型

```json
{
  "id": 1,
  "first_name": "张",
  "last_name": "三",
  "email": "zhangsan@example.com",
  "class": "A"
}
```

### 1. 获取所有学生

**端点**: `GET /students`

**查询参数**:
- `page`: 页码（默认 1）
- `limit`: 每页数量（默认 10）
- `first_name`: 按名字过滤
- `last_name`: 按姓氏过滤
- `email`: 按邮箱过滤
- `class`: 按班级过滤
- `sortby`: 排序字段（格式：`field:ASC` 或 `field:DESC`）

**响应**:
```json
{
  "status": "success",
  "count": 50,
  "page": 1,
  "page_size": 10,
  "data": [
    {
      "id": 1,
      "first_name": "张",
      "last_name": "三",
      "email": "zhangsan@example.com",
      "class": "A"
    }
  ]
}
```

**cURL 示例**:
```bash
# 基础查询
curl -k "https://localhost:8443/students"

# 带分页
curl -k "https://localhost:8443/students?page=2&limit=20"

# 带过滤和排序
curl -k "https://localhost:8443/students?class=A&sortby=last_name:ASC"
```

---

### 2. 获取单个学生

**端点**: `GET /students/{id}`

**路径参数**:
- `id`: 学生 ID

**响应**:
```json
{
  "id": 1,
  "first_name": "张",
  "last_name": "三",
  "email": "zhangsan@example.com",
  "class": "A"
}
```

**cURL 示例**:
```bash
curl -k "https://localhost:8443/students/1"
```

---

### 3. 创建学生（批量）

**端点**: `POST /students`

**请求体**:
```json
[
  {
    "first_name": "张",
    "last_name": "三",
    "email": "zhangsan@example.com",
    "class": "A"
  },
  {
    "first_name": "李",
    "last_name": "四",
    "email": "lisi@example.com",
    "class": "B"
  }
]
```

**响应**:
```json
{
  "status": "success",
  "count": 2,
  "data": [
    {
      "id": 1,
      "first_name": "张",
      "last_name": "三",
      "email": "zhangsan@example.com",
      "class": "A"
    },
    {
      "id": 2,
      "first_name": "李",
      "last_name": "四",
      "email": "lisi@example.com",
      "class": "B"
    }
  ]
}
```

**cURL 示例**:
```bash
curl -k -X POST https://localhost:8443/students \
  -H "Content-Type: application/json" \
  -d '[
    {
      "first_name": "张",
      "last_name": "三",
      "email": "zhangsan@example.com",
      "class": "A"
    }
  ]'
```

---

### 4. 更新学生（完整更新）

**端点**: `PUT /students/{id}`

**路径参数**:
- `id`: 学生 ID

**请求体**:
```json
{
  "first_name": "张",
  "last_name": "三",
  "email": "zhangsan_new@example.com",
  "class": "B"
}
```

**响应**:
```json
{
  "id": 1,
  "first_name": "张",
  "last_name": "三",
  "email": "zhangsan_new@example.com",
  "class": "B"
}
```

**cURL 示例**:
```bash
curl -k -X PUT https://localhost:8443/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "张",
    "last_name": "三",
    "email": "zhangsan_new@example.com",
    "class": "B"
  }'
```

---

### 5. 更新学生（部分更新）

**端点**: `PATCH /students/{id}`

**路径参数**:
- `id`: 学生 ID

**请求体**（只需提供要更新的字段）:
```json
{
  "email": "newemail@example.com"
}
```

**响应**:
```json
{
  "id": 1,
  "first_name": "张",
  "last_name": "三",
  "email": "newemail@example.com",
  "class": "A"
}
```

**cURL 示例**:
```bash
curl -k -X PATCH https://localhost:8443/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'
```

---

### 6. 批量部分更新

**端点**: `PATCH /students`

**请求体**:
```json
[
  {
    "id": 1,
    "class": "C"
  },
  {
    "id": 2,
    "email": "updated@example.com"
  }
]
```

**响应**: `204 No Content`

**cURL 示例**:
```bash
curl -k -X PATCH https://localhost:8443/students \
  -H "Content-Type: application/json" \
  -d '[
    {"id": 1, "class": "C"},
    {"id": 2, "email": "updated@example.com"}
  ]'
```

---

### 7. 删除单个学生

**端点**: `DELETE /students/{id}`

**路径参数**:
- `id`: 学生 ID

**响应**:
```json
{
  "status": "Student successfully deleted",
  "id": 1
}
```

**cURL 示例**:
```bash
curl -k -X DELETE https://localhost:8443/students/1
```

---

### 8. 批量删除学生

**端点**: `DELETE /students`

**请求体**:
```json
[1, 2, 3, 4, 5]
```

**响应**:
```json
{
  "status": "students successfully deleted",
  "deleted_ids": [1, 2, 3, 4, 5]
}
```

**cURL 示例**:
```bash
curl -k -X DELETE https://localhost:8443/students \
  -H "Content-Type: application/json" \
  -d '[1, 2, 3]'
```

---

## 👨‍🏫 教师管理 API

### 数据模型

```json
{
  "id": 1,
  "first_name": "王",
  "last_name": "老师",
  "email": "wang@example.com",
  "class": "A",
  "subject": "数学"
}
```

### 1. 获取所有教师

**端点**: `GET /teachers`

**查询参数**: 同学生 API，额外支持：
- `subject`: 按科目过滤

**响应**:
```json
{
  "status": "success",
  "count": 20,
  "page": 1,
  "page_size": 10,
  "data": [
    {
      "id": 1,
      "first_name": "王",
      "last_name": "老师",
      "email": "wang@example.com",
      "class": "A",
      "subject": "数学"
    }
  ]
}
```

---

### 2. 获取单个教师

**端点**: `GET /teachers/{id}`

**响应**:
```json
{
  "id": 1,
  "first_name": "王",
  "last_name": "老师",
  "email": "wang@example.com",
  "class": "A",
  "subject": "数学"
}
```

---

### 3. 获取教师的所有学生

**端点**: `GET /teachers/{id}/students`

**路径参数**:
- `id`: 教师 ID

**响应**:
```json
{
  "status": "success",
  "count": 25,
  "data": [
    {
      "id": 1,
      "first_name": "张",
      "last_name": "三",
      "email": "zhangsan@example.com",
      "class": "A"
    }
  ]
}
```

**cURL 示例**:
```bash
curl -k "https://localhost:8443/teachers/1/students"
```

---

### 4. 获取教师的学生数量

**端点**: `GET /teachers/{id}/studentcount`

**路径参数**:
- `id`: 教师 ID

**响应**:
```json
{
  "teacher_id": 1,
  "student_count": 25
}
```

**cURL 示例**:
```bash
curl -k "https://localhost:8443/teachers/1/studentcount"
```

---

### 5. 创建教师（批量）

**端点**: `POST /teachers`

**请求体**:
```json
[
  {
    "first_name": "王",
    "last_name": "老师",
    "email": "wang@example.com",
    "class": "A",
    "subject": "数学"
  }
]
```

**响应**:
```json
{
  "status": "success",
  "count": 1,
  "data": [
    {
      "id": 1,
      "first_name": "王",
      "last_name": "老师",
      "email": "wang@example.com",
      "class": "A",
      "subject": "数学"
    }
  ]
}
```

---

### 6. 更新教师（完整更新）

**端点**: `PUT /teachers/{id}`

**请求体**:
```json
{
  "first_name": "王",
  "last_name": "老师",
  "email": "wang_new@example.com",
  "class": "B",
  "subject": "物理"
}
```

---

### 7. 更新教师（部分更新）

**端点**: `PATCH /teachers/{id}`

**请求体**:
```json
{
  "subject": "化学"
}
```

---

### 8. 批量部分更新教师

**端点**: `PATCH /teachers`

**请求体**:
```json
[
  {
    "id": 1,
    "subject": "生物"
  }
]
```

**响应**: `204 No Content`

---

### 9. 删除单个教师

**端点**: `DELETE /teachers/{id}`

**响应**:
```json
{
  "status": "Teacher successfully deleted",
  "id": 1
}
```

---

### 10. 批量删除教师

**端点**: `DELETE /teachers`

**请求体**:
```json
[1, 2, 3]
```

**响应**:
```json
{
  "status": "teachers successfully deleted",
  "deleted_ids": [1, 2, 3]
}
```

---

## 👤 管理员 API

### 数据模型

```json
{
  "id": 1,
  "first_name": "管理",
  "last_name": "员",
  "email": "admin@example.com",
  "username": "admin",
  "role": "admin",
  "inactive_status": false,
  "user_created_at": "2024-01-01T00:00:00Z"
}
```

### 1. 获取所有管理员

**端点**: `GET /execs`

**查询参数**: 支持分页、过滤、排序

**响应**:
```json
{
  "status": "success",
  "count": 5,
  "data": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  ]
}
```

---

### 2. 获取单个管理员

**端点**: `GET /execs/{id}`

**响应**:
```json
{
  "id": 1,
  "first_name": "管理",
  "last_name": "员",
  "email": "admin@example.com",
  "username": "admin",
  "role": "admin"
}
```

---

### 3. 创建管理员

**端点**: `POST /execs`

**请求体**:
```json
[
  {
    "first_name": "新",
    "last_name": "管理员",
    "email": "newadmin@example.com",
    "username": "newadmin",
    "password": "securepassword123",
    "role": "user"
  }
]
```

**响应**:
```json
{
  "status": "success",
  "count": 1,
  "data": [
    {
      "id": 2,
      "username": "newadmin",
      "email": "newadmin@example.com",
      "role": "user"
    }
  ]
}
```

---

### 4. 更新管理员（部分更新）

**端点**: `PATCH /execs/{id}`

**请求体**:
```json
{
  "email": "updated@example.com"
}
```

---

### 5. 批量部分更新管理员

**端点**: `PATCH /execs`

**请求体**:
```json
[
  {
    "id": 1,
    "role": "superadmin"
  }
]
```

**响应**: `204 No Content`

---

### 6. 删除管理员

**端点**: `DELETE /execs/{id}`

**响应**:
```json
{
  "status": "Exec successfully deleted",
  "id": 1
}
```

---

## ❌ 错误处理

### 常见错误响应

**400 Bad Request**
```json
{
  "error": "Invalid request payload"
}
```

**401 Unauthorized**
```json
{
  "error": "Invalid credentials"
}
```

**404 Not Found**
```json
{
  "error": "Resource not found"
}
```

**500 Internal Server Error**
```json
{
  "error": "Internal server error"
}
```

---

## 🔍 查询参数

### 分页参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | int | 1 | 页码 |
| `limit` | int | 10 | 每页数量 |

**示例**:
```
GET /students?page=2&limit=20
```

---

### 过滤参数

支持按任意字段过滤，多个过滤条件使用 AND 逻辑：

**示例**:
```
GET /students?first_name=张&class=A
GET /teachers?subject=数学&class=B
```

---

### 排序参数

使用 `sortby` 参数，格式为 `field:order`：

- `order`: `ASC`（升序）或 `DESC`（降序）
- 支持多字段排序

**示例**:
```
GET /students?sortby=last_name:ASC
GET /students?sortby=class:DESC&sortby=last_name:ASC
```

---

### 组合查询示例

```bash
# 查询 A 班的学生，按姓氏升序排列，第 2 页，每页 20 条
curl -k "https://localhost:8443/students?class=A&sortby=last_name:ASC&page=2&limit=20"

# 查询数学老师，按班级降序排列
curl -k "https://localhost:8443/teachers?subject=数学&sortby=class:DESC"
```

---

## 📝 注意事项

1. **HTTPS 必需**: 所有请求必须使用 HTTPS
2. **Content-Type**: POST/PUT/PATCH 请求必须设置 `Content-Type: application/json`
3. **认证**: 部分端点需要在请求头中包含 `Authorization: Bearer <token>`
4. **批量操作**: 创建、更新、删除操作都支持批量处理
5. **字段验证**: 所有必填字段不能为空，邮箱必须唯一
6. **SQL 注入防护**: 所有查询使用参数化语句，ORDER BY 字段经过白名单验证

---

## 🔗 相关链接

- [README.md](./README.md) - 项目概述和快速开始
- [CLAUDE.md](./CLAUDE.md) - 项目架构和开发指南

---

**最后更新**: 2026-03-06

# REST API - Go

一个基于 Go 标准库构建的 RESTful API 项目，提供学生、教师和管理员的完整 CRUD 操作，支持 HTTPS、身份认证和密码管理功能。

## 📋 目录

- [特性](#特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [环境配置](#环境配置)
- [API 文档](#api-文档)
- [开发指南](#开发指南)

## ✨ 特性

- **RESTful API 设计** - 遵循 REST 规范的 API 端点
- **HTTPS 支持** - TLS 1.2+ 加密传输
- **身份认证系统** - JWT Token 认证，支持登录/登出
- **密码管理** - 密码加密存储、修改密码、忘记密码/重置密码流程
- **查询功能** - 支持过滤、排序、分页
- **批量操作** - 支持批量创建、更新、删除
- **安全中间件** - 安全头、CORS、速率限制、压缩等
- **参数验证** - 自动字段验证和 SQL 注入防护

## 🛠 技术栈

- **语言**: Go 1.19+
- **数据库**: MariaDB/MySQL
- **认证**: JWT (JSON Web Token)
- **加密**: bcrypt 密码哈希
- **HTTP**: Go 标准库 `net/http`
- **TLS**: 自签名证书或 CA 签发证书

## 📁 项目结构

```
rest-api-go/
├── cmd/
│   └── api/
│       └── server.go          # 应用入口
├── internal/
│   ├── api/
│   │   ├── handlers/          # HTTP 请求处理器
│   │   ├── middlewares/       # HTTP 中间件
│   │   └── router/            # 路由配置
│   ├── models/                # 数据模型
│   └── repository/
│       └── sqlconnect/        # 数据库操作
├── pkg/
│   └── utils/                 # 工具函数
├── cert.pem                   # TLS 证书
├── key.pem                    # TLS 私钥
├── .env                       # 环境变量配置
└── go.mod                     # Go 模块依赖
```

## 🚀 快速开始

### 前置要求

- Go 1.19 或更高版本
- MariaDB 或 MySQL 数据库
- OpenSSL (用于生成证书)

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd rest-api-go
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**

创建数据库并导入表结构：
```sql
CREATE DATABASE your_database;
USE your_database;

-- 创建学生表
CREATE TABLE students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    class VARCHAR(50)
);

-- 创建教师表
CREATE TABLE teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    class VARCHAR(50),
    subject VARCHAR(100)
);

-- 创建管理员表
CREATE TABLE execs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255) UNIQUE,
    username VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    password_changed_at DATETIME,
    user_created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    password_reset_token VARCHAR(255),
    password_token_expires DATETIME,
    inactive_status BOOLEAN DEFAULT FALSE,
    role VARCHAR(50) DEFAULT 'user'
);
```

4. **配置环境变量**

创建 `.env` 文件：
```env
API_PORT=:8443
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=your_database
HOST=127.0.0.1
DB_PORT=3306
JWT_SECRET=your_jwt_secret_key
```

5. **生成 TLS 证书**

使用 OpenSSL 生成自签名证书：
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

6. **运行服务器**
```bash
go run cmd/api/server.go
```

服务器将在 `https://localhost:8443` 启动。

## ⚙️ 环境配置

### 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `API_PORT` | API 服务端口 | `:8443` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_NAME` | 数据库名称 | - |
| `HOST` | 数据库主机地址 | `127.0.0.1` |
| `DB_PORT` | 数据库端口 | `3306` |
| `JWT_SECRET` | JWT 签名密钥 | - |

### 中间件配置

项目包含多个可选中间件（在 `cmd/api/server.go` 中配置）：

- **SecurityHeaders** - 安全响应头（默认启用）
- **CORS** - 跨域资源共享
- **RateLimiter** - 速率限制
- **HPP** - HTTP 参数污染防护
- **Compression** - Gzip 响应压缩
- **ResponseTime** - 响应时间记录

## 📖 API 文档

详细的 API 文档请参阅 [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)

### 快速示例

**获取所有学生（带分页）**
```bash
curl -k "https://localhost:8443/students?page=1&limit=10"
```

**创建学生**
```bash
curl -k -X POST https://localhost:8443/students \
  -H "Content-Type: application/json" \
  -d '[{
    "first_name": "张",
    "last_name": "三",
    "email": "zhangsan@example.com",
    "class": "A"
  }]'
```

**用户登录**
```bash
curl -k -X POST https://localhost:8443/execs/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

## 🔧 开发指南

### 代码格式化
```bash
go fmt ./...
```

### 静态检查
```bash
go vet ./...
```

### 运行测试
```bash
go test ./...
```

### 构建生产版本
```bash
go build -o api-server cmd/api/server.go
```

### 使用 Docker（可选）
```bash
# 构建镜像
docker build -t rest-api-go .

# 运行容器
docker run -p 8443:8443 --env-file .env rest-api-go
```

## 🔒 安全注意事项

1. **生产环境使用 CA 签发的证书**，不要使用自签名证书
2. **JWT_SECRET 必须使用强随机字符串**，建议 32 字符以上
3. **数据库密码不要使用弱密码**
4. **启用所有安全中间件**（CORS、RateLimiter、HPP 等）
5. **定期更新依赖**以修复安全漏洞

## 📝 数据库设计

### 核心表结构

- **students** - 学生信息表
- **teachers** - 教师信息表
- **execs** - 管理员/执行者表（包含认证信息）

### 关系说明

- 教师和学生通过 `class` 字段关联
- 支持查询某个教师的所有学生
- 支持统计某个教师的学生数量

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add some amazing feature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 📧 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件至 [your-email@example.com]

---

**注意**: 本项目仅供学习和开发使用，生产环境部署前请进行充分的安全审计和性能测试。

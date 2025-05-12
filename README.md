# 博客后端 API

这是一个使用 Go 语言和 Gin 框架开发的博客后端 API。

## 项目结构

```
blog-go/
├── config/          # 配置相关
├── controllers/     # 控制器
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由
├── database/        # 数据库相关
├── ent/             # Ent ORM生成的代码
├── .env             # 环境变量
├── go.mod           # Go模块文件
└── main.go          # 入口文件
```

## 功能

- 文章管理（增删改查）
- 标签管理
- 评论系统
- 用户身份验证
- 图片上传
- 友链管理
- 图书管理
- 一言管理
- 统计数据

## 开发环境

- Go 1.16+
- PostgreSQL 13+
- Gin
- Ent ORM

## 安装和运行

1. 克隆仓库

```bash
git clone <repository-url>
cd blog-go
```

2. 安装依赖

```bash
go mod download
```

3. 创建并配置 .env 文件

```
# 环境配置
ENV=development

# 服务器配置
PORT=8080

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=blog

# JWT配置
JWT_SECRET=your-secret-key-change-in-production
```

4. 运行项目

```bash
go run main.go
```

## API 文档

服务启动后，可以通过以下接口访问：

- `GET /api/posts`: 获取所有文章
- `GET /api/posts/:id`: 获取特定文章
- `POST /api/posts`: 创建文章
- `PUT /api/posts/:id`: 更新文章
- `DELETE /api/posts/:id`: 删除文章

更多接口详情请参考 [API 文档](./API.md)。

## 数据库

项目使用 PostgreSQL 数据库，需要预先创建名为 `blog` 的数据库。

## 贡献

欢迎提交 Issues 和 Pull Requests。 
# 校园二手交易平台后端

## 技术栈
- Gin 作为 web 框架
- Gorm 作为 ORM 框架
- PostgreSQL 作为数据库
- JWT 作为认证
- 配置直接写死在代码中

## 项目结构
- handler 负责处理前端请求, 校验参数, 定义路由
- service 负责业务逻辑, 调用 gorm 进行数据库操作
- model 定义数据库模型, 包含 gorm 的字段标签
- middleware 负责 JWT 认证

## 启动
1. 确保 PostgreSQL 运行中，创建数据库 `campus_secondhand`
2. 配置 main.go 中的 dsn 连接信息
3. 运行 `go run main.go`

## API 列表

### 认证
- POST /api/auth/register - 注册
- POST /api/auth/login - 登录
- POST /api/auth/wx-login - 微信登录

### 商品
- GET /api/goods - 商品列表
- GET /api/goods/:id - 商品详情
- POST /api/goods - 发布商品（需认证）
- PUT /api/goods/:id - 更新商品（需认证）
- DELETE /api/goods/:id - 删除商品（需认证）

### 订单
- POST /api/orders - 创建订单（需认证）
- GET /api/orders - 我的订单（需认证）
- GET /api/orders/:id - 订单详情（需认证）
- PUT /api/orders/:id/status - 更新订单状态（需认证）

### 用户
- GET /api/users/me - 个人信息（需认证）
- PUT /api/users/me - 更新个人信息（需认证）

# 校园二手交易平台后端设计

## 技术栈

- **Gin** — Web 框架
- **Gorm** — ORM
- **PostgreSQL** — 数据库
- **JWT** — 认证
- **phuslu/log** — 日志
- 配置写死在代码中，无配置文件

## 项目结构

```
backend/
├── main.go
├── go.mod
├── handler/          # HTTP 处理，参数校验，路由定义
│   ├── user.go
│   ├── goods.go
│   └── order.go
├── service/         # 业务逻辑
│   ├── user.go
│   ├── goods.go
│   └── order.go
├── model/           # 数据库模型
│   ├── user.go
│   ├── goods.go
│   └── order.go
└── middleware/      # JWT 认证中间件
    └── auth.go
```

## 数据模型

### User
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| openid | string | 微信 openid（可为空） |
| nickname | string | 昵称 |
| avatar | string | 头像 URL |
| phone | string | 手机号 |
| role | string | buyer/seller/admin |
| password | string | 密码 |
| created_at | time | 创建时间 |

### Goods
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| title | string | 商品标题 |
| description | string | 商品描述 |
| price | float64 | 售价 |
| original_price | float64 | 原价 |
| category | string | 分类 |
| condition | string | 新旧程度 |
| images | string (JSON) | 图片 URL 数组 |
| seller_id | uint | 卖家 ID |
| status | string | on_sale/sold/off_shelf |
| view_count | int | 浏览次数 |
| created_at | time | 创建时间 |

### Order
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| goods_id | uint | 商品 ID |
| buyer_id | uint | 买家 ID |
| seller_id | uint | 卖家 ID |
| status | string | pending/confirmed/cancelled/completed |
| remark | string | 备注 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

## API 设计

### 认证
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 注册 | 否 |
| POST | /api/auth/login | 登录 | 否 |

### 商品
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/goods | 列表（关键词/分类/价格筛选） | 否 |
| GET | /api/goods/:id | 商品详情 | 否 |
| POST | /api/goods | 发布商品 | 是 |
| PUT | /api/goods/:id | 更新商品 | 是（仅卖家本人） |
| DELETE | /api/goods/:id | 删除商品 | 是（仅卖家本人或admin） |

### 订单
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/orders | 创建订单 | 是 |
| GET | /api/orders | 查我的订单（buyer或seller） | 是 |
| PUT | /api/orders/:id/status | 更新订单状态 | 是 |

### 用户
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/users/me | 个人信息 | 是 |
| PUT | /api/users/me | 更新个人信息 | 是 |

## 认证流程

- 注册/登录返回 JWT token
- 客户端在 Header 携带 `Authorization: Bearer <token>`
- JWT 包含 user_id 和 role
- role 标识：buyer/seller/admin

## 错误处理

- 简单返回错误信息，HTTP 状态码区分错误类型
- 400: 参数错误
- 401: 未认证
- 403: 无权限
- 404: 资源不存在
- 500: 服务器错误

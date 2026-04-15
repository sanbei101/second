# 校园二手交易平台后端实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个精简的校园二手交易平台 Go 后端，支持用户认证、商品管理和订单管理。

**Architecture:** 使用 Gin + Gorm + PostgreSQL + JWT。handler 层处理 HTTP 请求和参数校验，service 层封装业务逻辑，model 层定义数据库模型，middleware 层处理 JWT 认证。

**Tech Stack:** Go 1.26, Gin, Gorm, PostgreSQL, JWT (golang-jwt/jwt/v5), phuslu/log

---

## 文件结构

```
backend/
├── main.go                        # 入口，初始化数据库、注册路由、启动服务
├── go.mod                         # 依赖：gin, gorm, pg, jwt, log
├── model/
│   ├── user.go                    # User 模型
│   ├── goods.go                   # Goods 模型
│   └── order.go                   # Order 模型
├── service/
│   ├── user.go                    # 用户注册/登录/信息更新
│   ├── goods.go                   # 商品 CRUD
│   └── order.go                   # 订单 CRUD
├── handler/
│   ├── user.go                    # 用户 HTTP 处理
│   ├── goods.go                   # 商品 HTTP 处理
│   └── order.go                   # 订单 HTTP 处理
└── middleware/
    └── auth.go                    # JWT 认证中间件
```

---

## Task 1: 初始化 go.mod 依赖

**Files:**
- Modify: `backend/go.mod`

- [ ] **Step 1: 更新 go.mod 添加依赖**

```go
module github.com/sanbei101/second

go 1.26.2

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/phuslu/log v1.0.124
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)
```

Run: `cd backend && go mod tidy`
Expected: 下载所有依赖，无报错

- [ ] **Step 2: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "deps: add gin, gorm, pg, jwt dependencies"
```

---

## Task 2: 创建数据模型

**Files:**
- Create: `backend/model/user.go`
- Create: `backend/model/goods.go`
- Create: `backend/model/order.go`

- [ ] **Step 1: 创建 User 模型**

`backend/model/user.go`:

```go
package model

import "time"

type UserRole string

const (
	RoleBuyer  UserRole = "buyer"
	RoleSeller UserRole = "seller"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Openid    string    `gorm:"size:128;index" json:"openid"`    // 微信 openid
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Avatar    string    `gorm:"size:256" json:"avatar"`
	Phone     string    `gorm:"size:16;index" json:"phone"`
	Role      UserRole  `gorm:"size:16;default:buyer" json:"role"`
	Password  string    `gorm:"size:64" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}
```

- [ ] **Step 2: 创建 Goods 模型**

`backend/model/goods.go`:

```go
package model

import "time"

type GoodsStatus string

const (
	GoodsOnSale   GoodsStatus = "on_sale"
	GoodsSold     GoodsStatus = "sold"
	GoodsOffShelf GoodsStatus = "off_shelf"
)

type Goods struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Title        string      `gorm:"size:128" json:"title"`
	Description  string      `gorm:"type:text" json:"description"`
	Price        float64     `gorm:"type:decimal(10,2)" json:"price"`
	OriginalPrice float64    `gorm:"type:decimal(10,2)" json:"originalPrice"`
	Category     string      `gorm:"size:32" json:"category"`
	Condition    string      `gorm:"size:16" json:"condition"`
	Images       string      `gorm:"type:text" json:"images"` // JSON array of URLs
	SellerID     uint        `gorm:"index" json:"sellerId"`
	Seller       *User       `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Status       GoodsStatus `gorm:"size:16;default:on_sale" json:"status"`
	ViewCount    int         `gorm:"default:0" json:"viewCount"`
	CreatedAt    time.Time   `json:"createdAt"`
}
```

- [ ] **Step 3: 创建 Order 模型**

`backend/model/order.go`:

```go
package model

import "time"

type OrderStatus string

const (
	OrderPending    OrderStatus = "pending"
	OrderConfirmed  OrderStatus = "confirmed"
	OrderCancelled OrderStatus = "cancelled"
	OrderCompleted OrderStatus = "completed"
)

type Order struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	GoodsID   uint        `gorm:"index" json:"goodsId"`
	Goods     *Goods      `gorm:"foreignKey:GoodsID" json:"goods,omitempty"`
	BuyerID   uint        `gorm:"index" json:"buyerId"`
	Buyer     *User       `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	SellerID  uint        `gorm:"index" json:"sellerId"`
	Seller    *User       `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Status    OrderStatus `gorm:"size:16;default:pending" json:"status"`
	Remark    string      `gorm:"size:256" json:"remark"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/model/user.go backend/model/goods.go backend/model/order.go
git commit -m "model: add User, Goods, Order models"
```

---

## Task 3: 创建 JWT 认证中间件

**Files:**
- Create: `backend/middleware/auth.go`

- [ ] **Step 1: 创建 JWT 中间件**

`backend/middleware/auth.go`:

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("campus-secondhand-secret-key")

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	id, _ := c.Get("userID")
	return id.(uint)
}

func GetRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/middleware/auth.go
git commit -m "middleware: add JWT auth middleware"
```

---

## Task 4: 创建 Service 层

**Files:**
- Create: `backend/service/user.go`
- Create: `backend/service/goods.go`
- Create: `backend/service/order.go`

- [ ] **Step 1: 创建 user service**

`backend/service/user.go`:

```go
package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"campus-secondhand/model"
	"campus-secondhand/middleware"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(phone, password, nickname string, role model.UserRole) (*model.User, string, error) {
	var exist model.User
	if err := s.db.Where("phone = ?", phone).First(&exist).Error; err == nil {
		return nil, "", errors.New("phone already registered")
	}

	user := model.User{
		Phone:    phone,
		Password: password, // 直接存储密码，生产环境应哈希
		Nickname: nickname,
		Role:     role,
		Avatar:   "https://img.yzcdn.cn/vant/cat.jpeg",
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(&user)
	return &user, token, err
}

func (s *UserService) Login(phone, password string) (*model.User, string, error) {
	var user model.User
	if err := s.db.Where("phone = ? AND password = ?", phone, password).First(&user).Error; err != nil {
		return nil, "", errors.New("invalid phone or password")
	}

	token, err := s.generateToken(&user)
	return &user, token, err
}

func (s *UserService) WxLogin(openid string, role model.UserRole) (*model.User, string, error) {
	var user model.User
	err := s.db.Where("openid = ?", openid).First(&user).Error
	if err == nil {
		token, err := s.generateToken(&user)
		return &user, token, err
	}

	user = model.User{
		Openid:   openid,
		Nickname: "微信用户",
		Role:     role,
		Avatar:   "https://img.yzcdn.cn/vant/cat.jpeg",
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(&user)
	return &user, token, err
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(id uint, nickname, avatar, phone string) error {
	updates := map[string]interface{}{}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}
	if phone != "" {
		updates["phone"] = phone
	}
	return s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (s *UserService) generateToken(user *model.User) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTSecret)
}
```

- [ ] **Step 2: 创建 goods service**

`backend/service/goods.go`:

```go
package service

import (
	"errors"

	"gorm.io/gorm"

	"campus-secondhand/model"
)

type GoodsService struct {
	db *gorm.DB
}

func NewGoodsService(db *gorm.DB) *GoodsService {
	return &GoodsService{db: db}
}

func (s *GoodsService) Create(goods *model.Goods) error {
	return s.db.Create(goods).Error
}

func (s *GoodsService) GetByID(id uint) (*model.Goods, error) {
	var goods model.Goods
	if err := s.db.Preload("Seller").First(&goods, id).Error; err != nil {
		return nil, err
	}
	// 增加浏览次数
	s.db.Model(&goods).Update("view_count", gorm.Expr("view_count + 1"))
	return &goods, nil
}

func (s *GoodsService) List(keyword, category string, minPrice, maxPrice *float64) ([]model.Goods, error) {
	query := s.db.Model(&model.Goods{}).Preload("Seller").Where("status = ?", model.GoodsOnSale)

	if keyword != "" {
		query = query.Where("title ILIKE ?", "%"+keyword+"%")
	}
	if category != "" && category != "全部" {
		query = query.Where("category = ?", category)
	}
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}

	var goods []model.Goods
	if err := query.Order("created_at DESC").Find(&goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

func (s *GoodsService) GetBySeller(sellerID uint) ([]model.Goods, error) {
	var goods []model.Goods
	if err := s.db.Where("seller_id = ?", sellerID).Order("created_at DESC").Find(&goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

func (s *GoodsService) Update(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.Goods{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("goods not found")
	}
	return nil
}

func (s *GoodsService) Delete(id uint) error {
	return s.db.Delete(&model.Goods{}, id).Error
}

func (s *GoodsService) UpdateStatus(id uint, status model.GoodsStatus) error {
	return s.Update(id, map[string]interface{}{"status": status})
}
```

- [ ] **Step 3: 创建 order service**

`backend/service/order.go`:

```go
package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"campus-secondhand/model"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) Create(goodsID, buyerID, sellerID uint, remark string) (*model.Order, error) {
	// 检查商品是否存在且在售
	var goods model.Goods
	if err := s.db.First(&goods, goodsID).Error; err != nil {
		return nil, errors.New("goods not found")
	}
	if goods.Status != model.GoodsOnSale {
		return nil, errors.New("goods not available")
	}
	if buyerID == sellerID {
		return nil, errors.New("cannot buy your own goods")
	}

	// 创建订单
	order := &model.Order{
		GoodsID:   goodsID,
		BuyerID:   buyerID,
		SellerID:  sellerID,
		Status:    model.OrderPending,
		Remark:    remark,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.db.Create(order).Error; err != nil {
		return nil, err
	}

	// 将商品状态改为已售
	s.db.Model(&goods).Update("status", model.GoodsSold)

	return order, nil
}

func (s *OrderService) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := s.db.Preload("Goods").Preload("Buyer").Preload("Seller").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetByUser(userID uint, asBuyer bool) ([]model.Order, error) {
	var orders []model.Order
	query := s.db.Preload("Goods").Preload("Buyer").Preload("Seller")
	if asBuyer {
		query = query.Where("buyer_id = ?", userID)
	} else {
		query = query.Where("seller_id = ?", userID)
	}
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) UpdateStatus(id uint, status model.OrderStatus) error {
	result := s.db.Model(&model.Order{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/service/user.go backend/service/goods.go backend/service/order.go
git commit -m "service: add user, goods, order services"
```

---

## Task 5: 创建 Handler 层

**Files:**
- Create: `backend/handler/user.go`
- Create: `backend/handler/goods.go`
- Create: `backend/handler/order.go`

- [ ] **Step 1: 创建 user handler**

`backend/handler/user.go`:

```go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"campus-secondhand/middleware"
	"campus-secondhand/model"
	"campus-secondhand/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type RegisterReq struct {
	Phone    string           `json:"phone" binding:"required"`
	Password string           `json:"password" binding:"required"`
	Nickname string           `json:"nickname"`
	Role     model.UserRole   `json:"role" binding:"required"`
}

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type WxLoginReq struct {
	Openid string           `json:"openid" binding:"required"`
	Role   model.UserRole   `json:"role"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = model.RoleBuyer
	}

	user, token, err := h.svc.Register(req.Phone, req.Password, req.Nickname, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.svc.Login(req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *UserHandler) WxLogin(c *gin.Context) {
	var req WxLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = model.RoleBuyer
	}

	user, token, err := h.svc.WxLogin(req.Openid, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.svc.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.svc.Update(userID, req.Nickname, req.Avatar, req.Phone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, _ := h.svc.GetByID(userID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup, svc *service.UserService) {
	h.svc = svc
	auth := rg.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/wx-login", h.WxLogin)

	users := rg.Group("/users")
	users.Use(middleware.AuthRequired())
	users.GET("/me", h.Me)
	users.PUT("/me", h.UpdateProfile)
}
```

- [ ] **Step 2: 创建 goods handler**

`backend/handler/goods.go`:

```go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-secondhand/middleware"
	"campus-secondhand/model"
	"campus-secondhand/service"
)

type GoodsHandler struct {
	svc *service.GoodsService
}

func NewGoodsHandler(svc *service.GoodsService) *GoodsHandler {
	return &GoodsHandler{svc: svc}
}

type CreateGoodsReq struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	Price        float64  `json:"price" binding:"required"`
	OriginalPrice float64 `json:"originalPrice"`
	Category     string   `json:"category" binding:"required"`
	Condition    string   `json:"condition" binding:"required"`
	Images       []string `json:"images"`
}

func (h *GoodsHandler) Create(c *gin.Context) {
	var req CreateGoodsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imagesJSON, _ := json.Marshal(req.Images)
	goods := &model.Goods{
		Title:         req.Title,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Category:      req.Category,
		Condition:     req.Condition,
		Images:        string(imagesJSON),
		SellerID:      middleware.GetUserID(c),
		Status:        model.GoodsOnSale,
		ViewCount:     0,
	}

	if err := h.svc.Create(goods); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goods})
}

func (h *GoodsHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	var minPrice, maxPrice *float64

	if mp := c.Query("minPrice"); mp != "" {
		v, _ := strconv.ParseFloat(mp, 64)
		minPrice = &v
	}
	if mp := c.Query("maxPrice"); mp != "" {
		v, _ := strconv.ParseFloat(mp, 64)
		maxPrice = &v
	}

	goods, err := h.svc.List(keyword, category, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 解析 images JSON
	type GoodsResponse struct {
		model.Goods
		Images []string `json:"images"`
	}
	result := make([]GoodsResponse, len(goods))
	for i, g := range goods {
		result[i] = GoodsResponse{Goods: g}
		json.Unmarshal([]byte(g.Images), &result[i].Images)
	}

	c.JSON(http.StatusOK, gin.H{"goods": result})
}

func (h *GoodsHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	goods, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goods not found"})
		return
	}

	type GoodsResponse struct {
		model.Goods
		Images []string `json:"images"`
	}
	var result GoodsResponse
	result.Goods = *goods
	json.Unmarshal([]byte(goods.Images), &result.Images)

	c.JSON(http.StatusOK, gin.H{"goods": result})
}

func (h *GoodsHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	goods, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goods not found"})
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	if goods.SellerID != userID && role != string(model.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req CreateGoodsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imagesJSON, _ := json.Marshal(req.Images)
	updates := map[string]interface{}{
		"title":          req.Title,
		"description":    req.Description,
		"price":          req.Price,
		"original_price": req.OriginalPrice,
		"category":       req.Category,
		"condition":      req.Condition,
		"images":         string(imagesJSON),
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *GoodsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	goods, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goods not found"})
		return
	}

	userID := middleware.GetUserID(c)
	role := middleware.GetRole(c)
	if goods.SellerID != userID && role != string(model.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *GoodsHandler) RegisterRoutes(rg *gin.RouterGroup, svc *service.GoodsService) {
	h.svc = svc
	goods := rg.Group("/goods")
	goods.GET("", h.List)
	goods.GET("/:id", h.GetByID)
	goods.POST("", middleware.AuthRequired(), h.Create)
	goods.PUT("/:id", middleware.AuthRequired(), h.Update)
	goods.DELETE("/:id", middleware.AuthRequired(), h.Delete)
}
```

- [ ] **Step 3: 创建 order handler**

`backend/handler/order.go`:

```go
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"campus-secondhand/middleware"
	"campus-secondhand/model"
	"campus-secondhand/service"
)

type OrderHandler struct {
	svc       *service.OrderService
	goodsSvc  *service.GoodsService
}

func NewOrderHandler(svc *service.OrderService, goodsSvc *service.GoodsService) *OrderHandler {
	return &OrderHandler{svc: svc, goodsSvc: goodsSvc}
}

type CreateOrderReq struct {
	GoodsID  uint   `json:"goodsId" binding:"required"`
	Remark   string `json:"remark"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buyerID := middleware.GetUserID(c)

	// 获取商品信息以确定卖家
	goods, err := h.goodsSvc.GetByID(req.GoodsID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "goods not found"})
		return
	}

	order, err := h.svc.Create(req.GoodsID, buyerID, goods.SellerID, req.Remark)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *OrderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := c.DefaultQuery("role", "buyer") // buyer or seller

	asBuyer := role == "buyer"
	orders, err := h.svc.GetByUser(userID, asBuyer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	order, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	userID := middleware.GetUserID(c)
	if order.BuyerID != userID && order.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

type UpdateStatusReq struct {
	Status model.OrderStatus `json:"status" binding:"required"`
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	userID := middleware.GetUserID(c)
	if order.BuyerID != userID && order.SellerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := h.svc.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *OrderHandler) RegisterRoutes(rg *gin.RouterGroup, svc *service.OrderService, goodsSvc *service.GoodsService) {
	h.svc = svc
	h.goodsSvc = goodsSvc
	orders := rg.Group("/orders")
	orders.Use(middleware.AuthRequired())
	orders.POST("", h.Create)
	orders.GET("", h.List)
	orders.GET("/:id", h.GetByID)
	orders.PUT("/:id/status", h.UpdateStatus)
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/handler/user.go backend/handler/goods.go backend/handler/order.go
git commit -m "handler: add user, goods, order HTTP handlers"
```

---

## Task 6: 编写 main.go 并启动服务

**Files:**
- Modify: `backend/main.go`

- [ ] **Step 1: 编写 main.go**

`backend/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"campus-secondhand/handler"
	"campus-secondhand/middleware"
	"campus-secondhand/model"
	"campus-secondhand/service"
)

func main() {
	log.DefaultLogger = log.Logger{
		Level:  log.InfoLevel,
		Writer: &log.IOWriter{Writer: os.Stderr},
	}

	// 数据库配置（写死在代码中）
	dsn := "host=localhost user=postgres password=postgres dbname=campus_secondhand port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}, &model.Goods{}, &model.Order{}); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	log.Info().Msg("database connected and migrated")

	// 初始化服务
	userSvc := service.NewUserService(db)
	goodsSvc := service.NewGoodsService(db)
	orderSvc := service.NewOrderService(db)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userSvc)
	goodsHandler := handler.NewGoodsHandler(goodsSvc)
	orderHandler := handler.NewOrderHandler(orderSvc, goodsSvc)

	// Gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册路由
	api := r.Group("/api")
	userHandler.RegisterRoutes(api, userSvc)
	goodsHandler.RegisterRoutes(api, goodsSvc)
	orderHandler.RegisterRoutes(api, orderSvc, goodsSvc)

	// 静态文件（可选，用于图片等）
	// r.Static("/uploads", "./uploads")

	port := 8080
	log.Info().Msgf("server starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
```

- [ ] **Step 2: 更新 README.md**

`backend/README.md`:

```markdown
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
```

- [ ] **Step 3: Commit**

```bash
git add backend/main.go backend/README.md
git commit -m "main: wire up all components and start server"
```

---

## Task 7: 验证编译

- [ ] **Step 1: 运行 go mod tidy 和编译**

Run: `cd backend && go mod tidy && go build ./...`
Expected: 无错误输出

- [ ] **Step 2: Commit**

```bash
git add backend/go.sum
git commit -m "build: ensure all dependencies downloaded"
```

---

## 依赖汇总

| 包 | 用途 |
|----|------|
| github.com/gin-gonic/gin v1.10.0 | Web 框架 |
| github.com/golang-jwt/jwt/v5 v5.2.1 | JWT 认证 |
| github.com/phuslu/log v1.0.124 | 日志 |
| gorm.io/driver/postgres v1.5.9 | PostgreSQL 驱动 |
| gorm.io/gorm v1.25.12 | ORM |

## 自检清单

- [ ] go.mod 包含所有依赖
- [ ] model/user.go, goods.go, order.go 三个模型完整
- [ ] middleware/auth.go JWT 中间件正确解析 token
- [ ] service 层调用 db 操作数据库
- [ ] handler 层正确绑定参数和返回 JSON
- [ ] main.go 初始化 db、services、handlers 并注册路由
- [ ] 所有文件通过 go build ./... 编译通过

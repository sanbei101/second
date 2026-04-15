package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sanbei101/second/middleware"
	"github.com/sanbei101/second/model"
	"github.com/sanbei101/second/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type RegisterReq struct {
	Phone    string         `json:"phone" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Nickname string         `json:"nickname"`
	Role     model.UserRole `json:"role" binding:"required"`
}

type LoginReq struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type WxLoginReq struct {
	Openid string         `json:"openid" binding:"required"`
	Role   model.UserRole `json:"role"`
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

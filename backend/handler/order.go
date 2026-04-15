package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"

	"github.com/sanbei101/second/middleware"
	"github.com/sanbei101/second/model"
	"github.com/sanbei101/second/service"
)

type OrderHandler struct {
	svc      *service.OrderService
	goodsSvc *service.GoodsService
}

func NewOrderHandler(svc *service.OrderService, goodsSvc *service.GoodsService) *OrderHandler {
	return &OrderHandler{svc: svc, goodsSvc: goodsSvc}
}

type CreateOrderReq struct {
	GoodsID uint   `json:"goodsId" binding:"required"`
	Remark  string `json:"remark"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buyerID := middleware.GetUserID(c)

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

	log.Info().Uint("order_id", order.ID).Uint("buyer_id", buyerID).Msg("order created")
	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *OrderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := c.DefaultQuery("role", "buyer")

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

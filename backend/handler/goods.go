package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"

	"github.com/sanbei101/second/middleware"
	"github.com/sanbei101/second/model"
	"github.com/sanbei101/second/service"
)

type GoodsHandler struct {
	svc *service.GoodsService
}

func NewGoodsHandler(svc *service.GoodsService) *GoodsHandler {
	return &GoodsHandler{svc: svc}
}

type CreateGoodsReq struct {
	Title         string   `json:"title" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required"`
	OriginalPrice float64  `json:"originalPrice"`
	Category      string   `json:"category" binding:"required"`
	Condition     string   `json:"condition" binding:"required"`
	Images        []string `json:"images"`
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

	log.Info().Str("title", goods.Title).Uint("seller_id", goods.SellerID).Msg("goods created")
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

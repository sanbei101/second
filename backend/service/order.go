package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/sanbei101/second/model"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) Create(goodsID, buyerID, sellerID uint, remark string) (*model.Order, error) {
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

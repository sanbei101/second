package service

import (
	"errors"
	"time"

	"github.com/phuslu/log"
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
		log.Warn().Err(err).Uint("goods_id", goodsID).Msg("create order failed: goods not found")
		return nil, errors.New("goods not found")
	}
	if goods.Status != model.GoodsOnSale {
		log.Warn().Uint("goods_id", goodsID).Str("status", string(goods.Status)).Msg("create order failed: goods not available")
		return nil, errors.New("goods not available")
	}
	if buyerID == sellerID {
		log.Warn().Uint("buyer_id", buyerID).Uint("seller_id", sellerID).Msg("create order failed: self-buy")
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
		log.Error().Err(err).Uint("goods_id", goodsID).Uint("buyer_id", buyerID).Msg("create order failed: db error")
		return nil, err
	}

	s.db.Model(&goods).Update("status", model.GoodsSold)
	log.Info().Uint("order_id", order.ID).Uint("goods_id", goodsID).Uint("buyer_id", buyerID).Msg("order created")
	return order, nil
}

func (s *OrderService) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := s.db.Preload("Goods").Preload("Buyer").Preload("Seller").First(&order, id).Error; err != nil {
		log.Warn().Err(err).Uint("order_id", id).Msg("get order failed")
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
		log.Error().Err(err).Uint("user_id", userID).Msg("get orders failed")
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
		log.Error().Err(result.Error).Uint("order_id", id).Str("status", string(status)).Msg("update order status failed")
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Warn().Uint("order_id", id).Msg("update order status: not found")
		return errors.New("order not found")
	}
	log.Info().Uint("order_id", id).Str("status", string(status)).Msg("order status updated")
	return nil
}

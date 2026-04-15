package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/sanbei101/second/model"
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

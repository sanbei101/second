package service

import (
	"errors"

	"github.com/phuslu/log"
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
	if err := s.db.Create(goods).Error; err != nil {
		log.Error().Err(err).Str("title", goods.Title).Msg("create goods failed")
		return err
	}
	log.Info().Uint("goods_id", goods.ID).Str("title", goods.Title).Uint("seller_id", goods.SellerID).Msg("goods created")
	return nil
}

func (s *GoodsService) GetByID(id uint) (*model.Goods, error) {
	var goods model.Goods
	if err := s.db.Preload("Seller").First(&goods, id).Error; err != nil {
		log.Warn().Err(err).Uint("goods_id", id).Msg("get goods failed")
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
		log.Error().Err(err).Msg("list goods failed")
		return nil, err
	}
	log.Info().Int("count", len(goods)).Str("keyword", keyword).Str("category", category).Msg("goods listed")
	return goods, nil
}

func (s *GoodsService) GetBySeller(sellerID uint) ([]model.Goods, error) {
	var goods []model.Goods
	if err := s.db.Where("seller_id = ?", sellerID).Order("created_at DESC").Find(&goods).Error; err != nil {
		log.Error().Err(err).Uint("seller_id", sellerID).Msg("get goods by seller failed")
		return nil, err
	}
	return goods, nil
}

func (s *GoodsService) Update(id uint, updates map[string]any) error {
	result := s.db.Model(&model.Goods{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		log.Error().Err(result.Error).Uint("goods_id", id).Msg("update goods failed")
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Warn().Uint("goods_id", id).Msg("update goods: not found")
		return errors.New("goods not found")
	}
	log.Info().Uint("goods_id", id).Msg("goods updated")
	return nil
}

func (s *GoodsService) Delete(id uint) error {
	if err := s.db.Delete(&model.Goods{}, id).Error; err != nil {
		log.Error().Err(err).Uint("goods_id", id).Msg("delete goods failed")
		return err
	}
	log.Info().Uint("goods_id", id).Msg("goods deleted")
	return nil
}

func (s *GoodsService) UpdateStatus(id uint, status model.GoodsStatus) error {
	return s.Update(id, map[string]any{"status": status})
}

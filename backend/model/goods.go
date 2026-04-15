package model

import "time"

type GoodsStatus string

const (
	GoodsOnSale   GoodsStatus = "on_sale"
	GoodsSold     GoodsStatus = "sold"
	GoodsOffShelf GoodsStatus = "off_shelf"
)

type Goods struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	Title         string      `gorm:"size:128" json:"title"`
	Description   string      `gorm:"type:text" json:"description"`
	Price         float64     `gorm:"type:decimal(10,2)" json:"price"`
	OriginalPrice float64     `gorm:"type:decimal(10,2)" json:"originalPrice"`
	Category      string      `gorm:"size:32" json:"category"`
	Condition     string      `gorm:"size:16" json:"condition"`
	Images        string      `gorm:"type:text" json:"images"`
	SellerID      uint        `gorm:"index" json:"sellerId"`
	Seller        *User       `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Status        GoodsStatus `gorm:"size:16;default:on_sale" json:"status"`
	ViewCount     int         `gorm:"default:0" json:"viewCount"`
	CreatedAt     time.Time   `json:"createdAt"`
}

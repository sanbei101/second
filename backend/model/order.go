package model

import "time"

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderConfirmed OrderStatus = "confirmed"
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

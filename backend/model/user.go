package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Openid    string    `gorm:"size:128;index" json:"openid"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Avatar    string    `gorm:"size:256" json:"avatar"`
	Phone     string    `gorm:"size:16;index" json:"phone"`
	Password  string    `gorm:"size:64" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

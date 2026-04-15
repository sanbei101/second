package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/phuslu/log"
	"gorm.io/gorm"

	"github.com/sanbei101/second/model"
)

var (
	categories  = []string{"数码", "书籍", "服装", "美妆", "运动", "家居", "票务", "其他"}
	conditions  = []string{"全新", "九九新", "九成新", "八成新", "七成新"}
	goodsTitles = []string{
		"iPhone 14 Pro", "MacBook Air M2", "AirPods Pro 2", "iPad mini 6", "Kindle Paperwhite",
		"高等数学教材", "考研英语词汇", "C++ Primer", "数据结构", "算法导论",
		"优衣库羽绒服", "Nike 运动鞋", "Adidas 背包", "ZARA 外套", "HM 牛仔裤",
		"SK-II 神仙水", "YSL 口红", "Dior 粉底液", "Lancome 眼霜", "Estee Lauder 小棕瓶",
		"瑜伽垫", "哑铃套装", "羽毛球拍", "篮球", "跳绳",
		"宜家台灯", "收纳箱", "抱枕", "热水壶", "吹风机",
		"音乐节门票", "电影票", "健身房月卡", "游泳馆次卡", "演唱会门票",
		"二手自行车", "充电宝", "数据线", "鼠标键盘", "显示器",
	}
	images = []string{
		"https://images.unsplash.com/photo-1621381571163-18b89a7bd543?w=800",
		"https://images.unsplash.com/photo-1604364184592-82356ca9d7a3?w=800",
		"https://images.unsplash.com/photo-1602324865142-f408d45769ae?w=800",
		"https://images.unsplash.com/photo-1677096603108-2721b9642b0a?w=800",
		"https://images.unsplash.com/photo-1653548410454-0d2b4098c4da?w=800",
		"https://images.unsplash.com/photo-1653548410467-85a75d900b0c?w=800",
		"https://images.unsplash.com/photo-1552990608-cfe5eb330b89?w=800",
		"https://images.unsplash.com/photo-1582078926742-1d84405a9283?w=800",
		"https://images.unsplash.com/photo-1636262513953-2351e8d1671b?w=800",
		"https://images.unsplash.com/photo-1629201973727-8aabeed95aad?w=800",
		"https://images.unsplash.com/photo-1654939438953-58dca903c976?w=800",
		"https://images.unsplash.com/photo-1673726803845-a76083e94712?w=800",
		"https://images.unsplash.com/photo-1627435601550-406e8683be82?w=800",
		"https://images.unsplash.com/photo-1611682011252-21667a85e3ac?w=800",
		"https://images.unsplash.com/photo-1599447291786-724cfd131568?w=800",
		"https://images.unsplash.com/photo-1759978257038-ff90be507a3d?w=800",
		"https://images.unsplash.com/photo-1758471995115-81c662cf949f?w=800",
		"https://images.unsplash.com/photo-1765608262875-2b97700ec852?w=800",
		"https://images.unsplash.com/photo-1559348349-86f1f65817fe?w=800",
		"https://images.unsplash.com/photo-1591334770599-d5df9b063fb9?w=800",
		"https://images.unsplash.com/photo-1637176472598-46133e5d47ce?w=800",
		"https://images.unsplash.com/photo-1592335509190-ac997442fbef?w=800",
		"https://images.unsplash.com/photo-1687096447510-a3301cfe7dd8?w=800",
		"https://images.unsplash.com/photo-1525825691042-e14d9042fc70?w=800",
		"https://images.unsplash.com/photo-1770064319289-ef144cde0055?w=800",
		"https://images.unsplash.com/photo-1746010387602-d22b9d6401f8?w=800",
		"https://images.unsplash.com/photo-1603809666145-1fecc56f49a5?w=800",
		"https://images.unsplash.com/photo-1677245359931-0bd5993ad965?w=800",
		"https://images.unsplash.com/photo-1668395099117-8ba402a86dc6?w=800",
		"https://images.unsplash.com/photo-1626589618584-677eacfcd265?w=800",
	}
)

func Migrate(db *gorm.DB) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 清空旧数据
	db.Exec("TRUNCATE TABLE orders, goods, users RESTART IDENTITY CASCADE")
	log.Info().Msg("old data truncated")

	// 1. 插入 10 个用户
	users := make([]model.User, 10)
	for i := range 10 {
		users[i] = model.User{
			Phone:    fmt.Sprintf("138%08d", 10000000+i),
			Password: "123456",
			Nickname: fmt.Sprintf("用户%d", i+1),
			Avatar:   "https://img.yzcdn.cn/vant/cat.jpeg",
		}
	}
	if err := db.Create(&users).Error; err != nil {
		log.Fatal().Err(err).Msg("create users failed")
	}
	log.Info().Int("count", len(users)).Msg("users created")

	userIDs := make([]uint, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	goodsBatch := make([]model.Goods, 1000)
	for i := range 1000 {
		title := goodsTitles[r.Intn(len(goodsTitles))]
		if r.Intn(2) == 0 {
			title += fmt.Sprintf(" %d", r.Intn(100))
		}
		price := float64(r.Intn(10000) + 10)
		origPrice := price * (1 + float64(r.Intn(50)+10)/100)

		imgJSON, _ := json.Marshal([]string{images[r.Intn(len(images))]})

		statuses := []model.GoodsStatus{model.GoodsOnSale, model.GoodsSold, model.GoodsOffShelf}
		goodsBatch[i] = model.Goods{
			Title:         title,
			Description:   fmt.Sprintf("这是一段商品描述,编号 %d,品相不错,欢迎咨询。", i+1),
			Price:         price,
			OriginalPrice: origPrice,
			Category:      categories[r.Intn(len(categories))],
			Condition:     conditions[r.Intn(len(conditions))],
			Images:        string(imgJSON),
			SellerID:      userIDs[r.Intn(len(userIDs))],
			Status:        statuses[r.Intn(len(statuses))],
			ViewCount:     r.Intn(500),
			CreatedAt:     time.Now().Add(-time.Duration(r.Intn(720)) * time.Hour),
		}
	}
	if err := db.Create(&goodsBatch).Error; err != nil {
		log.Fatal().Err(err).Msg("create goods failed")
	}
	log.Info().Int("count", len(goodsBatch)).Msg("goods created")

	goodsIDs := make([]uint, len(goodsBatch))
	for i, g := range goodsBatch {
		goodsIDs[i] = g.ID
	}

	// 3. 插入 100 个订单
	orders := make([]model.Order, 100)
	for i := range 100 {
		sellerID := userIDs[r.Intn(len(userIDs))]
		var buyerID uint
		for {
			buyerID = userIDs[r.Intn(len(userIDs))]
			if buyerID != sellerID {
				break
			}
		}
		statuses := []model.OrderStatus{model.OrderPending, model.OrderConfirmed, model.OrderCancelled, model.OrderCompleted}
		orders[i] = model.Order{
			GoodsID:   goodsIDs[r.Intn(len(goodsIDs))],
			BuyerID:   buyerID,
			SellerID:  sellerID,
			Status:    statuses[r.Intn(len(statuses))],
			Remark:    fmt.Sprintf("备注 %d:希望能尽快发货", i+1),
			CreatedAt: time.Now().Add(-time.Duration(r.Intn(240)) * time.Hour),
			UpdatedAt: time.Now().Add(-time.Duration(r.Intn(120)) * time.Hour),
		}
	}
	if err := db.Create(&orders).Error; err != nil {
		log.Fatal().Err(err).Msg("create orders failed")
	}
	log.Info().Int("count", len(orders)).Msg("orders created")

	log.Info().Msg("mock data migration completed")
}

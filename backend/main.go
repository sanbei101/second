package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sanbei101/second/handler"
	"github.com/sanbei101/second/model"
	"github.com/sanbei101/second/service"
)

func main() {

	log.DefaultLogger = log.Logger{
		Level:  log.InfoLevel,
		Writer: &log.IOWriter{Writer: os.Stderr},
	}

	dsn := "host=localhost user=sanbei password=123456 dbname=sanbei port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	if err := db.AutoMigrate(&model.User{}, &model.Goods{}, &model.Order{}); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	log.Info().Msg("database connected and migrated")

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		Migrate(db)
		return
	}

	userSvc := service.NewUserService(db)
	goodsSvc := service.NewGoodsService(db)
	orderSvc := service.NewOrderService(db)

	userHandler := handler.NewUserHandler(userSvc)
	goodsHandler := handler.NewGoodsHandler(goodsSvc)
	orderHandler := handler.NewOrderHandler(orderSvc, goodsSvc)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")
	userHandler.RegisterRoutes(api, userSvc)
	goodsHandler.RegisterRoutes(api, goodsSvc)
	orderHandler.RegisterRoutes(api, orderSvc, goodsSvc)

	port := 8848
	log.Info().Msgf("server starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	dsn := "host=154.8.213.38 user=myuser password=mypassword dbname=second port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := db.AutoMigrate(&model.User{}, &model.Goods{}, &model.Order{}); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
		log.Info().Msg("database connected and migrated")
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

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	userHandler.RegisterRoutes(api)
	goodsHandler.RegisterRoutes(api)
	orderHandler.RegisterRoutes(api)

	port := 8848
	log.Info().Msgf("server starting on :%d", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

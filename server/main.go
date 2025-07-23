package main

import (
	"context"
	"os"
	"os/signal"
	"tg-stars/bot"
	"tg-stars/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("tg-stars.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	bot := bot.NewBot(os.Getenv("BOT_TOKEN"), db)

	r := gin.Default()
	r.GET("/health", handler.HealthHandler(bot))

	go r.Run(":7000")
	bot.Start(ctx)
}

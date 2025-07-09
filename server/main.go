package main

import (
	"context"
	"os"
	"os/signal"
	"tg-stars/bot"
	"tg-stars/handler"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	bot := bot.NewBot(os.Getenv("BOT_TOKEN"))

	r := gin.Default()
	r.GET("/health", handler.HealthHandler(bot))

	go r.Run(":7000")
	bot.Start(ctx)
}

package handler

import (
	"net/http"
	"os"
	"strconv"
	"tg-stars/bot"

	"github.com/gin-gonic/gin"
)

func HealthHandler(b bot.Bot) gin.HandlerFunc {
	return func(c *gin.Context) {
		chatId, err := strconv.ParseInt(os.Getenv("HEALTH_CHECK_CHAT_ID"), 10, 64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid HEALTH_CHECK_CHAT_ID",
			})
			return
		}
		b.SendMessage(chatId, "Health check")
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

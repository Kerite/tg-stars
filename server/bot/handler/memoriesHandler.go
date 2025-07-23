package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"

	botModels "tg-stars/models"
)

func (bh *BotHandler) MemoriesHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// fmt.Println("[Bot] MemoriesHandler called")
	memories, err := gorm.G[botModels.Memory](bh.db).Find(ctx)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Error fetching memories: " + err.Error(),
		})
		return
	}

	text := "Here are your memories:\n"
	for _, memory := range memories {
		text += "- " + memory.Description + " Price: " + fmt.Sprint(memory.Price) + "\n"
	}
	if len(memories) == 0 {
		text = "No memories yet. Use /share to add one!"
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeMarkdown,
		Text:      text,
	})
}

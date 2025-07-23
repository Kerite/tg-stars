package handler

import (
	"bytes"
	"context"
	"fmt"
	"tg-stars/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (bh *BotHandler) ExportHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	user_id := update.Message.From.Username

	data, err := utils.ExportMemory(user_id)
	if err != nil {
		fmt.Println("Error exporting memory:", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to export memory. Please try again later.",
		})
		return
	}

	b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: update.Message.Chat.ID,
		Document: &models.InputFileUpload{
			Filename: "exported_memory_" + user_id + ".snapshot",
			Data:     bytes.NewReader(data),
		},
		Caption: "Your current memory.",
	})
}

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	botModels "tg-stars/models"
	"tg-stars/utils"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/gorm"
)

type PendingOperation = int

const (
	PendingOperationNone PendingOperation = iota
	PendingOperationImport
	PendingOperationExport
)

type BotHandler struct {
	db         *gorm.DB
	backendURL string
	pendingUsers map[int64]PendingOperation
}

func NewBotHandler(db *gorm.DB) *BotHandler {
	db.AutoMigrate(&botModels.Memory{})
	backend_url := os.Getenv("BACKEND_URL")
	return &BotHandler{
		db:         db,
		backendURL: backend_url,
		pendingUsers: make(map[int64]PendingOperation),
	}
}

func (bh *BotHandler) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// fmt.Println("Received Message:", update.Message.Text)
	message_bytes, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}
	fmt.Println(string(message_bytes))

	if update.PreCheckoutQuery != nil {
		fmt.Println("Pre-checkout query received:", update.PreCheckoutQuery)
		success, err := b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{
			PreCheckoutQueryID: update.PreCheckoutQuery.ID,
			OK:                 true,
		})
		if !success {
			fmt.Println("Failed to answer pre-checkout query")
		}
		if err != nil {
			fmt.Println("Error answering pre-checkout query:", err)
		}
		return
	} else if update.Message != nil {
		if update.Message.SuccessfulPayment != nil {
			fmt.Println("Successful payment received:", update.Message.SuccessfulPayment)
			return
		}

		if update.Message.Document != nil {
			if bh.pendingUsers[update.Message.From.ID] == PendingOperationImport {
				fmt.Println("ImportHandler called with document:", update.Message.Document.FileID)
				file, err := b.GetFile(ctx, &bot.GetFileParams{
					FileID: update.Message.Document.FileID,
				})
				if err != nil {
					fmt.Println("Error getting file:", err)
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Failed to retrieve the file. Please try again later.",
					})
					return
				}
				fileDownloadLink := b.FileDownloadLink(file)
				fmt.Println("File download link:", fileDownloadLink)
				fileData, err := utils.GetFileData(fileDownloadLink)
				if err != nil {
					fmt.Println("Error downloading file:", err)
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Failed to download the file. Please try again later.",
					})
					return
				}

				err = utils.ImportMemory(update.Message.From.Username, fileData)
				if err != nil {
					fmt.Println("Error importing memory:", err)
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Failed to import memory. Please ensure the file is valid.",
					})
				} else {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   "Memory imported successfully!",
					})
				}
				
				bh.pendingUsers[update.Message.From.ID] = PendingOperationNone
			}
			return
		}

		// Normal chat
		message, err := utils.Chat(update.Message.From.Username, update.Message.Text)
		if err != nil {
			fmt.Println("Error processing chat message:", err)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "An error occurred while processing your message. Please try again later.",
			})
			return
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   message,
		})
	}
}

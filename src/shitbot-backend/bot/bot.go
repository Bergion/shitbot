package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/matterbridge/telegram-bot-api/v6"
	"go.uber.org/fx"
)

type TelegramBot struct {
	telegramBot *tgbotapi.BotAPI
	appURL      string
}

func NewTelegramBot(lc fx.Lifecycle) *TelegramBot {
	token := os.Getenv("TELEGRAM_TOKEN")
	appURL := os.Getenv("APP_URL")
	telegramBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		if err.Error() == "Not Found" {
			panic("Bot not found. Token is invalid")
		}
		panic(err)
	}

	bot := &TelegramBot{
		telegramBot: telegramBot,
		appURL:      appURL,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting Telegram bot")
			go bot.Start()
			return nil
		},
	})

	return bot
}

func (b *TelegramBot) Start() {
	// ctx := context.Background()
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.telegramBot.GetUpdatesChan(updateConfig)
	for update := range updates {
		fmt.Println(update.Message.Text)
		if update.Message != nil {
			if strings.Contains(update.Message.Text, "/start") {
				appURL := b.appURL
				_, refCode, found := strings.Cut(update.Message.Text, " ")
				if found {
					appURL = appURL + "?ref_code=" + refCode
				}
				fmt.Println(appURL)
				setMenuConfig := tgbotapi.SetChatMenuButtonConfig{
					ChatID: update.Message.Chat.ID,
					MenuButton: &tgbotapi.MenuButton{
						Type: "web_app",
						Text: "Play",
						WebApp: &tgbotapi.WebAppInfo{
							URL: appURL,
						},
					},
				}
				response, err := b.telegramBot.Request(setMenuConfig)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(response)

			}
		}
	}
}

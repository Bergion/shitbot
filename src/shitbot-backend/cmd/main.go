package main

import (
	"net/http"
	"shitbot/bot"
	"shitbot/internal/config"
	"shitbot/internal/db"
	"shitbot/internal/handlers"
	"shitbot/internal/router"
	"shitbot/server"

	"go.uber.org/fx"
)

func main() {
	config.LoadConfig()
	fx.New(
		fx.Provide(
			server.NewHTTPServer,
			router.NewRouter,
			db.NewMongoClient,
			handlers.NewUserHandler,
			bot.NewTelegramBot,
		),
		fx.Invoke(
			func(*http.Server) {},
			func(*bot.TelegramBot) {},
		),
	).Run()
}

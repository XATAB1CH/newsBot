package main

import (
	"time"

	tg "gopkg.in/telebot.v3"
	"./tags"
)

func main() {
	bot, err := tg.NewBot(tg.Settings{
		Token:  "7006639507:AAEVXwFmksp027JhYPqvA5oh2B4U4htoBS8",
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	bot.Handle("/start", func(c tg.Context) error {
		return c.Send(SendMenu(c))
	})

	bot.Handle(&tg.Btn{Unique: tags.menuHandler, Text: "Меню"}, HandleMenu)
	bot.Handle(&tg.Btn{Unique: tags.menuHandler, Text: "Рассылка"}, HandleMailMenu)
	bot.Handle(&tg.Btn{Unique: tags.menuHandler, Text: "Регистрация"}, HandleRegistryMenu)

	bot.Start()
}

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/XATAB1CH/news-bot/news"
	tg "gopkg.in/telebot.v3"
)

const (
	Id = 1403958448
)

var (
	t1 = time.Date(2024, time.May, 7, -1, 33, 0, 0, time.UTC)
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

	bot.Handle(&tg.Btn{Unique: "menuHandler", Text: "Новости"}, HandleNewsMenu)

	go NewsUpdater(bot)

	bot.Start()
}

func SendMenu(c tg.Context) error {
	menu := &tg.ReplyMarkup{ResizeKeyboard: true}

	newsMenu := menu.Data("Новости", "menuHandler", "newsMenu")

	menu.Inline(
		menu.Row(newsMenu),
	)

	p := &tg.Photo{
		File:    tg.FromDisk("photo.png"),
		Caption: "Выберите пункт меню",
	}

	return c.Send(p, menu)
}

// func HandleMenu(c tg.Context) error {
// 	c.Send("хуй")
// 	return SendNewsMenu(c)
// }

// func SendNewsMenu(c tg.Context) error {
// 	newsMenu := &tg.ReplyMarkup{ResizeKeyboard: true}

// 	start := newsMenu.Data("Начать", "menuHandler", "start")
// 	back := newsMenu.Data("Назад", "menuHandler", "back")

// 	newsMenu.Inline(
// 		newsMenu.Row(start),
// 		newsMenu.Row(back),
// 	)

// 	c.Edit(newsMenu)
// 	return c.Respond()
// }

func HandleNewsMenu(c tg.Context) error {
	switch c.Data() {
	case "back":
		return SendMenu(c)
	}

	t2 := time.Now()

	t := t1.Sub(t2)
	days := t / (time.Hour * 24)

	t -= days * time.Hour * 24
	hours := t / time.Hour

	t -= hours * time.Hour
	minutes := t / time.Minute

	c.Send(fmt.Sprintf("Отправка новостей начнётся через %s %s !", fhours(int(hours)), fminutes(int(minutes))))

	return c.Respond()
}

func fminutes(n int) string {
	switch {
	case n%10 == 1 && n%100 != 11:
		return strconv.Itoa(n) + " минута"
	case n%10 > 1 && n%10 < 5 && !(n >= 11 && n <= 19):
		return strconv.Itoa(n) + " минуты"
	}

	return strconv.Itoa(n) + " минут"
}

func fhours(n int) string {
	switch {
	case n%10 == 1 && n%100 != 11:
		return strconv.Itoa(n) + " час"
	case n%10 > 1 && n%10 < 5 && !(n >= 11 && n <= 19):
		return strconv.Itoa(n) + " часа"
	}

	return strconv.Itoa(n) + " часов"
}

func NewsUpdater(bot *tg.Bot) {
	var t2 time.Time

	for {
		time.Sleep(time.Second * 10)
		t2 = time.Now()
		if checkTime(t1, t2) {
			for _, v := range news.NewsList {
				time.Sleep(10 * time.Second)
				bot.Send(&tg.Chat{ID: Id}, v.Title+"\n"+v.Text)
			}
			break
		}
	}
}

func checkTime(t1, t2 time.Time) bool {
	if t1.Sub(t2) <= 0 {
		return true
	} else {
		return false
	}
}

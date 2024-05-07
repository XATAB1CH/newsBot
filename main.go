package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/XATAB1CH/news-bot/news"
	"github.com/XATAB1CH/news-bot/user"
	_ "github.com/lib/pq"
	tg "gopkg.in/telebot.v3"
)

var (
	t1 = time.Date(2024, time.May, 7, 22, 44, 0, 0, time.UTC)
)

func main() {

	// подключаем бд
	connStr := "user=postgres password=anton132 dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	user.UpdateUserArr(db)
	defer db.Close()

	// подключаем бота
	bot, err := tg.NewBot(tg.Settings{
		Token:  "7006639507:AAEVXwFmksp027JhYPqvA5oh2B4U4htoBS8",
		Poller: &tg.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	bot.Handle("/start", func(c tg.Context) error {
		tgId := c.Sender().ID
		result, err := db.Exec("insert into users (id, rate) values ($1, $2) on conflict (id) do nothing", tgId, 10)
		if err != nil {
			panic(err)
		}
		fmt.Println(result.RowsAffected())

		return c.Send(SendNewMenu(c))
	})

	bot.Handle(&tg.Btn{Unique: "startHandler", Text: "Меню"}, HandleMenu)
	bot.Handle(&tg.Btn{Unique: "menuHandler", Text: "Новости"}, HandleNewsMenu)

	go NewsUpdater(bot)

	bot.Start()
}

func SendMenu(c tg.Context) error {
	menu := &tg.ReplyMarkup{ResizeKeyboard: true}

	newsMenu := menu.Data("Новости", "startHandler", "newsMenu")

	menu.Inline(
		menu.Row(newsMenu),
	)

	p := &tg.Photo{
		File:    tg.FromDisk("./assets/menu.jpg"),
		Caption: "Выберите пункт меню",
	}

	return c.Edit(p, menu)
}

func SendNewMenu(c tg.Context) error {
	menu := &tg.ReplyMarkup{ResizeKeyboard: true}

	newsMenu := menu.Data("Новости", "startHandler", "newsMenu")

	menu.Inline(
		menu.Row(newsMenu),
	)

	p := &tg.Photo{
		File:    tg.FromDisk("./assets/menu.jpg"),
		Caption: "Выберите пункт меню",
	}

	return c.Send(p, menu)
}

func HandleMenu(c tg.Context) error {
	switch c.Data() {
	case "newsMenu":
		return SendNewsMenu(c)
	}
	return c.Respond()
}

func SendNewsMenu(c tg.Context) error {
	newsMenu := &tg.ReplyMarkup{ResizeKeyboard: true}

	start := newsMenu.Data("Начать", "menuHandler", "start")
	back := newsMenu.Data("Назад", "menuHandler", "back")

	newsMenu.Inline(
		newsMenu.Row(start),
		newsMenu.Row(back),
	)

	p := &tg.Photo{
		File:    tg.FromDisk("./assets/news.png"),
		Caption: "Нажмите чтобы начать отправку новостей",
	}

	c.Edit(p, newsMenu)
	return c.Respond()
}

func HandleNewsMenu(c tg.Context) error {
	switch c.Data() {
	case "back":
		c.Edit(SendMenu(c))
		return c.Respond()
	}

	t2 := time.Now()

	t := t1.Sub(t2)
	days := t / (time.Hour * 24)

	t -= days * time.Hour * 24
	hours := t / time.Hour

	t -= hours * time.Hour
	minutes := t / time.Minute

	c.Send(fmt.Sprintf("Отправка новостей начнётся через %s %s!", fhours(int(hours)), fminutes(int(minutes))))

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
				for _, u := range user.UserArr {
					bot.Send(&tg.Chat{ID: u.Id}, v.Title+"\n"+v.Text)
				}
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

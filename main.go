package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

const (
	telegramToken = "8070880491:AAHmubgiYH3zeo46Q7qijpHvBx11SuGzWs8"
	dbConnString  = "postgres://postgres:Geksagon-2018@localhost:5432/symbol_transport?search_path=perso"
)

func main() {
	// Подключение к базе данных
	conn, err := pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())

	// Создание Telegram бота
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v\n", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Обработка обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "date":
				err := updateDate(conn)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Помилка оновлення дати.")
					bot.Send(msg)
					fmt.Printf("Помилка оновлення дати о %s\n", time.Now())
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Дата успішно оновлена для серії 11, проєкта Чернігів. "+
						"Можна приступати до персоналізації")
					bot.Send(msg)
					fmt.Printf("Дата успішно оновлена о %s\n", time.Now())
				}
			}
		}
	}
}

func updateDate(conn *pgx.Conn) error {
	// Обновляем поле date в таблице test
	currentTime := time.Now()
	formattedDate := currentTime.Format("060102")
	_, err := conn.Exec(context.Background(), "UPDATE barcode_profile SET ticket_expiry_date = $1 WHERE project_id = 10 AND card_series = '11' AND d=FALSE", formattedDate)
	return err
}

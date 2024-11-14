package main

import (
	"context"
	"fmt"
	"os"
	"telegram/initializers"
	"telegram/queryes"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Подключение к базе данных
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())

	// Создание Telegram бота
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAMTOKEN"))
	if err != nil {
		fmt.Errorf("Failed to create bot: %v\n", err)
	}

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

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
				err := queryes.UpdateData(conn)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Помилка оновлення дати.")
					bot.Send(msg)
					fmt.Printf("Error: Помилка оновлення дати о %s\n", time.Now())
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Дата успішно оновлена для серії 11, проєкта Чернігів. "+
						"Можна приступати до персоналізації")
					bot.Send(msg)
					fmt.Printf("Info: Оновлення дати о %s\n", time.Now())
				}

			case "getdate":
				results, err := queryes.GetData(conn)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Помилка отримання дати.")
					bot.Send(msg)
					fmt.Printf("Error: Помилка отримання дати о %s\n", time.Now())
					fmt.Printf("Помилка отримання дати: %s\n", err)
				} else {
					var response string
					for _, row := range results {
						for key, value := range row {
							response += fmt.Sprintf("%s: %v\n", key, value)
						}
						response += "\n" // add empty row before print
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
					bot.Send(msg)
					fmt.Printf("Info: Отримання дати о %s\n", time.Now())
				}
			}
		}
	}
}

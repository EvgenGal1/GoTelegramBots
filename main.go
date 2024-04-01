package main

import (
	// imp библ.
	"fmt"
	"os"
	// telego
	"github.com/mymmrac/telego"
	// алиас > обработчик + предикаты(вместо путей)
	// th "github.com/mymmrac/telego/telegohandler"
	// алиас >  утилит
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	os.Setenv("TOKEN", "7164464409:AAFidAN6SFXkkbV1cGzw9YjEMFxZLhd0Q9s")
	// хран.токен
	botToken := os.Getenv("TOKEN")

	// ч/з мтд.NewBot созд.экземпляр бота с обраб.ошб.Logger(`регистратор`) е/и Токен не верный
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// канал получения обновлений. Созд.executor(`исполнитель`). Есть 2 типа - Long Poling(`Длинный опрос`) | Webhook(`сеть-перехватчик`). Здесь 1 ч/з `обновления ч/з длительный опрос`
	updates, _ := bot.UpdatesViaLongPolling(nil)

	// прекратить получать обновления. В стек + fn StopLongPolling(`прекратить долгий опрос`). Вызов при возврате из main по LIFO(от последн.к первому)
	defer bot.StopLongPolling()

	// цикл получ.обнов.серверов Телеграма
	for updates := range updates {
		// е/и не пустое смс
		if updates.Message != nil {
			// ID смс
			chatID := tu.ID(updates.Message.Chat.ID)
			// клавиатура
			keyboard := tu.Keyboard(
				tu.KeyboardRow(
					tu.KeyboardButton("Старт"),
					tu.KeyboardButton("Помощь"),
					tu.KeyboardButton("Тех. поддержка"),
				),
				tu.KeyboardRow(
					tu.KeyboardButton("Отправить локацию").WithRequestLocation(),
					tu.KeyboardButton("Отправить контакт").WithRequestContact(),
					tu.KeyboardButton("Отмена"),
				),
			)

			// прикреп.клвт.к ответ.смс ч/з WithReplyMarkup
			message := tu.Message(
				chatID,
				"С смс придёт клава. ! но не пришла",
			).WithReplyMarkup(keyboard)

			// отправка клавиатуры
			_, _ = bot.SendMessage(message)

			// копир/отправка такого же смс
			_, _ = bot.CopyMessage(
				tu.CopyMessage(
					chatID,
					chatID,
					updates.Message.MessageID,
				),
			)

			// ответ стикером
			_, _ = bot.SendSticker(
				// мтд.Sticker(helper из telego) созд.парам.
				tu.Sticker(
					// куда отправ.стик.
					chatID,
					// мтд.FileFromID указ. ID стикер
					tu.FileFromID("CAACAgIAAxkBAAEL0-tmCu0FAkLYtjHbWr0vhKSCkxRbiQACpwcAAmMr4gmVkahaRG9rAjQE"),
				),
			)
		}
	}
}

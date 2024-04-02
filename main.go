package main

import (
	// imp библ.
	"fmt"
	"os"
	// telego
	"github.com/mymmrac/telego"
	// алиас th > обработчик + предикаты(вместо путей)
	th "github.com/mymmrac/telego/telegohandler"
	// алиас tu > утилит
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

	// обраб.бота c указ.получ.обновления
	bh, _ := th.NewBotHandler(bot, updates)

	// прекратить обработку обновлений
	defer bh.Stop()

	// прекратить получать обновления. В стек + fn StopLongPolling(`прекратить долгий опрос`). Вызов при возврате из main по LIFO(от последн.к первому)
	defer bot.StopLongPolling()

	// цикл получ.обнов.серверов Телеграма
	// for updates := range updates {
	// е/и не пустое смс
	// if updates.Message != nil {
	// переезд в обработчик. updates <> update
	// ID смс, клавиатура, прикреп.клвт.к ответ.смс, отправка клавиатуры, ответ стикером
	// }
	// }

	// регистр.нов.обраб.по КМД."/start". Узак.fn + тип.обраб.
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// обработчик
		// ID смс
		chatID := tu.ID(update.Message.Chat.ID)
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

		// копир/отправка такого же смс
		_, _ = bot.CopyMessage(
			tu.CopyMessage(
				chatID,
				chatID,
				update.Message.MessageID,
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

		// прикреп.клвт.к ответ.смс ч/з WithReplyMarkup
		message := tu.Message(
			chatID,
			"Выбирете пункт меню",
		).WithReplyMarkup(keyboard)

		// отправка клвт.
		_, _ = bot.SendMessage(message)
	}, th.CommandEqual("start"))

	// регистр.нов.обраб.для всех КМД.кроме `/start`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Отправить сообщение
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Неизвестная команда, используйте /start",
		))
	}, th.AnyCommand())

	// регистр.нов.обраб.КНП. Пока > Тех.поддержка и Помощь
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Проверяем, что полученное сообщение содержит кнопку "Тех. поддержка"
		if update.Message != nil {
			// ID чата, куда отправлять ответ
			chatID := tu.ID(update.Message.Chat.ID)
			// перем.хран.смс
			var smsCre string
			// общ.улов.и проверка
			if update.Message.Text == "Тех. поддержка" {
				// смс > кнп."Тех. поддержка"
				smsCre = "Вы нажали на кнопку 'Тех. поддержка'"
			}
			if update.Message.Text == "Помощь" {
				smsCre = "Какая вам нужна 'Помощь'"
			}
			// отправ.смс от кажд.кнп.
			_, _ = bot.SendMessage(tu.Message(
				chatID,
				smsCre,
			))
		}
	}, th.AnyMessage())

	// старт.обраб.обновлений
	bh.Start()
}

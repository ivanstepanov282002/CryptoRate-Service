package bot

import (
	"cryptorate-service/internal/api"
	"cryptorate-service/internal/models"
	"cryptorate-service/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	api       *tgbotapi.BotAPI
	updates   tgbotapi.UpdatesChannel
	apiClient *api.CoinGeckoClient
	repo      *repository.Repository
}

// Создаем нового бота
func NewBot(token string, db *sql.DB) (*TelegramBot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = true
	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0) //Запрашивает все письма с последнего непрочитанного
	u.Timeout = 60             //После 60 секунд бездействия начинается новый цикл

	updates := botAPI.GetUpdatesChan(u) //Получаем канал сообщений

	// Создаём API клиент и репозиторий
	apiClient := api.NewCoinGeckoClient()
	repo := repository.NewRepository(db)

	return &TelegramBot{
		api:       botAPI,
		updates:   updates,
		apiClient: apiClient,
		repo:      repo,
	}, nil
}

// Метод для запуски бота
func (b *TelegramBot) Start() {
	for update := range b.updates {
		if update.Message == nil {
			continue
		}

		err := b.repo.EnsureUser(update.Message.Chat.ID, update.Message.From.UserName)
		if err != nil {
			log.Printf("Error ensuring user: %v", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "") //Записываем ID диалогв

		//Читаем что прислал пользователь и сравниваем с возможными вариантами
		switch update.Message.Command() {
		case "start":
			msg.Text = "Привет! Я бот для отслеживания курсов криптовалют.\n\n" +
				"Доступные команды:\n" +
				"/rates - все курсы\n" +
				"/rates [валюта] - курс конкретной валюты\n" +
				"/currencies - список всех валют\n" +
				"/startauto [минуты] - автоотправка\n" +
				"/stopauto - остановить автоотправку"

		case "rates":
			args := update.Message.CommandArguments()
			if args == "" {
				rates, err := b.repo.GetLatestRates()
				if err != nil {
					log.Printf("Database error: %v", err)
					msg.Text = "Ошибка получения курсов"
				} else if len(rates) == 0 {
					msg.Text = "Курсов пока нет. Попробуйте позже."
				} else {
					log.Printf("[DEBUG] Получено %d курсов", len(rates))
					for i, rate := range rates {
						log.Printf("[DEBUG] Курс %d: %s, CurrencyID: %d",
							i+1, rate.NameCurrency, rate.CurrencyID)
					}

					var response strings.Builder
					response.WriteString("📊 Последние курсы:\n\n")
					for _, rate := range rates {
						// Теперь rate.CurrencyID доступен
						symbol, _ := b.repo.GetCurrencySymbolByID(rate.CurrencyID)
						timeStr := rate.RecordedAt.Format("15:04")
						response.WriteString(fmt.Sprintf("• %s (%s): $%.2f (%s)\n",
							rate.NameCurrency, symbol, rate.Price, timeStr))
					}
					response.WriteString("\n🔄 Обновляется каждые 5 минут")
					msg.Text = response.String()
				}
			} else {
				// Курс конкретной валюты
				currencyName := strings.ToLower(args)

				// Пробуем найти по символу (BTC, ETH)
				currencyID, err := b.repo.GetCurrencyIDBySymbol(currencyName)
				if err != nil {
					// Если не нашли по символу, ищем по имени
					currencyID, err = b.repo.GetCurrencyID(currencyName)
				}

				if err != nil {
					msg.Text = "Валюта не найдена. Используйте /currencies для списка"
				} else {
					rate, err := b.repo.GetCurrencyRate(currencyID)
					if err != nil {
						msg.Text = "Ошибка получения курса"
					} else {
						min, max, _ := b.repo.GetDailyMinMax(currencyID)
						change, _ := b.repo.GetHourlyChange(currencyID)

						// Получаем информацию о валюте
						symbol, _ := b.repo.GetCurrencySymbolByID(currencyID)
						displayName, _ := b.repo.GetCurrencyDisplayName(currencyID)

						msg.Text = fmt.Sprintf(
							"📊 %s (%s)\n"+
								"💵 Текущий курс: $%.2f\n"+
								"📈 День: $%.2f - $%.2f\n"+
								"🕐 Час: %.2f%%\n"+
								"⏰ Обновлено: %s",
							displayName,
							symbol,
							rate.Price,
							min,
							max,
							change,
							rate.RecordedAt.Format("15:04"),
						)
					}
				}
			}

		case "currencies":
			currencies, err := b.repo.GetAllCurrencies()
			if err != nil {
				msg.Text = "Ошибка получения списка валют"
			} else {
				var response strings.Builder
				response.WriteString("📋 Доступные валюты:\n\n")

				for _, currency := range currencies {
					response.WriteString(fmt.Sprintf("• %s (%s)\n",
						currency.DisplayName, currency.Symbol))
				}

				response.WriteString("\n💡 Используйте /rates [символ] для получения курса\n")
				response.WriteString("Пример: /rates BTC или /rates bitcoin")
				msg.Text = response.String()
			}

		case "startauto":
			args := update.Message.CommandArguments()
			if args == "" {
				msg.Text = "Укажите интервал в минутах. Пример: /start-auto 10"
			} else {
				interval, err := strconv.Atoi(args)
				if err != nil || interval <= 0 {
					msg.Text = "Интервал должен быть положительным числом (минуты)"
				} else if interval < 5 {
					msg.Text = "Минимальный интервал - 5 минут"
				} else {
					err := b.repo.SetUserInterval(update.Message.Chat.ID, interval)
					if err != nil {
						log.Printf("Error setting interval: %v", err)
						msg.Text = "Ошибка настройки автоотправки"
					} else {
						msg.Text = fmt.Sprintf(
							"✅ Автоотправка включена\n"+
								"📩 Курсы будут приходить каждые %d минут\n\n"+
								"❌ Используйте /stop-auto для отключения",
							interval,
						)
					}
				}
			}

		case "stopauto":
			err := b.repo.StopAuto(update.Message.Chat.ID)
			if err != nil {
				log.Printf("Error stopping auto: %v", err)
				msg.Text = "Ошибка отключения автоотправки"
			} else {
				msg.Text = "✅ Автоотправка отключена"
			}

		default:
			if update.Message.Text != "" {
				msg.Text = "Неизвестная команда. Используйте /start"
			}
		}

		if msg.Text != "" {
			b.api.Send(msg)
		}
	}
}

// autoSendWorker отправляет автоматические уведомления
func (b *TelegramBot) autoSendWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		users, err := b.repo.GetSubscribedUsers()
		if err != nil {
			log.Printf("Error getting subscribed users: %v", err)
			continue
		}

		currentTime := time.Now()

		for _, user := range users {
			// Проверяем, пора ли отправлять
			nextSendTime := user.LastSent.Add(time.Duration(user.Interval) * time.Minute)

			if currentTime.After(nextSendTime) {
				// Формируем сообщение
				message := b.buildAutoMessage(user.Currencies)
				if message != "" {
					msg := tgbotapi.NewMessage(user.UserID, message)

					// Отправляем с повторными попытками
					for i := 0; i < 3; i++ {
						_, err := b.api.Send(msg)
						if err == nil {
							// Успешно отправили, обновляем время
							b.repo.UpdateLastSent(user.UserID)
							break
						}

						if i < 2 {
							time.Sleep(2 * time.Second)
						} else {
							log.Printf("Failed to send to user %d after 3 attempts: %v",
								user.UserID, err)
						}
					}
				}
			}
		}
	}
}

// buildAutoMessage формирует сообщение для автоотправки
func (b *TelegramBot) buildAutoMessage(currencies []models.Currency) string {
	if len(currencies) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("🔄 Автообновление курсов:\n\n")

	// Ограничиваем количество валют в одном сообщении
	maxCurrencies := 3
	if len(currencies) > maxCurrencies {
		currencies = currencies[:maxCurrencies]
	}

	for _, currency := range currencies {
		currencyID, err := b.repo.GetCurrencyID(currency.NameCurrency)
		if err != nil {
			continue
		}

		rate, err := b.repo.GetCurrencyRate(currencyID)
		if err != nil {
			continue
		}

		min, max, _ := b.repo.GetDailyMinMax(currencyID)
		change, _ := b.repo.GetHourlyChange(currencyID)

		builder.WriteString(fmt.Sprintf(
			"• %s (%s): $%.2f\n"+
				"  📊 День: $%.2f - $%.2f\n"+
				"  📈 Час: %.2f%%\n\n",
			currency.DisplayName,
			currency.Symbol,
			rate.Price,
			min,
			max,
			change,
		))
	}

	builder.WriteString("⏰ " + time.Now().Format("15:04"))
	builder.WriteString("\n💡 /stop-auto для отключения")

	return builder.String()
}

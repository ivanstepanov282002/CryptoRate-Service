package bot

import (
	"cryptorate-service/internal/api"            
	"cryptorate-service/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"strings"

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

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "") //Записываем ID диалогв

		//Читаем что прислал пользователь и сравниваем с возможными вариантами
		switch update.Message.Text {
		case "/start":
			msg.Text = "Привет! Я бот для отслеживания курсов криптовалют.\nИспользуй /rates чтобы получить текущие курсы."
		case "/rates":
			rates, err := b.repo.GetLatestRates()
			if err != nil {
				log.Printf("Database error: %v", err)
				msg.Text = "Ошибка получения курсов из базы данных"
			} else if len(rates) == 0 {
				msg.Text = "В базе данных ещё нет курсов. Попробуйте позже."
			} else {
				// Формируем красивое сообщение
				var response strings.Builder
				response.WriteString("📊 Последние курсы:\n\n")

				for _, rate := range rates {
					// Форматируем время
					timeStr := rate.RecordedAt.Format("15:04:00")
					response.WriteString(fmt.Sprintf("• %s: $%.2f (%s)\n",
						rate.NameCurrency, rate.Price, timeStr))
				}

				//response.WriteString("\n🔄 Обновляются каждые 5 минут")
				msg.Text = response.String()
			}

		default:
			msg.Text = "Неизвестная команда. Используй /start или /rates"
		}

			b.api.Send(msg)
	}
}
package bot

import (
	"autoadmin/internal/client"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	client *client.Client
}

func New(apiURL string) (*Bot, error) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN environment variable is required")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("Bot authorized as @%s", api.Self.UserName)

	return &Bot{
		api:    api,
		client: client.New(apiURL),
	}, nil
}

func (b *Bot) Start() {
	log.Println("Bot starting...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go b.handleMessage(update.Message)
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		return
	}

	switch msg.Command() {
	case "start":
		b.handleStart(msg)
	case "help":
		b.handleHelp(msg)
	case "masters":
		b.handleMasters(msg)
	case "master":
		b.handleMaster(msg)
	case "clients":
		b.handleClients(msg)
	case "client":
		b.handleClient(msg)
	case "services":
		b.handleServices(msg)
	case "schedule":
		b.handleSchedule(msg)
	case "tariffs":
		b.handleTariffs(msg)
	case "available":
		b.handleAvailable(msg)
	case "dashboard":
		b.handleDashboard(msg)
	default:
		b.send(msg, "Неизвестная команда. Используйте /help")
	}
}

func (b *Bot) send(msg *tgbotapi.Message, text string) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "HTML"
	b.api.Send(reply)
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	text := `👋 <b>Добро пожаловать в AutoAdmin Bot!</b>

Я помогу вам управлять записями и данными.

Доступные команды:
/masters — Список мастеров
/master [id] — Информация о мастере
/clients — Список клиентов
/client [id] — Информация о клиенте
/services [master_id] — Список услуг
/schedule — Расписание записей
/tariffs — Тарифы
/available [master_id] [service_id] [дата] — Свободные слоты
/dashboard [id] — Панель управления мастера
/help — Справка

Пример: <code>/available 1 1 2026-06-15</code>`

	b.send(msg, text)
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	b.handleStart(msg)
}

func (b *Bot) handleMasters(msg *tgbotapi.Message) {
	masters, err := b.client.GetMasters()
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	if len(masters) == 0 {
		b.send(msg, "📭 Мастеров пока нет")
		return
	}

	var sb strings.Builder
	sb.WriteString("📋 <b>Список мастеров:</b>\n\n")

	for _, m := range masters {
		name := "Без имени"
		if m.Name != nil {
			name = *m.Name
		}
		status := "✅"
		if m.IsActive == 0 {
			status = "❌"
		}
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> | ID: %d | TG: %d\n", status, name, m.ID, m.TelegramID))
	}

	b.send(msg, sb.String())
}

func (b *Bot) handleMaster(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.send(msg, "Использование: /master [id]\nПример: /master 1")
		return
	}

	id, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil {
		b.send(msg, "❌ Неверный ID")
		return
	}

	m, err := b.client.GetMaster(id)
	if err != nil {
		b.send(msg, "❌ Мастер не найден")
		return
	}

	name := "Без имени"
	if m.Name != nil {
		name = *m.Name
	}
	phone := "не указан"
	if m.Phone != nil {
		phone = *m.Phone
	}
	status := "Активен"
	if m.IsActive == 0 {
		status = "Неактивен"
	}

	text := fmt.Sprintf(`👤 <b>Мастер #%d</b>

Имя: %s
Телефон: %s
Telegram ID: %d
Статус: %s
Создан: %s`, m.ID, name, phone, m.TelegramID, status, m.CreatedAt)

	b.send(msg, text)
}

func (b *Bot) handleClients(msg *tgbotapi.Message) {
	clients, err := b.client.GetClients()
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	if len(clients) == 0 {
		b.send(msg, "📭 Клиентов пока нет")
		return
	}

	var sb strings.Builder
	sb.WriteString("👥 <b>Список клиентов:</b>\n\n")

	for _, cl := range clients {
		name := "Без имени"
		if cl.Name != nil {
			name = *cl.Name
		}
		status := "✅"
		if cl.IsBlocked == 1 {
			status = "🚫"
		}
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> | ID: %d | No-show: %d\n", status, name, cl.ID, cl.NoShowCount))
	}

	b.send(msg, sb.String())
}

func (b *Bot) handleClient(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.send(msg, "Использование: /client [id]\nПример: /client 1")
		return
	}

	id, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil {
		b.send(msg, "❌ Неверный ID")
		return
	}

	cl, err := b.client.GetClient(id)
	if err != nil {
		b.send(msg, "❌ Клиент не найден")
		return
	}

	name := "Без имени"
	if cl.Name != nil {
		name = *cl.Name
	}
	phone := "не указан"
	if cl.Phone != nil {
		phone = *cl.Phone
	}
	status := "Активен"
	if cl.IsBlocked == 1 {
		status = "Заблокирован"
	}

	text := fmt.Sprintf(`👤 <b>Клиент #%d</b>

Имя: %s
Телефон: %s
Telegram ID: %v
No-show: %d
Статус: %s
Создан: %s`, cl.ID, name, phone, cl.TelegramID, cl.NoShowCount, status, cl.CreatedAt)

	b.send(msg, text)
}

func (b *Bot) handleServices(msg *tgbotapi.Message) {
	var masterID *int
	args := msg.CommandArguments()
	if args != "" {
		id, err := strconv.Atoi(strings.TrimSpace(args))
		if err == nil {
			masterID = &id
		}
	}

	services, err := b.client.GetServices(masterID)
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	if len(services) == 0 {
		b.send(msg, "📭 Услуг пока нет")
		return
	}

	var sb strings.Builder
	sb.WriteString("💇 <b>Список услуг:</b>\n\n")

	for _, s := range services {
		price := "бесплатно"
		if s.Price != nil {
			price = fmt.Sprintf("%.0f ₽", *s.Price)
		}
		sb.WriteString(fmt.Sprintf("<b>%s</b> | %d мин | %s | Мастер: %d\n",
			s.Name, s.DurationMinutes, price, s.MasterID))
	}

	b.send(msg, sb.String())
}

func (b *Bot) handleSchedule(msg *tgbotapi.Message) {
	var masterID *int
	var status *string

	args := strings.Fields(msg.CommandArguments())
	if len(args) > 0 {
		id, err := strconv.Atoi(args[0])
		if err == nil {
			masterID = &id
		}
	}
	if len(args) > 1 {
		status = &args[1]
	}

	slots, err := b.client.GetScheduleSlots(masterID, status)
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	if len(slots) == 0 {
		b.send(msg, "📭 Записей пока нет")
		return
	}

	var sb strings.Builder
	sb.WriteString("📅 <b>Расписание:</b>\n\n")

	for _, s := range slots {
		statusEmoji := "⏳"
		switch s.Status {
		case "confirmed":
			statusEmoji = "✅"
		case "cancelled":
			statusEmoji = "❌"
		case "completed":
			statusEmoji = "✔️"
		case "no_show":
			statusEmoji = "🚫"
		}

		clientID := "-"
		if s.ClientID != nil {
			clientID = strconv.Itoa(*s.ClientID)
		}

		sb.WriteString(fmt.Sprintf("%s #%d | %s → %s | Клиент: %s | %s\n",
			statusEmoji, s.ID, s.StartTime, s.EndTime, clientID, s.Status))
	}

	b.send(msg, sb.String())
}

func (b *Bot) handleTariffs(msg *tgbotapi.Message) {
	tariffs, err := b.client.GetTariffs()
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	if len(tariffs) == 0 {
		b.send(msg, "📭 Тарифов пока нет")
		return
	}

	var sb strings.Builder
	sb.WriteString("💰 <b>Тарифы:</b>\n\n")

	for _, t := range tariffs {
		status := "✅"
		if t.IsActive == 0 {
			status = "❌"
		}
		meetings := "∞"
		if t.MeetingLimit != nil {
			meetings = strconv.Itoa(*t.MeetingLimit)
		}
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> | %.0f ₽ | Встреч: %s | %d дней\n",
			status, t.Name, t.Price, meetings, t.DurationDays))
	}

	b.send(msg, sb.String())
}

func (b *Bot) handleAvailable(msg *tgbotapi.Message) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 3 {
		b.send(msg, "Использование: /available [master_id] [service_id] [дата]\nПример: /available 1 1 2026-06-15")
		return
	}

	masterID, err := strconv.Atoi(args[0])
	if err != nil {
		b.send(msg, "❌ Неверный master_id")
		return
	}

	serviceID, err := strconv.Atoi(args[1])
	if err != nil {
		b.send(msg, "❌ Неверный service_id")
		return
	}

	date := args[2]

	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		b.send(msg, "❌ Неверный формат даты. Используйте YYYY-MM-DD")
		return
	}

	slots, err := b.client.GetAvailableSlots(masterID, serviceID, date)
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	if len(slots) == 0 {
		b.send(msg, "📭 Свободных слотов на "+date+" нет")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🕐 <b>Свободные слоты на %s:</b>\n\n", date))

	for _, s := range slots {
		start := s.StartTime
		if idx := strings.Index(start, " "); idx != -1 {
			start = start[idx+1:]
		}
		end := s.EndTime
		if idx := strings.Index(end, " "); idx != -1 {
			end = end[idx+1:]
		}
		sb.WriteString(fmt.Sprintf("  %s — %s\n", start, end))
	}

	b.send(msg, sb.String())
}

func (b *Bot) handleDashboard(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.send(msg, "Использование: /dashboard [master_id]\nПример: /dashboard 1")
		return
	}

	id, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil {
		b.send(msg, "❌ Неверный ID")
		return
	}

	d, err := b.client.GetDashboard(id)
	if err != nil {
		b.send(msg, "❌ Ошибка: "+err.Error())
		return
	}

	name := "Без имени"
	if d.Master.Name != nil {
		name = *d.Master.Name
	}

	workDays := 0
	for _, h := range d.WorkingHours {
		if h.IsDayOff == 0 {
			workDays++
		}
	}

	text := fmt.Sprintf(`📊 <b>Панель управления — %s</b>

👤 Активных клиентов: %d
📅 Всего записей: %d
⏳ Ожидают подтверждения: %d
💇 Услуг: %d
⏰ Рабочих дней: %d`, name, d.ActiveClients, d.TotalBookings, d.PendingBookings, len(d.Services), workDays)

	b.send(msg, text)
}

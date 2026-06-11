package model

type Master struct {
	ID                 int     `json:"id"`
	TelegramID         int     `json:"telegram_id"`
	Name               *string `json:"name"`
	Phone              *string `json:"phone"`
	Description        *string `json:"description"`
	SubscriptionUntil  *string `json:"subscription_until"`
	TariffID           *int    `json:"tariff_id"`
	IsActive           int     `json:"is_active"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type Client struct {
	ID           int     `json:"id"`
	TelegramID   *int    `json:"telegram_id"`
	Name         *string `json:"name"`
	Phone        *string `json:"phone"`
	NoShowCount  int     `json:"no_show_count"`
	IsBlocked    int     `json:"is_blocked"`
	BlockedUntil *string `json:"blocked_until"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type Service struct {
	ID              int      `json:"id"`
	MasterID        int      `json:"master_id"`
	Name            string   `json:"name"`
	DurationMinutes int      `json:"duration_minutes"`
	Price           *float64 `json:"price"`
	CreatedAt       string   `json:"created_at"`
}

type ScheduleSlot struct {
	ID              int     `json:"id"`
	MasterID        int     `json:"master_id"`
	ClientID        *int    `json:"client_id"`
	ServiceID       *int    `json:"service_id"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	Status          string  `json:"status"`
	Details         *string `json:"details"`
	ConfirmCode     *string `json:"confirm_code"`
	ConfirmDeadline *string `json:"confirm_deadline"`
	CheckinTime     *string `json:"checkin_time"`
	CancelledBy     *string `json:"cancelled_by"`
	CancelledAt     *string `json:"cancelled_at"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type WorkingHour struct {
	ID         int     `json:"id"`
	MasterID   int     `json:"master_id"`
	DayOfWeek  int     `json:"day_of_week"`
	TimeStart  string  `json:"time_start"`
	TimeEnd    string  `json:"time_end"`
	BreakStart *string `json:"break_start"`
	BreakEnd   *string `json:"break_end"`
	IsDayOff   int     `json:"is_day_off"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type Tariff struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	MeetingLimit *int     `json:"meeting_limit"`
	ClientLimit  *int     `json:"client_limit"`
	DurationDays int      `json:"duration_days"`
	IsActive     int      `json:"is_active"`
	CreatedAt    string   `json:"created_at"`
}

type NoShowSetting struct {
	MasterID       int    `json:"master_id"`
	EnablePenalty  int    `json:"enable_penalty"`
	PenaltyPercent int    `json:"penalty_percent"`
	NoShowLimit    int    `json:"no_show_limit"`
	BlockDays      int    `json:"block_days"`
	ConfirmMinutes int    `json:"confirm_minutes"`
	CheckinMethod  string `json:"checkin_method"`
	RemindMinutes  int    `json:"remind_minutes"`
	UpdatedAt      string `json:"updated_at"`
}

type BlacklistEntry struct {
	ID        int     `json:"id"`
	MasterID  int     `json:"master_id"`
	ClientID  int     `json:"client_id"`
	Reason    *string `json:"reason"`
	CreatedAt string  `json:"created_at"`
}

type BlockedSlot struct {
	ID        int     `json:"id"`
	MasterID  int     `json:"master_id"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	Reason    *string `json:"reason"`
	CreatedAt string  `json:"created_at"`
}

type RefCode struct {
	ID        int     `json:"id"`
	MasterID  int     `json:"master_id"`
	ShortID   string  `json:"short_id"`
	QRCodeURL *string `json:"qr_code_url"`
	IsActive  int     `json:"is_active"`
	CreatedAt string  `json:"created_at"`
}

type SubscriptionPayment struct {
	ID            int      `json:"id"`
	MasterID      int      `json:"master_id"`
	TariffID      *int     `json:"tariff_id"`
	Amount        float64  `json:"amount"`
	PaymentMethod *string  `json:"payment_method"`
	PaymentID     *string  `json:"payment_id"`
	Status        string   `json:"status"`
	PaidAt        *string  `json:"paid_at"`
	ValidFrom     *string  `json:"valid_from"`
	ValidUntil    *string  `json:"valid_until"`
	CreatedAt     string   `json:"created_at"`
}

type NotificationLog struct {
	ID               int      `json:"id"`
	UserType         string   `json:"user_type"`
	UserID           int      `json:"user_id"`
	TelegramID       int      `json:"telegram_id"`
	NotificationType *string  `json:"notification_type"`
	Message          *string  `json:"message"`
	Status           string   `json:"status"`
	ErrorMessage     *string  `json:"error_message"`
	CreatedAt        string   `json:"created_at"`
}

type MasterClientBinding struct {
	ID       int     `json:"id"`
	MasterID int     `json:"master_id"`
	ClientID int     `json:"client_id"`
	BindedAt string  `json:"binded_at"`
	BindType *string `json:"bind_type"`
	IsActive int     `json:"is_active"`
}

type DashboardResponse struct {
	Master          Master         `json:"master"`
	Services        []Service      `json:"services"`
	WorkingHours    []WorkingHour  `json:"working_hours"`
	TodaySlots      []ScheduleSlot `json:"today_slots"`
	UpcomingSlots   []ScheduleSlot `json:"upcoming_slots"`
	ActiveClients   int            `json:"active_clients"`
	TotalBookings   int            `json:"total_bookings"`
	PendingBookings int            `json:"pending_bookings"`
}

type AvailableSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type LoginRequest struct {
	TelegramID int `json:"telegram_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Master Master `json:"master"`
}

type ConfirmRequest struct {
	Code string `json:"code"`
}

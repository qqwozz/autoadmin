package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) get(path string, out any) error {
	resp, err := c.HTTPClient.Get(c.BaseURL + path)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

type Master struct {
	ID                int     `json:"id"`
	TelegramID        int     `json:"telegram_id"`
	Name              *string `json:"name"`
	Phone             *string `json:"phone"`
	Description       *string `json:"description"`
	SubscriptionUntil *string `json:"subscription_until"`
	TariffID          *int    `json:"tariff_id"`
	IsActive          int     `json:"is_active"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type ClientModel struct {
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
	CreatedAt       string  `json:"created_at"`
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

type AvailableSlot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Dashboard struct {
	Master          Master         `json:"master"`
	Services        []Service      `json:"services"`
	WorkingHours    []WorkingHour  `json:"working_hours"`
	ActiveClients   int            `json:"active_clients"`
	TotalBookings   int            `json:"total_bookings"`
	PendingBookings int            `json:"pending_bookings"`
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
}

func (c *Client) GetMasters() ([]Master, error) {
	var masters []Master
	err := c.get("/api/masters", &masters)
	return masters, err
}

func (c *Client) GetMaster(id int) (*Master, error) {
	var m Master
	err := c.get(fmt.Sprintf("/api/masters/%d", id), &m)
	return &m, err
}

func (c *Client) GetMasterByTelegramID(telegramID int) (*Master, error) {
	var masters []Master
	err := c.get(fmt.Sprintf("/api/masters?telegram_id=%d", telegramID), &masters)
	if err != nil {
		return nil, err
	}
	if len(masters) == 0 {
		return nil, fmt.Errorf("master not found")
	}
	return &masters[0], nil
}

func (c *Client) GetClients() ([]ClientModel, error) {
	var clients []ClientModel
	err := c.get("/api/clients", &clients)
	return clients, err
}

func (c *Client) GetClient(id int) (*ClientModel, error) {
	var cl ClientModel
	err := c.get(fmt.Sprintf("/api/clients/%d", id), &cl)
	return &cl, err
}

func (c *Client) GetServices(masterID *int) ([]Service, error) {
	path := "/api/services"
	if masterID != nil {
		path += fmt.Sprintf("?master_id=%d", *masterID)
	}
	var services []Service
	err := c.get(path, &services)
	return services, err
}

func (c *Client) GetScheduleSlots(masterID *int, status *string) ([]ScheduleSlot, error) {
	path := "/api/schedule-slots"
	params := []string{}
	if masterID != nil {
		params = append(params, fmt.Sprintf("master_id=%d", *masterID))
	}
	if status != nil {
		params = append(params, fmt.Sprintf("status=%s", *status))
	}
	if len(params) > 0 {
		path += "?"
		for i, p := range params {
			if i > 0 {
				path += "&"
			}
			path += p
		}
	}
	var slots []ScheduleSlot
	err := c.get(path, &slots)
	return slots, err
}

func (c *Client) GetTariffs() ([]Tariff, error) {
	var tariffs []Tariff
	err := c.get("/api/tariffs", &tariffs)
	return tariffs, err
}

func (c *Client) GetAvailableSlots(masterID, serviceID int, date string) ([]AvailableSlot, error) {
	path := fmt.Sprintf("/api/available-slots?master_id=%d&service_id=%d&date=%s", masterID, serviceID, date)
	var slots []AvailableSlot
	err := c.get(path, &slots)
	return slots, err
}

func (c *Client) GetDashboard(masterID int) (*Dashboard, error) {
	var d Dashboard
	err := c.get(fmt.Sprintf("/api/master/%d/dashboard", masterID), &d)
	return &d, err
}

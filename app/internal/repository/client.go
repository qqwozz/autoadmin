package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetClientByID(id int) (*model.Client, error) {
	var c model.Client
	err := db.DB.QueryRow(
		"SELECT id, telegram_id, name, phone, no_show_count, is_blocked, blocked_until, created_at, updated_at FROM clients WHERE id = ?", id,
	).Scan(&c.ID, &c.TelegramID, &c.Name, &c.Phone, &c.NoShowCount, &c.IsBlocked, &c.BlockedUntil, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func GetClientByTelegramID(telegramID int) (*model.Client, error) {
	var c model.Client
	err := db.DB.QueryRow(
		"SELECT id, telegram_id, name, phone, no_show_count, is_blocked, blocked_until, created_at, updated_at FROM clients WHERE telegram_id = ?", telegramID,
	).Scan(&c.ID, &c.TelegramID, &c.Name, &c.Phone, &c.NoShowCount, &c.IsBlocked, &c.BlockedUntil, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func ListClients() ([]model.Client, error) {
	rows, err := db.DB.Query("SELECT id, telegram_id, name, phone, no_show_count, is_blocked, blocked_until, created_at, updated_at FROM clients ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []model.Client
	for rows.Next() {
		var c model.Client
		if err := rows.Scan(&c.ID, &c.TelegramID, &c.Name, &c.Phone, &c.NoShowCount, &c.IsBlocked, &c.BlockedUntil, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func CreateClient(c *model.Client) error {
	res, err := db.DB.Exec(
		"INSERT INTO clients (telegram_id, name, phone) VALUES (?, ?, ?)",
		c.TelegramID, c.Name, c.Phone,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	c.ID = int(id)
	return nil
}

func UpdateClient(id int, c *model.Client) error {
	_, err := db.DB.Exec(
		"UPDATE clients SET name=?, phone=?, no_show_count=?, is_blocked=?, blocked_until=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		c.Name, c.Phone, c.NoShowCount, c.IsBlocked, c.BlockedUntil, id,
	)
	return err
}

func DeleteClient(id int) error {
	_, err := db.DB.Exec("DELETE FROM clients WHERE id = ?", id)
	return err
}

func IncrementNoShow(clientID int) error {
	_, err := db.DB.Exec("UPDATE clients SET no_show_count = no_show_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", clientID)
	return err
}

func BlockClient(clientID int, blockedUntil string) error {
	_, err := db.DB.Exec("UPDATE clients SET is_blocked = 1, blocked_until = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", blockedUntil, clientID)
	return err
}

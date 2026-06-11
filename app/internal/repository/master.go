package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetMasterByID(id int) (*model.Master, error) {
	var m model.Master
	err := db.DB.QueryRow(
		"SELECT id, telegram_id, name, phone, description, subscription_until, tariff_id, is_active, created_at, updated_at FROM masters WHERE id = ?", id,
	).Scan(&m.ID, &m.TelegramID, &m.Name, &m.Phone, &m.Description, &m.SubscriptionUntil, &m.TariffID, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func GetMasterByTelegramID(telegramID int) (*model.Master, error) {
	var m model.Master
	err := db.DB.QueryRow(
		"SELECT id, telegram_id, name, phone, description, subscription_until, tariff_id, is_active, created_at, updated_at FROM masters WHERE telegram_id = ?", telegramID,
	).Scan(&m.ID, &m.TelegramID, &m.Name, &m.Phone, &m.Description, &m.SubscriptionUntil, &m.TariffID, &m.IsActive, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func ListMasters(telegramID *int) ([]model.Master, error) {
	query := "SELECT id, telegram_id, name, phone, description, subscription_until, tariff_id, is_active, created_at, updated_at FROM masters WHERE 1=1"
	args := []any{}

	if telegramID != nil {
		query += " AND telegram_id = ?"
		args = append(args, *telegramID)
	}
	query += " ORDER BY id DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var masters []model.Master
	for rows.Next() {
		var m model.Master
		if err := rows.Scan(&m.ID, &m.TelegramID, &m.Name, &m.Phone, &m.Description, &m.SubscriptionUntil, &m.TariffID, &m.IsActive, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		masters = append(masters, m)
	}
	return masters, nil
}

func CreateMaster(m *model.Master) error {
	res, err := db.DB.Exec(
		"INSERT INTO masters (telegram_id, name, phone, description, tariff_id, is_active) VALUES (?, ?, ?, ?, ?, ?)",
		m.TelegramID, m.Name, m.Phone, m.Description, m.TariffID, m.IsActive,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = int(id)
	return nil
}

func UpdateMaster(id int, m *model.Master) error {
	_, err := db.DB.Exec(
		"UPDATE masters SET name=?, phone=?, description=?, tariff_id=?, is_active=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		m.Name, m.Phone, m.Description, m.TariffID, m.IsActive, id,
	)
	return err
}

func DeleteMaster(id int) error {
	_, err := db.DB.Exec("DELETE FROM masters WHERE id = ?", id)
	return err
}

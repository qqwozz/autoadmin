package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetTariffByID(id int) (*model.Tariff, error) {
	var t model.Tariff
	err := db.DB.QueryRow(
		"SELECT id, name, price, meeting_limit, client_limit, duration_days, is_active, created_at FROM tariffs WHERE id = ?", id,
	).Scan(&t.ID, &t.Name, &t.Price, &t.MeetingLimit, &t.ClientLimit, &t.DurationDays, &t.IsActive, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ListTariffs() ([]model.Tariff, error) {
	rows, err := db.DB.Query("SELECT id, name, price, meeting_limit, client_limit, duration_days, is_active, created_at FROM tariffs ORDER BY price ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tariffs []model.Tariff
	for rows.Next() {
		var t model.Tariff
		if err := rows.Scan(&t.ID, &t.Name, &t.Price, &t.MeetingLimit, &t.ClientLimit, &t.DurationDays, &t.IsActive, &t.CreatedAt); err != nil {
			return nil, err
		}
		tariffs = append(tariffs, t)
	}
	return tariffs, nil
}

func CreateTariff(t *model.Tariff) error {
	if t.DurationDays == 0 {
		t.DurationDays = 30
	}
	res, err := db.DB.Exec(
		"INSERT INTO tariffs (name, price, meeting_limit, client_limit, duration_days, is_active) VALUES (?, ?, ?, ?, ?, ?)",
		t.Name, t.Price, t.MeetingLimit, t.ClientLimit, t.DurationDays, t.IsActive,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = int(id)
	return nil
}

func UpdateTariff(id int, t *model.Tariff) error {
	_, err := db.DB.Exec(
		"UPDATE tariffs SET name=?, price=?, meeting_limit=?, client_limit=?, duration_days=?, is_active=? WHERE id=?",
		t.Name, t.Price, t.MeetingLimit, t.ClientLimit, t.DurationDays, t.IsActive, id,
	)
	return err
}

func DeleteTariff(id int) error {
	_, err := db.DB.Exec("DELETE FROM tariffs WHERE id = ?", id)
	return err
}

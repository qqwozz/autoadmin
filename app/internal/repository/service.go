package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetServiceByID(id int) (*model.Service, error) {
	var s model.Service
	err := db.DB.QueryRow(
		"SELECT id, master_id, name, duration_minutes, price, created_at FROM services WHERE id = ?", id,
	).Scan(&s.ID, &s.MasterID, &s.Name, &s.DurationMinutes, &s.Price, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func ListServices(masterID *int) ([]model.Service, error) {
	query := "SELECT id, master_id, name, duration_minutes, price, created_at FROM services WHERE 1=1"
	args := []any{}

	if masterID != nil {
		query += " AND master_id = ?"
		args = append(args, *masterID)
	}
	query += " ORDER BY id DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(&s.ID, &s.MasterID, &s.Name, &s.DurationMinutes, &s.Price, &s.CreatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func CreateService(s *model.Service) error {
	res, err := db.DB.Exec(
		"INSERT INTO services (master_id, name, duration_minutes, price) VALUES (?, ?, ?, ?)",
		s.MasterID, s.Name, s.DurationMinutes, s.Price,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	s.ID = int(id)
	return nil
}

func UpdateService(id int, s *model.Service) error {
	_, err := db.DB.Exec(
		"UPDATE services SET name=?, duration_minutes=?, price=? WHERE id=?",
		s.Name, s.DurationMinutes, s.Price, id,
	)
	return err
}

func DeleteService(id int) error {
	_, err := db.DB.Exec("DELETE FROM services WHERE id = ?", id)
	return err
}

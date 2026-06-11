package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetWorkingHours(masterID *int) ([]model.WorkingHour, error) {
	query := "SELECT id, master_id, day_of_week, time_start, time_end, break_start, break_end, is_day_off, created_at, updated_at FROM working_hours WHERE 1=1"
	args := []any{}

	if masterID != nil {
		query += " AND master_id = ?"
		args = append(args, *masterID)
	}
	query += " ORDER BY master_id, day_of_week"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hours []model.WorkingHour
	for rows.Next() {
		var wh model.WorkingHour
		if err := rows.Scan(&wh.ID, &wh.MasterID, &wh.DayOfWeek, &wh.TimeStart, &wh.TimeEnd, &wh.BreakStart, &wh.BreakEnd, &wh.IsDayOff, &wh.CreatedAt, &wh.UpdatedAt); err != nil {
			return nil, err
		}
		hours = append(hours, wh)
	}
	return hours, nil
}

func GetWorkingHourByID(id int) (*model.WorkingHour, error) {
	var wh model.WorkingHour
	err := db.DB.QueryRow(
		"SELECT id, master_id, day_of_week, time_start, time_end, break_start, break_end, is_day_off, created_at, updated_at FROM working_hours WHERE id = ?", id,
	).Scan(&wh.ID, &wh.MasterID, &wh.DayOfWeek, &wh.TimeStart, &wh.TimeEnd, &wh.BreakStart, &wh.BreakEnd, &wh.IsDayOff, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func GetWorkingHoursByMasterAndDay(masterID, dayOfWeek int) (*model.WorkingHour, error) {
	var wh model.WorkingHour
	err := db.DB.QueryRow(
		"SELECT id, master_id, day_of_week, time_start, time_end, break_start, break_end, is_day_off, created_at, updated_at FROM working_hours WHERE master_id = ? AND day_of_week = ?", masterID, dayOfWeek,
	).Scan(&wh.ID, &wh.MasterID, &wh.DayOfWeek, &wh.TimeStart, &wh.TimeEnd, &wh.BreakStart, &wh.BreakEnd, &wh.IsDayOff, &wh.CreatedAt, &wh.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &wh, nil
}

func CreateWorkingHour(wh *model.WorkingHour) error {
	res, err := db.DB.Exec(
		"INSERT INTO working_hours (master_id, day_of_week, time_start, time_end, break_start, break_end, is_day_off) VALUES (?, ?, ?, ?, ?, ?, ?)",
		wh.MasterID, wh.DayOfWeek, wh.TimeStart, wh.TimeEnd, wh.BreakStart, wh.BreakEnd, wh.IsDayOff,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	wh.ID = int(id)
	return nil
}

func UpdateWorkingHour(id int, wh *model.WorkingHour) error {
	_, err := db.DB.Exec(
		"UPDATE working_hours SET time_start=?, time_end=?, break_start=?, break_end=?, is_day_off=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		wh.TimeStart, wh.TimeEnd, wh.BreakStart, wh.BreakEnd, wh.IsDayOff, id,
	)
	return err
}

func DeleteWorkingHour(id int) error {
	_, err := db.DB.Exec("DELETE FROM working_hours WHERE id = ?", id)
	return err
}

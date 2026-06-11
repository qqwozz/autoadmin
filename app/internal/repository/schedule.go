package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetScheduleSlotByID(id int) (*model.ScheduleSlot, error) {
	var s model.ScheduleSlot
	err := db.DB.QueryRow(
		"SELECT id, master_id, client_id, service_id, start_time, end_time, status, details, confirm_code, confirm_deadline, checkin_time, cancelled_by, cancelled_at, created_at, updated_at FROM schedule_slots WHERE id = ?", id,
	).Scan(&s.ID, &s.MasterID, &s.ClientID, &s.ServiceID, &s.StartTime, &s.EndTime, &s.Status, &s.Details, &s.ConfirmCode, &s.ConfirmDeadline, &s.CheckinTime, &s.CancelledBy, &s.CancelledAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func ListScheduleSlots(masterID *int, status *string) ([]model.ScheduleSlot, error) {
	query := "SELECT id, master_id, client_id, service_id, start_time, end_time, status, details, confirm_code, confirm_deadline, checkin_time, cancelled_by, cancelled_at, created_at, updated_at FROM schedule_slots WHERE 1=1"
	args := []any{}

	if masterID != nil {
		query += " AND master_id = ?"
		args = append(args, *masterID)
	}
	if status != nil {
		query += " AND status = ?"
		args = append(args, *status)
	}
	query += " ORDER BY start_time DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.ScheduleSlot
	for rows.Next() {
		var s model.ScheduleSlot
		if err := rows.Scan(&s.ID, &s.MasterID, &s.ClientID, &s.ServiceID, &s.StartTime, &s.EndTime, &s.Status, &s.Details, &s.ConfirmCode, &s.ConfirmDeadline, &s.CheckinTime, &s.CancelledBy, &s.CancelledAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func GetOccupiedSlots(masterID int, date string) ([]model.ScheduleSlot, error) {
	rows, err := db.DB.Query(
		"SELECT id, master_id, client_id, service_id, start_time, end_time, status, details, confirm_code, confirm_deadline, checkin_time, cancelled_by, cancelled_at, created_at, updated_at FROM schedule_slots WHERE master_id = ? AND date(start_time) = ? AND status NOT IN ('cancelled', 'no_show') ORDER BY start_time",
		masterID, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.ScheduleSlot
	for rows.Next() {
		var s model.ScheduleSlot
		if err := rows.Scan(&s.ID, &s.MasterID, &s.ClientID, &s.ServiceID, &s.StartTime, &s.EndTime, &s.Status, &s.Details, &s.ConfirmCode, &s.ConfirmDeadline, &s.CheckinTime, &s.CancelledBy, &s.CancelledAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func CreateScheduleSlot(s *model.ScheduleSlot) error {
	if s.Status == "" {
		s.Status = "pending_confirmation"
	}
	res, err := db.DB.Exec(
		"INSERT INTO schedule_slots (master_id, client_id, service_id, start_time, end_time, status, details) VALUES (?, ?, ?, ?, ?, ?, ?)",
		s.MasterID, s.ClientID, s.ServiceID, s.StartTime, s.EndTime, s.Status, s.Details,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	s.ID = int(id)
	return nil
}

func UpdateScheduleSlot(id int, s *model.ScheduleSlot) error {
	_, err := db.DB.Exec(
		"UPDATE schedule_slots SET status=?, details=?, checkin_time=?, cancelled_by=?, cancelled_at=?, confirm_code=?, confirm_deadline=?, updated_at=CURRENT_TIMESTAMP WHERE id=?",
		s.Status, s.Details, s.CheckinTime, s.CancelledBy, s.CancelledAt, s.ConfirmCode, s.ConfirmDeadline, id,
	)
	return err
}

func DeleteScheduleSlot(id int) error {
	_, err := db.DB.Exec("DELETE FROM schedule_slots WHERE id = ?", id)
	return err
}

package repository

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
)

func GetNoShowSetting(masterID int) (*model.NoShowSetting, error) {
	var ns model.NoShowSetting
	err := db.DB.QueryRow(
		"SELECT master_id, enable_penalty, penalty_percent, no_show_limit, block_days, confirm_minutes, checkin_method, remind_minutes, updated_at FROM no_show_settings WHERE master_id = ?", masterID,
	).Scan(&ns.MasterID, &ns.EnablePenalty, &ns.PenaltyPercent, &ns.NoShowLimit, &ns.BlockDays, &ns.ConfirmMinutes, &ns.CheckinMethod, &ns.RemindMinutes, &ns.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &ns, nil
}

func ListNoShowSettings() ([]model.NoShowSetting, error) {
	rows, err := db.DB.Query("SELECT master_id, enable_penalty, penalty_percent, no_show_limit, block_days, confirm_minutes, checkin_method, remind_minutes, updated_at FROM no_show_settings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []model.NoShowSetting
	for rows.Next() {
		var ns model.NoShowSetting
		if err := rows.Scan(&ns.MasterID, &ns.EnablePenalty, &ns.PenaltyPercent, &ns.NoShowLimit, &ns.BlockDays, &ns.ConfirmMinutes, &ns.CheckinMethod, &ns.RemindMinutes, &ns.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, ns)
	}
	return settings, nil
}

func UpsertNoShowSetting(ns *model.NoShowSetting) error {
	_, err := db.DB.Exec(
		"INSERT OR REPLACE INTO no_show_settings (master_id, enable_penalty, penalty_percent, no_show_limit, block_days, confirm_minutes, checkin_method, remind_minutes, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)",
		ns.MasterID, ns.EnablePenalty, ns.PenaltyPercent, ns.NoShowLimit, ns.BlockDays, ns.ConfirmMinutes, ns.CheckinMethod, ns.RemindMinutes,
	)
	return err
}

func ListBlacklist(masterID *int) ([]model.BlacklistEntry, error) {
	query := "SELECT id, master_id, client_id, reason, created_at FROM master_blacklist WHERE 1=1"
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

	var entries []model.BlacklistEntry
	for rows.Next() {
		var b model.BlacklistEntry
		if err := rows.Scan(&b.ID, &b.MasterID, &b.ClientID, &b.Reason, &b.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, b)
	}
	return entries, nil
}

func CreateBlacklistEntry(b *model.BlacklistEntry) error {
	res, err := db.DB.Exec(
		"INSERT INTO master_blacklist (master_id, client_id, reason) VALUES (?, ?, ?)",
		b.MasterID, b.ClientID, b.Reason,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	b.ID = int(id)
	return nil
}

func DeleteBlacklistEntry(id int) error {
	_, err := db.DB.Exec("DELETE FROM master_blacklist WHERE id = ?", id)
	return err
}

func ListBlockedSlots(masterID *int) ([]model.BlockedSlot, error) {
	query := "SELECT id, master_id, start_time, end_time, reason, created_at FROM blocked_slots WHERE 1=1"
	args := []any{}

	if masterID != nil {
		query += " AND master_id = ?"
		args = append(args, *masterID)
	}
	query += " ORDER BY start_time DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.BlockedSlot
	for rows.Next() {
		var bs model.BlockedSlot
		if err := rows.Scan(&bs.ID, &bs.MasterID, &bs.StartTime, &bs.EndTime, &bs.Reason, &bs.CreatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, bs)
	}
	return slots, nil
}

func GetBlockedSlotsByMasterAndDate(masterID int, date string) ([]model.BlockedSlot, error) {
	rows, err := db.DB.Query(
		"SELECT id, master_id, start_time, end_time, reason, created_at FROM blocked_slots WHERE master_id = ? AND date(start_time) = ? ORDER BY start_time",
		masterID, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.BlockedSlot
	for rows.Next() {
		var bs model.BlockedSlot
		if err := rows.Scan(&bs.ID, &bs.MasterID, &bs.StartTime, &bs.EndTime, &bs.Reason, &bs.CreatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, bs)
	}
	return slots, nil
}

func CreateBlockedSlot(bs *model.BlockedSlot) error {
	res, err := db.DB.Exec(
		"INSERT INTO blocked_slots (master_id, start_time, end_time, reason) VALUES (?, ?, ?, ?)",
		bs.MasterID, bs.StartTime, bs.EndTime, bs.Reason,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	bs.ID = int(id)
	return nil
}

func DeleteBlockedSlot(id int) error {
	_, err := db.DB.Exec("DELETE FROM blocked_slots WHERE id = ?", id)
	return err
}

func ListRefCodes(masterID *int) ([]model.RefCode, error) {
	query := "SELECT id, master_id, short_id, qr_code_url, is_active, created_at FROM master_ref_codes WHERE 1=1"
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

	var codes []model.RefCode
	for rows.Next() {
		var rc model.RefCode
		if err := rows.Scan(&rc.ID, &rc.MasterID, &rc.ShortID, &rc.QRCodeURL, &rc.IsActive, &rc.CreatedAt); err != nil {
			return nil, err
		}
		codes = append(codes, rc)
	}
	return codes, nil
}

func CreateRefCode(rc *model.RefCode) error {
	res, err := db.DB.Exec(
		"INSERT INTO master_ref_codes (master_id, short_id, qr_code_url, is_active) VALUES (?, ?, ?, ?)",
		rc.MasterID, rc.ShortID, rc.QRCodeURL, rc.IsActive,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	rc.ID = int(id)
	return nil
}

func DeleteRefCode(id int) error {
	_, err := db.DB.Exec("DELETE FROM master_ref_codes WHERE id = ?", id)
	return err
}

func ListPayments(masterID *int) ([]model.SubscriptionPayment, error) {
	query := "SELECT id, master_id, tariff_id, amount, payment_method, payment_id, status, paid_at, valid_from, valid_until, created_at FROM subscription_payments WHERE 1=1"
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

	var payments []model.SubscriptionPayment
	for rows.Next() {
		var p model.SubscriptionPayment
		if err := rows.Scan(&p.ID, &p.MasterID, &p.TariffID, &p.Amount, &p.PaymentMethod, &p.PaymentID, &p.Status, &p.PaidAt, &p.ValidFrom, &p.ValidUntil, &p.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func CreatePayment(p *model.SubscriptionPayment) error {
	if p.Status == "" {
		p.Status = "pending"
	}
	res, err := db.DB.Exec(
		"INSERT INTO subscription_payments (master_id, tariff_id, amount, payment_method, payment_id, status) VALUES (?, ?, ?, ?, ?, ?)",
		p.MasterID, p.TariffID, p.Amount, p.PaymentMethod, p.PaymentID, p.Status,
	)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.ID = int(id)
	return nil
}

func ListNotificationsLog(userType, userID *string) ([]model.NotificationLog, error) {
	query := "SELECT id, user_type, user_id, telegram_id, notification_type, message, status, error_message, created_at FROM notifications_log WHERE 1=1"
	args := []any{}

	if userType != nil {
		query += " AND user_type = ?"
		args = append(args, *userType)
	}
	if userID != nil {
		query += " AND user_id = ?"
		args = append(args, *userID)
	}
	query += " ORDER BY id DESC LIMIT 100"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.NotificationLog
	for rows.Next() {
		var nl model.NotificationLog
		if err := rows.Scan(&nl.ID, &nl.UserType, &nl.UserID, &nl.TelegramID, &nl.NotificationType, &nl.Message, &nl.Status, &nl.ErrorMessage, &nl.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, nl)
	}
	return logs, nil
}

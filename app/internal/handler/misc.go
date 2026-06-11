package handler

import (
	"autoadmin/internal/db"
	"autoadmin/internal/model"
	"autoadmin/internal/repository"
	"net/http"
	"strconv"
)

func HandleGetScheduleSlots(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, _ := parseOptionalInt(mid)
		masterID = &v
	}

	var statusStr *string
	if s := r.URL.Query().Get("status"); s != "" {
		statusStr = &s
	}

	slots, err := repository.ListScheduleSlots(masterID, statusStr)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, slots)
}

func HandleGetScheduleSlot(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	slot, err := repository.GetScheduleSlotByID(id)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, slot)
}

func HandleCreateScheduleSlot(w http.ResponseWriter, r *http.Request) {
	var s struct {
		MasterID  int     `json:"master_id"`
		ClientID  *int    `json:"client_id"`
		ServiceID *int    `json:"service_id"`
		StartTime string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
		Status    string  `json:"status"`
		Details   *string `json:"details"`
	}
	if err := decodeJSON(r, &s); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	slot := &model.ScheduleSlot{
		MasterID:  s.MasterID,
		ClientID:  s.ClientID,
		ServiceID: s.ServiceID,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Status:    s.Status,
		Details:   s.Details,
	}
	if err := repository.CreateScheduleSlot(slot); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, slot)
}

func HandleUpdateScheduleSlot(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var s struct {
		Status          *string `json:"status"`
		Details         *string `json:"details"`
		CheckinTime     *string `json:"checkin_time"`
		CancelledBy     *string `json:"cancelled_by"`
		CancelledAt     *string `json:"cancelled_at"`
		ConfirmCode     *string `json:"confirm_code"`
		ConfirmDeadline *string `json:"confirm_deadline"`
	}
	if err := decodeJSON(r, &s); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	slot := &model.ScheduleSlot{
		Status:          derefStr(s.Status),
		Details:         s.Details,
		CheckinTime:     s.CheckinTime,
		CancelledBy:     s.CancelledBy,
		CancelledAt:     s.CancelledAt,
		ConfirmCode:     s.ConfirmCode,
		ConfirmDeadline: s.ConfirmDeadline,
	}
	if err := repository.UpdateScheduleSlot(id, slot); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	slot.ID = id
	respondJSON(w, 200, slot)
}

func HandleDeleteScheduleSlot(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteScheduleSlot(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetWorkingHours(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, _ := parseOptionalInt(mid)
		masterID = &v
	}

	hours, err := repository.GetWorkingHours(masterID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, hours)
}

func HandleGetWorkingHour(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	wh, err := repository.GetWorkingHourByID(id)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, wh)
}

func HandleCreateWorkingHour(w http.ResponseWriter, r *http.Request) {
	var wh struct {
		MasterID   int     `json:"master_id"`
		DayOfWeek  int     `json:"day_of_week"`
		TimeStart  string  `json:"time_start"`
		TimeEnd    string  `json:"time_end"`
		BreakStart *string `json:"break_start"`
		BreakEnd   *string `json:"break_end"`
		IsDayOff   int     `json:"is_day_off"`
	}
	if err := decodeJSON(r, &wh); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	hour := &model.WorkingHour{
		MasterID:   wh.MasterID,
		DayOfWeek:  wh.DayOfWeek,
		TimeStart:  wh.TimeStart,
		TimeEnd:    wh.TimeEnd,
		BreakStart: wh.BreakStart,
		BreakEnd:   wh.BreakEnd,
		IsDayOff:   wh.IsDayOff,
	}
	if err := repository.CreateWorkingHour(hour); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, hour)
}

func HandleUpdateWorkingHour(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var wh struct {
		TimeStart  string  `json:"time_start"`
		TimeEnd    string  `json:"time_end"`
		BreakStart *string `json:"break_start"`
		BreakEnd   *string `json:"break_end"`
		IsDayOff   int     `json:"is_day_off"`
	}
	if err := decodeJSON(r, &wh); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	hour := &model.WorkingHour{
		TimeStart:  wh.TimeStart,
		TimeEnd:    wh.TimeEnd,
		BreakStart: wh.BreakStart,
		BreakEnd:   wh.BreakEnd,
		IsDayOff:   wh.IsDayOff,
	}
	if err := repository.UpdateWorkingHour(id, hour); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	hour.ID = id
	respondJSON(w, 200, hour)
}

func HandleDeleteWorkingHour(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteWorkingHour(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetTariffs(w http.ResponseWriter, r *http.Request) {
	tariffs, err := repository.ListTariffs()
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, tariffs)
}

func HandleGetTariff(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	tariff, err := repository.GetTariffByID(id)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, tariff)
}

func HandleCreateTariff(w http.ResponseWriter, r *http.Request) {
	var t struct {
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		MeetingLimit *int    `json:"meeting_limit"`
		ClientLimit  *int    `json:"client_limit"`
		DurationDays int     `json:"duration_days"`
		IsActive     int     `json:"is_active"`
	}
	if err := decodeJSON(r, &t); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	tariff := &model.Tariff{
		Name:         t.Name,
		Price:        t.Price,
		MeetingLimit: t.MeetingLimit,
		ClientLimit:  t.ClientLimit,
		DurationDays: t.DurationDays,
		IsActive:     t.IsActive,
	}
	if err := repository.CreateTariff(tariff); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, tariff)
}

func HandleUpdateTariff(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var t struct {
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		MeetingLimit *int    `json:"meeting_limit"`
		ClientLimit  *int    `json:"client_limit"`
		DurationDays int     `json:"duration_days"`
		IsActive     int     `json:"is_active"`
	}
	if err := decodeJSON(r, &t); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	tariff := &model.Tariff{
		Name:         t.Name,
		Price:        t.Price,
		MeetingLimit: t.MeetingLimit,
		ClientLimit:  t.ClientLimit,
		DurationDays: t.DurationDays,
		IsActive:     t.IsActive,
	}
	if err := repository.UpdateTariff(id, tariff); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	tariff.ID = id
	respondJSON(w, 200, tariff)
}

func HandleDeleteTariff(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteTariff(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetNoShowSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := repository.ListNoShowSettings()
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, settings)
}

func HandleGetNoShowSetting(w http.ResponseWriter, r *http.Request) {
	masterID, ok := getID(r, "masterId")
	if !ok {
		respondError(w, 400, "invalid master_id")
		return
	}

	setting, err := repository.GetNoShowSetting(masterID)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, setting)
}

func HandleUpdateNoShowSetting(w http.ResponseWriter, r *http.Request) {
	masterID, ok := getID(r, "masterId")
	if !ok {
		respondError(w, 400, "invalid master_id")
		return
	}

	var ns struct {
		EnablePenalty  int    `json:"enable_penalty"`
		PenaltyPercent int    `json:"penalty_percent"`
		NoShowLimit    int    `json:"no_show_limit"`
		BlockDays      int    `json:"block_days"`
		ConfirmMinutes int    `json:"confirm_minutes"`
		CheckinMethod  string `json:"checkin_method"`
		RemindMinutes  int    `json:"remind_minutes"`
	}
	if err := decodeJSON(r, &ns); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	setting := &model.NoShowSetting{
		MasterID:       masterID,
		EnablePenalty:  ns.EnablePenalty,
		PenaltyPercent: ns.PenaltyPercent,
		NoShowLimit:    ns.NoShowLimit,
		BlockDays:      ns.BlockDays,
		ConfirmMinutes: ns.ConfirmMinutes,
		CheckinMethod:  ns.CheckinMethod,
		RemindMinutes:  ns.RemindMinutes,
	}
	if err := repository.UpsertNoShowSetting(setting); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, setting)
}

func HandleGetBlacklist(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, _ := parseOptionalInt(mid)
		masterID = &v
	}

	entries, err := repository.ListBlacklist(masterID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, entries)
}

func HandleCreateBlacklistEntry(w http.ResponseWriter, r *http.Request) {
	var b struct {
		MasterID int     `json:"master_id"`
		ClientID int     `json:"client_id"`
		Reason   *string `json:"reason"`
	}
	if err := decodeJSON(r, &b); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	entry := &model.BlacklistEntry{
		MasterID: b.MasterID,
		ClientID: b.ClientID,
		Reason:   b.Reason,
	}
	if err := repository.CreateBlacklistEntry(entry); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, entry)
}

func HandleDeleteBlacklistEntry(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteBlacklistEntry(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetBlockedSlots(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, _ := parseOptionalInt(mid)
		masterID = &v
	}

	slots, err := repository.ListBlockedSlots(masterID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, slots)
}

func HandleCreateBlockedSlot(w http.ResponseWriter, r *http.Request) {
	var bs struct {
		MasterID  int     `json:"master_id"`
		StartTime string  `json:"start_time"`
		EndTime   string  `json:"end_time"`
		Reason    *string `json:"reason"`
	}
	if err := decodeJSON(r, &bs); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	slot := &model.BlockedSlot{
		MasterID:  bs.MasterID,
		StartTime: bs.StartTime,
		EndTime:   bs.EndTime,
		Reason:    bs.Reason,
	}
	if err := repository.CreateBlockedSlot(slot); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, slot)
}

func HandleDeleteBlockedSlot(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteBlockedSlot(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetRefCodes(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, _ := parseOptionalInt(mid)
		masterID = &v
	}

	codes, err := repository.ListRefCodes(masterID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, codes)
}

func HandleCreateRefCode(w http.ResponseWriter, r *http.Request) {
	var rc struct {
		MasterID  int     `json:"master_id"`
		ShortID   string  `json:"short_id"`
		QRCodeURL *string `json:"qr_code_url"`
		IsActive  int     `json:"is_active"`
	}
	if err := decodeJSON(r, &rc); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	code := &model.RefCode{
		MasterID:  rc.MasterID,
		ShortID:   rc.ShortID,
		QRCodeURL: rc.QRCodeURL,
		IsActive:  rc.IsActive,
	}
	if err := repository.CreateRefCode(code); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, code)
}

func HandleDeleteRefCode(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteRefCode(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetPayments(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, _ := parseOptionalInt(mid)
		masterID = &v
	}

	payments, err := repository.ListPayments(masterID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, payments)
}

func HandleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var p struct {
		MasterID      int      `json:"master_id"`
		TariffID      *int     `json:"tariff_id"`
		Amount        float64  `json:"amount"`
		PaymentMethod *string  `json:"payment_method"`
		PaymentID     *string  `json:"payment_id"`
		Status        string   `json:"status"`
	}
	if err := decodeJSON(r, &p); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	payment := &model.SubscriptionPayment{
		MasterID:      p.MasterID,
		TariffID:      p.TariffID,
		Amount:        p.Amount,
		PaymentMethod: p.PaymentMethod,
		PaymentID:     p.PaymentID,
		Status:        p.Status,
	}
	if err := repository.CreatePayment(payment); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, payment)
}

func HandleGetNotificationsLog(w http.ResponseWriter, r *http.Request) {
	var userType, userID *string
	if ut := r.URL.Query().Get("user_type"); ut != "" {
		userType = &ut
	}
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID = &uid
	}

	logs, err := repository.ListNotificationsLog(userType, userID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, logs)
}

func HandleMasterDashboard(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	master, err := repository.GetMasterByID(id)
	if err != nil {
		respondError(w, 404, "master not found")
		return
	}

	services, _ := repository.ListServices(&id)
	hours, _ := repository.GetWorkingHours(&id)

	var activeClients int
	db.DB.QueryRow("SELECT COUNT(DISTINCT client_id) FROM master_client_bindings WHERE master_id = ? AND is_active = 1", id).Scan(&activeClients)

	var totalBookings, pendingBookings int
	db.DB.QueryRow("SELECT COUNT(*) FROM schedule_slots WHERE master_id = ?", id).Scan(&totalBookings)
	db.DB.QueryRow("SELECT COUNT(*) FROM schedule_slots WHERE master_id = ? AND status = 'pending_confirmation'", id).Scan(&pendingBookings)

	respondJSON(w, 200, map[string]any{
		"master":           master,
		"services":         services,
		"working_hours":    hours,
		"active_clients":   activeClients,
		"total_bookings":   totalBookings,
		"pending_bookings": pendingBookings,
	})
}

func parseOptionalInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

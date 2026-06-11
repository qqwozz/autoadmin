package service

import (
	"autoadmin/internal/auth"
	"autoadmin/internal/model"
	"autoadmin/internal/repository"
	"errors"
	"fmt"
	"time"
)

var (
	ErrMasterNotFound = errors.New("master not found")
	ErrClientNotFound = errors.New("client not found")
	ErrClientBlocked  = errors.New("client is blocked")
	ErrInvalidCode    = errors.New("invalid confirmation code")
	ErrDeadlinePassed = errors.New("confirmation deadline passed")
	ErrSlotNotFound   = errors.New("slot not found")
	ErrNotPending     = errors.New("slot is not pending confirmation")
)

func LoginByTelegramID(telegramID int) (string, *model.Master, error) {
	master, err := repository.GetMasterByTelegramID(telegramID)
	if err != nil {
		return "", nil, ErrMasterNotFound
	}

	token, err := auth.GenerateToken(master.ID, master.TelegramID)
	if err != nil {
		return "", nil, err
	}

	return token, master, nil
}

func GetAvailableSlots(masterID int, date string, serviceID int) ([]model.AvailableSlot, error) {
	master, err := repository.GetMasterByID(masterID)
	if err != nil {
		return nil, ErrMasterNotFound
	}
	if master.IsActive == 0 {
		return nil, errors.New("master is not active")
	}

	svc, err := repository.GetServiceByID(serviceID)
	if err != nil {
		return nil, errors.New("service not found")
	}

	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}
	dayOfWeek := int(t.Weekday())

	wh, err := repository.GetWorkingHoursByMasterAndDay(masterID, dayOfWeek)
	if err != nil {
		return nil, errors.New("working hours not configured for this day")
	}
	if wh.IsDayOff == 1 {
		return []model.AvailableSlot{}, nil
	}

	occupied, err := repository.GetOccupiedSlots(masterID, date)
	if err != nil {
		return nil, err
	}

	blocked, err := repository.GetBlockedSlotsByMasterAndDate(masterID, date)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(svc.DurationMinutes) * time.Minute
	return calculateFreeSlots(wh, occupied, blocked, duration, date), nil
}

func calculateFreeSlots(wh *model.WorkingHour, occupied []model.ScheduleSlot, blocked []model.BlockedSlot, duration time.Duration, date string) []model.AvailableSlot {
	if wh.IsDayOff == 1 {
		return []model.AvailableSlot{}
	}

	var slots []model.AvailableSlot

	now := time.Now()
	today, _ := time.Parse("2006-01-02", date)

	workStartStr := fmt.Sprintf("%s %s", date, wh.TimeStart)
	workEndStr := fmt.Sprintf("%s %s", date, wh.TimeEnd)
	workStart, _ := time.Parse("2006-01-02 15:04", workStartStr)
	workEnd, _ := time.Parse("2006-01-02 15:04", workEndStr)

	var breakStart, breakEnd time.Time
	hasBreak := wh.BreakStart != nil && wh.BreakEnd != nil
	if hasBreak {
		bs, _ := time.Parse("15:04", *wh.BreakStart)
		be, _ := time.Parse("15:04", *wh.BreakEnd)
		breakStart = time.Date(today.Year(), today.Month(), today.Day(), bs.Hour(), bs.Minute(), 0, 0, time.UTC)
		breakEnd = time.Date(today.Year(), today.Month(), today.Day(), be.Hour(), be.Minute(), 0, 0, time.UTC)
	}

	type timeRange struct {
		start, end time.Time
	}

	var occupiedRanges []timeRange
	for _, s := range occupied {
		st, _ := time.Parse("2006-01-02 15:04:05", s.StartTime)
		et, _ := time.Parse("2006-01-02 15:04:05", s.EndTime)
		occupiedRanges = append(occupiedRanges, timeRange{st, et})
	}

	for _, b := range blocked {
		st, _ := time.Parse("2006-01-02 15:04:05", b.StartTime)
		et, _ := time.Parse("2006-01-02 15:04:05", b.EndTime)
		occupiedRanges = append(occupiedRanges, timeRange{st, et})
	}

	current := workStart
	for current.Add(duration).Before(workEnd) || current.Add(duration).Equal(workEnd) {
		slotEnd := current.Add(duration)

		if hasBreak && current.Before(breakEnd) && slotEnd.After(breakStart) {
			current = breakEnd
			continue
		}

		slotDateTime := time.Date(today.Year(), today.Month(), today.Day(), current.Hour(), current.Minute(), 0, 0, time.UTC)
		if today.Equal(now.Truncate(24*time.Hour)) && slotDateTime.Before(now) {
			current = current.Add(15 * time.Minute)
			continue
		}

		conflict := false
		conflictEnd := time.Time{}
		for _, r := range occupiedRanges {
			if current.Before(r.end) && slotEnd.After(r.start) {
				conflict = true
				if r.end.After(conflictEnd) {
					conflictEnd = r.end
				}
			}
		}

		if conflict {
			current = conflictEnd
			continue
		}

		slots = append(slots, model.AvailableSlot{
			StartTime: fmt.Sprintf("%s %s", date, current.Format("15:04:00")),
			EndTime:   fmt.Sprintf("%s %s", date, slotEnd.Format("15:04:00")),
		})
		current = current.Add(15 * time.Minute)
	}

	return slots
}

func ConfirmBooking(slotID int, code string) (*model.ScheduleSlot, error) {
	slot, err := repository.GetScheduleSlotByID(slotID)
	if err != nil {
		return nil, ErrSlotNotFound
	}

	if slot.Status != "pending_confirmation" {
		return nil, ErrNotPending
	}

	if slot.ConfirmCode == nil || *slot.ConfirmCode != code {
		return nil, ErrInvalidCode
	}

	if slot.ConfirmDeadline != nil {
		deadline, err := time.Parse("2006-01-02 15:04:05", *slot.ConfirmDeadline)
		if err == nil && time.Now().After(deadline) {
			slot.Status = "cancelled"
			slot.CancelledBy = strPtr("system")
			now := time.Now().Format("2006-01-02 15:04:05")
			slot.CancelledAt = &now
			repository.UpdateScheduleSlot(slotID, slot)
			return slot, ErrDeadlinePassed
		}
	}

	slot.Status = "confirmed"
	repository.UpdateScheduleSlot(slotID, slot)
	return slot, nil
}

func CancelBooking(slotID int, cancelledBy string) (*model.ScheduleSlot, error) {
	slot, err := repository.GetScheduleSlotByID(slotID)
	if err != nil {
		return nil, ErrSlotNotFound
	}

	if slot.Status == "cancelled" || slot.Status == "no_show" {
		return nil, errors.New("slot already cancelled or no-show")
	}

	slot.Status = "cancelled"
	slot.CancelledBy = &cancelledBy
	now := time.Now().Format("2006-01-02 15:04:05")
	slot.CancelledAt = &now

	repository.UpdateScheduleSlot(slotID, slot)
	return slot, nil
}

func MarkNoShow(slotID int) (*model.ScheduleSlot, error) {
	slot, err := repository.GetScheduleSlotByID(slotID)
	if err != nil {
		return nil, ErrSlotNotFound
	}

	if slot.Status != "confirmed" {
		return nil, errors.New("only confirmed slots can be marked as no-show")
	}

	slot.Status = "no_show"
	repository.UpdateScheduleSlot(slotID, slot)

	if slot.ClientID != nil {
		repository.IncrementNoShow(*slot.ClientID)

		ns, err := repository.GetNoShowSetting(slot.MasterID)
		if err == nil && ns.EnablePenalty == 1 {
			client, err := repository.GetClientByID(*slot.ClientID)
			if err == nil && client.NoShowCount >= ns.NoShowLimit {
				blockedUntil := time.Now().AddDate(0, 0, ns.BlockDays).Format("2006-01-02 15:04:05")
				repository.BlockClient(*slot.ClientID, blockedUntil)
			}
		}
	}

	return slot, nil
}

func strPtr(s string) *string {
	return &s
}

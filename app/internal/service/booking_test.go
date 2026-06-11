package service

import (
	"autoadmin/internal/model"
	"testing"
	"time"
)

func intPtr(i int) *int {
	return &i
}

func TestCalculateFreeSlots_NoOccupied(t *testing.T) {
	wh := &model.WorkingHour{
		TimeStart: "09:00",
		TimeEnd:   "17:00",
		IsDayOff:  0,
	}

	duration := 60 * time.Minute
	slots := calculateFreeSlots(wh, nil, nil, duration, "2026-06-15")

	if len(slots) == 0 {
		t.Fatal("expected free slots, got none")
	}

	for _, s := range slots {
		if s.StartTime == "" || s.EndTime == "" {
			t.Error("slot has empty time")
		}
	}
}

func TestCalculateFreeSlots_WithOccupied(t *testing.T) {
	wh := &model.WorkingHour{
		TimeStart: "09:00",
		TimeEnd:   "12:00",
		IsDayOff:  0,
	}

	occupied := []model.ScheduleSlot{
		{
			StartTime: "2026-06-15 10:00:00",
			EndTime:   "2026-06-15 11:00:00",
			Status:    "confirmed",
		},
	}

	duration := 60 * time.Minute
	slots := calculateFreeSlots(wh, occupied, nil, duration, "2026-06-15")

	for _, s := range slots {
		start, _ := time.Parse("2006-01-02 15:04:00", s.StartTime)
		end, _ := time.Parse("2006-01-02 15:04:00", s.EndTime)

		occStart, _ := time.Parse("2006-01-02 15:04:05", "2026-06-15 10:00:00")
		occEnd, _ := time.Parse("2006-01-02 15:04:05", "2026-06-15 11:00:00")

		if start.Before(occEnd) && end.After(occStart) {
			t.Errorf("slot %s-%s overlaps with occupied slot", s.StartTime, s.EndTime)
		}
	}
}

func TestCalculateFreeSlots_DayOff(t *testing.T) {
	wh := &model.WorkingHour{
		TimeStart: "09:00",
		TimeEnd:   "17:00",
		IsDayOff:  1,
	}

	duration := 60 * time.Minute
	slots := calculateFreeSlots(wh, nil, nil, duration, "2026-06-15")

	if len(slots) != 0 {
		t.Errorf("expected no slots for day off, got %d", len(slots))
	}
}

func TestCalculateFreeSlots_WithBreak(t *testing.T) {
	wh := &model.WorkingHour{
		TimeStart:  "09:00",
		TimeEnd:    "17:00",
		BreakStart: strPtr("12:00"),
		BreakEnd:   strPtr("13:00"),
		IsDayOff:   0,
	}

	duration := 60 * time.Minute
	slots := calculateFreeSlots(wh, nil, nil, duration, "2026-06-15")

	for _, s := range slots {
		start, _ := time.Parse("2006-01-02 15:04:00", s.StartTime)
		end, _ := time.Parse("2006-01-02 15:04:00", s.EndTime)

		breakStart, _ := time.Parse("15:04", "12:00")
		breakEnd, _ := time.Parse("15:04", "13:00")

		if start.Before(breakEnd) && end.After(breakStart) {
			t.Errorf("slot %s-%s overlaps with break", s.StartTime, s.EndTime)
		}
	}
}

func TestCalculateFreeSlots_WithBlocked(t *testing.T) {
	wh := &model.WorkingHour{
		TimeStart: "09:00",
		TimeEnd:   "17:00",
		IsDayOff:  0,
	}

	blocked := []model.BlockedSlot{
		{
			StartTime: "2026-06-15 14:00:00",
			EndTime:   "2026-06-15 15:00:00",
		},
	}

	duration := 60 * time.Minute
	slots := calculateFreeSlots(wh, nil, blocked, duration, "2026-06-15")

	for _, s := range slots {
		start, _ := time.Parse("2006-01-02 15:04:00", s.StartTime)
		end, _ := time.Parse("2006-01-02 15:04:00", s.EndTime)

		blkStart, _ := time.Parse("2006-01-02 15:04:05", "2026-06-15 14:00:00")
		blkEnd, _ := time.Parse("2006-01-02 15:04:05", "2026-06-15 15:00:00")

		if start.Before(blkEnd) && end.After(blkStart) {
			t.Errorf("slot %s-%s overlaps with blocked slot", s.StartTime, s.EndTime)
		}
	}
}

func TestCalculateFreeSlots_ShortService(t *testing.T) {
	wh := &model.WorkingHour{
		TimeStart: "09:00",
		TimeEnd:   "10:00",
		IsDayOff:  0,
	}

	duration := 30 * time.Minute
	slots := calculateFreeSlots(wh, nil, nil, duration, "2026-06-15")

	if len(slots) < 1 {
		t.Error("expected at least 1 slot for 30min service in 1hr window")
	}
}



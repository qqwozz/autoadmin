package handler

import (
	"autoadmin/internal/model"
	"autoadmin/internal/repository"
	"autoadmin/internal/service"
	"net/http"
	"strconv"
)

func HandleGetServices(w http.ResponseWriter, r *http.Request) {
	var masterID *int
	if mid := r.URL.Query().Get("master_id"); mid != "" {
		v, err := strconv.Atoi(mid)
		if err == nil {
			masterID = &v
		}
	}

	services, err := repository.ListServices(masterID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, services)
}

func HandleGetService(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	svc, err := repository.GetServiceByID(id)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, svc)
}

func HandleCreateService(w http.ResponseWriter, r *http.Request) {
	var s struct {
		MasterID        int      `json:"master_id"`
		Name            string   `json:"name"`
		DurationMinutes int      `json:"duration_minutes"`
		Price           *float64 `json:"price"`
	}
	if err := decodeJSON(r, &s); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	svc := &model.Service{
		MasterID:        s.MasterID,
		Name:            s.Name,
		DurationMinutes: s.DurationMinutes,
		Price:           s.Price,
	}
	if err := repository.CreateService(svc); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, svc)
}

func HandleUpdateService(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var s struct {
		Name            string   `json:"name"`
		DurationMinutes int      `json:"duration_minutes"`
		Price           *float64 `json:"price"`
	}
	if err := decodeJSON(r, &s); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	svc := &model.Service{
		Name:            s.Name,
		DurationMinutes: s.DurationMinutes,
		Price:           s.Price,
	}
	if err := repository.UpdateService(id, svc); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	svc.ID = id
	respondJSON(w, 200, svc)
}

func HandleDeleteService(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteService(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

func HandleGetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	masterIDStr := r.URL.Query().Get("master_id")
	date := r.URL.Query().Get("date")
	serviceIDStr := r.URL.Query().Get("service_id")

	masterID, err := strconv.Atoi(masterIDStr)
	if err != nil {
		respondError(w, 400, "invalid master_id")
		return
	}
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		respondError(w, 400, "invalid service_id")
		return
	}
	if date == "" {
		respondError(w, 400, "date is required (YYYY-MM-DD)")
		return
	}

	slots, err := service.GetAvailableSlots(masterID, date, serviceID)
	if err != nil {
		respondError(w, 400, err.Error())
		return
	}
	respondJSON(w, 200, slots)
}

func HandleConfirmBooking(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	slot, err := service.ConfirmBooking(id, req.Code)
	if err != nil {
		respondError(w, 400, err.Error())
		return
	}
	respondJSON(w, 200, slot)
}

func HandleCancelBooking(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	cancelledBy := r.URL.Query().Get("by")
	if cancelledBy == "" {
		cancelledBy = "user"
	}

	slot, err := service.CancelBooking(id, cancelledBy)
	if err != nil {
		respondError(w, 400, err.Error())
		return
	}
	respondJSON(w, 200, slot)
}

func HandleMarkNoShow(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	slot, err := service.MarkNoShow(id)
	if err != nil {
		respondError(w, 400, err.Error())
		return
	}
	respondJSON(w, 200, slot)
}

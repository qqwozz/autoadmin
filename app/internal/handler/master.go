package handler

import (
	"autoadmin/internal/middleware"
	"autoadmin/internal/model"
	"autoadmin/internal/repository"
	"autoadmin/internal/service"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TelegramID int `json:"telegram_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	token, master, err := service.LoginByTelegramID(req.TelegramID)
	if err != nil {
		respondError(w, 401, "master not found")
		return
	}

	respondJSON(w, 200, map[string]any{
		"token":  token,
		"master": master,
	})
}

func HandleGetMasters(w http.ResponseWriter, r *http.Request) {
	var telegramID *int
	if tid := r.URL.Query().Get("telegram_id"); tid != "" {
		v, err := parseOptionalInt(tid)
		if err == nil {
			telegramID = &v
		}
	}

	masters, err := repository.ListMasters(telegramID)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, masters)
}

func HandleGetMasterByTelegram(w http.ResponseWriter, r *http.Request) {
	tid, ok := middleware.GetTelegramID(r)
	if !ok {
		respondError(w, 401, "unauthorized")
		return
	}

	master, err := repository.GetMasterByTelegramID(tid)
	if err != nil {
		respondError(w, 404, "master not found")
		return
	}
	respondJSON(w, 200, master)
}

func HandleGetMaster(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	master, err := repository.GetMasterByID(id)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, master)
}

func HandleCreateMaster(w http.ResponseWriter, r *http.Request) {
	var m struct {
		TelegramID  int     `json:"telegram_id"`
		Name        *string `json:"name"`
		Phone       *string `json:"phone"`
		Description *string `json:"description"`
		TariffID    *int    `json:"tariff_id"`
		IsActive    int     `json:"is_active"`
	}
	if err := decodeJSON(r, &m); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	master := &model.Master{
		TelegramID:  m.TelegramID,
		Name:        m.Name,
		Phone:       m.Phone,
		Description: m.Description,
		TariffID:    m.TariffID,
		IsActive:    m.IsActive,
	}
	if err := repository.CreateMaster(master); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, master)
}

func HandleUpdateMaster(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var m struct {
		Name        *string `json:"name"`
		Phone       *string `json:"phone"`
		Description *string `json:"description"`
		TariffID    *int    `json:"tariff_id"`
		IsActive    int     `json:"is_active"`
	}
	if err := decodeJSON(r, &m); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	master := &model.Master{
		Name:        m.Name,
		Phone:       m.Phone,
		Description: m.Description,
		TariffID:    m.TariffID,
		IsActive:    m.IsActive,
	}
	if err := repository.UpdateMaster(id, master); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	master.ID = id
	respondJSON(w, 200, master)
}

func HandleDeleteMaster(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteMaster(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

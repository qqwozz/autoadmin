package handler

import (
	"autoadmin/internal/model"
	"autoadmin/internal/repository"
	"net/http"
	"strconv"
)

func HandleGetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := repository.ListClients()
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, clients)
}

func HandleGetClientByTelegram(w http.ResponseWriter, r *http.Request) {
	tidStr := r.PathValue("telegramId")
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		respondError(w, 400, "invalid telegram_id")
		return
	}

	client, err := repository.GetClientByTelegramID(tid)
	if err != nil {
		respondError(w, 404, "client not found")
		return
	}
	respondJSON(w, 200, client)
}

func HandleGetClient(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	client, err := repository.GetClientByID(id)
	if err != nil {
		respondError(w, 404, "not found")
		return
	}
	respondJSON(w, 200, client)
}

func HandleCreateClient(w http.ResponseWriter, r *http.Request) {
	var c struct {
		TelegramID *int    `json:"telegram_id"`
		Name       *string `json:"name"`
		Phone      *string `json:"phone"`
	}
	if err := decodeJSON(r, &c); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	client := &model.Client{
		TelegramID: c.TelegramID,
		Name:       c.Name,
		Phone:      c.Phone,
	}
	if err := repository.CreateClient(client); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, client)
}

func HandleUpdateClient(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}

	var c struct {
		Name         *string `json:"name"`
		Phone        *string `json:"phone"`
		NoShowCount  int     `json:"no_show_count"`
		IsBlocked    int     `json:"is_blocked"`
		BlockedUntil *string `json:"blocked_until"`
	}
	if err := decodeJSON(r, &c); err != nil {
		respondError(w, 400, "invalid json")
		return
	}

	client := &model.Client{
		Name:         c.Name,
		Phone:        c.Phone,
		NoShowCount:  c.NoShowCount,
		IsBlocked:    c.IsBlocked,
		BlockedUntil: c.BlockedUntil,
	}
	if err := repository.UpdateClient(id, client); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	client.ID = id
	respondJSON(w, 200, client)
}

func HandleDeleteClient(w http.ResponseWriter, r *http.Request) {
	id, ok := getID(r, "id")
	if !ok {
		respondError(w, 400, "invalid id")
		return
	}
	if err := repository.DeleteClient(id); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, map[string]string{"deleted": "ok"})
}

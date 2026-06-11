package main

import (
	"autoadmin/internal/db"
	"autoadmin/internal/handler"
	"autoadmin/internal/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	db.Init()
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/login", handler.HandleLogin)

	mux.HandleFunc("GET /api/masters", handler.HandleGetMasters)
	mux.HandleFunc("GET /api/masters/{id}", handler.HandleGetMaster)
	mux.HandleFunc("POST /api/masters", handler.HandleCreateMaster)
	mux.HandleFunc("PUT /api/masters/{id}", handler.HandleUpdateMaster)
	mux.HandleFunc("DELETE /api/masters/{id}", handler.HandleDeleteMaster)

	mux.HandleFunc("GET /api/me/master", handler.HandleGetMasterByTelegram)

	mux.HandleFunc("GET /api/clients", handler.HandleGetClients)
	mux.HandleFunc("GET /api/clients/by-telegram/{telegramId}", handler.HandleGetClientByTelegram)
	mux.HandleFunc("GET /api/clients/{id}", handler.HandleGetClient)
	mux.HandleFunc("POST /api/clients", handler.HandleCreateClient)
	mux.HandleFunc("PUT /api/clients/{id}", handler.HandleUpdateClient)
	mux.HandleFunc("DELETE /api/clients/{id}", handler.HandleDeleteClient)

	mux.HandleFunc("GET /api/services", handler.HandleGetServices)
	mux.HandleFunc("GET /api/services/{id}", handler.HandleGetService)
	mux.HandleFunc("POST /api/services", handler.HandleCreateService)
	mux.HandleFunc("PUT /api/services/{id}", handler.HandleUpdateService)
	mux.HandleFunc("DELETE /api/services/{id}", handler.HandleDeleteService)

	mux.HandleFunc("GET /api/available-slots", handler.HandleGetAvailableSlots)

	mux.HandleFunc("GET /api/schedule-slots", handler.HandleGetScheduleSlots)
	mux.HandleFunc("GET /api/schedule-slots/{id}", handler.HandleGetScheduleSlot)
	mux.HandleFunc("POST /api/schedule-slots", handler.HandleCreateScheduleSlot)
	mux.HandleFunc("PUT /api/schedule-slots/{id}", handler.HandleUpdateScheduleSlot)
	mux.HandleFunc("DELETE /api/schedule-slots/{id}", handler.HandleDeleteScheduleSlot)

	mux.HandleFunc("POST /api/schedule-slots/{id}/confirm", handler.HandleConfirmBooking)
	mux.HandleFunc("POST /api/schedule-slots/{id}/cancel", handler.HandleCancelBooking)
	mux.HandleFunc("POST /api/schedule-slots/{id}/no-show", handler.HandleMarkNoShow)

	mux.HandleFunc("GET /api/working-hours", handler.HandleGetWorkingHours)
	mux.HandleFunc("GET /api/working-hours/{id}", handler.HandleGetWorkingHour)
	mux.HandleFunc("POST /api/working-hours", handler.HandleCreateWorkingHour)
	mux.HandleFunc("PUT /api/working-hours/{id}", handler.HandleUpdateWorkingHour)
	mux.HandleFunc("DELETE /api/working-hours/{id}", handler.HandleDeleteWorkingHour)

	mux.HandleFunc("GET /api/tariffs", handler.HandleGetTariffs)
	mux.HandleFunc("GET /api/tariffs/{id}", handler.HandleGetTariff)
	mux.HandleFunc("POST /api/tariffs", handler.HandleCreateTariff)
	mux.HandleFunc("PUT /api/tariffs/{id}", handler.HandleUpdateTariff)
	mux.HandleFunc("DELETE /api/tariffs/{id}", handler.HandleDeleteTariff)

	mux.HandleFunc("GET /api/no-show-settings", handler.HandleGetNoShowSettings)
	mux.HandleFunc("GET /api/no-show-settings/{masterId}", handler.HandleGetNoShowSetting)
	mux.HandleFunc("PUT /api/no-show-settings/{masterId}", handler.HandleUpdateNoShowSetting)

	mux.HandleFunc("GET /api/blacklist", handler.HandleGetBlacklist)
	mux.HandleFunc("POST /api/blacklist", handler.HandleCreateBlacklistEntry)
	mux.HandleFunc("DELETE /api/blacklist/{id}", handler.HandleDeleteBlacklistEntry)

	mux.HandleFunc("GET /api/blocked-slots", handler.HandleGetBlockedSlots)
	mux.HandleFunc("POST /api/blocked-slots", handler.HandleCreateBlockedSlot)
	mux.HandleFunc("DELETE /api/blocked-slots/{id}", handler.HandleDeleteBlockedSlot)

	mux.HandleFunc("GET /api/ref-codes", handler.HandleGetRefCodes)
	mux.HandleFunc("POST /api/ref-codes", handler.HandleCreateRefCode)
	mux.HandleFunc("DELETE /api/ref-codes/{id}", handler.HandleDeleteRefCode)

	mux.HandleFunc("GET /api/subscription-payments", handler.HandleGetPayments)
	mux.HandleFunc("POST /api/subscription-payments", handler.HandleCreatePayment)

	mux.HandleFunc("GET /api/notifications-log", handler.HandleGetNotificationsLog)

	mux.HandleFunc("GET /api/master/{id}/dashboard", handler.HandleMasterDashboard)

	h := middleware.CorsMiddleware(middleware.LoggingMiddleware(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("API server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}

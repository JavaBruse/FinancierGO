package main

import (
	"database/sql"
	"financierGo/config"
	"financierGo/internal/repositories"
	"financierGo/pkg/scheduler"
	"financierGo/routes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()
	db, _ := sql.Open("postgres", cfg.DBUrl)

	// Инициализация маршрутов
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Запуск сервера
	log.Printf("Server running on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))

	sched := scheduler.CreditScheduler{
		CreditRepo:   &repositories.CreditRepository{DB: db},
		ScheduleRepo: &repositories.PaymentScheduleRepository{DB: db},
		AccountRepo:  &repositories.AccountRepository{DB: db},
	}
	sched.Start(1 * time.Hour)

}

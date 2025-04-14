package main

import (
	"database/sql"
	"financierGo/config"
	"financierGo/internal/repositories"
	"financierGo/pkg/migrations"
	"financierGo/pkg/scheduler"
	"financierGo/routes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Подключение к базе данных
	connStr := cfg.Database.URL + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Выполнение миграций
	err = migrations.Migrate(db, "migrations")
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Инициализация маршрутов
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Инициализация и запуск планировщика
	sched := scheduler.CreditScheduler{
		CreditRepo:   &repositories.CreditRepository{DB: db},
		ScheduleRepo: &repositories.PaymentScheduleRepository{DB: db},
		AccountRepo:  &repositories.AccountRepository{DB: db},
	}
	go sched.Start(1 * time.Hour)

	// Запуск сервера
	log.Printf("Server running on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}

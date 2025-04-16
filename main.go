package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"financierGo/config"
	"financierGo/internal/repositories"
	"financierGo/pkg/migrations"
	"financierGo/pkg/scheduler"
	"financierGo/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Настройка логгера
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		DisableColors:   false,
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout) // Явно указываем вывод в консоль

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к базе данных
	connStr := cfg.Database.URL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		logrus.Fatal("Failed to ping database:", err)
	}

	// Выполнение миграций
	err = migrations.Migrate(db, "migrations")
	if err != nil {
		logrus.Fatal("Failed to run migrations:", err)
	}

	// Инициализация маршрутов
	router := mux.NewRouter()

	// Регистрируем маршруты
	routes.RegisterRoutes(router)

	// Инициализация и запуск планировщика
	sched := scheduler.CreditScheduler{
		CreditRepo:   &repositories.CreditRepository{DB: db},
		ScheduleRepo: &repositories.PaymentScheduleRepository{DB: db},
		AccountRepo:  &repositories.AccountRepository{DB: db},
	}
	go sched.Start(12 * time.Hour)

	// Запуск сервера
	logrus.Infof("Server running on port %s", cfg.Server.Port)
	logrus.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))

}

package routes

import (
	"database/sql"
	"financierGo/config"
	"financierGo/internal/handlers"
	"financierGo/internal/middleware"
	"financierGo/internal/repositories"
	"financierGo/internal/services"
	"financierGo/internal/utils"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	cfg := config.Load()
	connStr := cfg.Database.URL
	db, _ := sql.Open("postgres", connStr)
	userRepo := &repositories.UserRepository{DB: db}
	userService := &services.UserService{Repo: userRepo}
	authHandler := &handlers.AuthHandler{Service: userService}

	utils.SetJWTSecret(cfg.JWT.Secret)

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)
	accountRepo := &repositories.AccountRepository{DB: db}
	accountService := &services.AccountService{Repo: accountRepo}
	accountHandler := &handlers.AccountHandler{Service: accountService}

	api.HandleFunc("/accounts", accountHandler.Create).Methods("POST")
	api.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")

	cardRepo := &repositories.CardRepository{DB: db}
	cardService := &services.CardService{
		Repo: cardRepo, AccountRepo: accountRepo,
	}
	cardHandler := &handlers.CardHandler{Service: cardService}

	api.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")

	scheduleRepo := &repositories.PaymentScheduleRepository{DB: db}
	creditRepo := &repositories.CreditRepository{DB: db}
	creditService := &services.CreditService{
		Repo:         creditRepo,
		ScheduleRepo: scheduleRepo,
		AccountRepo:  accountRepo,
	}
	creditHandler := &handlers.CreditHandler{Service: creditService}

	api.HandleFunc("/credits", creditHandler.Create).Methods("POST")
	api.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetSchedule).Methods("GET")

	txRepo := &repositories.TransactionRepository{DB: db}
	analyticsService := &services.AnalyticsService{
		TxRepo:       txRepo,
		CreditRepo:   creditRepo,
		ScheduleRepo: scheduleRepo,
	}
	analyticsHandler := &handlers.AnalyticsHandler{Service: analyticsService}

	api.HandleFunc("/analytics", analyticsHandler.Stats).Methods("GET")
	api.HandleFunc("/analytics/credit", analyticsHandler.CreditLoad).Methods("GET")
	api.HandleFunc("/accounts/{accountId}/predict", analyticsHandler.Predict).Methods("GET")

	cbrHandler := &handlers.CBRHandler{Service: &services.CBRService{}}
	api.HandleFunc("/cbr/key-rate", cbrHandler.GetKeyRate).Methods("GET")

}

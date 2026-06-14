package routes

import (
	handler "bank/internal/handlers"
	"bank/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	authHandler *handler.AuthHandler,
	accountHandler *handler.AccountHandler,
	cardHandler *handler.CardHandler,
	transferHandler *handler.TransferHandler,
	creditHandler *handler.CreditHandler,
	financesHandler *handler.FinancesHandler,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	protected.HandleFunc("/accounts", accountHandler.GetAccounts).Methods("GET")
	protected.HandleFunc("/accounts/deposit", accountHandler.Deposit).Methods("POST")

	protected.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	protected.HandleFunc("/cards", cardHandler.GetCards).Methods("GET")

	protected.HandleFunc("/transfer", transferHandler.Transfer).Methods("POST")

	protected.HandleFunc("/credits", creditHandler.CreateCredit).Methods("POST")
	protected.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetSchedule).Methods("GET")

	protected.HandleFunc("/analytics", financesHandler.GetAnalytics).Methods("GET")
	protected.HandleFunc("/accounts/{accountId}/predict/{days}", financesHandler.PredictBalance).Methods("GET")

	return r
}

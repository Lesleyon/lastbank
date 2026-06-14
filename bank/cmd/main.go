package main

import (
	"database/sql"
	"fmt"
	"os"

	handler "bank/internal/handlers"
	"bank/internal/middleware"
	"bank/internal/routes"
	"bank/internal/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "BankBD")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	jwtSecret := getEnv("JWT_SECRET", "jwt-secret-key")
	hmacSecret := getEnv("HMAC_SECRET", "hmac-secret-key")
	serverPort := getEnv("SERVER_PORT", "8080")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}
	println("Connected to BankBD database")

	authHandler := handler.NewAuthHandler(db, jwtSecret)
	accountHandler := handler.NewAccountHandler(db)
	cardHandler := handler.NewCardHandler(db, hmacSecret)
	transferHandler := handler.NewTransferHandler(db)
	creditHandler := handler.NewCreditHandler(db)
	financesHandler := handler.NewFinancesHandler(db)

	middleware.SetJWTSecret(jwtSecret)

	router := routes.SetupRoutes(
		authHandler, accountHandler, cardHandler,
		transferHandler, creditHandler, financesHandler,
	)

	srv := server.NewServer(router, ":"+serverPort)

	println("Server starting on port", serverPort)
	if err := srv.Start(); err != nil {
		panic(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

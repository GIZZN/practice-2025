package main

import (
	"delivery-service/db"
	"delivery-service/handlers"
	"delivery-service/middleware"
	"delivery-service/repository"
	"delivery-service/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	// Инициализация базы данных
	if err := db.InitDB(); err != nil {
		
		log.Fatal("Error initializing database:", err)
	}
	defer db.DB.Close()

	// Инициализация репозиториев, сервисов и обработчиков
	userRepo := repository.NewUserRepository(db.DB)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Create router
	router := mux.NewRouter()

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigins := []string{
				"http://localhost:3000",
				"https://practice-2025.vercel.app",
				"https://practice-2025-git-main.vercel.app",
				"https://practice-2025-*.vercel.app",
			}
			origin := r.Header.Get("Origin")

			for _, allowedOrigin := range allowedOrigins {

				if origin == allowedOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// публичные роуты
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// защищенные роуты
	router.HandleFunc("/api/profile", authMiddleware.Authenticate(authHandler.GetProfile)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/profile", authMiddleware.Authenticate(authHandler.UpdateProfile)).Methods("PUT", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

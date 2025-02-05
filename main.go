package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "onlinebooking/handlers"
    "onlinebooking/middleware"
    "onlinebooking/models"
    "os"
    "os/signal"
    "syscall"
    "time"
    "fmt"
    "html/template"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/markbates/goth"
    "github.com/markbates/goth/gothic"
    "github.com/markbates/goth/providers/github"
    "github.com/sirupsen/logrus"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var appLogger = logrus.New()

func initLogger() {
    logFile, err := os.OpenFile("server_logs.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        appLogger.SetOutput(os.Stdout)
        appLogger.Warn("Failed to log to file, using default stdout")
    } else {
        appLogger.SetOutput(logFile)
    }

    appLogger.SetFormatter(&logrus.JSONFormatter{})
    appLogger.SetLevel(logrus.InfoLevel)
}

func loadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func connectDB() *gorm.DB {
    // dsn := os.Getenv("DATABASE_URL")
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        appLogger.WithField("error", err).Fatal("Failed to connect to PostgreSQL")
    }

    err = db.AutoMigrate(&models.User{}, &models.Apartment{})
    if err != nil {
        appLogger.WithField("error", err).Fatal("Failed to migrate database schema")
    }

    appLogger.Info("Connected to the database successfully")
    return db
}

func setupRouter(db *gorm.DB) *mux.Router {
    r := mux.NewRouter()

    // Serve static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    r.HandleFunc("/signup", handlers.SignUpHandler(db)).Methods("GET", "POST")
    r.HandleFunc("/verify", handlers.VerifyEmailHandler(db)).Methods("GET")
    r.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")
    r.HandleFunc("/login", Login).Methods("GET")
    r.HandleFunc("/logout", handlers.LogoutHandler()).Methods("GET")

    // OAuth routes
    r.HandleFunc("/auth/{provider}/callback", handlers.OAuthCallbackHandler(db)).Methods("GET")
    r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler).Methods("GET")

    // API routes
    r.HandleFunc("/api/check_auth", checkAuthHandler).Methods("GET")
    r.HandleFunc("/api/apartments", apartmentsHandler(db)).Methods("GET")

    // Protected routes
    r.Handle("/index", middleware.AuthMiddleware(http.HandlerFunc(handlers.IndexHandler(db)))).Methods("GET")

    r.Handle("/add_apartment", middleware.AuthMiddleware(http.HandlerFunc(handlers.AddApartmentHandler(db)))).Methods("GET", "POST")

    return r
}
func Login(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("static/login.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
    return
}
func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
    isAuthenticated := false
    if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
        isAuthenticated = true
    }

    response := map[string]bool{"isAuthenticated": isAuthenticated}
    json.NewEncoder(w).Encode(response)
}

func apartmentsHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var apartments []models.Apartment
        if err := db.Find(&apartments).Error; err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(apartments)
    }
}

func main() {
    initLogger()
    logrus.Info("Logger initialized")

    loadEnv()
    db := connectDB()

    // Initialize OAuth providers
    goth.UseProviders(
        
        github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_CLIENT_SECRET"), os.Getenv("GITHUB_REDIRECT_URL")),
    )

    r := setupRouter(db)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

    go func() {
        appLogger.Info("Server is running on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            appLogger.WithField("error", err).Fatal("Server error")
        }
    }()

    <-quit
    appLogger.Info("Server is shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        appLogger.WithField("error", err).Fatal("Server forced to shutdown")
    }

    appLogger.Info("Server exited gracefully")
}
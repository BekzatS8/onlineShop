package handlers

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/smtp"
    "onlinebooking/models"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

func generateVerificationToken() (string, error) {
    token := make([]byte, 16)
    if _, err := rand.Read(token); err != nil {
        return "", err
    }
    return hex.EncodeToString(token), nil
}

func sendVerificationEmail(email, token string) error {
    from := os.Getenv("EMAIL_USER")
    password := os.Getenv("EMAIL_PASSWORD")
    to := email
    subject := "Email Verification"
    body := fmt.Sprintf("Please verify your email by clicking the following link: http://localhost:8080/verify?token=%s", token)

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
    err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
    if err != nil {
        log.Printf("Failed to send email: %v", err)
        return err
    }

    return nil
}

func SignUpHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            http.ServeFile(w, r, "static/signup.html")
            return
        }

        if r.Method == http.MethodPost {
            email := r.FormValue("email")
            password := r.FormValue("password")

            // Check if the user already exists
            var existingUser models.User
            if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
                http.Error(w, "Email already in use", http.StatusBadRequest)
                return
            }

            hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
            if err != nil {
                log.Printf("Error hashing password: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            token, err := generateVerificationToken()
            if err != nil {
                log.Printf("Error generating verification token: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            // Store the hashed password and token temporarily
            tempUser := models.User{Email: email, Password: string(hashedPassword), VerificationToken: token}
            if err := db.Create(&tempUser).Error; err != nil {
                log.Printf("Error creating temporary user: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }

            if err := sendVerificationEmail(email, token); err != nil {
                log.Printf("Failed to send verification email: %v", err)
                http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
                return
            }

            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(map[string]string{"message": "Verification email sent. Please check your email to verify your account."})
        }
    }
}

func VerifyEmailHandler(db *gorm.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.URL.Query().Get("token")
        if token == "" {
            http.Error(w, "Invalid token", http.StatusBadRequest)
            return
        }

        var user models.User
        if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
            http.Error(w, "Invalid token", http.StatusBadRequest)
            return
        }

        if user.EmailVerified {
            http.Error(w, "Email already verified", http.StatusBadRequest)
            return
        }

        user.EmailVerified = true
        user.VerificationToken = ""
        if err := db.Save(&user).Error; err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        // Generate JWT token
        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Email: user.Email,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:    "token",
            Value:   tokenString,
            Expires: expirationTime,
        })

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}
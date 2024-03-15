package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "campus-api/utils"
    "campus-api/internal/config"
)

var cfg *config.Config

func RegisterHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
            Fname    string `json:"fname"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("RegisterHandler: Error decoding request body: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        verificationToken, err := utils.RegisterUser(db, req.Email, req.Password, req.Fname)
        if err != nil {
            log.Printf("RegisterHandler: Error registering user %s: %v", req.Email, err)
            http.Error(w, "Failed to register user", http.StatusInternalServerError)
            return
        }

        err = utils.SendVerificationEmail(cfg, req.Email, verificationToken)
        if err != nil {
            log.Printf("RegisterHandler: Error sending verification email to %s: %v", req.Email, err)
            http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
            return
        }

        log.Printf("RegisterHandler: Registration and verification email sent successfully to %s", req.Email)
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("User registered successfully. Please check your email to verify your account."))
    }
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("LoginHandler called for user with id")
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("Login request decoding error for email %s: %v", req.Email, err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        userAuthInfo, err := utils.VerifyUser(db, req.Email, req.Password)
        if err != nil {
            log.Printf("Verification error during login for user %s: %v", req.Email, err)
            http.Error(w, "Failed to login", http.StatusInternalServerError)
            return
        }
        if !userAuthInfo.Authenticated {
            log.Printf("Invalid credentials provided for user %s", req.Email)
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        if userAuthInfo.UserID == 0 {
            log.Printf("User ID came back as zero, indicating user does not exist for email: %s", req.Email)
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        response := map[string]interface{}{
            "success":         true,
            "isVerified":      userAuthInfo.IsVerified,
            "hasSelectedRole": userAuthInfo.HasSelectedRole,
            "userId":          userAuthInfo.UserID,
        }

        log.Printf("User %d logged in successfully. Has selected role: %t", userAuthInfo.UserID, userAuthInfo.HasSelectedRole)
        responseJSON, err := json.Marshal(response)
        if err != nil {
            log.Printf("Error marshalling response: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(responseJSON)
    }
}

func VerifyEmailHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.URL.Query().Get("token")
        if token == "" {
            log.Println("Verification token is missing.")
            http.Error(w, "Verification token is required", http.StatusBadRequest)
            return
        }

        var userID int
        err := db.QueryRow("UPDATE users SET is_verified = TRUE WHERE verification_token = $1 RETURNING id", token).Scan(&userID)
        if err != nil {
            if err == sql.ErrNoRows {
                log.Printf("Invalid or expired verification token.")
                http.Error(w, "Invalid or expired token", http.StatusBadRequest)
            } else {
                log.Printf("Failed to verify email: %v", err)
                http.Error(w, "Failed to verify email", http.StatusInternalServerError)
            }
            return
        }

        log.Printf("Email verification successful for user ID %d", userID)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Email verified successfully."))
    }
}

func ForgotPasswordHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Email string `json:"email"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("ForgotPasswordHandler: Error decoding request body: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        resetToken, err := utils.GenerateResetToken(db, req.Email)
        if err != nil {
            log.Printf("ForgotPasswordHandler: Error generating reset token for %s: %v", req.Email, err)
            http.Error(w, "Failed to initiate password reset", http.StatusInternalServerError)
            return
        }

        resetLink := cfg.FrontendURL + "/change_password?token=" + resetToken // Ensure FrontendURL is defined in your Config struct
        err = utils.SendPasswordResetEmail(cfg, req.Email, resetLink)
        if err != nil {
            log.Printf("ForgotPasswordHandler: Error sending password reset email to %s: %v", req.Email, err)
            http.Error(w, "Failed to send password reset email", http.StatusInternalServerError)
            return
        }

        log.Printf("ForgotPasswordHandler: Password reset email sent to %s successfully", req.Email)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Password reset email sent successfully. Please check your email."))
    }
}

func ChangePasswordHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Token       string `json:"token"`
            NewPassword string `json:"newPassword"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("ChangePasswordHandler: Error decoding request: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        email, err := utils.ResetPassword(db, req.Token, req.NewPassword)
        if err != nil {
            log.Printf("ChangePasswordHandler: Error resetting password: %v", err)
            http.Error(w, "Failed to reset password", http.StatusInternalServerError)
            return
        }

        log.Printf("ChangePasswordHandler: Password reset successfully for %s", email)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Password reset successfully."))
    }
}

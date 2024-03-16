package utils

import (
    "database/sql"
    "fmt"
    "log"
    "net/smtp"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"

    "campus-api/internal/config"
)

func RegisterUser(db *sql.DB, email, password, fname string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return "", err
    }

    verificationToken := uuid.NewString()

    _, err = db.Exec("INSERT INTO users (email, password_hash, fname, verification_token) VALUES ($1, $2, $3, $4)", email, string(hashedPassword), fname, verificationToken)
    if err != nil {
        log.Printf("Error inserting new user into database: %v", err)
        return "", err
    }

    return verificationToken, nil
}

// UserAuthInfo represents the needed authentication and role selection information for a user.
type UserAuthInfo struct {
    Authenticated  bool
    HasSelectedRole bool
    UserID         int
    IsVerified bool
}

// VerifyUser checks if a user exists with the given email and password, and returns their role selection status.
func VerifyUser(db *sql.DB, email, password string) (UserAuthInfo, error) {
    var info UserAuthInfo
    var hashedPassword string

    query := "SELECT id, password_hash, has_selected_role, is_verified FROM users WHERE email = $1"
    err := db.QueryRow(query, email).Scan(&info.UserID, &hashedPassword, &info.HasSelectedRole, &info.IsVerified)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No user found with email: %v", email)
            return info, err  // User not found
        } else {
            log.Printf("Error retrieving user from database: %v", err)
            return info, err  // Other error
        }
    }

    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        log.Printf("Password mismatch for user: %v", email)
        return info, err  // Password does not match
    }

    // If we reach here, authentication was successful
    info.Authenticated = true
    log.Printf("User authenticated: %v", email)

    return info, nil // Return the authentication info including the user ID
}

func SendVerificationEmail(cfg *config.Config, recipient, token string) error {
    from := cfg.Email
    password := cfg.EmailPassword
    smtpHost := cfg.SMTPHost
    smtpPort := cfg.SMTPPort

    verificationLink := cfg.FrontendURL + "/verify_email?token=" + token

    message := []byte("To: " + recipient + "\r\n" +
        "Subject: Verify Your Email\r\n\r\n" +
        "Click on the following link to verify your email: " + verificationLink)

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, message)
    if err != nil {
        log.Printf("Error sending verification email to %s: %v", recipient, err)
        return err
    }

    log.Printf("Verification email sent to %s successfully", recipient)
    return nil
}

func GenerateResetToken(db *sql.DB, email string) (string, error) {
    resetToken := uuid.NewString()
    expirationTime := time.Now().Add(1 * time.Hour) // Token expires in 1 hour

    result, err := db.Exec("UPDATE users SET reset_token = $1, reset_token_expires = $2 WHERE email = $3", resetToken, expirationTime, email)
    if err != nil {
        return "", fmt.Errorf("error updating reset token for user %s: %v", email, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return "", fmt.Errorf("error checking affected rows for user %s: %v", email, err)
    }

    if rowsAffected == 0 {
        return "", fmt.Errorf("no user found with email %s", email)
    }

    return resetToken, nil
}

func SendPasswordResetEmail(cfg *config.Config, recipient, resetLink string) error {
    from := cfg.Email
    password := cfg.EmailPassword
    smtpHost := cfg.SMTPHost
    smtpPort := cfg.SMTPPort

    message := []byte("To: " + recipient + "\r\n" +
        "Subject: Password Reset Request\r\n\r\n" +
        "You have requested to reset your password. Click on the link below to set a new password:\n" + resetLink +
        "\n\nIf you did not request this, please ignore this email.")

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, message)
    if err != nil {
        log.Printf("Error sending password reset email to %s: %v", recipient, err)
        return fmt.Errorf("error sending password reset email to %s: %v", recipient, err)
    }

    log.Printf("Password reset email sent to %s successfully", recipient)
    return nil
}

func ResetPassword(db *sql.DB, token, newPassword string) (string, error) {
    var email string
    expirationCheck := "SELECT email FROM users WHERE reset_token = $1 AND reset_token_expires > NOW()"
    err := db.QueryRow(expirationCheck, token).Scan(&email)
    if err != nil {
        return "", err // Token not found or expired
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }

    updatePassword := "UPDATE users SET password_hash = $1, reset_token = NULL, reset_token_expires = NULL WHERE email = $2"
    _, err = db.Exec(updatePassword, string(hashedPassword), email)
    if err != nil {
        return "", err
    }

    return email, nil
}

func SendEmail(cfg *config.Config, recipient, subject, body string) error {
    from := cfg.Email
    password := cfg.EmailPassword
    smtpHost := cfg.SMTPHost
    smtpPort := cfg.SMTPPort

    message := []byte("To: " + recipient + "\r\n" +
        "Subject: " + subject + "\r\n\r\n" +
        body)

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{recipient}, message)
    if err != nil {
        log.Printf("Error sending email to %s: %v", recipient, err)
        return err
    }

    log.Printf("Email sent to %s successfully", recipient)
    return nil
}
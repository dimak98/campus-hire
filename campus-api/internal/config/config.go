package config

import (
    "os"
    "strconv"
)

type Config struct {
    DBHost         string
    DBPort         int
    DBUser         string
    DBPassword     string
    DBName         string
    SMTPHost       string
    SMTPPort       string
    Email          string
    EmailPassword  string
    FrontendURL    string
}

func LoadConfig() (*Config, error) {
    port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, err
    }

    return &Config{
        DBHost:        getEnv("DB_HOST", "localhost"),
        DBPort:        port,
        DBUser:        getEnv("DB_USER", "user"),
        DBPassword:    getEnv("DB_PASSWORD", "password"),
        DBName:        getEnv("DB_NAME", "dbname"),
        SMTPHost:      getEnv("SMTP_HOST", "smtp.gmail.com"),
        SMTPPort:      getEnv("SMTP_PORT", "587"),
        Email:         getEnv("EMAIL", "email@gmail.com"),
        EmailPassword: getEnv("EMAIL_PASSWORD", "passwd"),
        FrontendURL:   getEnv("FRONTEND_URL", "url"), 
    }, nil
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}

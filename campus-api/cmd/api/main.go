package main

import (
    "campus-api/internal/config"
    "campus-api/internal/database"
    "campus-api/internal/handlers"
    "log"
    "net/http"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    db, err := database.Initialize(cfg)
    if err != nil {
        log.Fatal("Failed to initialize database:", err)
    }

    // Auth routes
    http.HandleFunc("/register", handlers.RegisterHandler(db, cfg))
    http.Handle("/login", handlers.LoginHandler(db))
    http.HandleFunc("/verify_email", handlers.VerifyEmailHandler(db))
    http.HandleFunc("/forgot_password", handlers.ForgotPasswordHandler(db, cfg))
    http.HandleFunc("/change_password", handlers.ChangePasswordHandler(db, cfg))

    // Registration routes
    http.HandleFunc("/student_registration", handlers.StudentRegistrationHandler(db))
    http.HandleFunc("/company_registration", handlers.CompanyRegistrationHandler(db))

    // Get routes
    http.HandleFunc("/user_details", handlers.GetUserDetailsHandler(db))
    http.HandleFunc("/student", handlers.GetStudentDetailsHandler(db))
    http.HandleFunc("/company", handlers.GetCompanyDetailsHandler(db))
    http.HandleFunc("/students", handlers.GetAllStudentsHandler(db))
    http.HandleFunc("/user_role", handlers.GetUserRoleHandler(db))
    http.HandleFunc("/jobs", handlers.GetAllJobsHandler(db))
    http.HandleFunc("/job", handlers.GetJobByIDHandler(db))

    // Post routes
    http.HandleFunc("/post_job", handlers.CreateJobPostHandler(db))

    // Update Routes
    http.HandleFunc("/apply_for_job", handlers.ApplyForJobHandler(db, cfg))

    log.Println("Starting server on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "campus-api/pkg/models"
)

func StudentRegistrationHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req models.StudentRegistrationRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("Error decoding registration request: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        log.Printf("Attempting to register new user with id %s as a student", req.UserID)
        tx, err := db.Begin()
        if err != nil {
            log.Printf("Error starting transaction: %v", err)
            http.Error(w, "Database transaction error", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        studentQuery := `INSERT INTO students (user_id, description, profile_image) VALUES ($1, $2, $3)`
        if _, err := tx.Exec(studentQuery, req.UserID, req.Description, req.ImagePath); err != nil {
            log.Printf("Error inserting student details: %v", err)
            http.Error(w, "Failed to insert student details", http.StatusInternalServerError)
            return
        }

        // Only insert jobs if there are any
        if len(req.Jobs) > 0 {
            jobQuery := `INSERT INTO jobs (user_id, title, company, start_date, end_date, description) VALUES ($1, $2, $3, $4, $5, $6)`
            for _, job := range req.Jobs {
                if _, err := tx.Exec(jobQuery, req.UserID, job.Title, job.Company, job.StartDate, job.EndDate, job.Description); err != nil {
                    log.Printf("Error inserting job: %v", err)
                    http.Error(w, "Failed to insert job", http.StatusInternalServerError)
                    return
                }
            }
        }

        // Only insert education if there are any
        if len(req.Education) > 0 {
            educationQuery := `INSERT INTO education (user_id, school, degree, field_of_study, start_date, end_date, description) VALUES ($1, $2, $3, $4, $5, $6, $7)`
            for _, edu := range req.Education {
                if _, err := tx.Exec(educationQuery, req.UserID, edu.School, edu.Degree, edu.FieldOfStudy, edu.StartDate, edu.EndDate, edu.Description); err != nil {
                    log.Printf("Error inserting education: %v", err)
                    http.Error(w, "Failed to insert education", http.StatusInternalServerError)
                    return
                }
            }
        }

        updateUserQuery := `UPDATE users SET has_selected_role = TRUE, role = 'student' WHERE id = $1`
        if _, err := tx.Exec(updateUserQuery, req.UserID); err != nil {
            log.Printf("Error updating user's role selection status: %v", err)
            http.Error(w, "Failed to update user's role selection status", http.StatusInternalServerError)
            return
        }        

        if err := tx.Commit(); err != nil {
            log.Printf("Error committing transaction: %v", err)
            http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
            return
        }

        log.Println("Student registration successful")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
    }
}

func CompanyRegistrationHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req models.CompanyRegistrationRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("Error decoding company registration request: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Verify user exists
        var exists bool
        err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", req.UserID).Scan(&exists)
        if err != nil || !exists {
            log.Printf("User with ID %d does not exist: %v", req.UserID, err)
            http.Error(w, "User does not exist", http.StatusBadRequest)
            return
        }

        // Begin transaction
        tx, err := db.Begin()
        if err != nil {
            log.Printf("Error starting transaction: %v", err)
            http.Error(w, "Database transaction error", http.StatusInternalServerError)
            return
        }
        defer tx.Rollback()

        // Insert company details
        companyQuery := `INSERT INTO companies (user_id, name, size, address, description, image_path, video_path) VALUES ($1, $2, $3, $4, $5, $6, $7)`
        if _, err := tx.Exec(companyQuery, req.UserID, req.Name, req.Size, req.Address, req.Description, req.ImagePath, req.VideoPath); err != nil {
            log.Printf("Error inserting company details: %v", err)
            http.Error(w, "Failed to insert company details", http.StatusInternalServerError)
            return
        }

        // Update user's role selection status
        updateUserQuery := `UPDATE users SET has_selected_role = TRUE, role = 'company' WHERE id = $1`
        if _, err := tx.Exec(updateUserQuery, req.UserID); err != nil {
            log.Printf("Error updating user's role selection status: %v", err)
            http.Error(w, "Failed to update user's role selection status", http.StatusInternalServerError)
            return
        }

        // Commit transaction
        if err := tx.Commit(); err != nil {
            log.Printf("Error committing transaction: %v", err)
            http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
            return
        }

        log.Println("Company registration successful")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
    }
}
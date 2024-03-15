package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "campus-api/pkg/models"
)

func CreateJobPostHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var jobPost models.JobPost
        if err := json.NewDecoder(r.Body).Decode(&jobPost); err != nil {
            log.Printf("Error decoding job post request: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        query := `INSERT INTO job_posts (user_id, title, description, requirements, status, salary, address)
                  VALUES ($1, $2, $3, $4, $5, $6, $7)`
        _, err := db.Exec(query, jobPost.UserID, jobPost.Title, jobPost.Description, jobPost.Requirements, jobPost.Status, jobPost.Salary, jobPost.Address)
        if err != nil {
            log.Printf("Error inserting job post for user with id %d: %v", jobPost.UserID, err)
            http.Error(w, "Failed to create job post", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Job post created successfully"})
    }
}

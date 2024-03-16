package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
	"fmt"

	"campus-api/internal/config"
	"campus-api/utils"

)

type JobApplicationRequest struct {
    UserID          int    `json:"user_id"`
    JobID           int    `json:"job_id"`
    CompanyEmail    string `json:"company_email"`
    CvDownloadLink  string `json:"cv_download_link"`
}

func ApplyForJobHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req JobApplicationRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Printf("Error decoding request for job application: %v", err)
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }

        // Log the decoded request
        log.Printf("Received apply job request for user_id: %d, job_id: %d", req.UserID, req.JobID)

        // Insert the job application into the database
        _, err := db.Exec("INSERT INTO student_job_applications (user_id, job_id) VALUES ($1, $2)", req.UserID, req.JobID)
        if err != nil {
            log.Printf("Error applying for job with user_id: %d, job_id: %d - %v", req.UserID, req.JobID, err)
            http.Error(w, "Failed to apply for job", http.StatusInternalServerError)
            return
        }

        // Send an email to the company with the CV download link
        emailSubject := "New Job Application"
        emailBody := fmt.Sprintf("A student with ID %d has applied for your job. You can download their CV here: %s", req.UserID, req.CvDownloadLink)
        err = utils.SendEmail(cfg, req.CompanyEmail, emailSubject, emailBody)
        if err != nil {
            log.Printf("Error sending email to company: %v", err)
            // Optionally, you can handle this error differently
        }

        // Log successful application submission
        log.Printf("Application submitted successfully for user_id: %d, job_id: %d", req.UserID, req.JobID)

        // Respond with a success message
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Application submitted successfully"})
    }
}
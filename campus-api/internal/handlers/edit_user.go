package handlers

import (
    "campus-api/pkg/models"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
)

func EditStudentJobHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var job models.StudentJob
        if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
            log.Printf("Error decoding student job data: %v", err)
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }

        log.Printf("Received request to update student job: %+v", job)

        query := `UPDATE jobs SET title=$1, company=$2, start_date=$3, end_date=$4, description=$5 WHERE id=$6 AND user_id=$7`
        res, err := db.Exec(query, job.Title, job.Company, job.StartDate, job.EndDate, job.Description, job.ID, job.UserID)
        if err != nil {
            log.Printf("Error updating student job: %v", err)
            http.Error(w, "Failed to update student job", http.StatusInternalServerError)
            return
        }

        rowsAffected, err := res.RowsAffected()
        if err != nil {
            log.Printf("Error getting rows affected: %v", err)
            http.Error(w, "Failed to update student job", http.StatusInternalServerError)
            return
        }

        if rowsAffected == 0 {
            log.Printf("No rows affected, check if the job ID and user ID are correct")
            http.Error(w, "No rows affected, check if the job ID and user ID are correct", http.StatusBadRequest)
            return
        }

        log.Printf("Student job updated successfully: %+v", job)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Student job updated successfully"})
    }
}

func EditEducationHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var education models.Education
        if err := json.NewDecoder(r.Body).Decode(&education); err != nil {
            log.Printf("Error decoding education data: %v", err)
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }

        log.Printf("Received request to update education: %+v", education)

        query := `UPDATE education SET school=$1, degree=$2, field_of_study=$3, start_date=$4, end_date=$5, description=$6 WHERE id=$7 AND user_id=$8`
        res, err := db.Exec(query, education.School, education.Degree, education.FieldOfStudy, education.StartDate, education.EndDate, education.Description, education.ID, education.UserID)
        if err != nil {
            log.Printf("Error updating education: %v", err)
            http.Error(w, "Failed to update education", http.StatusInternalServerError)
            return
        }

        rowsAffected, err := res.RowsAffected()
        if err != nil {
            log.Printf("Error getting rows affected: %v", err)
            http.Error(w, "Failed to update education", http.StatusInternalServerError)
            return
        }

        if rowsAffected == 0 {
            log.Printf("No rows affected, check if the education ID and user ID are correct")
            http.Error(w, "No rows affected, check if the education ID and user ID are correct", http.StatusBadRequest)
            return
        }

        log.Printf("Education updated successfully: %+v", education)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Education updated successfully"})
    }
}

func EditCompanyHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var company models.Company
        if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
            log.Printf("Error decoding company data: %v", err)
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }

        log.Printf("Received request to update company: %+v", company)

        query := `UPDATE companies SET size=$2, address=$3, description=$4, fname=$5, email=$6 WHERE user_id=$1`
        _, err := db.Exec(query, company.UserID, company.Size, company.Address, company.Description, company.Fname, )
        if err != nil {
            log.Printf("Error updating company: %v", err)
            http.Error(w, "Failed to update company", http.StatusInternalServerError)
            return
        }

        log.Printf("Company updated successfully: %+v", company)

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Company updated successfully"})
    }
}

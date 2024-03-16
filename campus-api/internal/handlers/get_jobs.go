package handlers

import (
    "campus-api/pkg/models"
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "strconv"
)

func GetAllJobsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var jobs []models.JobPost
        titleFilter := r.URL.Query().Get("title")
        latest := r.URL.Query().Get("latest") == "true"
        userID := r.URL.Query().Get("user_id")

        // Start building the query
        query := `SELECT jp.id, jp.user_id, jp.title, jp.description, jp.requirements, jp.status, jp.salary, jp.address, COALESCE(c.name, '') AS company_name, u.email AS company_email
                  FROM job_posts jp
                  LEFT JOIN companies c ON jp.user_id = c.user_id
                  LEFT JOIN users u ON c.user_id = u.id`

        // Add conditions to the query
        var whereClauses []string
        if titleFilter != "" {
            titleFilter = strings.Replace(titleFilter, "'", "''", -1)
            whereClauses = append(whereClauses, "LOWER(jp.title) LIKE LOWER('%"+titleFilter+"%')")
        }
        if latest {
            whereClauses = append(whereClauses, "jp.created_at >= CURRENT_TIMESTAMP - INTERVAL '2 weeks'")
        }
        if userID != "" {
            whereClauses = append(whereClauses, "jp.id NOT IN (SELECT job_id FROM student_job_applications WHERE user_id = "+userID+")")
        }
        if len(whereClauses) > 0 {
            query += " WHERE " + strings.Join(whereClauses, " AND ")
        }

        rows, err := db.Query(query)
        if err != nil {
            log.Printf("Error querying job posts: %v", err)
            http.Error(w, "Error querying job posts", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        // Iterate through the returned rows
        for rows.Next() {
            var job models.JobPost
            if err := rows.Scan(&job.ID, &job.UserID, &job.Title, &job.Description, &job.Requirements, &job.Status, &job.Salary, &job.Address, &job.CompanyName, &job.CompanyEmail); err != nil {
                log.Printf("Error scanning job post: %v", err)
                http.Error(w, "Error reading job posts", http.StatusInternalServerError)
                return
            }
            jobs = append(jobs, job)
        }

        // Check for errors from iterating over rows
        if err := rows.Err(); err != nil {
            log.Printf("Error iterating over job posts: %v", err)
            http.Error(w, "Error iterating over job posts", http.StatusInternalServerError)
            return
        }

        if len(jobs) == 0 {
            log.Println("No job posts found")
            w.WriteHeader(http.StatusNotFound)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(jobs); err != nil {
            log.Printf("Error encoding job posts to JSON: %v", err)
            http.Error(w, "Error encoding job posts", http.StatusInternalServerError)
        }
    }
}

func GetJobByIDHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get the job ID from the URL path
        jobIDStr := r.URL.Query().Get("jobID")
        jobID, err := strconv.Atoi(jobIDStr)
        if err != nil {
            log.Printf("Invalid job ID: %v", err)
            http.Error(w, "Invalid job ID", http.StatusBadRequest)
            return
        }

        // Query the database for the job post
        var job models.JobPost
        query := `SELECT jp.id, jp.user_id, jp.title, jp.description, jp.requirements, jp.status, jp.salary, jp.address, COALESCE(c.name, '') AS company_name
                  FROM job_posts jp
                  LEFT JOIN companies c ON jp.user_id = c.user_id
                  WHERE jp.id = $1`
        err = db.QueryRow(query, jobID).Scan(&job.ID, &job.UserID, &job.Title, &job.Description, &job.Requirements, &job.Status, &job.Salary, &job.Address, &job.CompanyName)
        if err != nil {
            if err == sql.ErrNoRows {
                log.Printf("Job post not found: %v", err)
                http.Error(w, "Job post not found", http.StatusNotFound)
                return
            }
            log.Printf("Error querying job post: %v", err)
            http.Error(w, "Error querying job post", http.StatusInternalServerError)
            return
        }

        // Return the job post as JSON
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(job); err != nil {
            log.Printf("Error encoding job post to JSON: %v", err)
            http.Error(w, "Error encoding job post", http.StatusInternalServerError)
        }
    }
}
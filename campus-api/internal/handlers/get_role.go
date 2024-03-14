package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    _ "github.com/lib/pq"
)

func GetUserRoleHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get the user ID from the query parameters
        userID := r.URL.Query().Get("userID")
        if userID == "" {
            http.Error(w, "User ID is required", http.StatusBadRequest)
            return
        }

        // Query the database for the user's role
        var role string
        err := db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "User not found", http.StatusNotFound)
            } else {
                log.Printf("Error querying user role: %v", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
            }
            return
        }

        // Return the role as JSON
        json.NewEncoder(w).Encode(map[string]string{"role": role})
    }
}

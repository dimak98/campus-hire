package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
)

// StudentDetails represents the details of a student.
type StudentDetails struct {
    ID           string            `json:"id"`
    Email        string            `json:"email"`
    FName        string            `json:"fname"`
    Role         string            `json:"role"`
    IsCvCreated  string            `json:"is_cv_created"`
    CvPath       string            `json:"cv_path"`
    Description  string            `json:"description"`
    ProfileImage string            `json:"profileImage"`
    Jobs         []StudentJob      `json:"jobs"`
    Education    []StudentEducation `json:"education"`
}

// StudentJob represents a job entry for a student.
type StudentJob struct {
    Title       string `json:"title"`
    Company     string `json:"company"`
    StartDate   string `json:"startDate"`
    EndDate     string `json:"endDate"`
    Description string `json:"description"`
}

// StudentEducation represents an education entry for a student.
type StudentEducation struct {
    School       string `json:"school"`
    Degree       string `json:"degree"`
    FieldOfStudy string `json:"fieldOfStudy"`
    StartDate    string `json:"startDate"`
    EndDate      string `json:"endDate"`
    Description  string `json:"description"`
}

// CompanyDetails represents the details of a company.
type CompanyDetails struct {
    Email       string `json:"email"`
    FName       string `json:"fname"`
    Role        string `json:"role"`
    Name        string `json:"name"`
    Size        string `json:"size"`
    Address     string `json:"address"`
    Description string `json:"description"`
    ImagePath   string `json:"image_path"`
    VideoPath   string `json:"video_path"`
}

// GetStudentDetailsHandler handles fetching student details.
func GetStudentDetailsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.URL.Query().Get("userID")

        // Fetch student details
        studentDetails, err := fetchStudentDetails(db, userID)
        if err != nil {
            log.Printf("Error fetching student details: %v", err)
            http.Error(w, "Error fetching student details", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(studentDetails)
    }
}

// GetCompanyDetailsHandler handles fetching company details.
func GetCompanyDetailsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userID := r.URL.Query().Get("userID")

        // Fetch company details
        companyDetails, err := fetchCompanyDetails(db, userID)
        if err != nil {
            log.Printf("Error fetching company details: %v", err)
            http.Error(w, "Error fetching company details", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(companyDetails)
    }
}

// GetUserDetailsHandler handles fetching user details based on their role.
func GetUserDetailsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract user ID from query parameters or URL path
        userID := r.URL.Query().Get("userID")

        // First, get the role of the user
        var role string
        err := db.QueryRow("SELECT role FROM users WHERE id = $1", userID).Scan(&role)
        if err != nil {
            log.Printf("Error fetching user role: %v", err)
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        var response interface{}

        switch role {
        case "student":
            response, err = fetchStudentDetails(db, userID)
        case "company":
            response, err = fetchCompanyDetails(db, userID)
        default:
            http.Error(w, "Invalid user role", http.StatusBadRequest)
            return
        }

        if err != nil {
            log.Printf("Error fetching user details: %v", err)
            http.Error(w, "Error fetching user details", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}

// fetchStudentDetails fetches the details of a student from the database.
func fetchStudentDetails(db *sql.DB, userID string) (*StudentDetails, error) {
    var studentDetails StudentDetails

    // Fetch user details
    err := db.QueryRow("SELECT id, email, fname, role FROM users WHERE id = $1", userID).Scan(&studentDetails.ID, &studentDetails.Email, &studentDetails.FName, &studentDetails.Role)
    if err != nil {
        return nil, err
    }

    // Fetch student details
    err = db.QueryRow("SELECT description, profile_image, is_cv_created, cv_path FROM students WHERE user_id = $1", userID).Scan(&studentDetails.Description, &studentDetails.ProfileImage, &studentDetails.IsCvCreated, &studentDetails.CvPath)
    if err != nil {
        return nil, err
    }

    // Fetch jobs
    rows, err := db.Query("SELECT title, company, start_date, end_date, description FROM jobs WHERE user_id = $1", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var job StudentJob
        if err := rows.Scan(&job.Title, &job.Company, &job.StartDate, &job.EndDate, &job.Description); err != nil {
            return nil, err
        }
        studentDetails.Jobs = append(studentDetails.Jobs, job)
    }

    // Fetch education
    rows, err = db.Query("SELECT school, degree, field_of_study, start_date, end_date, description FROM education WHERE user_id = $1", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var edu StudentEducation
        if err := rows.Scan(&edu.School, &edu.Degree, &edu.FieldOfStudy, &edu.StartDate, &edu.EndDate, &edu.Description); err != nil {
            return nil, err
        }
        studentDetails.Education = append(studentDetails.Education, edu)
    }

    return &studentDetails, nil
}

// fetchCompanyDetails fetches the details of a company from the database.
func fetchCompanyDetails(db *sql.DB, userID string) (*CompanyDetails, error) {
    var companyDetails CompanyDetails

    // Fetch user details
    err := db.QueryRow("SELECT email, fname, role FROM users WHERE id = $1", userID).Scan(&companyDetails.Email, &companyDetails.FName, &companyDetails.Role)
    if err != nil {
        return nil, err
    }

    // Fetch company details
    err = db.QueryRow("SELECT name, size, address, description, image_path, video_path FROM companies WHERE user_id = $1", userID).Scan(&companyDetails.Name, &companyDetails.Size, &companyDetails.Address, &companyDetails.Description, &companyDetails.ImagePath, &companyDetails.VideoPath)
    if err != nil {
        return nil, err
    }

    return &companyDetails, nil
}

// GetAllStudentsHandler handles fetching details of all students.
func GetAllStudentsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Fetch all student details
        students, err := fetchAllStudentsDetails(db)
        if err != nil {
            log.Printf("Error fetching all student details: %v", err)
            http.Error(w, "Error fetching all student details", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(students)
    }
}

// fetchAllStudentsDetails fetches the details of all students from the database.
func fetchAllStudentsDetails(db *sql.DB) ([]StudentDetails, error) {
    var students []StudentDetails

    // Fetch all student user details
    rows, err := db.Query("SELECT id, email, fname, role FROM users WHERE role = 'student'")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var student StudentDetails
        if err := rows.Scan(&student.ID, &student.Email, &student.FName, &student.Role); err != nil {
            return nil, err
        }

        // Fetch additional student details
        err = db.QueryRow("SELECT description, profile_image, is_cv_created, cv_path FROM students WHERE user_id = $1", student.ID).Scan(&student.Description, &student.ProfileImage, &student.IsCvCreated, &student.CvPath)
        if err != nil {
            return nil, err
        }

        // Fetch jobs for the student
        jobRows, err := db.Query("SELECT title, company, start_date, end_date, description FROM jobs WHERE user_id = $1", student.ID)
        if err != nil {
            return nil, err
        }

        for jobRows.Next() {
            var job StudentJob
            if err := jobRows.Scan(&job.Title, &job.Company, &job.StartDate, &job.EndDate, &job.Description); err != nil {
                return nil, err
            }
            student.Jobs = append(student.Jobs, job)
        }
        jobRows.Close()

        // Fetch education for the student
        eduRows, err := db.Query("SELECT school, degree, field_of_study, start_date, end_date, description FROM education WHERE user_id = $1", student.ID)
        if err != nil {
            return nil, err
        }

        for eduRows.Next() {
            var edu StudentEducation
            if err := eduRows.Scan(&edu.School, &edu.Degree, &edu.FieldOfStudy, &edu.StartDate, &edu.EndDate, &edu.Description); err != nil {
                return nil, err
            }
            student.Education = append(student.Education, edu)
        }
        eduRows.Close()

        students = append(students, student)
    }

    return students, nil
}

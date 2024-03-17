package models

import "time"

/////////////////////////////////////////////////////////////////////////////////////////
//                                      Tables                                         //
/////////////////////////////////////////////////////////////////////////////////////////

// User represents a user's basic information.
type User struct {
    ID                int       `json:"id"`
    Email             string    `json:"email"`
    FName             string    `json:"fname"`
    PasswordHash      string    `json:"-"`
    IsVerified        bool      `json:"isVerified"`
    VerificationToken string    `json:"-"`
    ResetToken        string    `json:"-"`
    ResetTokenExpires time.Time `json:"resetTokenExpires,omitempty"`
    HasSelectedRole   bool      `json:"hasSelectedRole"`
    CreatedAt         time.Time `json:"createdAt"`
    UpdatedAt         time.Time `json:"updatedAt"`
}

// Student represents a student's extended information.
type Student struct {
    UserID       int       `json:"userId"`
    Description  string    `json:"description"`
    ProfileImage string    `json:"profileImage"`
    IsCvCreated  string    `json:"is_cv_created"`
    CvPath       string    `json:"cv_path"`
    Interests    []string  `json:"interests"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}

// Job represents a job experience of a user.
type Job struct {
    ID          int       `json:"id"`
    UserID      int       `json:"userId"`
    Title       string    `json:"title"`
    Company     string    `json:"company"`
    StartDate   string    `json:"startDate"`
    EndDate     string    `json:"endDate,omitempty"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

// Education represents the educational background of a user.
type Education struct {
    ID            int       `json:"id"`
    UserID        int       `json:"userId"`
    School        string    `json:"school"`
    Degree        string    `json:"degree,omitempty"`
    FieldOfStudy  string    `json:"fieldOfStudy,omitempty"`
    StartDate     string `json:"startDate"`
    EndDate       string `json:"endDate,omitempty"`
    Description   string    `json:"description,omitempty"`
    CreatedAt     time.Time `json:"createdAt"`
    UpdatedAt     time.Time `json:"updatedAt"`
}

type JobPost struct {
    ID           int    `json:"id,omitempty"`
    UserID       int    `json:"user_id"`
    Title        string `json:"title"`
    Description  string `json:"description"`
    Requirements string `json:"requirements"`
    Status       string `json:"status"`
    Salary       string `json:"salary"`
    Address      string `json:"address"`
    CompanyEmail string `json:"company_email"`
    CompanyName  string `json:"company_name,omitempty"`
}

/////////////////////////////////////////////////////////////////////////////////////////
//                                      Requests                                       //
/////////////////////////////////////////////////////////////////////////////////////////

type StudentRegistrationRequest struct {
    UserID      int         `json:"userId"`
    Description string      `json:"description"`
    ImagePath   string      `json:"imagePath"`
    Interests   []string    `json:"interests"`
    Jobs        []Job       `json:"jobs"`
    Education   []Education `json:"education"`
}

type CompanyRegistrationRequest struct {
    UserID          int      `json:"userId"`
    Name            string   `json:"name"`
    Size            string   `json:"size"`
    Address         string   `json:"address"`
    Description     string   `json:"description"`
    ImagePath       string   `json:"image_path"`
    VideoPath       string   `json:"video_path"`
}

// StudentJob represents a job experience of a student.
type StudentJob struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    Title       string    `json:"title"`
    Company     string    `json:"company"`
    StartDate   string    `json:"start_date"`
    EndDate     string    `json:"end_date"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Company represents a company's information.
type Company struct {
    ID          int       `json:"id"`
    UserID      int       `json:"user_id"`
    Fname       string    `json:"fname"`
    Name        string    `json:"name"`
    Email       string 
    Size        string    `json:"size"`
    Address     string    `json:"address"`
    Description string    `json:"description"`
    ImagePath   string    `json:"image_path"`
    VideoPath   string    `json:"video_path"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
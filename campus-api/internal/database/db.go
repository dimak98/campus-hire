package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"campus-api/internal/config"
)


func Initialize(cfg *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

    if err := createSchema(db); err != nil {
        return nil, err
    }

    fmt.Println("Successfully connected to PostgreSQL!")
    return db, nil
}

func createSchema(db *sql.DB) error {
    usersQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        fname VARCHAR(255) NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        is_verified BOOLEAN DEFAULT FALSE,
        verification_token VARCHAR(255),
        reset_token VARCHAR(255),
        reset_token_expires TIMESTAMP,
        has_selected_role BOOLEAN DEFAULT FALSE,
        role VARCHAR(255),
        instagram_url VARCHAR(255),
        facebook_url VARCHAR(255),
        twitter_url VARCHAR(255),
        linkedin_url VARCHAR(255),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    _, err := db.Exec(usersQuery)
    if err != nil {
        return fmt.Errorf("error creating users table: %v", err)
    }

    studentsQuery := `
    CREATE TABLE IF NOT EXISTS students (
        user_id INTEGER PRIMARY KEY REFERENCES users(id),
        description TEXT CHECK (CHAR_LENGTH(description) <= 1000),
        profile_image VARCHAR(255),
        is_cv_created BOOLEAN DEFAULT FALSE,
        cv_path VARCHAR(255) DEFAULT 'null',
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(studentsQuery)
    if err != nil {
        return fmt.Errorf("error creating students table: %v", err)
    }

    jobsQuery := `CREATE TABLE IF NOT EXISTS jobs (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        title VARCHAR(255) NOT NULL,
        company VARCHAR(255) NOT NULL,
        start_date VARCHAR(255) NOT NULL,
        end_date VARCHAR(255),
        description TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(jobsQuery)
    if err != nil {
        return fmt.Errorf("error creating jobs table: %v", err)
    }

    studentJobApplicationsQuery := `CREATE TABLE IF NOT EXISTS student_job_applications (
        user_id INTEGER REFERENCES users(id),
        job_id INTEGER REFERENCES job_posts(id),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`
    
    _, err = db.Exec(studentJobApplicationsQuery)
    if err != nil {
        return fmt.Errorf("error creating student_job_applications table: %v", err)
    }

    educationQuery := `CREATE TABLE IF NOT EXISTS education (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        school VARCHAR(255) NOT NULL,
        degree VARCHAR(255),
        field_of_study VARCHAR(255),
        start_date VARCHAR(255) NOT NULL,
        end_date VARCHAR(255),
        description TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(educationQuery)
    if err != nil {
        return fmt.Errorf("error creating education table: %v", err)
    }

    companiesQuery := `CREATE TABLE IF NOT EXISTS companies (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        name VARCHAR(255) NOT NULL,
        size VARCHAR(255),
        address TEXT,
        description TEXT,
        image_path VARCHAR(255),
        video_path VARCHAR(255),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(companiesQuery)
    if err != nil {
        return fmt.Errorf("error creating companies table: %v", err)
    }

    jobPostsQuery := `CREATE TABLE IF NOT EXISTS job_posts (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        requirements TEXT NOT NULL,
        status VARCHAR(50) DEFAULT 'Open',
        salary VARCHAR(50) NOT NULL,
        address VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(jobPostsQuery)
    if err != nil {
        return fmt.Errorf("error creating job_posts table: %v", err)
    }

    return nil
}
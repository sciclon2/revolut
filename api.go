package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "regexp"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type request struct {
    DateOfBirth string `json:"dateOfBirth"`
}

type response struct {
    Message string `json:"message"`
}

var db *sql.DB

func createDatabaseIfNotExists() {
    dbHost, dbUser, dbPass, port, dbName, dbSslMode := getDBConfig()

    // Connect to the "postgres" database
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s", dbHost, port, dbUser, dbPass, dbSslMode)
    tempDb, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database connection: %v", err)
    }
    defer tempDb.Close()

    // Check if the target database exists
    var exists bool
    err = tempDb.QueryRow(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", dbName)).Scan(&exists)
    if err != nil && err != sql.ErrNoRows {
        log.Fatalf("Error checking if database exists: %v", err)
    }

    if !exists {
        // Create the target database
        _, err = tempDb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
        if err != nil {
            log.Fatalf("Error creating database: %v", err)
        }
        log.Printf("Database %s created successfully", dbName)
    } else {
        log.Printf("Database %s already exists", dbName)
    }
}

func initDB() {
    createDatabaseIfNotExists() // Ensure the database exists

    dbHost, dbUser, dbPass, port, dbName, dbSslMode := getDBConfig()

    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbHost, port, dbUser, dbPass, dbName, dbSslMode)
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT NOT NULL,
        dateOfBirth DATE NOT NULL
    );`

    _, err = db.Exec(createTableQuery)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }

    log.Println("Database initialized successfully")
}

func getDBConfig() (string, string, string, int, string, string) {
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")
    dbSslMode := os.Getenv("DB_SSL_MODE")

    dbPort := os.Getenv("DB_PORT")
    if dbPort == "" {
        dbPort = "5432" // Default port
    }

    port, err := strconv.Atoi(dbPort)
    if err != nil {
        log.Fatalf("Invalid port: %v", err)
    }

    return dbHost, dbUser, dbPass, port, dbName, dbSslMode
}

func getUserDateOfBirth(username string) (string, error) {
    var dateOfBirth string
    query := `SELECT dateOfBirth FROM users WHERE username = $1 LIMIT 1`
    err := db.QueryRow(query, username).Scan(&dateOfBirth)
    if err != nil {
        return "", err
    }
    return dateOfBirth, nil
}

func upsertUser(username, dateOfBirth string) (bool, error) {
    var userExists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 LIMIT 1)`
    err := db.QueryRow(query, username).Scan(&userExists)
    if err != nil {
        return false, err
    }

    if userExists {
        updateQuery := `UPDATE users SET dateOfBirth = $1 WHERE username = $2`
        _, err = db.Exec(updateQuery, dateOfBirth, username)
        return false, err
    } else {
        insertQuery := `
        INSERT INTO users (username, dateOfBirth)
        VALUES ($1, $2)`
        _, err = db.Exec(insertQuery, username, dateOfBirth)
        return true, err
    }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    if !ValidUsername(username) {
        http.Error(w, "Invalid username, must contain only letters", http.StatusBadRequest)
        return
    }

    dateOfBirth, err := getUserDateOfBirth(username)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            log.Printf("Database error: %v", err)
            http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
        }
        return
    }

    daysUntilBirthday := daysUntilNextBirthday(dateOfBirth)
    var message string
    if daysUntilBirthday == 0 {
        message = fmt.Sprintf("Happy Birthday, %s!", username)
    } else {
        message = fmt.Sprintf("Hello, %s! Your birthday is in %d days.", username, daysUntilBirthday)
    }

    res := response{Message: message}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(res)
}

func ValidUsername(username string) bool {
    re := regexp.MustCompile(`^[a-zA-Z]+$`)
    return re.MatchString(username)
}

func parseDate(date string) (time.Time, error) {
    formats := []string{"2006-01-02", "2006-01-02T15:04:05Z"}
    var parsedDate time.Time
    var err error
    for _, format := range formats {
        parsedDate, err = time.Parse(format, date)
        if err == nil {
            return parsedDate, nil
        }
    }
    return time.Time{}, err
}

func ValidDate(date string) bool {
    parsedDate, err := parseDate(date)
    if err != nil {
        return false
    }

    today := time.Now().UTC().Truncate(24 * time.Hour)

    // Check this is not a time traveler ;)
    if parsedDate.After(today) {
        return false
    }

    return true
}

func daysUntilNextBirthday(date string) int {
    parsedDate, err := parseDate(date)
    if err != nil {
        log.Printf("Error parsing date: %v", err)
        return -1
    }

    today := time.Now().UTC().Truncate(24 * time.Hour)
    nextBirthday := time.Date(today.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
    if today.After(nextBirthday) {
        nextBirthday = nextBirthday.AddDate(1, 0, 0)
    }

    days := int(nextBirthday.Sub(today).Hours() / 24)
    return days
}

func putHelloHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    username := vars["username"]

    if !ValidUsername(username) {
        http.Error(w, "Invalid username, must contain only letters", http.StatusBadRequest)
        return
    }

    var req request
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if !ValidDate(req.DateOfBirth) {
        http.Error(w, "Invalid date format, must be YYYY-MM-DD", http.StatusBadRequest)
        return
    }

    created, err := upsertUser(username, req.DateOfBirth)
    if err != nil {
        log.Printf("Database error: %v", err)
        http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
        return
    }

    if created {
        w.WriteHeader(http.StatusCreated) // 201 Created
    } else {
        w.WriteHeader(http.StatusOK) // 200 OK
    }
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}



func main() {
    initDB()

    r := mux.NewRouter()
    r.HandleFunc("/hello/{username}", helloHandler).Methods("GET")
    r.HandleFunc("/hello/{username}", putHelloHandler).Methods("PUT")
    r.HandleFunc("/health", healthCheckHandler).Methods("GET")

    server := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    log.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("could not start server: %s\n", err)
    }
}


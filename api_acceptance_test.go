// +build acceptance


package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "testing"
    "time"
    "math/rand"
    "strings"
)

var (
    baseURL     = "http://localhost:8080/hello"
    userID      string
    dateOfBirth string
    putData     request
    putBody     []byte
)

func setup() {
    userID = generateRandomUsername(15)
    dateOfBirth = time.Now().AddDate(-25, 0, 0).Format("2006-01-02")
    putData = request{DateOfBirth: dateOfBirth}
    var err error
    putBody, err = json.Marshal(putData)
    if err != nil {
        panic("Failed to marshal putData: " + err.Error())
    }
}

func TestCreateUser(t *testing.T) {
    setup()

    resp, err := httpPut(baseURL+"/"+userID, putBody)
    if err != nil {
        t.Fatalf("PUT /hello/%s failed: %v", userID, err)
    }
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status code 201 Created for new user, got %d", resp.StatusCode)
    }
}

func TestUpdateUser(t *testing.T) {
    setup()

    // First create the user
    resp, err := httpPut(baseURL+"/"+userID, putBody)
    if err != nil {
        t.Fatalf("Initial PUT /hello/%s failed: %v", userID, err)
    }
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status code 201 Created for new user, got %d", resp.StatusCode)
    }

    // Then update the user
    resp, err = httpPut(baseURL+"/"+userID, putBody)
    if err != nil {
        t.Fatalf("PUT /hello/%s failed: %v", userID, err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status code 200 OK for existing user, got %d", resp.StatusCode)
    }
}

func TestGetUser(t *testing.T) {
    setup()

    // First create the user
    resp, err := httpPut(baseURL+"/"+userID, putBody)
    if err != nil {
        t.Fatalf("Initial PUT /hello/%s failed: %v", userID, err)
    }
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status code 201 Created for new user, got %d", resp.StatusCode)
    }

    // Test valid GET request
    resp, err = httpGet(baseURL + "/" + userID)
    if err != nil {
        t.Fatalf("GET /hello/%s failed: %v", userID, err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status code 200 OK for GET, got %d", resp.StatusCode)
    }
}

func httpGet(url string) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    return client.Do(req)
}

func httpPut(url string, body []byte) (*http.Response, error) {
    client := &http.Client{}
    req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")
    return client.Do(req)
}

// Function to generate a random string of letters
func generateRandomUsername(length int) string {
    const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    var result strings.Builder
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < length; i++ {
        result.WriteByte(letters[rand.Intn(len(letters))])
    }
    return result.String()
}

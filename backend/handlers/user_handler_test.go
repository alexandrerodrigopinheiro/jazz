// backend/api/handlers/user_handler_test.go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jazz/backend/models"
	"jazz/backend/pkg/database"
)

func TestRegisterUserHandler(t *testing.T) {
	database.InitializeDatabase()

	payload := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(RegisterUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, status)
	}
}

func TestLoginHandler(t *testing.T) {
	database.InitializeDatabase()

	// Primeiro registramos um usuário
	user := models.User{Username: "testuser", Password: "password123"}
	user.HashPassword()
	db := database.GetDBInstance()
	db.Create(&user)

	payload := map[string]string{
		"username": "testuser",
		"password": "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, status)
	}

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	if _, ok := response["token"]; !ok {
		t.Error("Expected a token in the response but got none")
	}
}

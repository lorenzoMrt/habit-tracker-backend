package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

var mockDB sqlmock.Sqlmock

func setupMockDB(t *testing.T) *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mockDB = mock
	return db
}

func TestCreateHabit(t *testing.T) {
	db = setupMockDB(t)
	defer db.Close()

	tests := []struct {
		name           string
		input          Habit
		expectedStatus int
		expectedBody   Habit
	}{
		{
			name: "Valid Habit",
			input: Habit{
				Name:        "Exercise",
				Description: "Daily exercise routine",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: Habit{
				ID:          1,
				Name:        "Exercise",
				Description: "Daily exercise routine",
				Completed:   false,
			},
		},
		{
			name:           "Invalid JSON",
			input:          Habit{}, // Empty habit to simulate invalid JSON
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Valid Habit" {
				mockDB.ExpectQuery("INSERT INTO habits").
					WithArgs(tt.input.Name, tt.input.Description, false).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			}

			var body bytes.Buffer
			if tt.name != "Invalid JSON" {
				json.NewEncoder(&body).Encode(tt.input)
			} else {
				body.WriteString("invalid json")
			}

			req, err := http.NewRequest("POST", "/habits", &body)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createHabit)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated {
				var habit Habit
				if err := json.NewDecoder(rr.Body).Decode(&habit); err != nil {
					t.Fatal(err)
				}

				if habit.ID != tt.expectedBody.ID || habit.Name != tt.expectedBody.Name || habit.Description != tt.expectedBody.Description || habit.Completed != tt.expectedBody.Completed {
					t.Errorf("handler returned unexpected body: got %v want %v", habit, tt.expectedBody)
				}
			}
		})
	}
}

func TestListHabits(t *testing.T) {
	db = setupMockDB(t)
	defer db.Close()

	mockDB.ExpectQuery("SELECT id, name, description, completed FROM habits").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "completed"}).
			AddRow(1, "Exercise", "Daily exercise routine", false).
			AddRow(2, "Read", "Read a book", true))

	req, err := http.NewRequest("GET", "/habits", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listHabits)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseHabits []Habit
	if err := json.NewDecoder(rr.Body).Decode(&responseHabits); err != nil {
		t.Fatal(err)
	}

	if len(responseHabits) != 2 {
		t.Errorf("handler returned unexpected number of habits: got %v want %v", len(responseHabits), 2)
	}
}

func TestCompleteHabit(t *testing.T) {
	db = setupMockDB(t)
	defer db.Close()

	tests := []struct {
		name           string
		habitID        string
		expectedStatus int
		expectedBody   Habit
	}{
		{
			name:           "Valid ID",
			habitID:        "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Invalid ID",
			habitID:        "invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent ID",
			habitID:        "3",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Valid ID" {
				mockDB.ExpectExec("UPDATE habits SET completed = true WHERE id = ?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else if tt.name == "Non-existent ID" {
				mockDB.ExpectExec("UPDATE habits SET completed = true WHERE id = ?").
					WithArgs(3).
					WillReturnResult(sqlmock.NewResult(1, 0))
			}

			req, err := http.NewRequest("PUT", "/habits/"+tt.habitID+"/complete", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/habits/{id}/complete", completeHabit)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
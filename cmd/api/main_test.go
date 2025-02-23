package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateHabit(t *testing.T) {
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
	// Prepare some initial data
	habits = []Habit{
		{ID: 1, Name: "Exercise", Completed: false},
		{ID: 2, Name: "Read", Completed: true},
	}

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

	if len(responseHabits) != len(habits) {
		t.Errorf("handler returned unexpected number of habits: got %v want %v", len(responseHabits), len(habits))
	}
}

func TestCompleteHabit(t *testing.T) {
	// Prepare some initial data
	habits = []Habit{
		{ID: 1, Name: "Exercise", Completed: false},
		{ID: 2, Name: "Read", Completed: false},
	}

	tests := []struct {
		name           string
		habitID        string
		expectedStatus int
		expectedBody   Habit
	}{
		{
			name:           "Valid ID",
			habitID:        "1",
			expectedStatus: http.StatusOK,
			expectedBody:   Habit{ID: 1, Name: "Exercise", Completed: true},
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

			if tt.expectedStatus == http.StatusOK {
				var habit Habit
				if err := json.NewDecoder(rr.Body).Decode(&habit); err != nil {
					t.Fatal(err)
				}

				if habit.ID != tt.expectedBody.ID || habit.Name != tt.expectedBody.Name || habit.Completed != tt.expectedBody.Completed {
					t.Errorf("handler returned unexpected body: got %v want %v", habit, tt.expectedBody)
				}
			}
		})
	}
}
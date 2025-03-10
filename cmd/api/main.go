package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Habit representa la estructura de un hábito.
type Habit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
}

var db *sql.DB

func main() {
	var err error
	// Get database connection string from environment variable or use default
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=password dbname=habit_tracker host=localhost sslmode=disable"
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crear un router usando Gorilla Mux
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	// Handle all OPTIONS requests
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
	}).Methods("OPTIONS")

	// Endpoint para crear un hábito
	router.HandleFunc("/habits", createHabit).Methods("POST")
	// Endpoint para listar hábitos
	router.HandleFunc("/habits", listHabits).Methods("GET")
	// Endpoint para completar un hábito (actualización)
	router.HandleFunc("/habits/{id}/complete", completeHabit).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor corriendo en http://0.0.0.0:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// createHabit permite registrar un nuevo hábito.
func createHabit(w http.ResponseWriter, r *http.Request) {
	var habit Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO habits(name, description, completed) VALUES($1, $2, $3) RETURNING id",
		habit.Name, habit.Description, false).Scan(&habit.ID)
	if err != nil {
		http.Error(w, "Error al insertar en la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

// listHabits lista todos los hábitos registrados.
func listHabits(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, description, completed FROM habits")
	if err != nil {
		http.Error(w, "Error al consultar la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		if err := rows.Scan(&habit.ID, &habit.Name, &habit.Description, &habit.Completed); err != nil {
			http.Error(w, "Error al escanear la fila: "+err.Error(), http.StatusInternalServerError)
			return
		}
		habits = append(habits, habit)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

// completeHabit marca un hábito como completado.
func completeHabit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE habits SET completed = true WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error al actualizar la base de datos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, "Hábito no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment variable or use default
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:5173"
		}

		// Check if the request origin is in the list of allowed origins
		origin := r.Header.Get("Origin")
		if origin != "" {
			// If ALLOWED_ORIGINS contains "*", allow all origins
			if allowedOrigins == "*" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// Check if the origin is in the list of allowed origins
				for _, allowedOrigin := range strings.Split(allowedOrigins, ",") {
					if strings.TrimSpace(allowedOrigin) == origin {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		} else {
			// Fallback to the first allowed origin if no origin header
			w.Header().Set("Access-Control-Allow-Origin", strings.Split(allowedOrigins, ",")[0])
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600") // Cache preflight response for 1 hour
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

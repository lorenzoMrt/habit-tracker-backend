package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Habit representa la estructura de un hábito.
type Habit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
}

var (
	habits      []Habit      // almacena los hábitos
	nextID      = 1          // id incremental para cada hábito
	habitsMutex sync.Mutex   // mutex para evitar condiciones de carrera
)

func main() {
	// Crear un router usando Gorilla Mux
	router := mux.NewRouter()

	// Endpoint para crear un hábito
	router.HandleFunc("/habits", createHabit).Methods("POST")
	// Endpoint para listar hábitos
	router.HandleFunc("/habits", listHabits).Methods("GET")
	// Endpoint para completar un hábito (actualización)
	router.HandleFunc("/habits/{id}/complete", completeHabit).Methods("PUT")

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// createHabit permite registrar un nuevo hábito.
func createHabit(w http.ResponseWriter, r *http.Request) {
	var habit Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		http.Error(w, "Error al decodificar JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	habitsMutex.Lock()
	habit.ID = nextID
	nextID++
	habit.Completed = false
	habits = append(habits, habit)
	habitsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(habit)
}

// listHabits lista todos los hábitos registrados.
func listHabits(w http.ResponseWriter, r *http.Request) {
	habitsMutex.Lock()
	defer habitsMutex.Unlock()

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

	habitsMutex.Lock()
	defer habitsMutex.Unlock()

	for i, habit := range habits {
		if habit.ID == id {
			habits[i].Completed = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(habits[i])
			return
		}
	}

	http.Error(w, "Hábito no encontrado", http.StatusNotFound)
}

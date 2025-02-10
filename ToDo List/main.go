package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	tasks  = []Task{}
	nextID = 1
	mu     sync.Mutex
)

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i, task := range tasks {
		if task.ID == req.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Раздаём HTML/CSS/JS
	http.HandleFunc("/tasks", getTasks)
	http.HandleFunc("/add", addTask)
	http.HandleFunc("/delete", deleteTask)

	http.ListenAndServe(":8080", nil)
}

package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/yourname/download-service/internal/downloader"
	"github.com/yourname/download-service/internal/model"
	"github.com/yourname/download-service/internal/store"
)

var st *store.Store
var dl *downloader.Worker

func RegisterHandlers(store *store.Store, worker *downloader.Worker) {
	st = store
	dl = worker

	http.HandleFunc("/tasks", createTaskHandler) //Создание таска с указанием урла
	http.HandleFunc("/tasks/", getTaskHandler)   //Запрос статуса задачи (сделал так что выводит всю информацию о задаче)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string   `json:"name"`
		URLs []string `json:"urls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := &model.Task{
		ID:     uuid.NewString(),
		Name:   req.Name,
		URLs:   req.URLs,
		Status: "pending",
	}

	st.SaveTask(task)
	dl.Enqueue(task)

	resp := map[string]string{"task_id": task.ID}
	json.NewEncoder(w).Encode(resp)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}
	id := parts[2]

	task := st.GetTask(id)
	if task == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

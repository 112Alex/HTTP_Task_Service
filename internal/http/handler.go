package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"tz/internal/task"
)

// Handler — структура, содержащая ссылки на менеджер задач
type Handler struct {
	TaskManager *task.TaskManager
}

// NewHandler — инициализация HTTP-обработчика
func NewHandler(tm *task.TaskManager) *Handler {
	return &Handler{TaskManager: tm}
}

// handleCreateTask — обработчик POST /tasks
func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	task := h.TaskManager.CreateTask()

	w.Header().Set("Location", "/tasks/"+task.ID)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(task)
}

// handleGetTask — обработчик GET /tasks/{id}
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		writeErr(w, http.StatusBadRequest, "не указан ID задачи")
		return
	}

	t, err := h.TaskManager.GetTask(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "задача не найдена")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// handleDeleteTask — обработчик DELETE /tasks/{id}
func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := getIDFromPath(r.URL.Path)
	if id == "" {
		writeErr(w, http.StatusBadRequest, "не указан ID задачи")
		return
	}

	err := h.TaskManager.DeleteTask(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "задача не найдена")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleListTasks — обработчик GET /tasks
func (h *Handler) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.TaskManager.ListTasks()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tasks)

	if len(tasks) == 0 {
		writeErr(w, http.StatusNotFound, "Задачи не найдены")
		return
	}
}

// getIDFromPath — утилита для извлечения ID из пути
func getIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2]
	}
	return ""
}

// writeErr возвращает JSON-ошибку единообразно
func writeErr(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

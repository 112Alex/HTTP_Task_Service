package http

import (
	"log"
	"net/http"
	"time"
)

// RegisterRoutes — регистрирует маршруты API на HTTP mux
func RegisterRoutes(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Middleware логирования
	logMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		})
	}

	// POST /tasks и GET /tasks
	mux.Handle("/tasks", logMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.handleCreateTask(w, r)
		case http.MethodGet:
			h.handleListTasks(w, r)
		default:
			writeErr(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		}
	})))

	// GET /tasks/{id}, DELETE /tasks/{id}
	mux.Handle("/tasks/", logMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.handleGetTask(w, r)
		case http.MethodDelete:
			h.handleDeleteTask(w, r)
		default:
			writeErr(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		}
	})))

	return mux
}

package main

import (
	"fmt"
	"log"
	"net/http"

	httpapi "tz/internal/http"
	"tz/internal/task"
)

func main() {
	// Инициализация менеджера задач
	taskManager := task.NewManager()

	// Создание обработчиков HTTP
	handler := httpapi.NewHandler(taskManager)

	// Регистрируем маршруты
	mux := httpapi.RegisterRoutes(handler)

	// Запуск HTTP-сервера на порту 8080
	const addr = ":8080"
	fmt.Printf("Сервер запущен на http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

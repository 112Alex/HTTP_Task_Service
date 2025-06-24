package task

import (
	"context"
	"time"
)

// TaskStatus — перечисление возможных статусов задачи
type TaskStatus string

const (
	StatusCreated  TaskStatus = "created"
	StatusRunning  TaskStatus = "running"
	StatusFinished TaskStatus = "completed"
	StatusFailed   TaskStatus = "failed"
	StatusCanceled TaskStatus = "canceled"
)

// Task — структура, описывающая задачу
type Task struct {
	ID         string             `json:"id"`         // Уникальный идентификатор задачи
	Status     TaskStatus         `json:"status"`     // Статус выполнения
	CreatedAt  time.Time          `json:"created_at"` // Время создания
	Duration   time.Duration      `json:"duration"`   // Время выполнения задачи
	Result     string             `json:"result"`     // Результат выполнения (можно оставить пустым)
	StartedAt  time.Time          `json:"-"`          // Внутреннее поле (не сериализуется)
	FinishedAt time.Time          `json:"-"`          // Внутреннее поле (не сериализуется)
	ctx        context.Context    `json:"-"`          // Контекст для отмены задачи
	cancel     context.CancelFunc `json:"-"`          // Функция отмены задачи
}

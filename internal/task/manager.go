package task

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// TaskManager — структура, управляющая задачами в памяти
type TaskManager struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

// NewManager — создаёт новый менеджер задач
func NewManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

// CreateTask — создаёт новую задачу и запускает её выполнение в фоне
func (m *TaskManager) CreateTask() *Task {
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		ID:        newID(),
		Status:    StatusCreated,
		CreatedAt: time.Now(),
		ctx:       ctx,
		cancel:    cancel,
	}

	m.mu.Lock()
	m.tasks[task.ID] = task
	m.mu.Unlock()

	go m.runTask(task)

	return task
}

// newID генерирует 32-символьный hex-идентификатор без сторонних зависимостей.
func newID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// runTask — выполняет задачу в фоне (эмуляция I/O bound операции)
func (m *TaskManager) runTask(task *Task) {
	m.updateStatus(task.ID, StatusRunning)
	task.StartedAt = time.Now()

	select {
	case <-time.After(3 * time.Minute):
		task.FinishedAt = time.Now()
		task.Duration = task.FinishedAt.Sub(task.StartedAt)
		m.mu.Lock()
		task.Status = StatusFinished
		task.Result = "Task completed successfully"
		m.mu.Unlock()
	case <-task.ctx.Done():
		task.FinishedAt = time.Now()
		task.Duration = task.FinishedAt.Sub(task.StartedAt)
		m.mu.Lock()
		task.Status = StatusCanceled
		m.mu.Unlock()
	}
}

// GetTask — возвращает задачу по ID
func (m *TaskManager) GetTask(id string) (*Task, error) {
	m.mu.RLock()
	t, exists := m.tasks[id]
	if !exists {
		m.mu.RUnlock()
		return nil, errors.New("задача не найдена")
	}

	// Создаём копию, чтобы избежать гонки при чтении/записи
	copyTask := *t
	// Если задача ещё выполняется, обновляем длительность на лету
	if copyTask.Status == StatusRunning {
		copyTask.Duration = time.Since(copyTask.StartedAt)
	}
	m.mu.RUnlock()
	return &copyTask, nil
}

// DeleteTask — удаляет задачу по ID
func (m *TaskManager) DeleteTask(id string) error {
	m.mu.Lock()
	task, exists := m.tasks[id]
	if !exists {
		m.mu.Unlock()
		return errors.New("задача не найдена")
	}

	// Отменяем выполнение, если ещё идёт
	task.cancel()
	delete(m.tasks, id)
	m.mu.Unlock()
	return nil
}

// updateStatus — служебная функция для обновления статуса задачи
func (m *TaskManager) updateStatus(id string, status TaskStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if task, ok := m.tasks[id]; ok {
		task.Status = status
	}
}

// ListTasks возвращает срез копий всех задач
func (m *TaskManager) ListTasks() []Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		copyTask := *t
		if copyTask.Status == StatusRunning {
			copyTask.Duration = time.Since(copyTask.StartedAt)
		}
		tasks = append(tasks, copyTask)
	}
	return tasks
}

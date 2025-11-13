package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Task representa uma tarefa no sistema
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var (
	allowedStatuses = map[string]bool{
		"todo":        true,
		"in_progress": true,
		"done":        true,
	}
)

// Store define a interface de persistência usada pelos handlers
type Store interface {
	Initialize() error
	GetAll() []Task
	Create(title, description, status string) (*Task, error)
	Update(id, title, description, status string) (*Task, error)
	Delete(id string) error
}

// Storage (file-based) gerencia a persistência em JSON
type Storage struct {
	tasks     []Task
	currentID int
	mu        sync.Mutex
	dataFile  string
}

// New cria uma nova instância file-based de Storage e implementa Store
func New(dataFilePath string) Store {
	return &Storage{
		tasks:     []Task{},
		currentID: 1,
		dataFile:  dataFilePath,
	}
}

// Initialize carrega as tarefas do disco
func (s *Storage) Initialize() error {
	if err := s.ensureDataDir(); err != nil {
		return err
	}

	b, err := os.ReadFile(s.dataFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	var loaded []Task
	if len(b) > 0 {
		if err := json.Unmarshal(b, &loaded); err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = loaded
	// Recalcula currentID com base nos IDs existentes
	maxID := 0
	for _, t := range s.tasks {
		if id, err := strconv.Atoi(t.ID); err == nil && id > maxID {
			maxID = id
		}
	}
	s.currentID = maxID + 1

	return nil
}

// GetAll retorna todas as tarefas
func (s *Storage) GetAll() []Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retorna uma cópia para evitar modificações concorrentes
	result := make([]Task, len(s.tasks))
	copy(result, s.tasks)
	return result
}

// Create cria uma nova tarefa
func (s *Storage) Create(title, description, status string) (*Task, error) {
	if title == "" {
		return nil, errors.New("título é obrigatório")
	}

	if status == "" {
		status = "todo"
	} else if !allowedStatuses[status] {
		return nil, errors.New("status inválido. Use: todo, in_progress, done")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{
		ID:          strconv.Itoa(s.currentID),
		Title:       title,
		Description: description,
		Status:      status,
	}

	s.currentID++
	s.tasks = append(s.tasks, task)

	if err := s.save(); err != nil {
		return nil, err
	}

	return &task, nil
}

// Update atualiza uma tarefa existente
func (s *Storage) Update(id, title, description, status string) (*Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			if title != "" {
				s.tasks[i].Title = title
			}
			if description != "" {
				s.tasks[i].Description = description
			}
			if status != "" {
				if !allowedStatuses[status] {
					return nil, errors.New("status inválido. Use: todo, in_progress, done")
				}
				s.tasks[i].Status = status
			}

			if err := s.save(); err != nil {
				return nil, err
			}

			result := s.tasks[i]
			return &result, nil
		}
	}

	return nil, errors.New("tarefa não encontrada")
}

// Delete remove uma tarefa
func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			if err := s.save(); err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("tarefa não encontrada")
}

// save persiste os dados no disco (deve ser chamado dentro de um lock)
func (s *Storage) save() error {
	if err := s.ensureDataDir(); err != nil {
		return err
	}

	b, err := json.MarshalIndent(s.tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.dataFile, b, 0644)
}

// ensureDataDir cria o diretório se não existir
func (s *Storage) ensureDataDir() error {
	dir := filepath.Dir(s.dataFile)
	return os.MkdirAll(dir, 0755)
}

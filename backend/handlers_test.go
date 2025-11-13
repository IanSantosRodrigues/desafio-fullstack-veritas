package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"desafio-fullstack-veritas/backend/storage"
)

// setupTestStorage cria um storage temporário para testes
func setupTestStorage(t *testing.T) storage.Store {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "tasks.json")
	st := storage.New(dataFile)
	if err := st.Initialize(); err != nil {
		t.Fatalf("Failed to initialize storage: %v", err)
	}
	return st
}

// TestGetTasksEmpty testa o endpoint GET /tasks quando não há tarefas
func TestGetTasksEmpty(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasks)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "[]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q",
			rr.Body.String(), expected)
	}
}

// TestCreateTask testa a criação de uma tarefa
func TestCreateTask(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	payload := map[string]string{
		"title":       "Test Task",
		"description": "Test Description",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var task storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if task.Title != "Test Task" {
		t.Errorf("handler returned unexpected title: got %v want %v",
			task.Title, "Test Task")
	}

	if task.Status != "todo" {
		t.Errorf("handler returned unexpected default status: got %v want %v",
			task.Status, "todo")
	}

	if task.ID == "" {
		t.Error("handler did not assign ID to task")
	}
}

// TestCreateTaskMissingTitle testa criação sem título obrigatório
func TestCreateTaskMissingTitle(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	payload := map[string]string{
		"description": "Test Description",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Errorf("Failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("handler did not return error message")
	}
}

// TestCreateTaskInvalidStatus testa criação com status inválido
func TestCreateTaskInvalidStatus(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	payload := map[string]string{
		"title":  "Test Task",
		"status": "invalid_status",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// TestGetTasksWithData testa GET /tasks com tarefas existentes
func TestGetTasksWithData(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	// Cria algumas tarefas
	st.Create("Task 1", "Description 1", "todo")
	st.Create("Task 2", "Description 2", "in_progress")

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasks)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var tasks []storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("handler returned wrong number of tasks: got %v want %v",
			len(tasks), 2)
	}
}

// TestDeleteTaskNotFound testa deleção de tarefa inexistente
func TestDeleteTaskNotFound(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	err := st.Delete("999")
	if err == nil {
		t.Error("Expected error when deleting non-existent task")
	}
}

// TestCreateTaskWithCustomStatus testa criação com status customizado
func TestCreateTaskWithCustomStatus(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	payload := map[string]string{
		"title":  "Test Task",
		"status": "in_progress",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTask)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var task storage.Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if task.Status != "in_progress" {
		t.Errorf("handler returned unexpected status: got %v want %v",
			task.Status, "in_progress")
	}
}

// TestErrorResponseFormat testa se as respostas de erro estão em JSON
func TestErrorResponseFormat(t *testing.T) {
	st := setupTestStorage(t)
	InitializeStorage(st)

	// Tenta criar tarefa sem título
	payload := map[string]string{
		"description": "Only description",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateTask)
	handler.ServeHTTP(rr, req)

	// Verifica se o Content-Type é JSON
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected JSON content-type, got %v", ct)
	}

	// Verifica se consegue fazer parse como JSON
	var errResp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}
}

// TestTaskPersistence testa se tarefas são persistidas no arquivo
func TestTaskPersistence(t *testing.T) {
	tempDir := t.TempDir()
	dataFile := filepath.Join(tempDir, "tasks.json")

	// Cria storage e adiciona tarefas
	st1 := storage.New(dataFile)
	st1.Initialize()
	st1.Create("Task 1", "Description 1", "todo")
	st1.Create("Task 2", "Description 2", "in_progress")

	// Cria um novo storage com o mesmo arquivo
	st2 := storage.New(dataFile)
	st2.Initialize()

	// Verifica se as tarefas foram carregadas
	tasks := st2.GetAll()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 persisted tasks, got %v", len(tasks))
	}

	// Verifica se os dados estão corretos
	if tasks[0].Title != "Task 1" {
		t.Errorf("Expected first task title 'Task 1', got %v", tasks[0].Title)
	}

	if tasks[1].Title != "Task 2" {
		t.Errorf("Expected second task title 'Task 2', got %v", tasks[1].Title)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

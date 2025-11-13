package models

// Task representa a estrutura de uma tarefa no sistema.
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// ErrorResponse é a estrutura padrão de resposta de erro
type ErrorResponse struct {
	Error string `json:"error"`
}

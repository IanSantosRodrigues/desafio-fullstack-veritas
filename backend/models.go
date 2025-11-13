package main

import (
	"encoding/json"
	"net/http"
)

// Task representa uma tarefa no sistema (mantido por compatibilidade, usar storage.Task)
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

// RespondWithError envia uma resposta de erro padronizada em JSON
func RespondWithError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// RespondWithJSON envia uma resposta padronizada em JSON
func RespondWithJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

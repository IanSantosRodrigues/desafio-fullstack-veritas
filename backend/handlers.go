package main

import (
	"encoding/json"
	"net/http"

	"desafio-fullstack-veritas/backend/storage"

	"github.com/gorilla/mux"
)

// taskStorage é a instância global de storage (interface)
var taskStorage storage.Store

// InitializeStorage inicializa o storage de tarefas
func InitializeStorage(st storage.Store) {
	taskStorage = st
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := taskStorage.GetAll()

	// Se não há tarefas, retorna um array vazio em vez de null
	if tasks == nil {
		tasks = []storage.Task{}
	}

	RespondWithJSON(w, tasks, http.StatusOK)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var incoming struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		RespondWithError(w, "Corpo da requisição inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	task, err := taskStorage.Create(incoming.Title, incoming.Description, incoming.Status)
	if err != nil {
		RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	RespondWithJSON(w, task, http.StatusCreated)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var incoming struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		RespondWithError(w, "Corpo da requisição inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	task, err := taskStorage.Update(id, incoming.Title, incoming.Description, incoming.Status)
	if err != nil {
		RespondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	RespondWithJSON(w, task, http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := taskStorage.Delete(id)
	if err != nil {
		RespondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

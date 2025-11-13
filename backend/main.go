package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"desafio-fullstack-veritas/backend/storage"

	"github.com/gorilla/mux"
)

func main() {
	// Configuração via variáveis de ambiente
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080" // Valor padrão
	}

	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = filepath.Join("backend", "data", "tasks.json") // Valor padrão
	}

	// Inicializa o storage selecionado por ambiente
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "json" // padrão
	}

	var st storage.Store
	switch storageType {
	case "sqlite":
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			// valor padrão para sqlite
			dsn = filepath.Join("backend", "data", "tasks.db")
		}
		st = storage.NewSQLite(dsn)
	default:
		st = storage.New(dataFile)
	}

	if err := st.Initialize(); err != nil {
		log.Println("Aviso: falha ao inicializar storage:", err)
	}
	InitializeStorage(st)

	// Configura as rotas
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")

	log.Println("Servidor rodando na porta " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With, Origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

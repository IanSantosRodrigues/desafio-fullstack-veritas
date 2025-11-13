//go:build !sqlite
// +build !sqlite

package storage

import "errors"

// NewSQLite (stub)
// Esta versão é compilada quando a build NÃO habilita a tag 'sqlite'.
// Retorna um Store que sempre falha com uma mensagem explicativa.
func NewSQLite(dsn string) Store {
	return &unsupportedSQLite{reason: "SQLite não está habilitado na build. Compile com a tag 'sqlite' e instale o driver (ex: github.com/mattn/go-sqlite3)."}
}

type unsupportedSQLite struct {
	reason string
}

func (s *unsupportedSQLite) Initialize() error { return errors.New(s.reason) }
func (s *unsupportedSQLite) GetAll() []Task    { return []Task{} }
func (s *unsupportedSQLite) Create(title, description, status string) (*Task, error) {
	return nil, errors.New(s.reason)
}
func (s *unsupportedSQLite) Update(id, title, description, status string) (*Task, error) {
	return nil, errors.New(s.reason)
}
func (s *unsupportedSQLite) Delete(id string) error { return errors.New(s.reason) }

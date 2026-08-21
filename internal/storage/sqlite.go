package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Storage wraps a *sql.DB and provides message-saving methods.
type Storage struct {
	db *sql.DB
}

// Close fecha a conexão com o banco de dados.
func (s *Storage) Close() error {
	return s.db.Close()
}

// InitDB abre (ou cria) o arquivo .db e garante que o schema existe.
// Cria o diretório pai caso não exista (ex: ./data/mensageforia.db).
func InitDB(path string) (*Storage, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("criar diretório do banco: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}

	// Simple migration: create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			theme TEXT NOT NULL,
			prompt TEXT NOT NULL,
			response TEXT NOT NULL,
			timestamp TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("criar tabela: %w", err)
	}

	return &Storage{db: db}, nil
}

// SaveMessage insere um novo registro de mensagem gerada.
// Se falhar, apenas loga e retorna nil (não quebra o fluxo principal).
func (s *Storage) SaveMessage(theme, prompt, response, timestamp string) error {
	_, err := s.db.Exec(
		`INSERT INTO messages (theme, prompt, response, timestamp) VALUES (?, ?, ?, ?)`,
		theme, prompt, response, timestamp,
	)
	if err != nil {
		log.Printf("aviso: falha ao salvar no SQLite (continuando mesmo assim): %v", err)
		return nil // não quebra o fluxo principal
	}
	return nil
}

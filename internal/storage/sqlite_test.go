package storage_test

import (
	"path/filepath"
	"testing"

	"mensageforia/internal/storage"
)

// Regressão: InitDB deve criar o diretório pai se não existir
// (ex: ./data/mensageforia.db dentro do container).
func TestInitDBCriaDiretorioPai(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "pasta-nova", "mensageforia.db")

	s, err := storage.InitDB(path)
	if err != nil {
		t.Fatalf("InitDB com diretório inexistente: %v", err)
	}
	defer s.Close()

	// Deve conseguir inserir e ler de volta
	if err := s.SaveMessage("foco", "prompt", "resposta", "2026-08-21-08h"); err != nil {
		t.Fatalf("SaveMessage: %v", err)
	}
}

func TestInitDBReusaBancoExistente(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mensageforia.db")

	s1, err := storage.InitDB(path)
	if err != nil {
		t.Fatalf("primeiro InitDB: %v", err)
	}
	if err := s1.SaveMessage("coragem", "p", "r", "t1"); err != nil {
		t.Fatalf("SaveMessage #1: %v", err)
	}
	if err := s1.Close(); err != nil {
		t.Fatalf("Close #1: %v", err)
	}

	// Reabrir não pode falhar (schema já existe)
	s2, err := storage.InitDB(path)
	if err != nil {
		t.Fatalf("reabrir banco existente: %v", err)
	}
	defer s2.Close()

	if err := s2.SaveMessage("foco", "p2", "r2", "t2"); err != nil {
		t.Fatalf("SaveMessage #2: %v", err)
	}
}

package main

import (
	"fmt"
	"os"

	"mensageforia/internal/ollama"
	"mensageforia/internal/sheduler"
	"mensageforia/internal/storage"
)

func main() {
	ollamaURL := "http://ollama:11434"
	modelName := "llama3.2:1b"

	db, err := storage.InitDB("./.db")
	if err != nil {
		fmt.Printf("falha ao inicializar SQLite: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	client := ollama.NewClient(ollamaURL, modelName)

	if err := sheduler.SetupCron(client, db); err != nil {
		fmt.Printf("falha ao setup cron: %v\n", err)
		os.Exit(1)
	}
}


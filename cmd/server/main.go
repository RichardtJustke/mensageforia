package main

import (
	"fmt"
	"os"

	// Embute os dados de timezone no binário: o alpine da imagem final
	// não tem tzdata instalado, e o scheduler usa America/Sao_Paulo.
	_ "time/tzdata"

	"mensageforia/internal/message"
	"mensageforia/internal/ollama"
	"mensageforia/internal/scheduler"
	"mensageforia/internal/storage"
)

func main() {
	ollamaURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	modelName := os.Getenv("OLLAMA_MODEL")
	if modelName == "" {
		modelName = "llama3.2:1b"
	}

	db, err := storage.InitDB("./data/mensageforia.db")
	if err != nil {
		fmt.Printf("falha ao inicializar SQLite: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	client := ollama.NewClient(ollamaURL, modelName)

	// Modo --once: gera uma mensagem agora e sai. Útil pra validar o deploy
	// sem esperar o horário do cron (ex: docker compose run --rm app ./mensageforia --once).
	if len(os.Args) > 1 && os.Args[1] == "--once" {
		if err := message.GenerateAndCommit(client, db); err != nil {
			fmt.Printf("falha ao gerar mensagem: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := scheduler.SetupCron(client, db); err != nil {
		fmt.Printf("falha ao setup cron: %v\n", err)
		os.Exit(1)
	}
}

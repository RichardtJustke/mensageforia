package main

import (
	"os"

	"mensageforia/internal/ollama"
	"mensageforia/internal/sheduler"
)

func main() {
	client := ollama.NewClient("http://ollama:11434", "llama3.2:3b") // ajusta host/model

	if err := sheduler.setupCron(client); err != nil {
		os.Exit(1)
	}
}

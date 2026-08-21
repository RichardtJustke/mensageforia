package message

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"

	"mensageforia/internal/git"
	"mensageforia/internal/storage"
	"mensageforia/internal/ollama"
)

var temas = []string{
	"superação",
	"foco",
	"gratidão",
	"humor",
	"disciplina",
	"determinação",
	"sucesso",
	"coragem",
	"recomeço",
	"persistência",
	"determinação",
	"autoconfiança",
	"sonhos",
	"mudança",
	"resiliência",
	"propósito",
	"crescimento pessoal",
	"conquista",
}

const promptTemplate = "Gere uma mensagem motivacional curta sobre %s, em português, tom inspirador."

func randTema(temas []string) string {
	indice := rand.IntN(len(temas))
	sorteado := temas[indice]
	return sorteado
}
func temaPrompt() string {
	temaEscolhido := randTema(temas)
	promptFinal := fmt.Sprintf(promptTemplate, temaEscolhido)

	return promptFinal

}

// GenerateAndCommit gera a mensagem via Ollama e persiste em SQLite + arquivos + git.
// Se algum passo falhar, apenas loga e continua (não trava o fluxo principal).
func GenerateAndCommit(client *ollama.Client, s *storage.Storage) error {
	prompt := temaPrompt()

	resposta, err := client.Generate(prompt)
	if err != nil {
		return fmt.Errorf("gerar mensagem: %w", err)
	}

	now := time.Now()
	timestamp := now.Format("2006-01-02-15h04")
	tema := randTema(temas)

	// 1. Salvar no SQLite
	if err := s.SaveMessage(tema, prompt, resposta, timestamp); err != nil {
		fmt.Printf("aviso: não foi possível salvar no SQLite: %v\n", err)
		// continua mesmo assim
	}

	// 2. Salvar arquivo .md em messages/
	if err := saveMarkdown(tema, resposta, timestamp); err != nil {
		fmt.Printf("aviso: não foi possível salvar .md: %v\n", err)
	}

	// 3. Commit e push no git
	if err := git.CommitAndPush(tema, timestamp); err != nil {
		fmt.Printf("aviso: falha no git commit/push: %v\n", err)
		// continua mesmo assim
	}

	fmt.Println("Resposta gerada e persistida:", resposta)
	return nil
}

// saveMarkdown escreve o arquivo .md na pasta messages/.
// Nome: YYYY-MM-DD-HH.md (ex: 2026-08-21-08h.md)
// Conteúdo: resposta + metadado (tema, horário) no formato YAML no início + resposta.
func saveMarkdown(tema, resposta, timestamp string) error {
	// Garante que a pasta existe
	if err := os.MkdirAll("messages", 0755); err != nil {
		return fmt.Errorf("criar pasta messages: %w", err)
	}

	// Nome do arquivo com data + hora para evitar colisão
	filename := fmt.Sprintf("%s.md", timestamp)
	filePath := filepath.Join("messages", filename)

	// Conteúdo: metadado opcional no início + a resposta pura
	content := fmt.Sprintf("---\\ntheme: %s\\ntimestamp: %s\\n---\\n%s", tema, timestamp, resposta)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("escrever arquivo %s: %w", filename, err)
	}

	fmt.Printf("Arquivo salvo: %s\n", filePath)
	return nil
}
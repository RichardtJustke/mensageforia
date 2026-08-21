package message

import (
	"fmt"
	"math/rand/v2"
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
func GenerateAndCommit(client *ollama.Client) error {
	prompt := temaPrompt()

	resposta, err := client.Generate(prompt)
	if err != nil {
		return fmt.Errorf("gerar mensagem: %w", err)
	}

	fmt.Println("Resposta gerada:", resposta)

	// próximos passos aqui: storage.Save(...), salvar .md, git commit

	return nil
}

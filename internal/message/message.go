package message

import (
	"fmt"
	"math/rand/v2"
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

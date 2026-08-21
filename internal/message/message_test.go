package message

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Garante que a lista de temas não tem duplicatas (regressão do tema repetido).
func TestTemasSemDuplicatas(t *testing.T) {
	vistos := make(map[string]bool)
	for _, tema := range temas {
		if vistos[tema] {
			t.Errorf("tema duplicado na lista: %q", tema)
		}
		vistos[tema] = true
	}
}

func TestRandTemaRetornaTemaValido(t *testing.T) {
	validos := make(map[string]bool)
	for _, tema := range temas {
		validos[tema] = true
	}
	for range 100 {
		tema := randTema(temas)
		if !validos[tema] {
			t.Fatalf("randTema retornou tema inválido: %q", tema)
		}
	}
}

// Regressão do bug do \\n: o .md precisa ter quebras de linha reais,
// nunca o texto literal "\n".
func TestSaveMarkdownConteudo(t *testing.T) {
	dir := t.TempDir()
	wdOriginal, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wdOriginal) })

	testes := []struct {
		nome      string
		tema      string
		resposta  string
		timestamp string
	}{
		{"mensagem simples", "foco", "Mantenha o foco no que importa.", "2026-08-21-08h"},
		{"resposta multilinha", "gratidão", "primeira linha\nsegunda linha", "2026-08-21-09h"},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			if err := saveMarkdown(tt.tema, tt.resposta, tt.timestamp); err != nil {
				t.Fatalf("saveMarkdown: %v", err)
			}

			raw, err := os.ReadFile(filepath.Join("messages", tt.timestamp+".md"))
			if err != nil {
				t.Fatalf("ler arquivo gerado: %v", err)
			}
			got := string(raw)

			quer := "---\ntheme: " + tt.tema + "\ntimestamp: " + tt.timestamp + "\n---\n" + tt.resposta
			if got != quer {
				t.Errorf("conteúdo inesperado:\ngot:  %q\nwant: %q", got, quer)
			}
			if strings.Contains(got, `\n`) {
				t.Error("arquivo contém \\n literal (newline duplo-escapado)")
			}
		})
	}
}

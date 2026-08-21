package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// CommitAndPush adds the changed files, commits with a descriptive message, and pushes to remote.
// tema: o tema da mensagem escolhido
// timestamp: horário formatado usado no nome do arquivo
func CommitAndPush(tema, timestamp string) error {
	// 1. git add messages/ (ou . para incluir tudo)
	cmdAdd := exec.Command("git", "add", "messages/")
	if output, err := cmdAdd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %s: %w", string(output), err)
	}

	// 2. git commit -m "..."
	commitMsg := fmt.Sprintf("mensagem automática: %s (%s)", tema, timestamp)
	cmdCommit := exec.Command("git", "commit", "-m", commitMsg)
	if output, err := cmdCommit.CombinedOutput(); err != nil {
		// Se já estiver atualizado (nothing to commit), não é erro fatal
		if strings.Contains(string(output), "nothing to commit") {
			fmt.Println("Nenhuma mudança nova para commitar.")
			return nil
		}
		return fmt.Errorf("git commit: %s: %w", string(output), err)
	}
	fmt.Println("Commit criado:", commitMsg)

	// 3. git push
	// A autenticação depende do remote URL já estar configurado com o token (ex: https://token@github.com/...)
	cmdPush := exec.Command("git", "push")
	if output, err := cmdPush.CombinedOutput(); err != nil {
		return fmt.Errorf("git push: %s: %w", string(output), err)
	}

	fmt.Println("Push concluído.")
	return nil
}
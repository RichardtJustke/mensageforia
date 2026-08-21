package scheduler

import (
	"log/slog"
	"os"
	"time"

	"github.com/robfig/cron/v3"

	"mensageforia/internal/message"
	"mensageforia/internal/ollama"
)

func setLocation() (*time.Location, error) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func setupCron(client *ollama.Client) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	loc, err := setLocation()
	if err != nil {
		logger.Error("falha ao carregar localização", "detalhe", err)
		return err
	}

	c := cron.New(cron.WithLocation(loc))

	job := func() {
		if err := message.GenerateAndCommit(client); err != nil {
			logger.Error("falha ao gerar mensagem", "detalhe", err)
		}
	}

	_, err = c.AddFunc("00 08 * * *", job)
	if err != nil {
		logger.Error("falha ao agendar tarefa", "detalhe", err)
		return err
	}

	_, err = c.AddFunc("00 12 * * *", job)
	if err != nil {
		logger.Error("falha ao agendar tarefa", "detalhe", err)
		return err
	}

	_, err = c.AddFunc("00 18 * * *", job)
	if err != nil {
		logger.Error("falha ao agendar tarefa", "detalhe", err)
		return err
	}

	c.Start()
	select {}
}

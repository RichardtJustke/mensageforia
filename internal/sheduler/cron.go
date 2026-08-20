package sheduler

import (
	"github.com/robfig/cron/v3"
	"log/slog"
	"os"
	"time"
)

func setLocation() (*time.Location, error) {

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return nil, err
	}
	return loc, nil
}

func setupCron() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	loc, err := setLocation()
	if err != nil {
		logger.Error("falha ao carregar localização", "detalhe", err)
		return err
	}

	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("00 08 * * *", func() {
		//ativa o message as 8:00
	})
	if err != nil {
		logger.Error("falha ao agendar tarefa", "detalhe", err)
		return err
	}

	_, err := c.AddFunc("00 12 * * *", func() {
		//ativa o message as 12:00
	})
	if err != nil {
		logger.Error("falha ao agendar tarefa", "detalhe", err)
		return err
	}

	_, err := c.AddFunc("00 18 * * *", func() {
		//ativa o message as 18:00
	})
	if err != nil {
		logger.Error("falha ao agendar tarefa", "detalhe", err)
		return err
	}
	c.Start()
	defer c.Stop()
	select {}
}

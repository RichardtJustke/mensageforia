package sheduler

import (
	"github.com/robfig/cron/v3"
	"errors"
	"log/slog"
	"os"
)
func setLocation()(*time.Location, error){
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err!= nil{
		return nil, err
	}
	return loc, nil
}
func setupCron() error{	
	loc := setLocation()

	c := cron.New(cron.WithLocation(loc))

	logger:= slog.New(slog.NewTextHandler(os.Stdout, nil))
	_, err := c.AddFunc("00 08 * * *", func (){
		//ativa o message as 8:00
	})
	if err!= nil{
		logger.Error("falha ao agendar tarefa","detalhe", err)
		return
	}
	
	_, err := c.AddFunc("00 12 * * *", func (){
		//ativa o message as 12:00
	})
	if err!= nil{
		logger.Error("falha ao agendar tarefa","detalhe", err)
		return
	}

	_, err := c.AddFunc("00 18 * * *", func (){
		//ativa o message as 18:00
	})
	if err!= nil{
		logger.Error("falha ao agendar tarefa","detalhe", err)
		return
	}
	c.Start()
	defer c.Stop()
	select{}
}



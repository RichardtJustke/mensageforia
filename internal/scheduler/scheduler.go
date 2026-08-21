package scheduler

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"mensageforia/internal/message"
	"mensageforia/internal/ollama"
	"mensageforia/internal/storage"
)

const (
	tickEvery = 30 * time.Second
	drawHour  = 1 // sorteio todo dia às 01:00
	maxPerDay = 10
)

func setLocation() (*time.Location, error) {
	return time.LoadLocation("America/Sao_Paulo")
}

func parseHHMM(s string, fallbackH, fallbackM int) (int, int) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) != 2 {
		return fallbackH, fallbackM
	}
	h, errH := strconv.Atoi(parts[0])
	m, errM := strconv.Atoi(parts[1])
	if errH != nil || errM != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return fallbackH, fallbackM
	}
	return h, m
}

// spreadTimes gera n horários igualmente espaçados entre início e fim (inclusive).
// Exemplo: n=5, janela 08:00–20:00 → 08:00, 11:00, 14:00, 17:00, 20:00
func spreadTimes(n int, day time.Time, loc *time.Location, startH, startM, endH, endM int) []time.Time {
	if n <= 0 {
		return nil
	}
	y, mo, d := day.Date()
	start := time.Date(y, mo, d, startH, startM, 0, 0, loc)
	end := time.Date(y, mo, d, endH, endM, 0, 0, loc)

	if n == 1 || !end.After(start) {
		return []time.Time{start}
	}

	step := end.Sub(start) / time.Duration(n-1)
	times := make([]time.Time, 0, n)
	for i := range n {
		times = append(times, start.Add(step*time.Duration(i)))
	}
	return times
}

// drawDailyCount sorteia quantas mensagens o dia vai ter (1 a maxPerDay).
func drawDailyCount() int {
	return rand.IntN(maxPerDay) + 1
}

// SetupCron roda um loop com ticker que dispara mensagens conforme o sorteio diário.
// Às 01:00 sortea quantas mensagens o dia terá; os horários são distribuídos
// igualmente na janela configurável (padrão 08:00–20:00).
// Horários que já passaram ao iniciar o scheduler não são recuperados.
func SetupCron(client *ollama.Client, db *storage.Storage) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	loc, err := setLocation()
	if err != nil {
		logger.Error("falha ao carregar localização", "detalhe", err)
		return err
	}

	startH, startM := parseHHMM(os.Getenv("MSG_WINDOW_START"), 8, 0)
	endH, endM := parseHHMM(os.Getenv("MSG_WINDOW_END"), 20, 0)

	job := func() {
		if err := message.GenerateAndCommit(client, db); err != nil {
			logger.Error("falha ao gerar mensagem", "detalhe", err)
		}
	}

	drawDate := ""
	todayCount := 0
	prev := time.Now().In(loc)

	logger.Info("scheduler iniciado",
		"janela", fmt.Sprintf("%02d:%02d–%02d:%02d", startH, startM, endH, endM),
		"sorteio", fmt.Sprintf("01:00 (1 a %d mensagens)", maxPerDay),
		"tick", tickEvery,
	)

	ticker := time.NewTicker(tickEvery)
	defer ticker.Stop()

	for now := range ticker.C {
		now = now.In(loc)
		today := now.Format("2006-01-02")

		// Sorteio diário às 01:00 (ou no primeiro tick depois disso)
		if drawDate != today && now.Hour() >= drawHour {
			drawDate = today
			todayCount = drawDailyCount()
			logger.Info("🎲 sorteio do dia", "data", today, "mensagens", todayCount)
			for _, t := range spreadTimes(todayCount, now, loc, startH, startM, endH, endM) {
				logger.Info("horário previsto", "às", t.Format("15:04"))
			}
		}

		// Dispara mensagens nos horários sorteados
		for _, t := range spreadTimes(todayCount, now, loc, startH, startM, endH, endM) {
			if t.After(prev) && !t.After(now) {
				logger.Info("⏰ disparando mensagem agendada", "horário", t.Format("15:04"))
				job()
			}
		}

		prev = now
	}

	return nil
}

package scheduler

import (
	"testing"
	"time"
)

func TestSpreadTimesZero(t *testing.T) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	day := time.Date(2026, 8, 21, 0, 0, 0, 0, loc)
	got := spreadTimes(0, day, loc, 8, 0, 20, 0)
	if got != nil {
		t.Errorf("n=0 deve retornar nil, got %v", got)
	}
}

func TestSpreadTimesSingle(t *testing.T) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	day := time.Date(2026, 8, 21, 0, 0, 0, 0, loc)
	got := spreadTimes(1, day, loc, 8, 0, 20, 0)
	if len(got) != 1 || got[0].Hour() != 8 || got[0].Minute() != 0 {
		t.Errorf("n=1 deve retornar 08:00, got %v", got)
	}
}

func TestSpreadTimesMultiple(t *testing.T) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	day := time.Date(2026, 8, 21, 0, 0, 0, 0, loc)
	got := spreadTimes(5, day, loc, 8, 0, 20, 0)

	if len(got) != 5 {
		t.Fatalf("n=5 deve retornar 5 horários, got %d", len(got))
	}

	// Primeiro deve ser 08:00
	if got[0].Hour() != 8 || got[0].Minute() != 0 {
		t.Errorf("primeiro horário deve ser 08:00, got %s", got[0].Format("15:04"))
	}

	// Último deve ser 20:00
	last := got[len(got)-1]
	if last.Hour() != 20 || last.Minute() != 0 {
		t.Errorf("último horário deve ser 20:00, got %s", last.Format("15:04"))
	}

	// Ordem cronológica
	for i := 1; i < len(got); i++ {
		if !got[i].After(got[i-1]) {
			t.Errorf("horários devem estar em ordem: %s <= %s", got[i].Format("15:04"), got[i-1].Format("15:04"))
		}
	}
}

func TestSpreadTimesThree(t *testing.T) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	day := time.Date(2026, 8, 21, 0, 0, 0, 0, loc)
	got := spreadTimes(3, day, loc, 8, 0, 20, 0)

	if len(got) != 3 {
		t.Fatalf("n=3 deve retornar 3 horários, got %d", len(got))
	}

	// Esperado: 08:00, 14:00, 20:00
	expect := []struct{ h, m int }{{8, 0}, {14, 0}, {20, 0}}
	for i, e := range expect {
		if got[i].Hour() != e.h || got[i].Minute() != e.m {
			t.Errorf("horário[%d] deve ser %02d:%02d, got %s", i, e.h, e.m, got[i].Format("15:04"))
		}
	}
}

func TestDrawDailyCountBounds(t *testing.T) {
	for range 1000 {
		n := drawDailyCount()
		if n < 1 || n > maxPerDay {
			t.Fatalf("drawDailyCount()=%d, fora do range 1–%d", n, maxPerDay)
		}
	}
}

func TestParseHHMMValid(t *testing.T) {
	tests := []struct {
		input        string
		wantH, wantM int
	}{
		{"08:30", 8, 30},
		{"00:00", 0, 0},
		{"23:59", 23, 59},
		{"", 9, 0}, // fallback
		{"invalid", 9, 0},
		{"25:00", 9, 0}, // hora inválida
	}
	for _, tt := range tests {
		h, m := parseHHMM(tt.input, 9, 0)
		if h != tt.wantH || m != tt.wantM {
			t.Errorf("parseHHMM(%q) = %d:%d, want %d:%d", tt.input, h, m, tt.wantH, tt.wantM)
		}
	}
}

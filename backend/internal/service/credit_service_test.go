package service_test

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestMatchModelCreditRate_ExactPattern(t *testing.T) {
	rates := []*service.ModelCreditRate{
		{ModelPattern: "*", CreditsPer1kIn: 2, CreditsPer1kOut: 8, Priority: 0},
		{ModelPattern: "claude-opus-4*", CreditsPer1kIn: 15, CreditsPer1kOut: 75, Priority: 100},
	}
	got := service.MatchModelCreditRate("claude-opus-4-5", rates)
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.CreditsPer1kIn != 15 {
		t.Errorf("expected 15, got %d", got.CreditsPer1kIn)
	}
}

func TestMatchModelCreditRate_Fallback(t *testing.T) {
	rates := []*service.ModelCreditRate{
		{ModelPattern: "*", CreditsPer1kIn: 2, CreditsPer1kOut: 8, Priority: 0},
	}
	got := service.MatchModelCreditRate("unknown-model", rates)
	if got == nil {
		t.Fatal("expected non-nil result")
	}
	if got.CreditsPer1kIn != 2 {
		t.Errorf("expected 2, got %d", got.CreditsPer1kIn)
	}
}

func TestMatchModelCreditRate_NoMatch(t *testing.T) {
	rates := []*service.ModelCreditRate{
		{ModelPattern: "claude-*", CreditsPer1kIn: 5, CreditsPer1kOut: 15, Priority: 10},
	}
	got := service.MatchModelCreditRate("gpt-4o", rates)
	if got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}

func TestCalculateCreditDelta(t *testing.T) {
	rate := &service.ModelCreditRate{CreditsPer1kIn: 15, CreditsPer1kOut: 75}
	// (1000/1000)*15 + (500/1000)*75 = 15 + 37 = 52  (integer division)
	delta := service.CalculateCreditDelta(1000, 500, rate)
	if delta != 52 {
		t.Errorf("expected 52, got %d", delta)
	}
}

func TestCalculateCreditDelta_SubThousand(t *testing.T) {
	rate := &service.ModelCreditRate{CreditsPer1kIn: 15, CreditsPer1kOut: 75}
	// 999 tokens < 1000 → 0 charged (integer division floors to 0)
	delta := service.CalculateCreditDelta(999, 999, rate)
	if delta != 0 {
		t.Errorf("expected 0 for sub-1k tokens, got %d", delta)
	}
}

func TestMatchModelCreditRate_HighestPriorityWins(t *testing.T) {
	rates := []*service.ModelCreditRate{
		{ModelPattern: "*", CreditsPer1kIn: 2, CreditsPer1kOut: 8, Priority: 0},
		{ModelPattern: "claude-*", CreditsPer1kIn: 10, CreditsPer1kOut: 30, Priority: 50},
		{ModelPattern: "claude-opus-4*", CreditsPer1kIn: 15, CreditsPer1kOut: 75, Priority: 100},
	}
	got := service.MatchModelCreditRate("claude-opus-4-5", rates)
	if got.Priority != 100 {
		t.Errorf("expected priority 100, got %d", got.Priority)
	}
	if got.CreditsPer1kIn != 15 {
		t.Errorf("expected CreditsPer1kIn=15, got %d", got.CreditsPer1kIn)
	}
}

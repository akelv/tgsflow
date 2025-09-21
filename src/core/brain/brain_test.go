package brain

import (
	"context"
	"testing"

	"github.com/kelvin/tgsflow/src/core/config"
)

func TestNewTransport_SupportedModes(t *testing.T) {
	modes := []string{"shell", "proxy", "sdk", "mcp"}
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			cfg := config.Default()
			cfg.AI.Mode = mode
			tr, err := NewTransport(cfg)
			if err != nil {
				t.Fatalf("expected no error for mode %q, got %v", mode, err)
			}
			if tr == nil {
				t.Fatalf("expected transport instance for mode %q", mode)
			}
			// Even stub transports should implement Chat and return an error.
			if _, chatErr := tr.Chat(context.Background(), ChatReq{}); chatErr == nil {
				t.Fatalf("expected Chat to return error for stub transport in mode %q", mode)
			}
		})
	}
}

func TestNewTransport_UnsupportedMode(t *testing.T) {
	cfg := config.Default()
	cfg.AI.Mode = "unknown-mode"
	tr, err := NewTransport(cfg)
	if err == nil {
		t.Fatalf("expected error for unsupported mode, got nil")
	}
	if tr != nil {
		t.Fatalf("expected nil transport for unsupported mode")
	}
}

func TestBudget(t *testing.T) {
	cfg := config.Default()
	// Ensure a fresh map to control values precisely
	cfg.AI.Toolpack.Budgets = map[string]int{
		"positive": 5,
		"zero":     0,
		"negative": -3,
	}

	if got := Budget(cfg, "positive", 10); got != 5 {
		t.Fatalf("Budget positive: expected 5, got %d", got)
	}
	if got := Budget(cfg, "zero", 10); got != 10 {
		t.Fatalf("Budget zero uses default: expected 10, got %d", got)
	}
	if got := Budget(cfg, "negative", 10); got != 10 {
		t.Fatalf("Budget negative uses default: expected 10, got %d", got)
	}
	if got := Budget(cfg, "missing", 10); got != 10 {
		t.Fatalf("Budget missing key uses default: expected 10, got %d", got)
	}

	// Nil budgets map should also fall back to default
	cfgNil := config.Default()
	cfgNil.AI.Toolpack.Budgets = nil
	if got := Budget(cfgNil, "any", 42); got != 42 {
		t.Fatalf("Budget with nil map uses default: expected 42, got %d", got)
	}
}

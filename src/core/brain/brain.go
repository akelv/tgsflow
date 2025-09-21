package brain

import (
	"context"
	"fmt"

	"github.com/kelvin/tgsflow/src/core/config"
)

type ChatReq struct {
	System    string `json:"system"`
	Messages  []Msg  `json:"messages"`
	Tools     []Tool `json:"tools,omitempty"`
	MaxTokens int    `json:"max_tokens"`
}
type Msg struct {
	Role, Content string `json:"role","content"`
}
type Tool struct {
	Name, Description, JSONSchema string `json:"name","description","json_schema"`
}
type ToolCall struct {
	Name, ArgsJSON string `json:"name","args_json"`
}
type ChatResp struct {
	Text      string     `json:"text"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type Transport interface {
	Chat(ctx context.Context, req ChatReq) (ChatResp, error)
}

func NewTransport(cfg config.Config) (Transport, error) {
	switch cfg.AI.Mode {
	case "shell":
		return NewShellTransport(cfg), nil
	case "mcp":
		return NewMCPTransport(cfg) // small client; build-tag or optional dep
	case "proxy":
		return NewProxyTransport(cfg), nil
	case "sdk":
		return NewSDKTransport(cfg), nil
	default:
		return nil, fmt.Errorf("ai.mode %q unsupported", cfg.AI.Mode)
	}
}

// helper budget
func Budget(cfg config.Config, key string, def int) int {
	if v, ok := cfg.AI.Toolpack.Budgets[key]; ok && v > 0 {
		return v
	}
	return def
}

// Usage in tgs context pack:
// tr := NewTransport(cfg)
// resp, _ := tr.Chat(ctx, ChatReq{System:..., Messages:..., Tools:..., MaxTokens:Budget(cfg,"context_pack_tokens",1200)})

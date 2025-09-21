package brain

import (
	"context"
	"errors"

	"github.com/kelvin/tgsflow/src/core/config"
)

type noopTransport struct{}

func (n *noopTransport) Chat(ctx context.Context, req ChatReq) (ChatResp, error) {
	return ChatResp{}, errors.New("transport not implemented for this mode")
}

func NewShellTransport(cfg config.Config) Transport        { return &noopTransport{} }
func NewProxyTransport(cfg config.Config) Transport        { return &noopTransport{} }
func NewSDKTransport(cfg config.Config) Transport          { return &noopTransport{} }
func NewMCPTransport(cfg config.Config) (Transport, error) { return &noopTransport{}, nil }

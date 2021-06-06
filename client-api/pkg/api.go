package pkg

import (
	"context"
	"github.com/client-api/pkg/dto"
)

type ApiPortClient interface {
	Send(ctx context.Context, dto dto.PortDto) error
	GetPortByName(ctx context.Context, name string) (dto.PortDto, error)
}

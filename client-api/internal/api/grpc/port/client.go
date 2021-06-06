package port

import (
	"context"
	"errors"
	"fmt"
	grpcapi "github.com/client-api/internal/api/grpc/pb/port"
	"github.com/client-api/pkg"
	"github.com/client-api/pkg/dto"
	"google.golang.org/grpc"
	"net/http"
)

func NewClient(conn *grpc.ClientConn) pkg.ApiPortClient {
	return &t{client: grpcapi.NewPortDomainServiceClient(conn)}
}

type t struct {
	client grpcapi.PortDomainServiceClient
}

func (c t) Send(ctx context.Context, dto dto.PortDto) error {

	status, err := c.client.SendPort(
		ctx,
		&grpcapi.PortDto{
			Key:         dto.Key,
			Name:        dto.Name,
			City:        dto.City,
			Country:     dto.Country,
			Alias:       dto.Alias,
			Regions:     dto.Regions,
			Coordinates: dto.Coordinates,
			Province:    dto.Province,
			Timezone:    dto.Timezone,
			Unlocs:      dto.Unlocs,
			Code:        dto.Code,
		},
	)

	if err != nil || status.Status != http.StatusOK {
		return errors.New(fmt.Sprintf("could not send: %v", err))
	}

	return nil
}

func (c t) GetPortByName(ctx context.Context, name string) (dto.PortDto, error) {
	grpcdto, err := c.client.GetPortByName(
		ctx,
		&grpcapi.Name{
			Name: name,
		},
	)
	if err != nil {
		return dto.PortDto{}, err
	}

	return dto.PortDto{
		Key:         grpcdto.Key,
		Name:        grpcdto.Name,
		City:        grpcdto.City,
		Country:     grpcdto.Country,
		Alias:       grpcdto.Alias,
		Regions:     grpcdto.Regions,
		Coordinates: grpcdto.Coordinates,
		Province:    grpcdto.Province,
		Timezone:    grpcdto.Timezone,
		Unlocs:      grpcdto.Unlocs,
		Code:        grpcdto.Code,
	}, nil
}

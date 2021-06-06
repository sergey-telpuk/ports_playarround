package port

import (
	"context"
	"errors"
	grpc "github.com/port-domain/internal/api/grpc/pb/port"
	domain "github.com/port-domain/internal/domain/port"
	"github.com/port-domain/pkg"
	"net/http"
)

func NewApiService(repo pkg.Repository) pkg.PortApiService {
	return t{
		repo: repo,
	}
}

type t struct {
	repo pkg.Repository
}

func (s t) SendPort(ctx context.Context, grpcdto *grpc.PortDto) (*grpc.Response, error) {

	err := s.save(ctx, grpcdto)

	if err != nil {
		return &grpc.Response{
			Status:       http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &grpc.Response{
		Status:       http.StatusOK,
		ErrorMessage: "",
	}, nil
}

func (s t) GetPortByName(ctx context.Context, name *grpc.Name) (*grpc.PortDto, error) {
	grpcdto, err := s.getPortByName(ctx, name.Name)

	return &grpcdto, err
}

func (s t) save(ctx context.Context, dto *grpc.PortDto) error {
	return s.repo.Save(ctx,
		domain.Entity{
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
}

func (s t) getPortByName(ctx context.Context, name string) (grpc.PortDto, error) {
	entity, err := s.repo.GetByName(ctx, name)

	if err != nil {
		return grpc.PortDto{}, err
	}
	entityPort, ok := entity.(domain.Entity)

	if !ok {
		return grpc.PortDto{}, errors.New("cant convert to the PortDto")
	}

	return grpc.PortDto{
		Key:         entityPort.Key,
		Name:        entityPort.Name,
		City:        entityPort.City,
		Country:     entityPort.Country,
		Alias:       entityPort.Alias,
		Regions:     entityPort.Regions,
		Coordinates: entityPort.Coordinates,
		Province:    entityPort.Province,
		Timezone:    entityPort.Timezone,
		Unlocs:      entityPort.Unlocs,
		Code:        entityPort.Code,
	}, nil
}

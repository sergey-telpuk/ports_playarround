package pkg

import "github.com/client-api/pkg/dto"

type PortService interface {
	ReadFromJson() error
	ReadFromJsonWorkerPool(jobs int) error
	GetPortByName(name string) (dto.PortDto, error)
	Close()
}

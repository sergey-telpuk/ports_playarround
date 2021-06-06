package pkg

import "context"

type Repository interface {
	Save(cont context.Context, entity interface{}) error
	GetByName(cont context.Context, name string) (interface{}, error)
}

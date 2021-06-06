package port

import (
	"context"
	"errors"
	"github.com/port-domain/pkg"
	"sync"
)

func NewRepository() pkg.Repository {
	return t{
		storage: make(map[string]interface{}),
		mtx:     &sync.Mutex{},
		rwmtx:   &sync.RWMutex{},
	}
}

type t struct {
	storage map[string]interface{}
	mtx     *sync.Mutex
	rwmtx   *sync.RWMutex
}

func (r t) Save(ctx context.Context, entity interface{}) error {
	e, ok := entity.(Entity)
	if !ok {
		return errors.New("cant covert to Entity")
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.storage[e.Name] = entity

	return nil
}

func (r t) GetByName(ctx context.Context, name string) (interface{}, error) {
	r.rwmtx.Lock()
	defer r.rwmtx.Unlock()
	e, ok := r.storage[name]

	if !ok {
		return nil, errors.New("not found")
	}

	return e, nil
}

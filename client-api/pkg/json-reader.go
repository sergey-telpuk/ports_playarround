package pkg

import "context"

type JsonReader interface {
	Read(ctx context.Context) (out <-chan interface{}, done <-chan bool, err <-chan error)
}

package reader

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/client-api/pkg"
	"github.com/client-api/pkg/dto"
	"io"
)

func NewReader(reader io.Reader) pkg.JsonReader {
	return ports{
		reader: reader,
	}
}

type ports struct {
	reader io.Reader
}

func (p ports) Read(ctx context.Context) (<-chan interface{}, <-chan bool, <-chan error) {
	outCh := make(chan interface{})
	errCh := make(chan error)
	quitCh := make(chan bool)

	go func() {
		//defer close(outCh)
		//defer close(errCh)
		r := bufio.NewReader(p.reader)
		dec := json.NewDecoder(r)

		// We expect an object
		t, err := dec.Token()
		if delim, ok := t.(json.Delim); !ok || delim != '{' {
			errCh <- errors.New(fmt.Sprint("expected object"))
			return
		}

		// Read props
		for dec.More() {
			//read for key
			key, err := dec.Token()
			if err != nil {
				errCh <- err
				return
			}
			dtoPort := dto.PortDto{
				Key: key.(string),
			}

			dec.More()
			// read next item (large object)
			if err := dec.Decode(&dtoPort); err != nil {
				errCh <- err
				return
			}

			select {
			case outCh <- dtoPort:
			case <-ctx.Done():
				return
			}
		}

		// object closing delim
		t, err = dec.Token()
		if delim, ok := t.(json.Delim); !ok || delim != '}' {
			errCh <- err
		}

		quitCh <- true
	}()

	return outCh, quitCh, errCh
}

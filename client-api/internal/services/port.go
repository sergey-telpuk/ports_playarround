package services

import (
	"context"
	"errors"
	"github.com/client-api/pkg"
	"github.com/client-api/pkg/dto"
	"sync"
	"time"
)

type t struct {
	jsonReader pkg.JsonReader
	apiClient  pkg.ApiPortClient
	closeCh    chan bool
}

func NewPortService(jsonReader pkg.JsonReader, apiClient pkg.ApiPortClient) pkg.PortService {
	return t{
		jsonReader: jsonReader,
		apiClient:  apiClient,
		closeCh:    make(chan bool),
	}
}

func (s t) ReadFromJson() error {
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	outCh, quitCh, errCh := s.jsonReader.Read(cxt)

	for {
		select {
		case port := <-outCh:
			p, ok := port.(dto.PortDto)
			if !ok {
				return errors.New("cant convert to the dto")
			}
			if err := s.sendPort(p); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		case <-quitCh:
			return nil
		case <-cxt.Done():
			return cxt.Err()
		}
	}
}

// timeOutAfterQuit - need to use for waiting for all finished grpc calls
func (s t) ReadFromJsonWorkerPool(jobs int) error {
	taskCh := make(chan dto.PortDto, jobs)
	taskErrCh := make(chan error)
	errCh := make(chan error)
	done := make(chan bool)

	cxtP, cancelP := context.WithCancel(context.Background())
	defer cancelP()

	workerWg := &sync.WaitGroup{}
	go func() {
		// graceful shutdown of worker goroutines
		workerWg.Wait()
		close(done)
	}()

	for i := 0; i < jobs; i++ {
		workerWg.Add(1)
		go s.worker(cxtP, workerWg, taskCh, taskErrCh)
	}

	cxtRead, cancelRead := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelRead()

	outReadCh, quitReadCh, errReadCh := s.jsonReader.Read(cxtRead)

	go func() {
		defer close(taskCh) //finish workers
		for {
			select {
			case port := <-outReadCh:
				p, ok := port.(dto.PortDto)
				if !ok {
					errCh <- errors.New("cant convert to PortDto")
					return
				}
				//dont block if was error in worker
				if err := s.sendToTask(p, taskCh); err != nil {
					errCh <- err
					return
				}
			case err := <-errReadCh:
				if err != nil {
					errCh <- err
				}
				return
			case <-quitReadCh:
				return
			case <-cxtRead.Done():
				if cxtRead.Err() != nil {
					errCh <- cxtRead.Err()
				}
				return
			case <-s.closeCh:
				return
			}
		}
	}()

	for {
		select {
		case err := <-taskErrCh:
			return err
		case <-done:
			return nil
		}
	}

}

func (s t) GetPortByName(name string) (dto.PortDto, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	return s.apiClient.GetPortByName(cxt, name)
}

func (s t) Close() {
	s.closeCh <- true
}

func (s t) worker(ctx context.Context, wg *sync.WaitGroup, tasksCh <-chan dto.PortDto, taskErrCh chan<- error) {
	defer wg.Done()

	for {
		select {
		case port, ok := <-tasksCh:
			if !ok {
				return
			}
			if err := s.sendPort(port); err != nil {
				taskErrCh <- err
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s t) sendToTask(port dto.PortDto, taskCh chan<- dto.PortDto) error {
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errCh := make(chan error)
	done := make(chan bool)

	go func() {
		for {
			select {
			case taskCh <- port:
				done <- true
				return
			case <-cxt.Done():
				errCh <- cxt.Err()
				return
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-done:
		return nil
	}
}

func (s t) sendPort(port dto.PortDto) error {
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	return s.apiClient.Send(cxt, port)
}

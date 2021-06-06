package service

import (
	grpcapi "github.com/port-domain/internal/api/grpc/pb/port"
	"github.com/port-domain/pkg"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewPortDomainService(api pkg.PortApiService) pkg.PortDomainService {
	return t{
		api:     api,
		closeCh: make(chan bool),
	}
}

type t struct {
	api     pkg.PortApiService
	closeCh chan bool
}

func (s t) Run(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcapi.RegisterPortDomainServiceServer(grpcServer, s.api.(grpcapi.PortDomainServiceServer))

	go func() {
		<-s.closeCh
		grpcServer.Stop()
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s t) Close() {
	s.closeCh <- true
}

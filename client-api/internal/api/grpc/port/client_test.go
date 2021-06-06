package port_test

import (
	"context"
	apigrpc "github.com/client-api/internal/api/grpc/pb/port"
	"github.com/client-api/internal/api/grpc/port"
	"github.com/client-api/pkg/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"net/http"
	"testing"
)

type mockPortDomainServiceServer struct {
	apigrpc.UnimplementedPortDomainServiceServer
}

func (*mockPortDomainServiceServer) SendPort(ctx context.Context, dto *apigrpc.PortDto) (*apigrpc.Response, error) {
	if dto.Name != "test" {
		return &apigrpc.Response{Status: http.StatusBadRequest, ErrorMessage: "error"}, nil
	}
	return &apigrpc.Response{Status: http.StatusOK, ErrorMessage: ""}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	apigrpc.RegisterPortDomainServiceServer(server, &mockPortDomainServiceServer{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestSend(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := port.NewClient(conn)

	if client.Send(context.Background(), dto.PortDto{
		Key:         "test",
		Name:        "test",
		City:        "test",
		Country:     "test",
		Alias:       nil,
		Regions:     nil,
		Coordinates: nil,
		Province:    "",
		Timezone:    "",
		Unlocs:      nil,
		Code:        "",
	}) != nil {
		t.Error(err)
	}
}

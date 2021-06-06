package services_test

import (
	"bytes"
	"context"
	apigrpc "github.com/client-api/internal/api/grpc/pb/port"
	api "github.com/client-api/internal/api/grpc/port"
	"github.com/client-api/internal/file/reader"
	"github.com/client-api/internal/services"
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

func TestReadFromJson(t *testing.T) {

	var buffer bytes.Buffer

	buffer.WriteString(`
{
  "AEAJM": {
    "name": "test",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "test",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}
`,
	)

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	srv := services.NewPortService(
		reader.NewReader(&buffer),
		api.NewClient(conn),
	)

	if err := srv.ReadFromJson(); err != nil {
		t.Fatal(err)
	}

}

func TestReadFromJsonWorkerPool(t *testing.T) {
	var buffer bytes.Buffer

	buffer.WriteString(`
{
  "AEAJM": {
    "name": "test",
    "city": "Ajman",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "coordinates": [
      55.5136433,
      25.4052165
    ],
    "province": "Ajman",
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAJM"
    ],
    "code": "52000"
  },
  "AEAUH": {
    "name": "test",
    "coordinates": [
      54.37,
      24.47
    ],
    "city": "Abu Dhabi",
    "province": "Abu Z¸aby [Abu Dhabi]",
    "country": "United Arab Emirates",
    "alias": [],
    "regions": [],
    "timezone": "Asia/Dubai",
    "unlocs": [
      "AEAUH"
    ],
    "code": "52001"
  }
}
`,
	)

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	srv := services.NewPortService(
		reader.NewReader(&buffer),
		api.NewClient(conn),
	)

	if err := srv.ReadFromJsonWorkerPool(5); err != nil {
		t.Fatal(err)
	}
}

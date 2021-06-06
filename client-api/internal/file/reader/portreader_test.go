package reader_test

import (
	"bytes"
	"context"
	portreader "github.com/client-api/internal/file/reader"
	"github.com/client-api/pkg/dto"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	var buffer bytes.Buffer

	buffer.WriteString(`
{
  "AEAJM": {
    "name": "Ajman",
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
    "name": "Abu Dhabi",
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outCh, quitCh, errCh := portreader.NewReader(&buffer).Read(ctx)

	for {
		select {
		case err := <-errCh:
			t.Error(err)
			return
		case out := <-outCh:
			p, ok := out.(dto.PortDto)
			if !ok {
				t.Error("cant read the dto port")
			}
			// check the first record
			if p.Key == "AEAJM" && p.Name == "Ajman" {
				return
			}
		case <-quitCh:
			return
		case <-ctx.Done():
			if ctx.Err() != nil {
				t.Error(ctx.Err())
				return
			}
		}
	}
}

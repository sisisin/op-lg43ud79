package lg43client

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

var (
	ErrorDeviceNotFound = errors.New("device not found")
)

type ErrorCommandRejected struct {
	rawAck string
}

func (e ErrorCommandRejected) Error() string {
	return fmt.Sprintf("command rejected. response: %s", e.rawAck)
}

type ErrorUnknownResponse struct {
	rawAck string
}

func (e ErrorUnknownResponse) Error() string {
	return fmt.Sprintf("unknown response: %s", e.rawAck)
}

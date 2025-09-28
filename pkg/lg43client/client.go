package lg43client

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type LG43Client struct {
	port serial.Port
}

type vidType string

func VID(v string) vidType { return vidType(v) }

type pidType string

func PID(v string) pidType { return pidType(v) }

func NewLG43Client(ctx context.Context, vid vidType, pid pidType) (*LG43Client, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list serial ports")
	}
	var portName string
	for _, p := range ports {
		if p.PID == string(pid) && p.VID == string(vid) {
			portName = p.Name
			break
		}
	}
	if portName == "" {
		return nil, errors.Errorf("device not found")
	}

	mode := &serial.Mode{
		DataBits: 8,
		StopBits: serial.OneStopBit,
		Parity:   serial.NoParity,
		BaudRate: 9600,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open serial port")
	}
	// Read method will block until at least one byte is received, or the read timeout is reached.
	port.SetReadTimeout(1 * time.Second)
	return &LG43Client{
		port: port,
	}, nil
}

func (c *LG43Client) Close() {
	if err := c.port.Close(); err != nil {
		log.Printf("failed to close port: %v", err)
	}
}

func (c *LG43Client) PowerOn(ctx context.Context, setId string) (response string, err error) {
	w := fmt.Sprintf("ka %s 01\r", setId)
	return c.Write(ctx, []byte(w))
}

func (c *LG43Client) InputSelectToDP1(ctx context.Context, setId string) (response string, err error) {
	w := fmt.Sprintf("xb %s C0\r", setId)
	return c.Write(ctx, []byte(w))
}

func (c *LG43Client) InputSelectToHDMI4(ctx context.Context, setId string) (response string, err error) {
	w := fmt.Sprintf("xb %s 93\r", setId)
	return c.Write(ctx, []byte(w))
}

func (c *LG43Client) Write(ctx context.Context, buf []byte) (response string, err error) {
	writeTimeout := time.After(10 * time.Second)

	logDebug(ctx, "Write: %s", strings.TrimSpace(string(buf)))
	if _, err = c.port.Write(buf); err != nil {
		return "", errors.Wrap(err, "failed to write to serial port")
	}

	var res strings.Builder
	for {
		select {
		case <-ctx.Done():
			logInfo(ctx, "Cancelled, exiting read loop")
			return "", ctx.Err()
		case <-writeTimeout:
			return "", errors.New("Write timeout waiting for response")
		default:
		}

		buff := make([]byte, 1000)
		if n, err := c.port.Read(buff); err != nil {
			return "", errors.Wrap(err, "failed to read from serial port")
		} else if n == 0 {
			logDebug(ctx, "Read: 0 bytes. Maybe timeout, retrying...")
			continue
		} else if buff[n-1] == '\r' {
			logDebug(ctx, "Read: %d bytes, end of response detected", n)
			break
		} else {
			if isDebug(ctx) {
				data := strings.TrimSpace(string(buff[:n]))
				logDebug(ctx, "Read: %d bytes. data: %s, buffer: %v", n, data, buff[:n])
			}
			res.Write(buff[:n])
		}
	}

	ack := res.String()

	logDebug(ctx, "Read: %s", strings.TrimSpace(ack))

	if len(ack) != 10 {
		return ack, &ErrorUnknownResponse{rawAck: ack}
	}

	status := ack[5:7]
	switch status {
	case "OK":
		return ack, nil
	case "NG":
		return ack, &ErrorCommandRejected{rawAck: ack}
	default:
		return ack, &ErrorUnknownResponse{rawAck: ack}
	}
}

package duke

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Socket struct {
	ID	   string
	io	   *Stream
	header http.Header
	status uint16
}

// New hijacks the http request and returns Socket
func NewSocket(w http.ResponseWriter, req *http.Request) (*Socket, error) {
	hj, ok := w.(http.Hijacker)

	if !ok {
		return nil, errors.New("webserver doesn't support http hijacking")
	}

	conn, io, err := hj.Hijack()

	if err != nil {
		return nil, err
	}

	return &Socket{
		uuid.NewString(),
		NewStream(conn, io),
		req.Header,
		1000,
	}, nil
}

// Close sends close frame and closes the TCP connection
func (self *Socket) Close() error {
	f := NewCloseFrame(self.status)

	if err := f.Write(); err != nil {
		return err
	}

	return self.io.Close()
}

// Handshake performs the initial websocket handshake
func (self *Socket) handshake() error {
	lines := []string{
		"HTTP/1.1 101 Web Socket Protocol Handshake",
		"Server: go/echoserver",
		"Upgrade: WebSocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + self.getAcceptHash(),
		"", // required for extra CRLF
		"", // required for extra CRLF
	}

	return self.io.Write([]byte(strings.Join(lines, "\r\n")))
}

func (self *Socket) listen(fn func(frame *Frame)) error {
	for {
		frame, err := ReadNewFrame(self.io)

		if err != nil {
			return err
		}

		status, err := frame.Validate()

		if status != self.status {
			self.status = status
		}

		if err != nil || frame.IsClose() {
			return err
		}

		if frame.IsPing() {
			frame.Pong()
		}

		fn(frame)
	}
}

// Send sends a Frame
func (self *Socket) send(fr *Frame) error {
	return self.io.Write(fr.Buffer())
}

func (self *Socket) Emit(payload string) error {

}

func (self *Socket) getAcceptHash() string {
	h := sha1.New()
	h.Write([]byte(self.header.Get("Sec-WebSocket-Key")))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

package duke

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"strings"
)

const bufferSize = 4096

type Socket struct {
	conn   net.Conn
	io     *bufio.ReadWriter
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

	return &Socket{conn, io, req.Header, 1000}, nil
}

// Handshake performs the initial websocket handshake
func (self *Socket) Handshake() error {
	lines := []string{
		"HTTP/1.1 101 Web Socket Protocol Handshake",
		"Server: go/echoserver",
		"Upgrade: WebSocket",
		"Connection: Upgrade",
		"Sec-WebSocket-Accept: " + self.getAcceptHash(),
		"", // required for extra CRLF
		"", // required for extra CRLF
	}

	return self.write([]byte(strings.Join(lines, "\r\n")))
}

func (self *Socket) getAcceptHash() string {
	h := sha1.New()
	h.Write([]byte(self.header.Get("Sec-WebSocket-Key")))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (self *Socket) write(data []byte) error {
	if _, err := self.io.Write(data); err != nil {
		return err
	}

	return self.io.Flush()
}

func (self *Socket) read(size int) ([]byte, error) {
	data := make([]byte, 0)

	for {
		if len(data) == size {
			break
		}

		// Temporary slice to read chunk
		sz := bufferSize
		remaining := size - len(data)

		if sz > remaining {
			sz = remaining
		}

		temp := make([]byte, sz)
		n, err := self.io.Read(temp)

		if err != nil && err != io.EOF {
			return data, err
		}

		data = append(data, temp[:n]...)
	}

	return data, nil
}

// Recv receives data and returns a Frame
func (self *Socket) Recv() (Frame, error) {
	frame := Frame{}
	head, err := self.read(2)

	if err != nil {
		return frame, err
	}

	frame.IsFragment = (head[0] & 0x80) == 0x00
	frame.Opcode = head[0] & 0x0F
	frame.Reserved = (head[0] & 0x70)
	frame.IsMasked = (head[1] & 0x80) == 0x80
	length := uint64(head[1] & 0x7F)

	if length == 126 {
		data, err := self.read(2)

		if err != nil {
			return frame, err
		}

		length = uint64(binary.BigEndian.Uint16(data))
	} else if length == 127 {
		data, err := self.read(8)

		if err != nil {
			return frame, err
		}

		length = uint64(binary.BigEndian.Uint64(data))
	}

	mask, err := self.read(4)

	if err != nil {
		return frame, err
	}

	frame.Length = length
	payload, err := self.read(int(length)) // possible data loss

	if err != nil {
		return frame, err
	}

	for i := uint64(0); i < length; i++ {
		payload[i] ^= mask[i%4]
	}

	frame.Payload = payload
	status, err := frame.Validate(self.status)

	if status != self.status {
		self.status = status
	}

	return frame, err
}

// Send sends a Frame
func (self *Socket) Send(fr Frame) error {
	data := make([]byte, 2)
	data[0] = 0x80 | fr.Opcode

	if fr.IsFragment {
		data[0] &= 0x7F
	}

	if fr.Length <= 125 {
		data[1] = byte(fr.Length)
		data = append(data, fr.Payload...)
	} else if fr.Length > 125 && float64(fr.Length) < math.Pow(2, 16) {
		data[1] = byte(126)
		size := make([]byte, 2)
		binary.BigEndian.PutUint16(size, uint16(fr.Length))
		data = append(data, size...)
		data = append(data, fr.Payload...)
	} else if float64(fr.Length) >= math.Pow(2, 16) {
		data[1] = byte(127)
		size := make([]byte, 8)
		binary.BigEndian.PutUint64(size, fr.Length)
		data = append(data, size...)
		data = append(data, fr.Payload...)
	}

	return self.write(data)
}

// Close sends close frame and closes the TCP connection
func (self *Socket) Close() error {
	f := NewCloseFrame(self.status)

	if err := self.Send(*f); err != nil {
		return err
	}

	return self.conn.Close()
}

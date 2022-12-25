package duke

import (
	"bufio"
	"io"
	"net"
)

type Stream struct {
	conn net.Conn
	io	 *bufio.ReadWriter
}

func NewStream(conn net.Conn, rw *bufio.ReadWriter) *Stream {
	v := Stream{conn, rw}
	return &v
}

func (self *Stream) Write(data []byte) error {
	if _, err := self.io.Write(data); err != nil {
		return err
	}

	return self.io.Flush()
}

func (self *Stream) Read(size int) ([]byte, error) {
	data := make([]byte, 0)

	for {
		if len(data) == size {
			break
		}

		// Temporary slice to read chunk
		sz := 4096
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

func (self *Stream) Close() error {
	return self.conn.Close()
}

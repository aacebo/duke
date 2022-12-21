package duke

import (
	"log"
	"net/http"
)

type Server struct {
	cors func(origin string) bool
}

func Listen(r *http.ServeMux, cors func(origin string) bool) *Server {
	self := Server{cors}

	if r != nil {
		r.HandleFunc("/ws", self.upgrade)
	} else {
		http.HandleFunc("/ws", self.upgrade)
	}

	return &self
}

func (self *Server) upgrade(w http.ResponseWriter, r *http.Request) {
	socket, err := NewSocket(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	err = socket.Handshake()
	if err != nil {
		log.Println(err)
		return
	}

	defer socket.Close()

	for {
		frame, err := socket.Recv()

		if err != nil {
			log.Println("Error Decoding", err)
			return
		}

		switch frame.Opcode {
		case 8: // Close
			return
		case 9: // Ping
			frame.Opcode = 10
			fallthrough
		case 0: // Continuation
			fallthrough
		case 1: // Text
			fallthrough
		case 2: // Binary
			if err = socket.Send(frame); err != nil {
				log.Println("Error sending", err)
				return
			}
		}
	}
}

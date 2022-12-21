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
		r.HandleFunc("/ws", self.onHandshake)
	} else {
		http.HandleFunc("/ws", self.onHandshake)
	}

	return &self
}

func (self *Server) onHandshake(w http.ResponseWriter, r *http.Request) {
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

		if frame.IsClose() {
			return
		}

		if frame.IsPing() {
			frame.Pong()
		}

		if err = socket.Send(frame); err != nil {
			log.Println("Error sending", err)
			return
		}
	}
}

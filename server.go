package duke

import (
	"log"
	"net/http"
)

type Server struct {
	listeners map[string]func(socket *Socket)
}

func New(r *http.ServeMux) *Server {
	self := Server{
		listeners: make(map[string]func(socket *Socket)),
	}

	if r != nil {
		r.HandleFunc("/ws", self.onHandshake)
	} else {
		http.HandleFunc("/ws", self.onHandshake)
	}

	return &self
}

func (self *Server) On(event string, fn func(socket *Socket)) {
	self.listeners[event] = fn
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

	if self.listeners["connection"] != nil {
		self.listeners["connection"](socket)
	}

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

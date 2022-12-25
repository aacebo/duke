package duke

import (
	"net/http"
)

type Server struct {
	listeners map[string]func(socket *Socket)
	sockets   map[string]*Socket
}

func NewServer(r *http.ServeMux) *Server {
	self := Server{
		listeners: make(map[string]func(socket *Socket)),
		sockets: make(map[string]*Socket),
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

func (self *Server) GetSocket(id string) *Socket {
	return self.sockets[id]
}

func (self *Server) onHandshake(w http.ResponseWriter, r *http.Request) {
	socket, err := NewSocket(w, r)

	if err != nil {
		return
	}

	err = socket.handshake()

	if err != nil {
		return
	}

	defer self.onDisconnect(socket)
	self.onConnect(socket)
	socket.listen(func(frame *Frame) {
		if err = socket.send(frame); err != nil {
			return
		}
	})
}

func (self *Server) onConnect(socket *Socket) {
	self.sockets[socket.ID] = socket

	if self.listeners["connect"] != nil {
		self.listeners["connect"](socket)
	}
}

func (self *Server) onDisconnect(socket *Socket) {
	socket.Close()
	delete(self.sockets, socket.ID)

	if self.listeners["disconnect"] != nil {
		self.listeners["disconnect"](socket)
	}
}

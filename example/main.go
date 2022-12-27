package main

import (
	"duke"
	"fmt"
	"net/http"
)

func main() {
	ws := duke.NewServer(nil)
	ws.On("connect", func(socket *duke.Socket) {
		fmt.Println(socket.ID + " connected")

		socket.Listen(func(f *duke.Frame) {
			socket.Emit(f.String())
		})
	})

	ws.On("disconnect", func(socket *duke.Socket) {
		fmt.Println(socket.ID + " disconnected")
	})

	http.ListenAndServe(":3000", nil)
}

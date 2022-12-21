package main

import (
	"duke"
	"fmt"
	"net/http"
)

func main() {
	ws := duke.New(nil)
	ws.On("connection", func(socket *duke.Socket) {
		fmt.Println(socket.ID)
	})

	http.ListenAndServe(":3000", nil)
}

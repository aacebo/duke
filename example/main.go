package main

import (
	"duke"
	"net/http"
)

func main() {
	duke.Listen(nil, nil)
	http.ListenAndServe(":3000", nil)
}

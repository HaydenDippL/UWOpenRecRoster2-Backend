package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", hello_world)
}

func hello_world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

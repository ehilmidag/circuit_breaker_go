package main

import (
	"fmt"
	"net/http"
)

func example(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "example")
}

func main() {
	http.HandleFunc("/", example)
	http.ListenAndServe(":8080", nil)
}

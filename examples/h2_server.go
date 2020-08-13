package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, h2!")
	})
	server := &http.Server{
		Addr:    "localhost:9999",
		Handler: http.DefaultServeMux,
	}

	http2.ConfigureServer(server, nil)
	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		fmt.Printf("err:%+v\n", err)
		return
	}
}

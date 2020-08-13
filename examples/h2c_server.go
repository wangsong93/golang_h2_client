package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, h2c!")
	})
	server := &http.Server{
		Addr:    "localhost:9999",
		Handler: h2c.NewHandler(http.DefaultServeMux, &http2.Server{}),
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("err:%+v\n", err)
		return
	}
}

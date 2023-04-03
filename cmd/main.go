package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	c int
}

func (w *ResponseWriterWrapper) WriteHeader(c int) {
	w.c = c
	w.ResponseWriter.WriteHeader(c)
}

func main() {
	ctx := context.Background()
	fs := http.FileServer(http.Dir("./"))
	l, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ngrok ingress url: ", l.URL())
	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] --> %s %s \n", r.RemoteAddr, r.Method, r.URL)

		w.Header().Add("x-ngrok-file-server", "trendev")
		ww := &ResponseWriterWrapper{w, http.StatusOK}

		fs.ServeHTTP(ww, r)

		log.Printf("[%s] <-- %d %s", r.RemoteAddr, ww.c, http.StatusText(ww.c))
	}))
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	c int
}

func (w *ResponseWriterWrapper) WriteHeader(c int) {
	w.c = c
	w.ResponseWriter.WriteHeader(c)
}

func methodColor(m string) string {
	switch m {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

func logRequest(r http.Request) {
	log.Printf("[%s] --> %s%s%s %s", r.RemoteAddr, methodColor(r.Method), r.Method, reset, r.URL)
}

func logResponse(w ResponseWriterWrapper, r *http.Request) {
	log.Printf("[%s] <-- %d %s", r.RemoteAddr, w.c, http.StatusText(w.c))
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
		logRequest(*r)
		w.Header().Add("x-ngrok-file-server", "trendev")
		ww := &ResponseWriterWrapper{w, http.StatusOK}
		fs.ServeHTTP(ww, r)
		logResponse(*ww, r)
	}))
}

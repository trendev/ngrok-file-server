package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/trendev/ngrok-file-server/pkg/colorlog"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

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
		colorlog.LogRequest(*r)
		w.Header().Add("x-ngrok-file-server", "trendev")
		ww := colorlog.NewResponseWriterWrapper(w, http.StatusOK)
		fs.ServeHTTP(ww, r)
		colorlog.LogResponse(*ww, r)
	}))
}

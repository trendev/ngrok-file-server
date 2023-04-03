package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
		w.Header().Add("x-ngrok-file-server", "trendev")
		fs.ServeHTTP(w, r)
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	}))
}

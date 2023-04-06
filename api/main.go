package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/trendev/ngrok-file-server/pkg/colorlog"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func setConfigHTTPEndpoint() config.Tunnel {
	p := flag.String("provider", "", "oauth2 provider")
	d := flag.String("domain", "", "oauth2 authorized domain")
	flag.Parse()
	if *p != "" && *d == "" {
		return config.HTTPEndpoint(config.WithOAuth(*p),
			config.WithRequestHeader("email", "${.oauth.user.email}"))
	}
	if *p != "" && *d != "" {
		return config.HTTPEndpoint(config.WithOAuth(*p, config.WithAllowOAuthDomain(*d)),
			config.WithRequestHeader("email", "${.oauth.user.email}"))
	}
	return config.HTTPEndpoint()
}

func main() {
	ctx := context.Background()
	fs := http.FileServer(http.Dir("./shared"))

	l, err := ngrok.Listen(ctx,
		setConfigHTTPEndpoint(),
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

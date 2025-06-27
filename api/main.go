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
	sd := flag.String("static_domain", "", "ngrok static domain")
	p := flag.String("provider", "", "oauth2 provider")
	o2d := flag.String("oauth2_domain", "", "oauth2 authorized oauth2_domain")
	flag.Parse()

	var opts []config.HTTPEndpointOption

	if *sd != "" {
		opts = append(opts, config.WithDomain(*sd))
	}
	if *p != "" {
		if *o2d != "" {
			opts = append(opts, config.WithOAuth(*p, config.WithAllowOAuthDomain(*o2d)),
				config.WithRequestHeader("email", "${.oauth.user.email}"))
		} else {
			opts = append(opts, config.WithOAuth(*p), config.WithRequestHeader("email", "${.oauth.user.email}"))
		}
	}
	return config.HTTPEndpoint(opts...)
}

func main() {
	fs := http.FileServer(http.Dir("./shared"))

	l, err := ngrok.Listen(context.Background(),
		setConfigHTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ngrok ingress url: ", l.URL())
	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		colorlog.LogRequest(*r)
		w.Header().Add("x-ngrok-file-server", "trendev.fr")
		ww := colorlog.NewResponseWriterWrapper(w)
		fs.ServeHTTP(ww, r)
		colorlog.LogResponse(*ww, r)
	}))
}

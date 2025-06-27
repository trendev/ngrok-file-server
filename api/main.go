package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/trendev/ngrok-file-server/pkg/colorlog"
	"golang.ngrok.com/ngrok/v2"
)

const config = `
{
  "on_http_response": [
    {
      "actions": [
        {
          "type": "add-headers",
          "config": {
            "headers": {
              "x-ngrok-file-server": "trendev.fr",
              "x-endpoint-id": "${endpoint.id}",
              "x-client-ip": "${conn.client_ip}",
              "x-client-conn-start": "${conn.ts.start}",
              "x-client-loc": "${conn.geo.country}",
              "x-client-path": "${req.url.path}"
            }
          }
        }
      ]
    }
  ],
  "on_http_request": [
    {
      "actions": [
        {
          "type": "oauth",
          "config": {
            "provider": "google"
          }
        }
      ]
    }
  ]
}`

func setConfigHTTPEndpoint() []ngrok.EndpointOption {
	// sd := flag.String("static_domain", "", "ngrok static domain")
	// p := flag.String("provider", "", "oauth2 provider")
	// o2d := flag.String("oauth2_domain", "", "oauth2 authorized oauth2_domain")
	// flag.Parse()
	var opts []ngrok.EndpointOption
	var c map[string]interface{}

	err := json.Unmarshal([]byte(config), &c)
	if err != nil {
		log.Fatal(err)
	}

	// if *sd != "" {
	// 	opts = append(opts, config.WithDomain(*sd))
	// }
	// if *p != "" {
	// 	if *o2d != "" {
	// 		opts = append(opts, config.WithOAuth(*p, config.WithAllowOAuthDomain(*o2d)),
	// 			config.WithRequestHeader("email", "${.oauth.user.email}"))
	// 	} else {
	// 		opts = append(opts, config.WithOAuth(*p), config.WithRequestHeader("email", "${.oauth.user.email}"))
	// 	}
	// }
	tp, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	opts = append(opts, ngrok.WithTrafficPolicy(string(tp)))
	return opts
}

func main() {
	fs := http.FileServer(http.Dir("./shared"))

	l, err := ngrok.Listen(context.Background(), setConfigHTTPEndpoint()...)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ngrok ingress url: ", l.URL())
	http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		colorlog.LogRequest(*r)
		// w.Header().Add("x-ngrok-file-server", "trendev.fr")
		ww := colorlog.NewResponseWriterWrapper(w)
		fs.ServeHTTP(ww, r)
		colorlog.LogResponse(*ww, r)
	}))
}

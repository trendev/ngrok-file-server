package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Jeffail/gabs/v2"
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
  ]
}`

func setConfigHTTPEndpoint() []ngrok.EndpointOption {
	var opts []ngrok.EndpointOption
	p := flag.String("provider", "", "oauth2 provider")
	o2d := flag.String("oauth2_domain", "", "oauth2 authorized oauth2_domain")
	flag.Parse()

	c, err := gabs.ParseJSON([]byte(config))
	if err != nil {
		log.Fatal(err)
	}

	if *p != "" {
		// Create OAuth action structure
		a := gabs.New()
		a.Set("oauth", "type")
		a.Set(*p, "config", "provider")

		// Create request entry
		r := gabs.New()
		r.ArrayAppend(a.Data(), "actions")
		if *o2d != "" {
			// e := fmt.Sprintf("actions.ngrok.oauth.identity.email.endsWith('%s')", *o2d)
			// r.Set(e, "expressions")
		}

		// Append to on_http_request array
		c.ArrayAppend(r.Data(), "on_http_request")

	}

	opts = append(opts, ngrok.WithTrafficPolicy(c.String()))
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

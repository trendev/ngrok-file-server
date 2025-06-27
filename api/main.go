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
	p := flag.String("provider", "", "oauth2 provider (google, github, etc.)")
	o2d := flag.String("oauth2_domain", "", "oauth2 authorized oauth2_domain")
	flag.Parse()

	c, err := gabs.ParseJSON([]byte(config))
	if err != nil {
		log.Fatal(err)
	}

	if *p != "" {
		// Create OAuth action structure
		a1 := gabs.New()
		a1.Set("oauth", "type")
		a1.Set(*p, "config", "provider")

		// Create request entry
		r := gabs.New()
		r.ArrayAppend(a1.Data(), "actions")
		// Append to on_http_request array
		c.ArrayAppend(r.Data(), "on_http_request")

		if *o2d != "" {
			dr := gabs.New() // domain rule

			da := gabs.New() // deny action
			da.Set("deny", "type")
			da.Set(401, "config", "status_code")

			expression := fmt.Sprintf(
				"!actions.ngrok.oauth.identity.email.endsWith('%s')",
				*o2d,
			)

			dr.ArrayAppend(expression, "expressions")
			dr.ArrayAppend(da.Data(), "actions")

			// Ajouter la deuxième règle
			c.ArrayAppend(dr.Data(), "on_http_request")
		}

		log.Println(c.StringIndent("", "  "))
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

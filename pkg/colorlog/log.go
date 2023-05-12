package colorlog

import (
	"log"
	"net/http"
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

type responseWriterWrapper struct {
	http.ResponseWriter
	c int
}

func (w *responseWriterWrapper) WriteHeader(c int) {
	w.c = c
	w.ResponseWriter.WriteHeader(c)
}

func NewResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{w, http.StatusOK}
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

func statusCodeColor(c int) string {
	switch {
	case c >= http.StatusOK && c < http.StatusMultipleChoices:
		return green
	case c >= http.StatusMultipleChoices && c < http.StatusBadRequest:
		return white
	case c >= http.StatusBadRequest && c < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func LogRequest(r http.Request) {
	log.Printf("[%s] -->> %s%s%s %s", r.RemoteAddr, methodColor(r.Method), r.Method, reset, r.URL)
}

func LogResponse(w responseWriterWrapper, r *http.Request) {
	log.Printf("[%s] <<-- %s%d%s %s", r.RemoteAddr, statusCodeColor(w.c), w.c, reset, http.StatusText(w.c))
}

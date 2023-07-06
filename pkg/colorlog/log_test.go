package colorlog

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpStatusText(t *testing.T) {
	tt := []struct {
		s int    // status code
		t string // text
	}{
		{http.StatusOK, "OK"},
		{http.StatusAccepted, "Accepted"},
		{http.StatusCreated, "Created"},
		{http.StatusMovedPermanently, "Moved Permanently"},
		{http.StatusGone, "Gone"},
		{http.StatusBadRequest, "Bad Request"},
		{http.StatusNotFound, "Not Found"},
		{http.StatusUnauthorized, "Unauthorized"},
		{http.StatusInternalServerError, "Internal Server Error"},
		{http.StatusGatewayTimeout, "Gateway Timeout"},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d", tc.s), func(t *testing.T) {
			got := http.StatusText(tc.s)
			if got != tc.t {
				t.Errorf("http.StatusText(%d) == %s, got %s", tc.s, tc.t, got)
				t.FailNow()
			}
		})
	}
}

func TestStatusCodeColor(t *testing.T) {
	tt := []struct {
		s int    // status code
		c string // color
	}{
		{http.StatusOK, green},
		{http.StatusAccepted, green},
		{http.StatusCreated, green},
		{http.StatusMovedPermanently, white},
		{http.StatusGone, yellow},
		{http.StatusBadRequest, yellow},
		{http.StatusNotFound, yellow},
		{http.StatusUnauthorized, yellow},
		{http.StatusInternalServerError, red},
		{http.StatusGatewayTimeout, red},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%d", tc.s), func(t *testing.T) {
			got := statusCodeColor(tc.s)
			if got != tc.c {
				t.Errorf("statusCodeColor(%d) == %s, got %s", tc.s, tc.c, got)
				t.FailNow()
			}
		})
	}
}

func TestMethodColor(t *testing.T) {
	tt := []struct {
		m, c string // method and color
	}{
		{"GET", blue},
		{"POST", cyan},
		{"PUT", yellow},
		{"DELETE", red},
		{"PATCH", green},
		{"HEAD", magenta},
		{"OPTIONS", white},
		{"", reset},
		{"FOO", reset},
	}

	for _, tc := range tt {
		t.Run(tc.m, func(t *testing.T) {
			got := methodColor(tc.m)
			if got != tc.c {
				t.Errorf("methodColor(%s) == %s, got %s", tc.m, tc.c, got)
				t.FailNow()
			}
		})
	}
}

type fakeHttpResponseWriter struct{}

func (f fakeHttpResponseWriter) Header() http.Header {
	return http.Header{}
}

func (f fakeHttpResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (f fakeHttpResponseWriter) WriteHeader(statusCode int) {}

func TestWriteHeader(t *testing.T) {
	s := http.StatusAccepted
	rw := responseWriterWrapper{fakeHttpResponseWriter{}, http.StatusOK}
	rw.WriteHeader(s)

	if rw.c != s {
		t.Errorf("incorrect status, got %d, want %d", rw.c, s)
		t.FailNow()
	}

}

func TestNewResponseWriterWrapper(t *testing.T) {
	rw := NewResponseWriterWrapper(fakeHttpResponseWriter{})
	if rw == nil {
		t.Errorf("NewResponseWriterWrapper cannot be nil")
		t.FailNow()
	}
}

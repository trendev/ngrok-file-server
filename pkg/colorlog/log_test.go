package colorlog

import (
	"testing"
)

func TestMethodColor(t *testing.T) {
	tt := []struct {
		m, c string
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

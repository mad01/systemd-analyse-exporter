package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeFromSystemdAnalyse(t *testing.T) {
	data := []struct {
		osType   string
		input    string
		expected float64
	}{
		{"coreos", "Startup finished in 1.238s (kernel) + 1.677s (initrd) + 7.406s (userspace) = 10.322s", 10.322},
		{"ubuntu", "Startup finished in 5.425s (kernel) + 3min 35.449s (userspace) = 3min 40.875s", 220.875},
	}

	for _, d := range data {
		out, err := getTimeFromSystemdAnalyse(d.input)
		assert.Nil(t, err)
		assert.Equal(t, d.expected, out.Seconds())
	}

	assert.True(t, true, true)

}

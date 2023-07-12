package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	string1 := GenerateRandomString(10)
	// log.Printf("%s", string1)
	string2 := GenerateRandomString(10)
	// log.Printf("%s", string2)

	assert.NotEqual(t, string1, string2)
}

func TestParseEndpoint(t *testing.T) {
	testCases := []struct {
		desc     string
		endpoint string
		want     string
	}{
		{
			desc:     "path with prefix",
			endpoint: "unix:///@path",
			want:     "@path",
		},
		{
			desc:     "path without prefix",
			endpoint: "unix:///path",
			want:     "/path",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := ParseEndpoint(tt.endpoint)
			if err != nil {
				t.Errorf("ParseEndpoint(%q) error: %v", tt.endpoint, err)
			}
			if got != tt.want {
				t.Errorf("ParseEndpoint(%q) = %q, want %q", tt.endpoint, got, tt.want)
			}
		})
	}
}

func TestParseEndpointError(t *testing.T) {
	testCases := []struct {
		desc     string
		endpoint string
		wantErr  string
	}{
		{
			desc:     "empty endpoint",
			endpoint: "",
			wantErr:  "remote KMS provider can't use empty string as endpoint",
		},
		{
			desc:     "invalid scheme",
			endpoint: "http:///path",
			wantErr:  "unsupported scheme \"http\" for remote KMS provider",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := ParseEndpoint(tt.endpoint)
			if err == nil {
				t.Errorf("ParseEndpoint(%q) error: %v", tt.endpoint, err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ParseEndpoint(%q) = %q, want %q", tt.endpoint, err, tt.wantErr)
			}
		})
	}
}

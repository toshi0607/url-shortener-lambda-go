package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		url, method string
		status      int
	}{
		{"https://github.com/toshi0607/url-shortener-lambda-go", http.MethodPost, http.StatusOK},
		{"invalid URL", http.MethodPost, http.StatusBadRequest},
		{"invalid method", http.MethodGet, http.StatusBadRequest},
	}

	for _, te := range tests {
		res, _ := handler(events.APIGatewayProxyRequest{
			HTTPMethod: te.method,
			Body:       "{\"url\": \"" + te.url + "\"}",
		})

		if res.StatusCode != te.status {
			t.Errorf("ExitStatus=%d, want %d", res.StatusCode, te.status)
		}
	}
}

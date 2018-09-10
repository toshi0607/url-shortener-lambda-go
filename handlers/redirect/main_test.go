package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		path, method string
		status       int
	}{
		{"xKlNKGomg", http.MethodGet, http.StatusPermanentRedirect},
		{"xKlNKGomg", http.MethodPost, http.StatusBadRequest},
		{"invalid path", http.MethodGet, http.StatusInternalServerError},
	}

	for _, te := range tests {
		res, _ := handler(events.APIGatewayProxyRequest{
			PathParameters: map[string]string{"shorten_resource": te.path},
			HTTPMethod:     te.method,
		})

		if res.StatusCode != te.status {
			t.Errorf("ExitStatus=%d, want %d", res.StatusCode, te.status)
		}
	}
}

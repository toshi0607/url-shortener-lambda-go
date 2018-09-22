package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/toshi0607/url-shortner-lambda-go/db"
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

func prepare() {
	DynamoDB = db.TestNew()

	if err := DynamoDB.CreateLinkTable(); err != nil {
		panic(err)
	}
}

func cleanUp() {
	if err := DynamoDB.DeleteLinkTable(); err != nil {
		panic(err)
	}
	DynamoDB = db.DB{}
}

func TestMain(m *testing.M) {
	prepare()
	exitCode := m.Run()
	cleanUp()
	os.Exit(exitCode)
}

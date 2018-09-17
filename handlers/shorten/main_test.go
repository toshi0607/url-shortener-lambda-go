package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/toshi0607/url-shortner-lambda-go/db"
)

func TestHandler(t *testing.T) {
	prepare()

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
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(db.Region),
		Endpoint: aws.String("http://localhost:8000")}),
	)

	DynamoDB = db.DB{Instance: dynamodb.New(sess)}
}

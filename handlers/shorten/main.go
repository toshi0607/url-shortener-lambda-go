package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"github.com/toshi0607/url-shortner-lambda-go/db"
)

type request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortenResource string `json:"shorten_resource"`
}

type Link struct {
	ShortenResource string `json:"shorten_resource"`
	OriginalURL     string `json:"original_url"`
}

var DynamoDB db.DB

func init() {
	DynamoDB = db.New()
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	shortenResource := shortid.MustGenerate()
	for shortenResource == "shorten" {
		shortenResource = shortid.MustGenerate()
	}
	link := &Link{
		ShortenResource: shortenResource,
		OriginalURL:     p.URL,
	}

	_, err = DynamoDB.PutItem(link)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}

	b, err := responseBody(shortenResource)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}
	return response(http.StatusOK, b), nil
}

func parseRequest(req events.APIGatewayProxyRequest) (*request, error) {
	if req.HTTPMethod != http.MethodPost {
		return nil, fmt.Errorf("use POST request")
	}

	var r request
	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse request")
	}

	_, err = url.ParseRequestURI(r.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid URL")
	}

	return &r, nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func responseBody(shortenResource string) (string, error) {
	resp, err := json.Marshal(Response{ShortenResource: shortenResource})
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

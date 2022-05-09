package main

import (
	"bytes"
	"encoding/json"
	"io"
)

type APIGResult struct {
	Base64     bool                `json:"isBase64Encoded"`
	StatusCode int                 `json:"statusCode"`
	Headers    map[string]string   `json:"headers"`
	MHeaders   map[string][]string `json:"multiValueHeaders"`
	Body       string              `json:"body"`
}

func APIGRedirect(code int, url string) io.Reader {
	r := APIGResult{
		Base64:     false,
		StatusCode: code,
		Headers: map[string]string{
			"Location":     url,
			"Content-Type": "text/plain",
		},
		Body: "Redirecting to " + url,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return nil
	}
	return bytes.NewReader(b)
}

func APIGText(code int, msg string) io.Reader {
	r := APIGResult{
		Base64:     false,
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Body: msg,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return nil
	}
	return bytes.NewReader(b)
}

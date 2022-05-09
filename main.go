package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ironiridis/gistova"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

type Req struct {
	Base64     bool              `json:"isBase64Encoded"`
	RouteKey   string            `json:"routeKey"`
	Headers    map[string]string `json:"headers,omitempty"`
	Parameters map[string]string `json:"pathParameters,omitempty"`
	Body       string
	Req        struct {
		Host string `json:"domainName"`
	} `json:"requestContext"`
}

func Handle(ctx context.Context, p *gistova.Payload) error {
	var r Req
	err := json.NewDecoder(&p.JSON).Decode(&r)
	if err != nil {
		return fmt.Errorf("Cannot parse payload json: %w", err)
	}
	switch r.RouteKey {
	case "ANY /{x}":
		if url, err := GetToken(ctx, r.Req.Host, r.Parameters["x"]); err != nil {
			return p.Respond(APIGText(400, err.Error()))
		} else {
			return p.Respond(APIGRedirect(301, url))
		}
	case "POST /":
		tkn, err := ProcessPost(ctx, r.Req.Host, r.Base64, r.Body, r.Headers["content-type"])
		if err != nil {
			return p.Respond(APIGText(400, err.Error()))
		}
		return p.Respond(APIGText(202, "Created: https://"+r.Req.Host+"/"+tkn))
	}
	return p.Respond(APIGText(403, "Forbidden"))
}

func main() {
	r := gistova.DefaultRuntime()
	r.LoopFunc(Handle, nil)
}

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
)

func ProcessPost(ctx context.Context, base string, b64 bool, body, ctype string) (tkn string, err error) {
	if b64 {
		if dec, err := base64.StdEncoding.DecodeString(body); err != nil {
			return "", fmt.Errorf("Cannot decode API Gateway payload: %w", err)
		} else {
			body = string(dec)
		}
	}
	var u string
	switch ctype {
	case "application/x-www-form-urlencoded":
		if kv, err := url.ParseQuery(body); err != nil {
			return "", fmt.Errorf("Cannot decode POST data: %w", err)
		} else {
			u = kv.Get("url")
		}
	case "application/json":
		var postJSON struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal([]byte(body), &postJSON); err != nil {
			return "", fmt.Errorf("Cannot decode JSON data: %w", err)
		}
		u = postJSON.URL
	default:
		return "", fmt.Errorf("Unknown content type %q", ctype)
	}
	if u == "" {
		return "", errors.New("No URL provided")
	}
	var urlParsed *url.URL
	if urlParsed, err = url.Parse(u); err != nil {
		return "", fmt.Errorf("Cannot parse URL: %w", err)
	}
	if urlParsed.Scheme != "https" {
		return "", errors.New("URL must have https scheme")
	}
	if urlParsed.User != nil {
		return "", errors.New("URL must not have user authentication data")
	}
	if res, err := net.DefaultResolver.LookupHost(ctx, urlParsed.Hostname()); err != nil {
		return "", fmt.Errorf("Cannot look up %q: %w", urlParsed.Hostname(), err)
	} else if len(res) == 0 {
		return "", fmt.Errorf("Cannot look up %q: no results", urlParsed.Hostname())
	}

	if tkn, err = MakeToken(ctx, base); err != nil {
		return "", fmt.Errorf("Cannot make token: %w", err)
	}
	if err = StoreToken(ctx, base, tkn, urlParsed.String()); err != nil {
		return "", fmt.Errorf("Cannot store token: %w", err)
	}
	return
}

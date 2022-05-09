package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"regexp"
	"strings"
)

const MinTokenLen = 3
const MaxTokenLen = 7
const TokenChars = "0123456789abcdefghijklmnopqrstuvwxyz"

var ValidToken = regexp.MustCompile("^[a-z0-9]{3,7}$")

func GetToken(ctx context.Context, base, t string) (string, error) {
	if !ValidToken.MatchString(t) {
		return "", errors.New("Invalid token")
	}
	srcdns := t + "." + base
	cachefn := "/tmp/cache_" + srcdns
	if url, err := os.ReadFile(cachefn); err == nil {
		return string(url), nil
	}

	if recs, err := net.DefaultResolver.LookupTXT(ctx, srcdns); err != nil {
		return "", fmt.Errorf("Cannot lookup TXT for %q: %w", srcdns, err)
	} else if len(recs) >= 1 {
		os.WriteFile(cachefn, []byte(recs[0]), 0644)
		return recs[0], nil
	}

	return "", errors.New("Not found")
}

func generateTokenStr() string {
	charmaplen := big.NewInt(int64(len(TokenChars)))
	var s strings.Builder
	s.Grow(MaxTokenLen)
	for s.Len() < MaxTokenLen {
		if i, err := rand.Int(rand.Reader, charmaplen); err == nil {
			s.WriteByte(TokenChars[i.Int64()])
		} else {
			return ""
		}
	}
	return s.String()
}

func MakeToken(ctx context.Context, base string) (string, error) {
	attempts := 3
	for attempts > 0 {
		t := generateTokenStr()
		for l := MinTokenLen; l <= MaxTokenLen; l++ {
			if u, _ := GetToken(ctx, base, t[0:l]); u == "" {
				return t[0:l], nil
			}
		}
	}
	return "", errors.New("Exhausted retries")
}

package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

var httpClient http.Client

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"HttpGet": func(url string) (string, error) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return "", err
			}

			res, err := httpClient.Do(req)
			if err != nil {
				return "", err
			}
			defer res.Body.Close()

			if res.StatusCode < 200 || res.StatusCode >= 300 {
				return "", fmt.Errorf("unexpected error code when reaching %s, got %d", url, res.StatusCode)
			}

			payload, err := io.ReadAll(res.Body)
			if err != nil {
				return "", err
			}

			return string(payload), nil
		},

		"QueryString": func(input string) (url.Values, error) {
			return url.ParseQuery(input)
		},

		"URL": func(input string) (*url.URL, error) {
			return url.Parse(input)
		},
	}
}

package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func HTTPGet(url string) (string, error) {
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
}

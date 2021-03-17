package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPGetRet struct {
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Payload    string      `json:"payload"`
}

func (out HTTPGetRet) String() string {
	return out.Payload
}

func HTTPGet(url string) (*HTTPGetRet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected error code when reaching %s, got %d", url, res.StatusCode)
	}

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPGetRet{
		StatusCode: res.StatusCode,
		Header:     res.Header,
		Payload:    string(payload),
	}, nil
}

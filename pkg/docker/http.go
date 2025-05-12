package docker

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

func httpClient(dial func(context.Context) (net.Conn, error)) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (conn net.Conn, err error) {
				return dial(ctx)
			},
		},
	}
}

func get(ctx context.Context, httpClient *http.Client, endpoint string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost"+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-DockerDesktop-Host", "vm.docker.internal")

	response, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, &v); err != nil {
		return err
	}

	return nil
}

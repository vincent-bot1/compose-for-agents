package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Info struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

func getRegistryAuth(ctx context.Context) (string, error) {
	if _, err := os.Stat("/run/host-services/backend.sock"); err != nil {
		return "", nil
	}

	var info Info
	if err := get(ctx, httpClient(dialHostSideBackend), "/registry/info", &info); err != nil {
		return "", fmt.Errorf("getting auth token: %w", err)
	}

	auth_config := map[string]string{
		"username": info.Id,
		"password": info.Token,
	}
	buf, err := json.Marshal(auth_config)
	if err != nil {
		return "", fmt.Errorf("marshalling auth config: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func dialHostSideBackend(ctx context.Context) (net.Conn, error) {
	dialer := net.Dialer{}

	return dialer.DialContext(ctx, "unix", "/run/host-services/backend.sock")
}

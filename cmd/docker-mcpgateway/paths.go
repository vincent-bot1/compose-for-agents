package main

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

func dialJFS(ctx context.Context) (net.Conn, error) {
	return dial(ctx, "Library/Containers/com.docker.docker/Data/jfs.sock")
}

func dial(ctx context.Context, path string) (net.Conn, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dialer := net.Dialer{}
	return dialer.DialContext(ctx, "unix", filepath.Join(home, path))
}

package docker

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

func dialVolumeContents(ctx context.Context) (net.Conn, error) {
	return dial(ctx, "Library/Containers/com.docker.docker/Data/volume-contents.sock")
}

func dialJFS(ctx context.Context) (net.Conn, error) {
	return dial(ctx, "Library/Containers/com.docker.docker/Data/jfs.sock")
}

func dialHostSideBackend(ctx context.Context) (net.Conn, error) {
	dialer := net.Dialer{}
	return dialer.DialContext(ctx, "unix", "/run/host-services/backend.sock")
}

func dial(ctx context.Context, path string) (net.Conn, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dialer := net.Dialer{}
	return dialer.DialContext(ctx, "unix", filepath.Join(home, path))
}

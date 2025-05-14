package docker

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
)

func dialVolumeContents(ctx context.Context) (net.Conn, error) {
	return winio.DialPipeContext(ctx, "dockerVolumeContents")
}

func dialJFS(ctx context.Context) (net.Conn, error) {
	return winio.DialPipeContext(ctx, "dockerJfs")
}

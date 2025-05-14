package docker

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func RunOnDockerDesktop(ctx context.Context, args ...string) (string, error) {
	var host string
	if runtime.GOOS == "windows" {
		host = "npipe:////./pipe/docker_engine_linux"
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		host = "unix://" + filepath.Join(home, "Library/Containers/com.docker.docker/Data/docker.raw.sock")
	}

	args = append([]string{"-H", host, "run", "--rm"}, args...)

	out, err := exec.CommandContext(ctx, "docker", args...).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

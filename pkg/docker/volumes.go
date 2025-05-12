package docker

import (
	"context"
)

type FileContent struct {
	VolumeId   string `json:"volumeId"`
	TargetPath string `json:"targetPath"`
	Contents   string `json:"contents"`
}

func ReadPromptFile(ctx context.Context, name string) (string, error) {
	var content FileContent
	if err := get(ctx, httpClient(dialVolumeContents), "/volume-file-content?volumeId=docker-prompts&targetPath="+name, &content); err != nil {
		return "", err
	}

	return content.Contents, nil
}

package main

import (
	"context"

	"github.com/docker/gateway/cmd/gateway_provider/docker"
	"github.com/spf13/cobra"
)

func NewDownCmd(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "called during compose down",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			containerID := flags.ContainerName(args[0])

			if err := stopGateway(cmd.Context(), containerID); err != nil {
				errorMessage("could not stop the gateway", err)
			} else {
				infoMessage("stopped the gateway")
			}
		},
	}
}

func stopGateway(ctx context.Context, containerID string) error {
	client, err := docker.NewClient(ctx)
	if err != nil {
		return err
	}

	exists, _, err := client.Exists(ctx, containerID)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	return client.RemoveContainer(ctx, containerID, true)
}

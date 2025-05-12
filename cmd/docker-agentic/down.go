package main

import (
	"context"

	"github.com/docker/compose-agents-demo/pkg/compose"
	"github.com/docker/compose-agents-demo/pkg/docker"
	"github.com/spf13/cobra"
)

func NewDownCmd(flags *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "called during compose down",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if flags.UIPort != "" {
				uiContainerID := flags.UIContainerName(args[0])
				if err := stopContainer(cmd.Context(), uiContainerID); err != nil {
					compose.ErrorMessage("could not stop the UI container", err)
				} else {
					compose.InfoMessage("stopped the UI container")
				}
			}

			agentsContainerID := flags.AgentsContainerName(args[0])

			if err := stopContainer(cmd.Context(), agentsContainerID); err != nil {
				compose.ErrorMessage("could not stop the agents container", err)
			} else {
				compose.InfoMessage("stopped the agents container")
			}
		},
	}
}

func stopContainer(ctx context.Context, containerID string) error {
	client, err := docker.NewClient(ctx)
	if err != nil {
		return err
	}

	exists, _, err := client.ContainerExists(ctx, containerID)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	return client.RemoveContainer(ctx, containerID, true)
}

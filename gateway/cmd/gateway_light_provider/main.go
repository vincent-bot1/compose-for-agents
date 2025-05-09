package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
)

func main() {
	plugin.Run(func(command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Use: "compose",
		}
		cmd.PersistentFlags().String("project-name", "", "the project name to use for the compose project")
		cmd.AddCommand(&cobra.Command{
			Use: "up",
			Run: func(*cobra.Command, []string) {
				fmt.Println(`{"type": "setenv","message": "HOST=host.docker.internal:8811"}`)
			},
		})
		cmd.AddCommand(&cobra.Command{
			Use: "down",
			Run: func(*cobra.Command, []string) {},
		})

		originalPreRun := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if originalPreRun != nil {
				return originalPreRun(cmd, args)
			}
			return nil
		}
		return cmd
	},
		manager.Metadata{
			SchemaVersion:    "0.1.0",
			Vendor:           "Docker Inc.",
			Version:          "0.0.1",
			ShortDescription: "Docker MCP Lightweight Gateway Provider",
		})
}

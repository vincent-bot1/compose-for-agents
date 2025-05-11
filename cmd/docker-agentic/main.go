package main

import (
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose-agents-demo/pkg/compose"
)

func main() {
	plugin.Run(func(command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Use: "compose",
		}

		var flags Flags
		cmd.PersistentFlags().StringVar(&flags.Project, "project-name", "", "the project name to use for the compose project")
		cmd.PersistentFlags().StringVar(&flags.Network, "network", "default", "Which docker network to use")
		cmd.PersistentFlags().StringVar(&flags.Config, "config", "", "Configuration for the agents")
		cmd.PersistentFlags().StringVar(&flags.OpenAIAPIKey, "openai_api_key", "", "API Key for OpenAI")
		cmd.PersistentFlags().StringVar(&flags.APIPort, "api_port", "7777", "Port to use for the API")
		cmd.PersistentFlags().StringVar(&flags.UIPort, "ui_port", "", "Port to use for the UI")
		cmd.AddCommand(NewUpCmd(&flags))
		cmd.AddCommand(NewDownCmd(&flags))

		// Don't return an error. Instead, send an error message to compose.
		cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
			compose.ErrorMessage("Error parsing flags", err)
			return nil
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
			ShortDescription: "Docker Agentic Compose Provider",
		})

}

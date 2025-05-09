package main

import (
	"github.com/spf13/cobra"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
)

type Flags struct {
	Project     string
	Tools       string
	Network     string
	LogCalls    string // Should be a bool but compose provider mechanism doesn't like that
	ScanSecrets string // Should be a bool but compose provider mechanism doesn't like that
}

func (f *Flags) ContainerName(providerName string) string {
	return f.Project + "-" + providerName + "-" + f.Network
}

func (f *Flags) NetworkName() string {
	return f.Project + "_" + f.Network
}

func main() {
	plugin.Run(func(command.Cli) *cobra.Command {
		var flags Flags

		cmd := &cobra.Command{
			Use: "compose",
		}
		cmd.PersistentFlags().StringVar(&flags.Project, "project-name", "", "the project name to use for the compose project")
		cmd.PersistentFlags().StringVar(&flags.Network, "network", "default", "Which docker network to use")
		cmd.PersistentFlags().StringVar(&flags.Tools, "tools", "", "Which tools to expose, comma separated list of tools")
		cmd.PersistentFlags().StringVar(&flags.LogCalls, "logCalls", "", "Log the tool calls?")
		cmd.PersistentFlags().StringVar(&flags.ScanSecrets, "scanSecrets", "", "Verify that secrets are not passed to tools")
		cmd.AddCommand(NewUpCmd(&flags))
		cmd.AddCommand(NewDownCmd(&flags))

		// Don't return an error. Instead, send an error message to compose.
		cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
			errorMessage("Error parsing flags", err)
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
			ShortDescription: "Docker MCP Gateway Provider",
		})

}

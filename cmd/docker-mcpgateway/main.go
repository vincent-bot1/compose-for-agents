package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose-agents-demo/pkg/compose"
)

func main() {
	if plugin.RunningStandalone() {
		os.Args = append([]string{os.Args[0], "compose"}, os.Args[1:]...)
	}

	plugin.Run(func(command.Cli) *cobra.Command {
		cmd := &cobra.Command{
			Use: "compose",
		}

		var flags Flags
		cmd.PersistentFlags().StringVar(&flags.Project, "project-name", "", "the project name to use for the compose project")
		cmd.PersistentFlags().StringVar(&flags.Image, "image", "docker/agents_gateway", "Which docker image to use for the gateway")
		cmd.PersistentFlags().StringVar(&flags.Network, "network", "default", "Which docker network to use")
		cmd.PersistentFlags().StringVar(&flags.Tools, "tools", "", "Which tools to expose, comma separated list of tools")
		cmd.PersistentFlags().StringVar(&flags.LogCalls, "log_calls", "", "Log the tool calls?")
		cmd.PersistentFlags().StringVar(&flags.ScanSecrets, "scan_secrets", "", "Verify that secrets are not passed to tools")
		cmd.PersistentFlags().StringVar(&flags.VerifySignatures, "verify_signatures", "", "Verify the image signatures")
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
			ShortDescription: "Docker MCP Gateway Provider",
		})

}

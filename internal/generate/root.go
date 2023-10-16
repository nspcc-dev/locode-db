package generate

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "locode-db",
	Short: "UN/LOCODE database CLI",
	Long:  `locode db CLI is a tool for working with UN/LOCODE database.`,
	Args:  cobra.NoArgs,
	Run:   entryPoint,
}

func init() {
	rootCmd.AddCommand(locodeGenerateCmd)
	initUtilLocodeGenerateCmd()
}

func entryPoint(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.PrintErrln(err)
	}
}

package root

import "github.com/spf13/cobra"

const (
	flagConfig = "config"
)

var configFilePath string

func ConfigFilePath() string {
	return configFilePath
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:              "paralyzer",
		Short:            "code analysis",
		TraverseChildren: true,
		Run:              func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.PersistentFlags().StringVar(&configFilePath, flagConfig, "", "config file path")
	if err := rootCmd.MarkPersistentFlagRequired(flagConfig); err != nil {
		panic(err)
	}
	return rootCmd
}

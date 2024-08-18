package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/madlabx/fs/common/cfg"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the config information",
	Run: func(cmd *cobra.Command, args []string) {
		checkFatalError(initConfigAndLog(cmd.Context()))
		fmt.Println(cfg.Get())
	},
}

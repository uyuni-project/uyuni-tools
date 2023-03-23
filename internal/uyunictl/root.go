package uyunictl

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "uyunictl",
	Short:   "Uyuni control tool",
	Version: "0.0.1",
}

var Verbose bool

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.Execute()
}

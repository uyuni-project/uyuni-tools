package uyunictl

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "uyunictl",
	Short:   "Uyuni control tool",
	Version: "0.0.1",
}

func Execute() {
	rootCmd.Execute()
}

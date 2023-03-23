package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "uyuniadm",
	Short:   "Uyuni setup tool",
	Version: "0.0.1",
}

func main() {
	rootCmd.Execute()
}

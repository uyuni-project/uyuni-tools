package main

import (
	"os"

	"github.com/uyuni-project/uyuni-tools/uyunictl/cmd"
)

// Run runs the `uyunictl` root command
func Run() error {
	return cmd.NewUyunictlCommand().Execute()
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

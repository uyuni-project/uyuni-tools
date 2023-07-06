package main

import (
	"os"

	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd"
)

// Run runs the `uyunictl` root command
func Run() error {
	return cmd.NewUyuniadmCommand().Execute()
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}
}

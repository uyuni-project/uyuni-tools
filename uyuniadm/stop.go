package main

//TODO required go version >= 18
//zypper in libbtrfs-devel libgpgme-devel device-mapper-devel gpgme libassuan >= 2.5.3
//systemct start podman

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/spf13/cobra"
)

func stopContainer(context context.Context, name string) {
	if err := containers.Stop(context, name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Container stopped.")
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop command",
	Run: func(cmd *cobra.Command, arg []string) {
		cmd.OutOrStdout()
		cmd.OutOrStderr()
		fmt.Println("Stopping container")
		stopContainer(ctx, uyuniContainer.name)
	},
}

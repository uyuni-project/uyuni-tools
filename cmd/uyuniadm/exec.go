package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

var fooCmd = &cobra.Command{
    Use: "foo",
    Short: "fooute commands inside the uyuni containers",
    Run: func(cmd *cobra.Command, arg []string) {
        fmt.Println("Running Foo command")
        // TODO run the command
    },
}

func init() {
    rootCmd.AddCommand(fooCmd)
}

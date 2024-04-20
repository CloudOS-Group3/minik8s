package main

import (
	"fmt"
	"minik8s/pkg/kubectl/cmd"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "kubectl",
		Long: "Welcome to use kubectl CLI tool!",
		Run:  nil,
	}

	rootCmd.AddCommand(cmd.GetCmd())
	rootCmd.AddCommand(cmd.ApplyCmd())
	rootCmd.AddCommand(cmd.DeleteCmd())
	rootCmd.AddCommand(cmd.DescribeCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

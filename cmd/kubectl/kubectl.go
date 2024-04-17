package main

import (
	"fmt"
	"os"
	
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "kubectl",
		Long: "Welcome to use kubectl CLI tool!",
		Run: nil,
	}

	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(applyCmd())
	rootCmd.AddCommand(deleteCmd())
	rootCmd.AddCommand(describeCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

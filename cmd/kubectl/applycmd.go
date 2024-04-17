package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func applyCmd() *cobra.Command {

	
	applyCmd := &cobra.Command{
		Use: "apply",
		Short: "apply a yaml file to create a resource",
		Run: applyCmdHandler,
	}

	applyCmd.Flags().StringP("file", "f", "", "specify a file name")

	return applyCmd
}

func applyCmdHandler(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("file")

	if (err != nil) {
		return
	}

	fmt.Println("the file you specified is:", path)
}
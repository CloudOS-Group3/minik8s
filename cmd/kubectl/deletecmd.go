package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {

	var deleteCmd, deletePodCmd, deleteDeploymentCmd, deleteServiceCmd *cobra.Command

	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete a resourse",
		Run:   nil,
	}

	deletePodCmd = &cobra.Command{
		Use:   "pod",
		Short: "delete pod",
		Run:   deletePodCmdHandler,
	}

	deleteDeploymentCmd = &cobra.Command{
		Use:   "deployment",
		Short: "delete deployment",
		Run:   deleteDeploymentCmdHandler,
	}

	deleteServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "delete service",
		Run:   deleteServiceCmdHandler,
	}

	deleteCmd.AddCommand(deletePodCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	deleteCmd.AddCommand(deleteServiceCmd)

	return deleteCmd
}

func deletePodCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("pod name:", args)
}

func deleteDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("deployment name:", args)
}

func deleteServiceCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("service name:", args)
}

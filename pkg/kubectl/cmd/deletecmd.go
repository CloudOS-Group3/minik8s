package cmd

import (
	"minik8s/util/log"

	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {

	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a resourse",
		Run:   nil,
	}

	deletePodCmd := &cobra.Command{
		Use:   "pod",
		Short: "delete pod",
		Run:   deletePodCmdHandler,
	}

	deleteDeploymentCmd := &cobra.Command{
		Use:   "deployment",
		Short: "delete deployment",
		Run:   deleteDeploymentCmdHandler,
	}

	deleteServiceCmd := &cobra.Command{
		Use:   "service",
		Short: "delete service",
		Run:   deleteServiceCmdHandler,
	}

	deleteHPACmd := &cobra.Command{
		Use:   "hpa",
		Short: "delete hpa",
		Run:   deleteHPACmdHandler,
	}

	deleteCmd.AddCommand(deletePodCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	deleteCmd.AddCommand(deleteServiceCmd)
	deleteCmd.AddCommand(deleteHPACmd)

	return deleteCmd
}

func deletePodCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("pod name: %+v", args)
}

func deleteDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("deployment name: %+v", args)
}

func deleteServiceCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("service name: %+v", args)
}

func deleteHPACmdHandler(cmd *cobra.Command, args []string) {
	log.Info("hpa name: %+v", args)
}

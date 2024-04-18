package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func DescribeCmd() *cobra.Command {

	var describeCmd, describePodCmd, describeDeploymentCmd, describeServiceCmd, describeNodeCmd *cobra.Command

	describeCmd = &cobra.Command{
		Use:   "describe",
		Short: "describe the config of some resourse",
		Run:   nil,
	}

	describePodCmd = &cobra.Command{
		Use:   "pod",
		Short: "describe a pod",
		Run:   describePodCmdHandler,
	}

	describeDeploymentCmd = &cobra.Command{
		Use:   "deployment",
		Short: "describe a deployment",
		Run:   describeDeploymentCmdHandler,
	}

	describeServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "describe a service",
		Run:   describeServiceCmdHandler,
	}

	describeNodeCmd = &cobra.Command{
		Use:   "node",
		Short: "describe a node",
		Run:   describeNodeCmdHandler,
	}

	describeCmd.AddCommand(describePodCmd)
	describeCmd.AddCommand(describeDeploymentCmd)
	describeCmd.AddCommand(describeServiceCmd)
	describeCmd.AddCommand(describeNodeCmd)

	return describeCmd
}

func describePodCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("pod name:", args)
}

func describeDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("deployment name:", args)
}

func describeServiceCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("service name:", args)
}

func describeNodeCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("node name:", args)
}

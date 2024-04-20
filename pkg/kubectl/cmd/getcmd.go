package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {

	var getCmd, getPodCmd, getDeploymentCmd, getServiceCmd, getNodeCmd *cobra.Command

	// getCmd is the root of the other four commands
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "get the infomation of resource",
		Run:   nil,
	}

	getPodCmd = &cobra.Command{
		Use:   "pod",
		Short: "get pod",
		Run:   getPodCmdHandler,
	}

	getDeploymentCmd = &cobra.Command{
		Use:   "deployment",
		Short: "get deployment",
		Run:   getDeploymentCmdHandler,
	}

	getServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "get service",
		Run:   getServiceCmdHandler,
	}

	getNodeCmd = &cobra.Command{
		Use:   "node",
		Short: "get node",
		Run:   getNodeCmdHandler,
	}

	// support -a flag, but the implementation could be troublesome
	getPodCmd.Flags().BoolP("all", "a", false, "get all pod")
	getDeploymentCmd.Flags().BoolP("all", "a", false, "get all deployment")
	getServiceCmd.Flags().BoolP("all", "a", false, "get all service")
	getNodeCmd.Flags().BoolP("all", "a", false, "get all node")

	getCmd.AddCommand(getPodCmd)
	getCmd.AddCommand(getDeploymentCmd)
	getCmd.AddCommand(getServiceCmd)
	getCmd.AddCommand(getNodeCmd)

	return getCmd
}

// all the handlers below should be replaced by real k8s logic later
func getPodCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("getting pod")
	all, _ := cmd.Flags().GetBool("all")
	if all {
		fmt.Println("getting all pod")
	}
}

func getDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("getting deployment")
	all, _ := cmd.Flags().GetBool("all")
	if all {
		fmt.Println("getting all deployment")
	}
}

func getServiceCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("getting service")
	all, _ := cmd.Flags().GetBool("all")
	if all {
		fmt.Println("getting all service")
	}
}

func getNodeCmdHandler(cmd *cobra.Command, args []string) {
	fmt.Println("getting node")
	all, _ := cmd.Flags().GetBool("all")
	if all {
		fmt.Println("getting all node")
	}
}

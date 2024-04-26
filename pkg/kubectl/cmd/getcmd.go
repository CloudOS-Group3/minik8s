package cmd

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/apiserver/config"
	"minik8s/util/log"
	"net/http"
	"strings"

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

	for _, podName := range args {
		log.Debug("%v", podName)
		URL := config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, podName, -1)
		response, err := http.Get("http://localhost:6443" + URL)
		if err != nil {
			log.Error("error http get, %s", err.Error())
			return
		}
		pod := &api.Pod{}
		log.Info("%+v", response)
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(pod)
		if err != nil {
			log.Error("error decode response body")
			return
		}
		log.Debug("%v", pod)
	}
}

func getDeploymentCmdHandler(cmd *cobra.Command, args []string) {

}

func getServiceCmdHandler(cmd *cobra.Command, args []string) {

}

func getNodeCmdHandler(cmd *cobra.Command, args []string) {

}

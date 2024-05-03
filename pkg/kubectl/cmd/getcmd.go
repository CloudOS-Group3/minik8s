package cmd

import (
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
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

	getCmd.AddCommand(getPodCmd)
	getCmd.AddCommand(getDeploymentCmd)
	getCmd.AddCommand(getServiceCmd)
	getCmd.AddCommand(getNodeCmd)

	return getCmd
}


// TODO: all of these handlers have got the data, but they dont show it in the terminal
func getPodCmdHandler(cmd *cobra.Command, args []string) {

	log.Debug("the length of args is: %v", len(args))

	matchPods := []api.Pod{}
	if len(args) == 0 {
		log.Info("getting all pods")
		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

		err := httputil.Get(URL, &matchPods, "data")
		if err != nil {
			log.Error("error get app pods: %s", err.Error())
			return
		}
	} else {
		for _, podName := range args {
			pod := &api.Pod{}

			log.Debug("getting pod: %v", podName)
			URL := config.GetUrlPrefix() + config.PodURL
			URL = strings.Replace(URL, config.NamePlaceholder, podName, -1)
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

			err := httputil.Get(URL, pod, "data")

			if err != nil {
				log.Error("error get pod: %s", err.Error())
				return
			}

			log.Debug("%+v", pod)
			matchPods = append(matchPods, *pod)
		}
	}
}

func getDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("the length of the args is: %v", len(args))
	matchDeployments := []api.Deployment{}

	if len(args) == 0 {
		log.Info("getting all deployments")
		URL := config.GetUrlPrefix() + config.DeploymentsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		err := httputil.Get(URL, matchDeployments, "data")
		if err != nil {
			log.Error("error getting all deployments: %s", err.Error())
			return
		}
	} else {
		for _, deploymentName := range args {
			deployment := &api.Deployment{}

			URL := config.GetUrlPrefix() + config.DeploymentURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
			URL = strings.Replace(URL, config.NamePlaceholder, deploymentName, -1)

			httputil.Get(URL, deployment, "data")

			matchDeployments = append(matchDeployments, *deployment)
		}
	}
}

func getServiceCmdHandler(cmd *cobra.Command, args []string) {

}

func getNodeCmdHandler(cmd *cobra.Command, args []string) {
	log.Debug("the length of args is: %v", len(args))

	matchNodes := []api.Node{}

	if len(args) == 0 {
		log.Debug("getting all nodes")
		URL := config.GetUrlPrefix() + config.NodesURL
		err := httputil.Get(URL, matchNodes, "data")
		if err != nil {
			log.Error("error getting all nodes: %s", err.Error())
			return
		}
	} else {
		for _, nodeName := range args {
			log.Debug("%v", nodeName)
			node := &api.Node{}
			URL := config.GetUrlPrefix() + config.NodeURL
			URL = strings.Replace(URL, config.NamePlaceholder, nodeName, -1)
			err := httputil.Get(URL, node, "data")
			if err != nil {
				log.Error("error get node: %s", err.Error())
				return
			}
			log.Debug("%+v", node)
			matchNodes = append(matchNodes, *node)
		}
	}
}

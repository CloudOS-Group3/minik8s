package cmd

import (
	"fmt"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/util"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/prettyprint"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {

	// getCmd is the root of the other four commands
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "get the infomation of resource",
		Run:   nil,
	}

	getPodCmd := &cobra.Command{
		Use:   "pod",
		Short: "get pod",
		Run:   getPodCmdHandler,
	}

	getNodeCmd := &cobra.Command{
		Use:   "node",
		Short: "get node",
		Run:   getNodeCmdHandler,
	}

	getServiceCmd := &cobra.Command{
		Use:   "service",
		Short: "get service",
		Run:   getServiceCmdHandler,
	}

	getDeploymentCmd := &cobra.Command{
		Use:   "deployment",
		Short: "get deployment",
		Run:   getDeploymentCmdHandler,
	}

	getHPACmd := &cobra.Command{
		Use:   "hpa",
		Short: "get hpa",
		Run:   getHPACmdHandler,
	}

	getPodCmd.Aliases = []string{"po", "pods"}
	getNodeCmd.Aliases = []string{"no", "nodes"}
	getServiceCmd.Aliases = []string{"svc", "service"}
	getDeploymentCmd.Aliases = []string{"deployments"}
	getHPACmd.Aliases = []string{"hpas"}

	getPodCmd.Flags().StringP("namespace", "n", "default", "namespace of the pod")
	getServiceCmd.Flags().StringP("namespace", "n", "default", "namespace of the pod")
	getCmd.AddCommand(getPodCmd)
	getCmd.AddCommand(getNodeCmd)
	getCmd.AddCommand(getDeploymentCmd)
	getCmd.AddCommand(getServiceCmd)

	return getCmd
}

// TODO: all of these handlers have got the data, but they dont show it in the terminal
func getPodCmdHandler(cmd *cobra.Command, args []string) {

	namespace := cmd.Flag("namespace").Value.String()

	matchPods := []api.Pod{}
	if len(args) == 0 {
		log.Info("getting all pods")
		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)

		err := httputil.Get(URL, &matchPods, "data")
		if err != nil {
			log.Error("error get app pods: %s", err.Error())
			return
		}
		log.Debug("match pods are: %+v", matchPods)
	} else {
		for _, podName := range args {
			pod := &api.Pod{}

			log.Debug("getting pod: %v", podName)
			URL := config.GetUrlPrefix() + config.PodURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
			URL = strings.Replace(URL, config.NamePlaceholder, podName, -1)

			err := httputil.Get(URL, pod, "data")

			if err != nil {
				log.Error("error get pod: %s", err.Error())
				return
			}

			log.Debug("%+v", pod)
			if pod.Metadata.Name == "" {
				continue
			}
			matchPods = append(matchPods, *pod)
		}
	}

	header := []string{"name", "status", "age", "usage", "ip"}
	data := [][]string{}

	for _, matchPod := range matchPods {
		age := time.Now().Sub(matchPod.Status.StartTime).Round(time.Second).String()
		metric_string := fmt.Sprintf("cpu: %.2f%%, memory: %.2f%%", matchPod.Status.CPUPercentage*100, matchPod.Status.MemoryPercentage*100)
		data = append(data, []string{matchPod.Metadata.Name, matchPod.Status.Phase, age, metric_string, matchPod.Status.PodIP})
	}

	prettyprint.PrintTable(header, data)
}

func getDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("the length of the args is: %v", len(args))
	matchDeployments := []api.Deployment{}

	if len(args) == 0 {
		log.Info("getting all deployments")
		URL := config.GetUrlPrefix() + config.DeploymentsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		err := httputil.Get(URL, &matchDeployments, "data")
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
	namespace := cmd.Flag("namespace").Value.String()
	matchServices := []api.Service{}

	if len(args) == 0 {
		log.Debug("getting all services")
		URL := config.GetUrlPrefix() + config.ServicesURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
		err := httputil.Get(URL, &matchServices, "data")
		if err != nil {
			log.Error("error get all services: %s", err.Error())
			return
		}
	} else {
		for _, serviceName := range args {
			service := &api.Service{}
			URL := config.GetUrlPrefix() + config.ServiceURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
			URL = strings.Replace(URL, config.NamePlaceholder, serviceName, -1)

			err := httputil.Get(URL, service, "data")
			if err != nil {
				log.Error("error get service: %s", err.Error())
				return
			}
			matchServices = append(matchServices, *service)
		}
	}

	header := []string{"name", "label", "ip"}
	data := [][]string{}
	for _, matchService := range matchServices {
		labelstring := util.ConvertLabelToString(matchService.Spec.Selector)
		data = append(data, []string{matchService.Metadata.Name, labelstring, matchService.Status.ClusterIP})

	}
	prettyprint.PrintTable(header, data)
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

func getHPACmdHandler(cmd *cobra.Command, args []string) {
	log.Debug("the length of args is: %v", len(args))

	matchHPAs := []api.HPA{}
	if len(args) == 0 {
		log.Debug("getting all HPAs")
		URL := config.GetUrlPrefix() + config.HPAsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		err := httputil.Get(URL, matchHPAs, "data")
		if err != nil {
			log.Error("Error http get hpa: %s", err.Error())
			return
		}
	} else {
		for _, hpaName := range args {
			log.Debug("getting hpa: %s", hpaName)

			hpa := &api.HPA{}
			URL := config.GetUrlPrefix() + config.HPAURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
			URL = strings.Replace(URL, config.NamePlaceholder, hpaName, -1)

			err := httputil.Get(URL, hpa, "data")

			if err != nil {
				log.Error("Error http get hpa: %s", err.Error())
				return
			}

			log.Debug("hpa is: %+v", hpa)
			matchHPAs = append(matchHPAs, *hpa)
		}
	}

}

package cmd

import (
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"

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

	deletePodCmd.Aliases = []string{"po", "pods"}
	deleteServiceCmd.Aliases = []string{"svc", "service"}
	deleteDeploymentCmd.Aliases = []string{"deployments"}
	deleteHPACmd.Aliases = []string{"hpas"}

	deleteCmd.AddCommand(deletePodCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	deleteCmd.AddCommand(deleteServiceCmd)
	deleteCmd.AddCommand(deleteHPACmd)

	return deleteCmd
}

func deletePodCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("pod name: %+v", args)
	matchPods := []api.Pod{}
	if len(args) == 0 {
		log.Info("getting all pods")
		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

		err := httputil.Get(URL, &matchPods, "data")
		if err != nil {
			log.Error("error getting all pods: %s", err.Error())
		}
	} else {
		for _, podName := range args {
			var pod api.Pod
			URL := config.GetUrlPrefix() + config.PodURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
			URL = strings.Replace(URL, config.NamePlaceholder, podName, -1)

			err := httputil.Get(URL, &pod, "data")
			if err != nil {
				log.Error("error getting pod %s with error %s", podName, err.Error())
				continue
			}

			matchPods = append(matchPods, pod)
		}
	}

	for _, pod := range matchPods {
		URL := config.GetUrlPrefix() + config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, pod.Metadata.Name, -1)

		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error deleting pod %s with error %s", pod.Metadata.Name, err.Error())
		}
	}

}

func deleteDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("deployment name: %+v", args)
	matchDeployments := []api.Deployment{}
	if len(args) == 0 {
		log.Info("getting all deployments")
		URL := config.GetUrlPrefix() + config.DeploymentsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

		err := httputil.Get(URL, &matchDeployments, "data")
		if err != nil {
			log.Error("error getting all deployments: %s", err.Error())
		}
	} else {
		for _, deploymentName := range args {
			var deployment api.Deployment
			URL := config.GetUrlPrefix() + config.DeploymentURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
			URL = strings.Replace(URL, config.NamePlaceholder, deploymentName, -1)

			err := httputil.Get(URL, &deployment, "data")
			if err != nil {
				log.Error("error getting deployment %s with error %s", deploymentName, err.Error())
				continue
			}

			matchDeployments = append(matchDeployments, deployment)
		}
	}

	for _, deployment := range matchDeployments {
		URL := config.GetUrlPrefix() + config.DeploymentURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, deployment.Metadata.Name, -1)

		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error deleting pod %s with error %s", deployment.Metadata.Name, err.Error())
		}
	}

}

func deleteServiceCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("service name: %+v", args)
}

func deleteHPACmdHandler(cmd *cobra.Command, args []string) {
	log.Info("hpa name: %+v", args)
	matchHPAs := []api.HPA{}
	if len(args) == 0 {
		log.Info("getting all HPAs")
		URL := config.GetUrlPrefix() + config.HPAsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

		err := httputil.Get(URL, &matchHPAs, "data")
		if err != nil {
			log.Error("error getting all HPAs: %s", err.Error())
		}
	} else {
		for _, HPAName := range args {
			var HPA api.HPA
			URL := config.GetUrlPrefix() + config.HPAURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
			URL = strings.Replace(URL, config.NamePlaceholder, HPAName, -1)

			err := httputil.Get(URL, &HPA, "data")
			if err != nil {
				log.Error("error getting HPA %s with error %s", HPAName, err.Error())
				continue
			}

			matchHPAs = append(matchHPAs, HPA)
		}
	}

	for _, HPA := range matchHPAs {
		URL := config.GetUrlPrefix() + config.HPAURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, HPA.Metadata.Name, -1)

		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error deleting HPA %s with error %s", HPA.Metadata.Name, err.Error())
		}
	}

}

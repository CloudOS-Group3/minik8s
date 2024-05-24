package cmd

import (
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
		Use:   "pod [pod name]",
		Short: "delete pod",
		Args:  cobra.MinimumNArgs(1),
		Run:   deletePodCmdHandler,
	}

	deleteDeploymentCmd := &cobra.Command{
		Use:   "deployment",
		Short: "delete deployment",
		Run:   deleteDeploymentCmdHandler,
	}

	deleteServiceCmd := &cobra.Command{
		Use:   "service [service name]",
		Short: "delete service",
		Run:   deleteServiceCmdHandler,
	}

	deleteHPACmd := &cobra.Command{
		Use:   "hpa",
		Short: "delete hpa",
		Run:   deleteHPACmdHandler,
	}

	deletePodCmd.Flags().StringP("namespace", "n", "default", "specify the namespace of the resource")
	deleteServiceCmd.Flags().StringP("namespace", "n", "default", "specify the namespace of the resource")

	deleteCmd.AddCommand(deletePodCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	deleteCmd.AddCommand(deleteServiceCmd)
	deleteCmd.AddCommand(deleteHPACmd)

	return deleteCmd
}

func deletePodCmdHandler(cmd *cobra.Command, args []string) {
	// delete pod name --namespace=default

	if len(args) == 0 {
		log.Debug("deleting all pods")
		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error delete all pods")
			return
		}
		return
	}

	name := args[0]
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		log.Error("Error getting flags: %s", err)
		return
	}
	path := strings.Replace(config.PodURL, config.NamespacePlaceholder, namespace, -1)
	path = strings.Replace(path, config.NamePlaceholder, name, -1)
	URL := config.GetUrlPrefix() + path
	err = httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}

	log.Info("pod name: %s, namespace: %s", name, namespace)

}

func deleteDeploymentCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("deployment name: %+v", args)
}

func deleteServiceCmdHandler(cmd *cobra.Command, args []string) {
	name := args[0]
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		log.Error("Error getting flags: %s", err)
		return
	}
	path := strings.Replace(config.ServiceURL, config.NamespacePlaceholder, namespace, -1)
	path = strings.Replace(path, config.NamePlaceholder, name, -1)
	URL := config.GetUrlPrefix() + path
	err = httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
}

func deleteHPACmdHandler(cmd *cobra.Command, args []string) {
	name := args[0]

	URL := config.GetUrlPrefix() + config.HPAURL
	URL := strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)

	err := httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
}

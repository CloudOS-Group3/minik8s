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

	deleteDNSCmd := &cobra.Command{
		Use:   "dns [dns name]",
		Short: "delete DNS",
		Run:   deleteDNSCmdHandler,
	}

	deletePVCmd := &cobra.Command{
		Use:   "pv",
		Short: "delete PV",
		Run:   deletePVCmdHandler,
	}

	deletePVCCmd := &cobra.Command{
		Use:   "pvc",
		Short: "delete PVC",
		Run:   deletePVCCmdHandler,
	}

	deleteJobCmd := &cobra.Command{
		Use:   "job [job name]",
		Short: "delete job",
		Run:   deleteJobCmdHandler,
	}

	deleteFunctionCmd := &cobra.Command{
		Use:   "function [function name]",
		Short: "delete function",
		Run:   deleteFunctionCmdHandler,
	}

	deleteTriggerCmd := &cobra.Command{
		Use:   "trigger [function name]",
		Short: "delete trigger",
		Run:   deleteTriggerCmdHandler,
	}

	deleteNodeCmd := &cobra.Command{
		Use:   "node [node name]",
		Short: "delete node",
		Run:   deleteNodeCmdHandler,
	}

	deletePodCmd.Flags().StringP("namespace", "n", "default", "specify the namespace of the resource")
	deleteServiceCmd.Flags().StringP("namespace", "n", "default", "specify the namespace of the resource")
	deleteFunctionCmd.Flags().StringP("namespace", "n", "default", "specify the namespace of the resource")
	deleteTriggerCmd.Flags().StringP("namespace", "n", "default", "specify the namespace of the resource")
	deleteTriggerCmd.Flags().BoolP("workflow", "w", false, "Indicates if the trigger is a workflow")

	deleteCmd.AddCommand(deletePodCmd)
	deleteCmd.AddCommand(deleteDeploymentCmd)
	deleteCmd.AddCommand(deleteServiceCmd)
	deleteCmd.AddCommand(deleteHPACmd)
	deleteCmd.AddCommand(deleteDNSCmd)
	deleteCmd.AddCommand(deletePVCmd)
	deleteCmd.AddCommand(deletePVCCmd)
	deleteCmd.AddCommand(deleteFunctionCmd)
	deleteCmd.AddCommand(deleteTriggerCmd)
	deleteCmd.AddCommand(deleteJobCmd)
	deleteCmd.AddCommand(deleteNodeCmd)

	return deleteCmd
}

func deleteFunctionCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Error("function name is required")
		return
	}
	name := args[0]
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		log.Error("Error getting flags: %s", err)
		return
	}
	path := strings.Replace(config.FunctionURL, config.NamespacePlaceholder, namespace, -1)
	path = strings.Replace(path, config.NamePlaceholder, name, -1)
	URL := config.GetUrlPrefix() + path
	err = httputil.Delete(URL)
	if err != nil {
		log.Error("Wait for all related pods to be deleted first: %s", err.Error())
		return
	}
	log.Info("function name: %s, namespace: %s", name, namespace)

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
	name := args[0]
	URL := config.GetUrlPrefix() + config.DeploymentURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)

	err := httputil.Delete(URL)
	if err != nil {
		log.Error("error delete deployment: %s", err.Error())
		return
	}

	log.Info("delete deployment %s successfully", name)
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
	log.Info("deleting hpa")
	name := args[0]

	URL := config.GetUrlPrefix() + config.HPAURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)

	err := httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
	log.Info("successfully deleted hpa")
}

func deleteDNSCmdHandler(cmd *cobra.Command, args []string) {
	name := args[0]
	URL := config.GetUrlPrefix() + config.DNSURL
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
	err := httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
}

func deletePVCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("deleting pv")
	for _, name := range args {
		URL := config.GetUrlPrefix() + config.PersistentVolumeURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error http post: %s", err.Error())
			continue
		}
	}
	log.Info("successfully deleted pv")
}

func deletePVCCmdHandler(cmd *cobra.Command, args []string) {
	log.Info("deleting pvc")
	for _, name := range args {
		URL := config.GetUrlPrefix() + config.PersistentVolumeClaimURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error http post: %s", err.Error())
			continue
		}
	}
	log.Info("successfully deleted pvc")
}

func deleteJobCmdHandler(cmd *cobra.Command, args []string) {
	name := args[0]
	URL := config.GetUrlPrefix() + config.JobURL
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
	err := httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
}

func deleteNodeCmdHandler(cmd *cobra.Command, args []string) {
	name := args[0]
	URL := config.GetUrlPrefix() + config.NodeURL
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
	err := httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
}

func deleteTriggerCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Error("function name is required")
		return
	}
	name := args[0]
	namespace, err := cmd.Flags().GetString("namespace")
	isWorkflow, err := cmd.Flags().GetBool("workflow")
	if err != nil {
		log.Error("Error getting flags: %s", err)
		return
	}
	if isWorkflow {
		path := strings.Replace(config.TriggerWorkflowURL, config.NamespacePlaceholder, namespace, -1)
		path = strings.Replace(path, config.NamePlaceholder, name, -1)
		URL := config.GetUrlPrefix() + path
		err = httputil.Delete(URL)
		if err != nil {
			log.Error("error http post: %s", err.Error())
			return
		}
		return
	}
	path := strings.Replace(config.TriggerURL, config.NamespacePlaceholder, namespace, -1)
	path = strings.Replace(path, config.NamePlaceholder, name, -1)
	URL := config.GetUrlPrefix() + path
	err = httputil.Delete(URL)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
	log.Info("function name: %s, namespace: %s", name, namespace)
}

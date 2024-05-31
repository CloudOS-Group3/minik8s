package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/serverless/function/function_util"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
)

func HttpCmd() *cobra.Command {

	httpCmd := &cobra.Command{
		Use:     "http",
		Short:   "send http trigger to function or workflow",
		Run:     nil,
		Example: "http function -n namespace name arg1 arg2...",
	}

	httpFunctionCmd := &cobra.Command{
		Use:   "function",
		Short: "http function",
		Run:   httpFuncHandler,
	}

	httpWorkflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "http workflow",
		Run:   httpWorkflowHandler,
	}

	httpFunctionCmd.Flags().StringP("namespace", "n", "default", "namespace of the pod")
	httpWorkflowCmd.Flags().StringP("namespace", "n", "default", "namespace of the service")
	httpCmd.AddCommand(httpFunctionCmd)
	httpCmd.AddCommand(httpWorkflowCmd)

	return httpCmd
}

func httpWorkflowHandler(cmd *cobra.Command, args []string) {
	//curl -X POST -H "Content-Type: application/json" -d '{"params": {"x": 8, "y": 9}}' localhost:6443/api/v1/namespaces/default/functions/matrix-calculate/run
	namespace := cmd.Flag("namespace").Value.String()
	if len(args) < 1 {
		log.Fatal("Usage: http function -n namespace name arg1 arg2 ...")
		return
	}
	wfName := args[0]

	// Get the function
	workflow := api.Workflow{}
	URL := config.GetUrlPrefix() + config.WorkflowURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, wfName, -1)
	err := httputil.Get(URL, &workflow, "data")
	if err != nil || workflow.Metadata.Name == "" {
		log.Error("Can't find function: %s", wfName)
		return
	}

	// Send the request
	URL = config.GetUrlPrefix() + config.WorkflowRunURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, wfName, -1)

	byteArr, err := json.Marshal(workflow)
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error sending request: %s", err)
		return
	}
}

func httpFuncHandler(cmd *cobra.Command, args []string) {
	//curl -X POST -H "Content-Type: application/json" -d '{"params": {"x": 8, "y": 9}}' localhost:6443/api/v1/namespaces/default/functions/matrix-calculate/run
	namespace := cmd.Flag("namespace").Value.String()
	if len(args) < 1 {
		log.Fatal("Usage: http function -n namespace name arg1 arg2 ...")
		return
	}
	functionName := args[0]

	// Get the function
	function := api.Function{}
	URL := config.GetUrlPrefix() + config.FunctionURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, functionName, -1)
	err := httputil.Get(URL, &function, "data")
	if err != nil || function.Metadata.Name == "" {
		log.Error("Can't find function: %s", functionName)
		return
	}

	params, err := function_util.CheckParams(function.Params, args[1:])
	if err != nil {
		log.Error("Error checking params: %s", err)
		return
	}

	//Create the payload
	jsonParams, err := json.Marshal(params)
	payload := map[string]interface{}{
		"params": string(jsonParams),
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	log.Info("Payload: %s", string(jsonPayload))

	// Send the request
	URL = config.GetUrlPrefix() + config.FunctionRunURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, functionName, -1)
	err = httputil.Post(URL, jsonPayload)
	if err != nil {
		log.Error("Error sending request: %s", err)
		return
	}
}

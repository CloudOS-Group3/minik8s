package cmd

import (
	"encoding/json"
	"io"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func ApplyCmd() *cobra.Command {

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply a yaml file to create a resource",
		Run:   applyCmdHandler,
	}

	applyCmd.Flags().StringP("file", "f", "", "specify a file name")

	return applyCmd
}

func getKindFromYaml(content []byte) string {
	var resource map[string]interface{}
	yaml.Unmarshal(content, &resource)
	if resource["kind"] == "" {
		log.Error("kind field is empty")
		return ""
	}
	return resource["kind"].(string)
}

func applyCmdHandler(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("file")
	if err != nil {
		log.Error("Error getting flags: %s", err)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Error("Error opening file: %s", err)
		return
	}

	defer file.Close()

	fd, err := os.Open(path)
	if err != nil {
		log.Error("error open file")
		return
	}

	defer fd.Close()

	content, err := io.ReadAll(fd)
	if err != nil {
		log.Error("error read all data")
		return
	}
	kind := getKindFromYaml(content)

	switch kind {
	case "Pod":
		applyPodHandler(content)
	case "Service":
		applyServiceHandler(content)
	case "Deployment":
		applyDeploymentHandler(content)
	case "HPA":
		applyHPAHandler(content)
	case "DNS":
		applyDNSHandler(content)
	case "PV":
		applyPVHandler(content)
	case "PVC":
		applyPVCHandler(content)
	case "Function":
		applyFunctionHandler(content)
	case "Trigger":
		applyTriggerHandler(content)
	case "Workflow":
		applyWorkflowHandler(content)
	case "Node":
		applyNodeHandler(content)
	case "GPU":
		applyGPUHandler(content)
	default:
		log.Warn("Unknown resource kind")
	}

}

func applyGPUHandler(content []byte) {
	log.Info("Creating or updating GPU")
	gpu := &api.GPUJob{}
	err := yaml.Unmarshal(content, gpu)
	if err != nil {
		log.Error("Error yaml unmarshal GPU")
		return
	}
	byteArr, err := json.Marshal(*gpu)
	if err != nil {
		log.Error("Error json marshal GPU")
		return
	}

	URL := config.GetUrlPrefix() + config.GPUJobURL
	URL = strings.Replace(URL, config.NamePlaceholder, gpu.Metadata.Name, -1)
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply GPU successed")
}

func applyWorkflowHandler(content []byte) {
	log.Info("Creating or updating workflow")
	workflow := &api.Workflow{}
	err := yaml.Unmarshal(content, workflow)
	if err != nil {
		log.Error("Error yaml unmarshal workflow, %s", err.Error())
		return
	}

	workflow.Metadata.UUID = uuid.NewString()
	if workflow.Metadata.NameSpace == "" {
		workflow.Metadata.NameSpace = "default"
	}

	byteArr, err := json.Marshal(*workflow)
	if err != nil {
		log.Error("Error json marshal workflow, %s", err.Error())
		return
	}
	URL := config.GetUrlPrefix() + config.WorkflowURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, workflow.Metadata.NameSpace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, workflow.Metadata.Name, -1)
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post workflow, %s", err.Error())
		return
	}
	log.Info("apply workflow successed, %v", workflow)

}

func applyFunctionHandler(content []byte) {
	log.Info("Creating or updating function")
	function := &api.Function{}
	err := yaml.Unmarshal(content, function)
	if err != nil {
		log.Error("Error yaml unmarshal function, %s", err.Error())
		return
	}

	// Check if exsit
	existFunction := &api.Function{}
	URL := config.GetUrlPrefix() + config.FunctionURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, function.Metadata.NameSpace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, function.Metadata.Name, -1)
	_ = httputil.Get(URL, existFunction, "data")
	if existFunction.Metadata.Name != "" {
		byteArr, err := json.Marshal(*function)
		if err != nil {
			log.Error("Error json marshal function, %s", err.Error())
			return
		}
		log.Warn("Function %s already exists", function.Metadata.Name)
		URL = config.GetUrlPrefix() + config.FunctionURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, function.Metadata.NameSpace, -1)
		URL = strings.Replace(URL, config.NamePlaceholder, function.Metadata.Name, -1)
		// delete old function
		err = httputil.Delete(URL)
		if err != nil {
			log.Error("Error http delete function, %s", err.Error())
			return
		}
		// create new function
		err = httputil.Post(URL, byteArr)
		if err != nil {
			log.Error("Error http put function, %s", err.Error())
			return
		}
		log.Info("update function successed, %v", function)
		return
	}

	//Generate
	function.Metadata.UUID = uuid.NewString()
	if function.Metadata.NameSpace == "" {
		function.Metadata.NameSpace = "default"
	}

	byteArr, err := json.Marshal(*function)
	if err != nil {
		log.Error("Error json marshal function, %s", err.Error())
		return
	}
	URL = config.GetUrlPrefix() + config.FunctionURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, function.Metadata.NameSpace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, function.Metadata.Name, -1)
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post function, %s", err.Error())
		return
	}
	log.Info("apply function successed, %v", function)
}
func applyPodHandler(content []byte) {
	log.Info("Creating or updating pod")
	pod := &api.Pod{}
	err := yaml.Unmarshal(content, pod)
	if err != nil {
		log.Error("error marshal yaml, %s", err.Error())
		return
	}
	pod.Metadata.UUID = uuid.NewString()
	log.Debug("%+v\n", pod)

	var namespace string
	if pod.Metadata.NameSpace != "" {
		namespace = pod.Metadata.NameSpace
	} else {
		namespace = "default"
	}
	pod.Metadata.UUID = uuid.NewString()

	URL := config.GetUrlPrefix() + config.PodsURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	byteArr, err := json.Marshal(*pod)
	if err != nil {
		log.Error("error marshal yaml")
		return
	}

	err = httputil.Post(URL, byteArr)

	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}

	log.Info("apply pod successed")
}

func applyServiceHandler(content []byte) {
	log.Info("creating or updating service")
	service := &api.Service{}
	err := yaml.Unmarshal(content, service)
	if err != nil {
		log.Error("Error yaml unmarshal service")
		return
	}

	var namespace string
	if service.Metadata.NameSpace != "" {
		namespace = service.Metadata.NameSpace
	} else {
		namespace = "default"
		service.Metadata.NameSpace = namespace
	}
	path := strings.Replace(config.ServiceURL, config.NamespacePlaceholder, namespace, -1)
	path = strings.Replace(path, config.NamePlaceholder, service.Metadata.Name, -1)
	byteArr, err := json.Marshal(*service)
	URL := config.GetUrlPrefix() + path

	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("error http post: %s", err.Error())
		return
	}
}

func applyDeploymentHandler(content []byte) {
	log.Info("creating or updating deployment")

	deployment := &api.Deployment{}
	err := yaml.Unmarshal(content, deployment)
	if err != nil {
		log.Error("unmarshal deployment falied")
		return
	}
	log.Debug("deployment is %+v", *deployment)

	if deployment.Metadata.NameSpace == "" {
		deployment.Metadata.NameSpace = "default"
	}

	deployment.Metadata.UUID = uuid.NewString()

	byteArr, err := json.Marshal(*deployment)

	if err != nil {
		log.Error("Error json marshal deployment")
		return
	}

	URL := config.GetUrlPrefix() + config.DeploymentsURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, deployment.Metadata.NameSpace, -1)

	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("error http post deployment: %s", err.Error())
		return
	}
	log.Info("apply deployment successed")
}

func applyHPAHandler(content []byte) {
	log.Info("creating or updating HPA")

	hpa := &api.HPA{}
	err := yaml.Unmarshal(content, hpa)
	if err != nil {
		log.Error("Error yaml unmarshal hpa %s", err.Error())
		return
	}

	byteArr, err := json.Marshal(*hpa)
	if err != nil {
		log.Error("Error json marshal hpa")
		return
	}

	URL := config.GetUrlPrefix() + config.HPAsURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, hpa.Metadata.NameSpace, -1)
	err = httputil.Post(URL, byteArr)

	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply hpa successed")
}

func applyDNSHandler(content []byte) {
	log.Info("creating or updating DNS")
	dns := &api.DNS{}
	err := yaml.Unmarshal(content, dns)
	if err != nil {
		log.Error("Error yaml unmarshal DNS")
		return
	}

	byteArr, err := json.Marshal(*dns)
	if err != nil {
		log.Error("Error json marshal DNS")
		return
	}

	URL := config.GetUrlPrefix() + config.DNSsURL
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply DNS successed")
}

func applyPVHandler(content []byte) {
	log.Info("creating or updating PV")
	pv := &api.PV{}
	err := yaml.Unmarshal(content, pv)
	if err != nil {
		log.Error("Error yaml unmarshal PV")
		return
	}
	byteArr, err := json.Marshal(*pv)
	if err != nil {
		log.Error("Error json marshal PV")
		return
	}
	URL := config.GetUrlPrefix() + config.PersistentVolumesURL
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply PV successed")
}

func applyPVCHandler(content []byte) {
	log.Info("creating or updating PV")
	pvc := &api.PVC{}
	err := yaml.Unmarshal(content, pvc)
	if err != nil {
		log.Error("Error yaml unmarshal PVC")
		return
	}
	byteArr, err := json.Marshal(*pvc)
	if err != nil {
		log.Error("Error json marshal PVC")
		return
	}
	URL := config.GetUrlPrefix() + config.PersistentVolumeClaimsURL
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply PVC successed")
}

func applyTriggerHandler(content []byte) {
	log.Info("creating or updating trigger")
	trigger := &api.Trigger{}
	err := yaml.Unmarshal(content, trigger)
	if err != nil {
		log.Error("Error yaml unmarshal trigger")
		return
	}
	byteArr, err := json.Marshal(*trigger)
	if err != nil {
		log.Error("Error json marshal trigger")
		return
	}
	URL := config.GetUrlPrefix() + config.TriggersURL
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply trigger successed")
}

func applyNodeHandler(content []byte) {
	log.Info("creating or updating node")
	node := &api.Node{}
	err := yaml.Unmarshal(content, node)
	if err != nil {
		log.Error("Error yaml unmarshal node")
		return
	}
	byteArr, err := json.Marshal(*node)
	if err != nil {
		log.Error("Error json marshal node")
		return
	}
	URL := config.GetUrlPrefix() + config.NodesURL
	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("Error http post: %s", err.Error())
		return
	}
	log.Info("apply node successed")
}

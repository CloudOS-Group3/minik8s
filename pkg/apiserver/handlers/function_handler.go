package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	function_manager "minik8s/pkg/serverless/function"
	"minik8s/pkg/serverless/function/function_util"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"os/exec"
)

func AddFunction(context *gin.Context) {
	// Add function
	log.Info("Add function")
	var function api.Function
	if err := context.ShouldBind(&function); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	// Build image

	err := function_manager.CreateImageFromFunction(&function)
	if err != nil {
		log.Error("Error create image: %s", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// Add function to etcd
	log.Info("Add function %s", function.Metadata.Name)

	URL := config.FunctionPath + function.Metadata.NameSpace + "/" + function.Metadata.Name
	functionByteArr, err := json.Marshal(function)
	if err != nil {
		log.Error("Error marshal function: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(functionByteArr))
}

func GetFunction(context *gin.Context) {
	// Get function
	log.Info("Get function")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	URL := config.FunctionPath + namespace + "/" + name
	function := etcdClient.GetEtcdPair(URL)
	var function_ api.Function
	if len(function) == 0 {
		log.Info("Function %s not found", name)
	} else {
		err := json.Unmarshal([]byte(function), &function_)
		if err != nil {
			log.Error("Error unmarshalling function json %v", err)
			return
		}
	}
	byteArr, err := json.Marshal(function_)
	if err != nil {
		log.Error("Error marshal function: %s", err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func GetFunctions(context *gin.Context) {
	log.Info("received get pods request")

	URL := config.FunctionPath
	functions := etcdClient.PrefixGet(URL)

	jsonString := stringutil.EtcdResEntryToJSON(functions)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}
func UpdateFunction(context *gin.Context) {
	// Update function
	log.Info("Update function")
	var function api.Function
	if err := context.ShouldBind(&function); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	URL := config.FunctionPath + function.Metadata.NameSpace + "/" + function.Metadata.Name
	functionByteArr, err := json.Marshal(function)
	if err != nil {
		log.Error("Error marshal function: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(functionByteArr))
}

func DeleteFunction(context *gin.Context) {
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	// Step1 : check running pod
	// get all related pods
	podname_prefix := function_util.GeneratePodName(name, namespace)
	URL := config.EtcdPodPath
	URL = URL + "/" + podname_prefix
	pods := etcdClient.PrefixGet(URL)
	if len(pods) != 0 {
		log.Error("Pods are still running, please delete them first")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// Step2: delete registry image
	// And cannot invoke function after delete
	cmd := exec.Command("docker", "rmi",
		config.Remotehost+":"+function_util.RegistryPort+"/"+function_util.GetImageName(name, namespace))
	//log.Info("cmd: %s", cmd.String())
	output, err := cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		log.Error("Error delete function image in master registry: %s", err.Error())
	}
	// The image was built on master, so delete it
	cmd = exec.Command("docker", "rmi", function_util.GetImageName(name, namespace))
	//log.Info("cmd: %s", cmd.String())
	output, err = cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		log.Error("Error delete function image in master registry: %s", err.Error())
	}

	// Step3: delete etcd key
	URL = config.FunctionPath + namespace + "/" + name
	etcdClient.DeleteEtcdPair(URL)

	// Step4: let all kubelet send local image
	msg := msg_type.DeleteImageMsg{
		ImageName: config.Remotehost + ":" + function_util.RegistryPort + "/" + function_util.GetImageName(name, namespace),
		Namespace: namespace,
	}
	msg_json, _ := json.Marshal(msg)
	publisher.Publish(msg_type.DeleteImageTopic, string(msg_json))
}

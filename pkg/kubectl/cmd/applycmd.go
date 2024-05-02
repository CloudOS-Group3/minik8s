package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
	"os"
	"strings"

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
	yaml.Unmarshal(content, resource)
	if resource["kind"] == "" {
		log.Error("kind field is empty")
		return ""
	}
	return resource["kind"].(string)
}

func applyCmdHandler(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("file")
	if err != nil {
		fmt.Println("Error getting flags:", err)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
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
	case "ReplicaSet":
		applyReplicaSetHandler(content)
	case "StatefulSet":
		applyStatefulSetHandler(content)
	default:
		fmt.Println("Unknown resource kind")
	}

}

func applyPodHandler(content []byte) {
	log.Info("Creating or updating pod")
	log.Debug("data is %v", content)
	pod := &api.Pod{}
	err := yaml.Unmarshal(content, pod)
	if err != nil {
		log.Error("error marshal yaml")
		return
	}
	log.Debug("%+v\n", pod)

	path := strings.Replace(config.PodsURL, config.NamespacePlaceholder, "default", -1)
	byteArr, err := json.Marshal(pod)
	log.Debug("path = %v", path)
	URL := config.GetUrlPrefix() + path
	response, err := http.Post(URL, config.JsonContent, bytes.NewBuffer(byteArr))
	if err != nil {
		log.Error("error http post")
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Warn("status code not ok")
		return
	}

	log.Info("apply pod successed")
}

func applyServiceHandler(content []byte) {
	log.Info("creating or updating service")
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
	path := strings.Replace(config.DeploymentsURL, config.NamespacePlaceholder, "default", -1)

	byteArr, err := json.Marshal(deployment)
	log.Debug("path = %+v", path)
	URL := config.GetUrlPrefix() + path

	response, err := http.Post(URL, config.JsonContent, bytes.NewBuffer(byteArr))
	if err != nil {
		log.Error("error http post deployment")
		return
	}
	if response.StatusCode != http.StatusOK {
		log.Warn("status code not ok")
		return
	}
	log.Info("apply deployment successed")
}

func applyReplicaSetHandler(content []byte) {
	log.Info("creating or updating replicaset")
}

func applyStatefulSetHandler(content []byte) {
	log.Info("creating or updating statefulset")
}

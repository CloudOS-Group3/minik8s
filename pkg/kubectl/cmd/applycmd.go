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

type Resource struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
	} `yaml:"spec"`
}

func ApplyCmd() *cobra.Command {

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "apply a yaml file to create a resource",
		Run:   applyCmdHandler,
	}

	applyCmd.Flags().StringP("file", "f", "", "specify a file name")

	return applyCmd
}

func parseYamlFileToResource(file *os.File) *Resource {
	decoder := yaml.NewDecoder(file)
	resource := &Resource{}

	err := decoder.Decode(resource)
	if err != nil {
		log.Error("Error decoding yaml: %s", err.Error())
		return nil
	}
	log.Debug("Decode yaml successfully, resource:%+v\n", resource)
	return resource
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

	resource := parseYamlFileToResource(file)

	fd, err := os.Open(path)

	if err != nil {
		log.Error("error open file")
		return
	}

	defer fd.Close()

	data, err := io.ReadAll(fd)

	if err != nil {
		log.Error("error read all data")
	}

	switch resource.Kind {
	case "Pod":
		applyPodHandler(data)
	case "Service":
		applyServiceHandler(data)
	case "Deployment":
		applyDeploymentHandler(data)
	case "ReplicaSet":
		applyReplicaSetHandler(data)
	case "StatefulSet":
		applyStatefulSetHandler(data)
	default:
		fmt.Println("Unknown resource kind")
	}

}

func applyPodHandler(data []byte) {
	log.Info("Creating or updating pod")
	log.Debug("data is %v", data)
	pod := &api.Pod{}
	err := yaml.Unmarshal(data, pod)
	if err != nil {
		log.Error("error marshal yaml")
		return
	}
	log.Debug("%+v\n", pod)

	URL := strings.Replace(config.PodsURL, config.NamespacePlaceholder, "default", -1)
	json, err := json.Marshal(pod)
	log.Debug("URL = %v", URL)
	response, err := http.Post(config.GetUrlPrefix()+URL, "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Error("error http post")
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Warn("post may have failed, because the status code is not ok")
		return
	}

	log.Info("apply pod ok")
}

func applyServiceHandler(data []byte) {
	fmt.Println("creating or updating service")
}

func applyDeploymentHandler(date []byte) {
	fmt.Println("creating or updating deployment")
}

func applyReplicaSetHandler(data []byte) {
	fmt.Println("creating or updating replicaset")
}

func applyStatefulSetHandler(data []byte) {
	fmt.Println("creating or updating statefulset")
}

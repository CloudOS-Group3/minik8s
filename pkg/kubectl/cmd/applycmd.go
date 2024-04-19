package cmd

import (
	"fmt"
	"os"

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
				App string `yaml:"nginx"`
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

	decoder := yaml.NewDecoder(file)

	resource := &Resource{}
	err = decoder.Decode(resource)
	if err != nil {
		fmt.Println("Error decoding yaml:", err)
	}

	fmt.Printf("Decode yaml successfully, resource:%+v\n", resource)

	switch resource.Kind {
	case "Pod":
		applyPodHandler(resource)
	case "Service":
		applyServiceHandler(resource)
	case "Deployment":
		applyDeploymentHandler(resource)
	case "ReplicaSet":
		applyReplicaSetHandler(resource)
	case "StatefulSet":
		applyStatefulSetHandler(resource)
	default:
		fmt.Println("Unknown resource kind")
	}

}

func applyPodHandler(resource *Resource) {
	fmt.Println("creating or updating pod")
}

func applyServiceHandler(resource *Resource) {
	fmt.Println("creating or updating service")
}

func applyDeploymentHandler(resource *Resource) {
	fmt.Println("creating or updating deployment")
}

func applyReplicaSetHandler(resource *Resource) {
	fmt.Println("creating or updating replicaset")
}

func applyStatefulSetHandler(resource *Resource) {
	fmt.Println("creating or updating statefulset")
}

package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestParseYamlFileToResource(t *testing.T) {

	file, err := os.Open("test.yaml")
	if err != nil {
		fmt.Println("err opening yaml file")
	}

	resource := parseYamlFileToResource(file)

	if resource.ApiVersion != "v1" ||
		resource.Kind != "Pod" ||
		resource.Metadata.Name != "nginx" ||
		resource.Spec.Replicas != 3 ||
		resource.Spec.Selector.MatchLabels.App != "nginx" {
		t.Errorf("err parse yaml file\n")
	}
}

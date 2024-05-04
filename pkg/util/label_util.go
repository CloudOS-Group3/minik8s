package util

import "strings"

func ConvertLabelToString(label map[string]string) string {
	labelString := ""
	for key, value := range label {
		labelString = labelString + key + "=" + value + ","
	}
	return labelString
}

func GetUniqueName(namespace string, name string) string {
	return namespace + "/" + name
}

func GetNamespaceAndName(uniqueName string) (namespace string, name string) {
	s := strings.Split(uniqueName, "/")
	return s[0], s[1]
}

func IsLabelEqual(label1 map[string]string, label2 map[string]string) bool {
	if len(label1) != len(label2) {
		return false
	}
	for key, value := range label1 {
		if label2[key] != value {
			return false
		}
	}
	return true
}

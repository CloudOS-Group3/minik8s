package util

func ConvertLabelToString(label map[string]string) string {
	labelString := ""
	for key, value := range label {
		labelString = labelString + key + "=" + value + ","
	}
	return labelString
}

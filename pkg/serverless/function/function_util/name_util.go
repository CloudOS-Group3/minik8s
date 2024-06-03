package function_util

func GeneratePodName(functionName string, namespace string) string {
	return namespace + "-" + functionName + "-functionPod"
}
func GetFunctionFilePath(functionName string, namespace string) string {
	return "~/function/" + "/" + namespace + "/" + functionName + "/"
}

func GetImageName(functionName string, namespace string) string {
	return "function-" + namespace + "-" + functionName
}

package function

import (
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/serverless/function/function_util"
	"minik8s/util/log"
)

// This file handles:
// 1. create pod from function

func CreatePodFromFunction(function *api.Function) *api.Pod {
	if function.Language == "python" {
		return CreatePythonPod(function)
	} else {
		log.Error("We only support python function now")
		return nil
	}
}

func CreatePythonPod(function *api.Function) *api.Pod {
	log.Info("Create python pod")
	imageName, err := function_util.CreateImage(function)
	if err != nil {
		log.Error("error create image: %s", err.Error())
		return nil
	}
	pod := &api.Pod{
		Metadata: api.ObjectMeta{
			Name:      function_util.GeneratePodName(function.Metadata.Name),
			NameSpace: function.Metadata.NameSpace,
			UUID:      uuid.NewString(),
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:            "python-" + function.Metadata.NameSpace + "-" + function.Metadata.Name,
					Image:           imageName,
					ImagePullPolicy: api.PullFromRegistry,
				},
			},
		},
	}
	log.Debug("%+v\n", pod)

	//URL := config.GetUrlPrefix() + config.PodsURL
	//URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
	//byteArr, err := json.Marshal(*pod)
	//if err != nil {
	//	log.Error("error marshal yaml")
	//	return nil
	//}
	//
	//err = httputil.Post(URL, byteArr)
	//
	//if err != nil {
	//	log.Error("error http post: %s", err.Error())
	//	return nil
	//}
	//pod_manager.CreatePod(pod)
	return pod
}

func DeleteFunction(name string, namespace string) {
	// Delete function
	log.Info("Delete function")
	// Step 1: delete all pods(replicas)

	// Step 2: delete images
	err := function_util.DeleteFunctionImage(name, namespace)
	if err != nil {
		log.Error("error delete function image: %s", err.Error())
		return
	}
}

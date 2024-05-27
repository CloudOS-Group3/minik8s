package function

import (
	"github.com/google/uuid"
	"minik8s/pkg/api"
	pod_manager "minik8s/pkg/kubelet/pod"
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
	err := function_util.CreateImage(function)
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
					Name:  "python",
					Image: function_util.GetImageName(function.Metadata.Name, function.Metadata.NameSpace),
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
	pod_manager.CreatePod(pod)
	return pod
}

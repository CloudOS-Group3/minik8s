package function

import (
	"errors"
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/serverless/function/function_util"
	"minik8s/util/log"
)

// This file handles:
// 1. create pod from function
func CreateImageFromFunction(function *api.Function) error {
	if function.Language == "python" {
		imageName, err := function_util.CreateImage(function)
		if err != nil {
			log.Error("error create image: %s", err.Error())
			return err
		}
		log.Info("Create image %s", imageName)
	} else {
		//log.Error("We only support python function now")
		return errors.New("We only support python function now")
	}
	return nil
}

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
	imageName := config.Remotehost + ":" + function_util.RegistryPort + "/" + function_util.GetImageName(function.Metadata.Name, function.Metadata.NameSpace)
	pod := &api.Pod{
		Metadata: api.ObjectMeta{
			Name:      function_util.GeneratePodName(function.Metadata.Name, function.Metadata.NameSpace),
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
	return pod
}

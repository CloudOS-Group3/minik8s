package gpu

import (
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/serverless/function/function_util"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"os/exec"
)

const GPUNamespace = "GPUJob"

func GetGPUPodName(gpu_config *api.GPUJob) string {
	return gpu_config.Metadata.Name + "-" + gpu_config.Metadata.UUID
}

func CreateGPUPod(gpu_config *api.GPUJob) *api.Pod {
	log.Info("Create gpu pod")
	imageName := config.Remotehost + ":" + function_util.RegistryPort + "/" + function_util.GetImageName(gpu_config.Metadata.Name, GPUNamespace)
	pod := &api.Pod{
		Metadata: api.ObjectMeta{
			Name:      GetGPUPodName(gpu_config),
			NameSpace: GPUNamespace,
			UUID:      uuid.NewString(),
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:            gpu_config.Metadata.Name + stringutil.GenerateRandomString(4),
					Image:           imageName,
					ImagePullPolicy: api.PullFromRegistry,
					VolumeMounts: []api.VolumeMount{
						{
							Name:      "shared-data",
							MountPath: "/src",
						},
					},
				},
			},
			Volumes: []api.Volume{
				{
					Name:     "shared-data",
					HostPath: gpu_config.SourcePath,
				},
			},
		},
	}
	log.Debug("%+v\n", pod)
	return pod
}

func CreateGpuImage(gpu_config *api.GPUJob) error {
	// Step 1: Build Image
	// docker build --build-arg SOURCE_DIR=/path/to/source -t my-python-app .
	cmd := exec.Command("docker", "build",
		"--build-arg", "job_name="+gpu_config.Metadata.Name+"-"+gpu_config.Metadata.UUID,
		"--build-arg", "partition="+gpu_config.Args["partition"],
		"--build-arg", "N="+gpu_config.Args["N"],
		"--build-arg", "ntasks_per_node="+gpu_config.Args["ntasks_per_node"],
		"--build-arg", "cpus_per_task="+gpu_config.Args["cpus_per_task"],
		"--build-arg", "gres="+gpu_config.Args["gres"],
		"-t",
		function_util.GetImageName(gpu_config.Metadata.Name, GPUNamespace), "/root/minik8s/pkg/gpu/image/")
	output, err := cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		log.Error("failed to run build: %s", string(output))
		return err
	}

	// Step 2: Tag Image
	// docker tag myimage:latest localhost:5000/myimage:latest
	cmd = exec.Command("docker", "tag", function_util.GetImageName(gpu_config.Metadata.Name, GPUNamespace),
		config.Remotehost+":"+function_util.RegistryPort+"/"+function_util.GetImageName(gpu_config.Metadata.Name, GPUNamespace))
	log.Info("cmd: %s", cmd.String())
	output, err = cmd.CombinedOutput()
	//log.Info("output: %s", string(output))
	if err != nil {
		log.Error("failed to run tag: %s", string(output))
		return err
	}

	// Step 3: Push Image
	// docker push localhost:5000/myimage:latest
	cmd = exec.Command("docker", "push",
		config.Remotehost+":"+function_util.RegistryPort+"/"+function_util.GetImageName(gpu_config.Metadata.Name, GPUNamespace))
	//log.Info("cmd: %s", cmd.String())
	output, err = cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		log.Error("failed to run push: %s", string(output))
		return err
	}
	return nil
}

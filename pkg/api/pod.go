package api

import (
	"time"
)

type Pod struct {
	APIVersion string     `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string     `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       PodSpec    `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     PodStatus  `json:"status" yaml:"status"`
}

type PodSpec struct {
	// NodeSelector: a map of labels.
	NodeSelector map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	NodeName     string            `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`
	Containers   []Container       `json:"containers,omitempty" yaml:"containers,omitempty"`
	Volumes      []Volume          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}

type PodStatus struct {
	Conditions        []PodCondition    `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty" yaml:"containerStatuses,omitempty"`
	HostIP            string            `json:"hostIP,omitempty" yaml:"hostIP,omitempty"`
	Message           string            `json:"message,omitempty" yaml:"message,omitempty"`
	Phase             string            `json:"phase,omitempty" yaml:"phase,omitempty"`
	PodIP             string            `json:"podIP,omitempty" yaml:"podIP,omitempty"`
	PauseId           string            `json:"pauseId,omitempty" yaml:"pauseId,omitempty"`
	StartTime         time.Time         `json:"startTime,omitempty" yaml:"startTime,omitempty"`
	Metrics           PodMetrics        `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	CPUPercentage     float64           `json:"cpuPercentage,omitempty" yaml:"cpuPercentage,omitempty"`
	MemoryPercentage  float64           `json:"memoryPercentage,omitempty" yaml:"memoryPercentage,omitempty"`
}

type PodMetrics struct {
	CpuUsage         float64
	MemoryUsage      float64
	ContainerMetrics []ContainerMetrics
}

type ContainerMetrics struct {
	CpuUsage      float64 `protobuf:"fixed64,1"`
	MemoryUsage   float64 `protobuf:"fixed64,2"`
	ProcessStatus string  `protobuf:"bytes,3"`
	//Running ProcessStatus = "running"
	//Created ProcessStatus = "created"
	//Stopped ProcessStatus = "stopped"
	//Paused ProcessStatus = "paused"
	//Pausing ProcessStatus = "pausing"
	//Unknown ProcessStatus = "unknown"
}

type PodCondition struct {
	LastProbeTime      time.Time `json:"lastProbeTime,omitempty" yaml:"lastProbeTime,omitempty"`
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty" yaml:"lastTransitionTime,omitempty"`
	Message            string    `json:"message,omitempty" yaml:"message,omitempty"`
	Reason             string    `json:"reason,omitempty" yaml:"reason,omitempty"`
	Status             string    `json:"status,omitempty" yaml:"status,omitempty"`
	Type               string    `json:"type,omitempty" yaml:"type,omitempty"`
}

type ContainerStatus struct {
	ContainerID string `json:"containerID,omitempty" yaml:"containerID,omitempty"`
	Image       string `json:"image,omitempty" yaml:"image,omitempty"`
	ImageID     string `json:"imageID,omitempty" yaml:"imageID,omitempty"`
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Ready       bool   `json:"ready,omitempty" yaml:"ready,omitempty"`
}

type Container struct {
	Name            string               `json:"name,omitempty" yaml:"name,omitempty"`
	Ports           []ContainerPort      `json:"ports,omitempty" yaml:"ports,omitempty"`
	Args            []string             `json:"args,omitempty" yaml:"args,omitempty"`
	Command         []string             `json:"command,omitempty" yaml:"command,omitempty"`
	Env             []EnvVar             `json:"env,omitempty" yaml:"env,omitempty"`
	Image           string               `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy string               `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	Resources       ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`
	VolumeMounts    []VolumeMount        `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
}

type ContainerPort struct {
	ContainerPort int32  `json:"containerPort,omitempty" yaml:"containerPort,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type ResourceRequirements struct {
	Claims   []ResourceClaim `json:"claims,omitempty" yaml:"claims,omitempty"`
	Limits   ComputeResource `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests ComputeResource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type ResourceClaim struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ComputeResource struct {
	Cpu     string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory  string `json:"memory,omitempty" yaml:"memory,omitempty"`
	Storage string `json:"storage,omitempty" yaml:"storage,omitempty"`
}

type VolumeMount struct {
	MountPath string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type Volume struct {
	Name                  string                            `json:"name,omitempty" yaml:"name,omitempty"`
	HostPath              string                            `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
	NFS                   NFSVolumeSource                   `json:"nfs,omitempty" yaml:"nfs,omitempty"`
	PersistentVolumeClaim PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty" yaml:"persistentVolumeClaim,omitempty"`
}

type PersistentVolumeClaimVolumeSource struct {
	ClaimName string `json:"claimName,omitempty" yaml:"claimName,omitempty"`
}

const (
	PullPolicyAlways       = "Always"
	PullPolicyIfNotPresent = "IfNotPresent"
	PullPolicyNever        = "Never"
)

type PodPhase string

const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
	PodUnknown   PodPhase = "Unknown"
)

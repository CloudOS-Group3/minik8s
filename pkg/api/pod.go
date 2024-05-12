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
	Conditions        []PodCondition    `json:"conditions" yaml:"conditions"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses" yaml:"containerStatuses"`
	HostIP            string            `json:"hostIP" yaml:"hostIP"`
	Message           string            `json:"message" yaml:"message"`
	Phase             string            `json:"phase" yaml:"phase"`
	PodIP             string            `json:"podIP" yaml:"podIP"`
	StartTime         time.Time         `json:"startTime" yaml:"startTime"`
	CPUPercentage     float64           `json:"cpuPercentage" yaml:"cpuPercentage"`
	MemoryPercentage  float64           `json:"memoryPercentage" yaml:"memoryPercentage"`
}

type PodCondition struct {
	LastProbeTime      time.Time `json:"lastProbeTime" yaml:"lastProbeTime"`
	LastTransitionTime time.Time `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	Message            string    `json:"message" yaml:"message"`
	Reason             string    `json:"reason" yaml:"reason"`
	Status             string    `json:"status" yaml:"status"`
	Type               string    `json:"type" yaml:"type"`
}

type ContainerStatus struct {
	ContainerID string `json:"containerID" yaml:"containerID"`
	Image       string `json:"image" yaml:"image"`
	ImageID     string `json:"imageID" yaml:"imageID"`
	Name        string `json:"name" yaml:"name"`
	Ready       bool   `json:"ready" yaml:"ready"`
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
	Limits   ComputeResource `json:"limits,omitempty" yaml:"limits,omitempty"`
	Requests ComputeResource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type ComputeResource struct {
	Cpu    string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type VolumeMount struct {
	MountPath string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type Volume struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
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

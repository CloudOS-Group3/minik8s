package api

import (
	"time"
)

// Pod is a collection of containers, created by clients and scheduled onto hosts.
/*
	This file defines the Pod struct and its associated types.
	Examples of a yaml / json can be found in exampleFile/
*/

type Pod struct {
	// APIVersion defines the versioned schema of this representation of an object.
	// eg: "v1"
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`

	// Kind is a string value representing the REST resource this object represents.
	// Values: "Pod"
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	// Metadata is data like name, UID, etc.
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the pod.
	Spec PodSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status PodStatus `json:"status" yaml:"status"`
}

type PodSpec struct {
	// NodeSelector must match labels in the pod template.
	NodeSelector map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`

	// NodeName is a request to schedule this pod onto a specific node.
	NodeName string `json:"nodeName,omitempty" yaml:"nodeName,omitempty"`

	// Containers is a list of containers belonging to the pod.
	Containers []Container `json:"containers,omitempty" yaml:"containers,omitempty"`

	// Volumes is a list of volumes mounted by containers.
	Volumes []Volume `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}

type PodStatus struct {
	Conditions        []PodCondition    `json:"conditions" yaml:"conditions"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses" yaml:"containerStatuses"`
	HostIP            string            `json:"hostIP" yaml:"hostIP"`
	Message           string            `json:"message" yaml:"message"`
	Phase             string            `json:"phase" yaml:"phase"`
	PodIP             string            `json:"podIP" yaml:"podIP"`
	StartTime         time.Time         `json:"startTime" yaml:"startTime"`
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

// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#container-v1-core
type Container struct {
	// Name must be unique within a pod.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Ports is to expose from the container.
	Ports []ContainerPort `json:"ports,omitempty" yaml:"ports,omitempty"`

	// Args is the arguments to the command.
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`

	// Command is the command to run in the container.
	Command []string `json:"command,omitempty" yaml:"command,omitempty"`

	// Env is a list of environment variables.
	Env []EnvVar `json:"env,omitempty" yaml:"env,omitempty"`

	// Image is the docker image to run.
	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	// ImagePullPolicy: Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified.
	ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`

	// Resources requirements.
	Resources ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`

	// VolumeMounts is a list of volumes mounted by the container.
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
}

type ContainerPort struct {
	// ContainerPort: 0 < x < 65536
	ContainerPort int32 `json:"containerPort,omitempty" yaml:"containerPort,omitempty"`

	// Name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Protocol must be UDP, TCP, or SCTP. Defaults to "TCP".
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

type EnvVar struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	Limits ComputeResource `json:"limits,omitempty" yaml:"limits,omitempty"`

	// Requests describes the minimum amount of compute resources required.
	Requests ComputeResource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

type ComputeResource struct {
	Cpu    string `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type VolumeMount struct {
	// MountPath: Must not contain ':'
	MountPath string `json:"mountPath,omitempty" yaml:"mountPath,omitempty"`

	// Name must match the name of a volume.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// ReadOnly: read-only if true, read-write otherwise. Defaults to false.
	ReadOnly bool `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
}

type Volume struct {
	// Name must be unique within the pod.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// HostPath, a directory on the host.
	HostPath string `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
}

const (
	PullPolicyAlways       = "Always"
	PullPolicyIfNotPresent = "IfNotPresent"
	PullPolicyNever        = "Never"
)

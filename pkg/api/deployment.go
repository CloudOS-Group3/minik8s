package api

type Deployment struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta       `json:"metadata" yaml:"metadata"`
	Spec       DeploymentSpec   `json:"spec" yaml:"spec"`
	Status     DeploymentStatus `json:"status" yaml:"status"`
}

type DeploymentSpec struct {
	Replicas int             `json:"replicas" yaml:"replicas"`
	Selector LabelSelector   `json:"selector" yaml:"selector"`
	Template PodTemplateSpec `json:"template" yaml:"template"`
}

type DeploymentStatus struct {
	AvailableReplicas   int `json:"availableReplicas" yaml:"availableReplicas"`
	ReadyReplicas       int `json:"readyReplicas" yaml:"readyReplicas"`
	Replicas            int `json:"replicas" yaml:"replicas"`
	UnAvailableReplicas int `json:"unavailableReplicas" yaml:"unavailableReplicas"`
	UpdatedReplicas     int `json:"updatedReplicas" yaml:"updatedReplicas"`
}

type LabelSelector struct {
	MatchLabels map[string]string `json:"matchLabels" yaml:"matchLabels"`
}

type PodTemplateSpec struct {
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     PodSpec    `json:"spec" yaml:"spec"`
}

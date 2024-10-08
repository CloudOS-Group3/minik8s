package api

import "time"

type Node struct {
	APIVersion string     `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	Kind       string     `yaml:"kind,omitempty" json:"kind,omitempty"`
	Metadata   ObjectMeta `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Spec       NodeSpec   `yaml:"spec,omitempty" json:"spec,omitempty"`
	Status     NodeStatus `yaml:"status,omitempty" json:"status,omitempty"`
}

type NodeSpec struct {
	NodeIP     string   `yaml:"nodeIP,omitempty" json:"nodeIP,omitempty"`
	PodCIDR    string   `yaml:"podCIDR,omitempty" json:"podCIDR,omitempty"`
	PodCIDRs   []string `yaml:"podCIDRs,omitempty" json:"podCIDRs,omitempty"`
	ProviderID string   `yaml:"providerID,omitempty" json:"providerID,omitempty"`
}

type NodeCondition struct {
	Status            string    `yaml:"status,omitempty" json:"status,omitempty"`
	LastHeartbeatTime time.Time `yaml:"lastHeartbeatTime,omitempty" json:"lastHeartbeatTime,omitempty"`
}

const (
	NodeReady   string = "Ready"
	NodeFailed  string = "Failed"
	NodeUnknown string = "Unknown"
)

type NodeStatus struct {
	Condition  NodeCondition `yaml:"condition,omitempty" json:"condition,omitempty"`
	PodsNumber int           `yaml:"podsNumber,omitempty" json:"podsNumber,omitempty"`
	Pods       []Pod         `yaml:"pods,omitempty" json:"pods,omitempty"`
}

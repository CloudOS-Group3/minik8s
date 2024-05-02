package api

type Node struct {
	APIVersion string     `yaml:"apiVersion" json:"apiVersion"`
	Kind       string     `yaml:"kind" json:"kind"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       NodeSpec   `yaml:"spec" json:"spec"`
	Status     NodeStatus `yaml:"status" json:"status"`
}

type NodeSpec struct {
	ExternalID string   `yaml:"externalID" json:"externalID"`
	PodCIDR    string   `yaml:"podCIDR" json:"podCIDR"`
	PodCIDRs   []string `yaml:"podCIDRs" json:"podCIDRs"`
	ProviderID string   `yaml:"providerID" json:"providerID"`
}

type NodeCondition = string

const (
	Ready   NodeCondition = "Ready"
	Failed  NodeCondition = "Failed"
	Unknown NodeCondition = "Unknown"
)

type NodeStatus struct {
	Hostname   string        `yaml:"hostname" json:"hostname"`
	Condition  NodeCondition `yaml:"condition" json:"condition"`
	PodsNumber int           `yaml:"podsNumber" json:"podsNumber"`
}

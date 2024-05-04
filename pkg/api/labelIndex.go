package api

// LabelIndex : A two-way mapping between service and pod, all by label.
// unique index by label
type LabelIndex struct {
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	// PodName is namespace / podname
	PodName []string `json:"podName,omitempty" yaml:"podName,omitempty"`
	// ServiceName is namespace / servicename
	ServiceName []string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
}

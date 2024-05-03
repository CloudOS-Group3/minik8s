package api

type Service struct {
	APIVersion string        `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string        `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   ObjectMeta    `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       ServiceSpec   `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     ServiceStatus `json:"status" yaml:"status"`
}

type ServiceSpec struct {
	// Type: how the service is exposed. only support ClusterIP now.
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Ports: a list of ports that are exposed.
	Ports []ServicePort `json:"ports,omitempty" yaml:"ports,omitempty"`

	// Selector: a map of labels that the service selects.
	Selector map[string]string `json:"selector,omitempty" yaml:"selector,omitempty"`
}

type ServicePort struct {
	// The port exposed.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// The port on each node that will be used by the service.
	TargetPort int `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`

	// The protocol: TCP(default) or UDP.
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`

	// The name of this port within the service.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type ServiceStatus struct {
	ClusterIP string     `json:"clusterIP,omitempty" yaml:"clusterIP,omitempty"`
	EndPoints []EndPoint `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
}

type EndPoint struct {
	// Don't support cross namespace search.
	NameSpace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	IP        string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Ports     []int  `json:"ports,omitempty" yaml:"ports,omitempty"`
	// EndPoint is a two-way mapping between service and pod, all by label.
	PodName     string `json:"podName,omitempty" yaml:"podName,omitempty"`
	ServiceName string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
}

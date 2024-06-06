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
	TargetPort string `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`

	// The protocol: TCP(default) or UDP.
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`

	// The name of this port within the service.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// NodePort: curl <NodeIP>:NodePort, 30000-32767
	NodePort int `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`
}

type ServiceStatus struct {
	ClusterIP string     `json:"clusterIP,omitempty" yaml:"clusterIP,omitempty"`
	EndPoints []EndPoint `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
}

// EndPoint : match a pod by label
type EndPoint struct {
	// ServicePort: bind PodIP:ContainerPort to ClusterIP:ServicePort
	ServicePort string `json:"targetPort,omitempty" yaml:"targetPort,omitempty"`
	// NodePort: bind PodIP:ContainerPort to NodeIP:NodePort
	NodePort int `json:"nodePort,omitempty" yaml:"nodePort,omitempty"`
	// IP is the IP address of the pod
	IP string `json:"ip,omitempty" yaml:"ip,omitempty"`
	// Ports is a list of ports that are exposed.
	Ports []ContainerPort `json:"ports,omitempty" yaml:"ports,omitempty"`
}

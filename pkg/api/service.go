package api

type Service struct {
	APIVersion string        `json:"apiVersion" yaml:"apiVersion"`
	Kind       string        `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta    `json:"metadata" yaml:"metadata"`
	Spec       ServiceSpec   `json:"spec" yaml:"spec"`
	Status     ServiceStatus `json:"status" yaml:"status"`
}

type ServiceSpec struct {
	Type     string            `json:"type" yaml:"type"`
	Selector map[string]string `json:"selector" yaml:"selector"`
	Ports    []ServicePort     `json:"ports" yaml:"ports"`
}

type ServiceStatus struct {
}

type ServicePort struct {
	AppProtocol string `json:"appProtocol" yaml:"appProtocol"`
	Name        string `json:"name" yaml:"name"`
	NodePort    int    `json:"nodePort" yaml:"nodePort"`
	Port        int    `json:"port" yaml:"port"`
	Protocol    string `json:"protocol" yaml:"protocol"`
	TargetPort  int    `json:"targetPort" yaml:"targetPort"`
}

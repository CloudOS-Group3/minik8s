package api

type DNS struct {
	APIVersion string    `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string    `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name       string    `json:"name,omitempty" yaml:"name,omitempty"`
	Host       string    `json:"host,omitempty" yaml:"host,omitempty"`
	Paths      []DNSPath `json:"paths,omitempty" yaml:"paths,omitempty"`
}

type DNSPath struct {
	Path             string `json:"path,omitempty" yaml:"path,omitempty"`
	ServiceNamespace string `json:"serviceNamespace,omitempty" yaml:"serviceNamespace,omitempty"`
	ServiceName      string `json:"serviceName,omitempty" yaml:"serviceName,omitempty"`
	ServiceIP        string `json:"serviceIP,omitempty" yaml:"serviceIP,omitempty"`
	ServicePort      string `json:"servicePort,omitempty" yaml:"servicePort,omitempty"`
}

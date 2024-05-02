package api

type ObjectMeta struct {
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	NameSpace       string            `json:"namespace,omitempty" yaml:"nameSpace,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	ResourceVersion string            `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

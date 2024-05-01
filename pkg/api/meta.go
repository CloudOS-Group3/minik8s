package api

type ObjectMeta struct {
	// Name is unique within a namespace to identify the object.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// NameSpace defines the space where name must be unique.
	NameSpace string `json:"namespace,omitempty" yaml:"nameSpace,omitempty"`

	// UUID is unique in time and space.
	UUID string `json:"uuid,omitempty" yaml:"uuid,omitempty"`

	// ResourceVersion is used by clients to determine when objects have changed.
	ResourceVersion string `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`

	// Labels are key value pairs that are attached to objects, such as pods.
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

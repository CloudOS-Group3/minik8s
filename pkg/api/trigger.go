package api

type Trigger struct {
	APIVersion string      `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string      `json:"kind,omitempty" yaml:"kind,omitempty"`
	Spec       TriggerSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type TriggerSpec struct {
	Type              string `json:"type,omitempty" yaml:"type,omitempty"`
	FunctionNamespace string `json:"functionNamespace,omitempty" yaml:"functionNamespace,omitempty"`
	FunctionName      string `json:"functionName,omitempty" yaml:"functionName,omitempty"`
}

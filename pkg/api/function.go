package api

type Function struct {
	// Metadata: name, namespace, uuid
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// Language: only python
	Language string `json:"language,omitempty" yaml:"language,omitempty"`
	// FilePath: the path of the function file
	FilePath string `json:"filePath,omitempty" yaml:"filePath,omitempty"`
	// Trigger: the trigger of the function
	Trigger TriggerType `json:"triggerType,omitempty" yaml:"triggerType,omitempty"`
	// Params: the parameters template
	Params []Template `json:"params,omitempty" yaml:"params,omitempty"`
	// Result: the result template
	Result []Template `json:"result,omitempty" yaml:"result,omitempty"`
}
type TriggerType struct {
	Http  bool `json:"http,omitempty" yaml:"http,omitempty"`
	Event bool `json:"event,omitempty" yaml:"event,omitempty"`
}

// Template eg. "a" int = 1
type Template struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`
	Default string `json:"default,omitempty" yaml:"default,omitempty"`
}

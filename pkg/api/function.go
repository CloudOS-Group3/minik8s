package api

type Function struct {
	// Metadata: name, namespace, uuid
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// Language: only python
	Language string `json:"language,omitempty" yaml:"language,omitempty"`
	// FilePath: the path of the function file
	FilePath string `json:"filePath,omitempty" yaml:"filePath,omitempty"`
	// ConfigPath: eg requirements.txt
	ConfigPath string `json:"configPath,omitempty" yaml:"configPath,omitempty"`
	// Trigger: the trigger of the function
	Trigger TriggerType `json:"triggerType,omitempty" yaml:"triggerType,omitempty"`
}
type TriggerType struct {
	Http  bool `json:"http,omitempty" yaml:"http,omitempty"`
	Event bool `json:"event,omitempty" yaml:"event,omitempty"`
}

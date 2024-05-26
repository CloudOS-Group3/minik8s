package api

type Function struct {
	// Metadata: name, namespace, uuid
	Metadata ObjectMeta `json:"metadata"`
	// Language: only python
	Language string `json:"language"`
	// FilePath: the path of the function file
	FilePath string `json:"filePath"`
	// ConfigPath: eg requirements.txt
	ConfigPath string `json:"configPath"`
	// Trigger: the trigger of the function
	Trigger TriggerType `json:"triggerType"`
}
type TriggerType struct {
	Http  bool `json:"http"`
	Event bool `json:"event"`
}

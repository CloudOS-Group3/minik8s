package api

type GPUJob struct {
	// Metadata: name, namespace, uuid
	Metadata ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Args map[string]string `json:"args,omitempty" yaml:"args,omitempty"`

	// SourcePath: the path of the source file
	SourcePath string `json:"sourcePath,omitempty" yaml:"sourcePath,omitempty"`

	// Result: the result of the job
	Result string `json:"result,omitempty" yaml:"result,omitempty"`

	// Status: Created, Running, Ended
	Status string `json:"status,omitempty" yaml:"status,omitempty"`

	// StartTime: the start time of the job
	StartTime string `json:"startTime,omitempty" yaml:"startTime,omitempty"`

	// EndTime: the end time of the job
	EndTime string `json:"endTime,omitempty" yaml:"endTime,omitempty"`
}

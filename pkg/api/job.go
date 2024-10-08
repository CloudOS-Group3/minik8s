package api

const (
	JOB_CREATED = "Created"
	JOB_RUNNING = "Running"
	JOB_ENDED   = "Ended"
)

type Job struct {
	// Metadata: name, namespace, uuid
	JobID      string `json:"jobID,omitempty" yaml:"jobID,omitempty"`
	Status     string `json:"status,omitempty" yaml:"status,omitempty"`
	Function   string `json:"function,omitempty" yaml:"function,omitempty"`
	Instance   Pod    `json:"instance,omitempty" yaml:"instance,omitempty"`
	Result     string `json:"result,omitempty" yaml:"result,omitempty"`
	CreateTime string `json:"createTime,omitempty" yaml:"createTime,omitempty"`
	Params     string `json:"params,omitempty" yaml:"params,omitempty"`
}

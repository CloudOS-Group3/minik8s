package api

const (
	JOB_CREATED = "CREATED"
	JOB_RUNNING = "RUNNING"
	JOB_DELETED = "DELETED"
)

type Job struct {
	// Metadata: name, namespace, uuid
	JobID      string `json:"jobID,omitempty" yaml:"jobID,omitempty"`
	Status     string `json:"status,omitempty" yaml:"status,omitempty"`
	Instance   Pod    `json:"instance,omitempty" yaml:"instance,omitempty"`
	Result     string `json:"result,omitempty" yaml:"result,omitempty"`
	CreateTime string `json:"createTime,omitempty" yaml:"createTime,omitempty"`
	Params     string `json:"params,omitempty" yaml:"params,omitempty"`
}

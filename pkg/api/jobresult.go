package api

type JobResult struct {
	UUID string `json:"uuid,omitempty"`
	Result string `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

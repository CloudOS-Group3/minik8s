package msg_type

import "minik8s/pkg/api"

const (
	// add operation
	Add = "add"
	// update operation
	Update = "update"
	// delete operation
	Delete = "delete"
)

const (
	// pod topic
	PodTopic = "endpoint_pod"
	// service topic
	ServiceTopic = "endpoint_service"
	// pod operation
	PodOperation = "opt"
	// pod new pod
	NewPod = "new_pod"
	// pod old pod
	OldPod = "old_pod"
	// service new service
	NewService = "new_service"
	// service old service
	OldService = "old_service"
)

type PodMsg struct {
	// operation type: add, update, delete
	Opt string `json:"opt"`
	// new pod
	NewPod api.Pod `json:"new_pod,omitempty"`
	// old pod
	OldPod api.Pod `json:"old_pod,omitempty"`
}

type ServiceMsg struct {
	// operation type: add, update, delete
	Opt string `json:"opt"`
	// new service
	NewService api.Service `json:"new_service,omitempty"`
	// old service
	OldService api.Service `json:"old_service,omitempty"`
}

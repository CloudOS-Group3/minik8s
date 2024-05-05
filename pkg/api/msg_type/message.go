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
	PodTopic = "pod"
	// pod operation
	PodOperation = "opt"
	// pod new pod
	NewPod = "new_pod"
	// pod old pod
	OldPod = "old_pod"
)

type PodMsg struct {
	// operation type: add, update, delete
	Opt string `json:"opt"`
	// new pod
	NewPod api.Pod `json:"new_pod"`
	// old pod
	OldPod api.Pod `json:"old_pod"`
}

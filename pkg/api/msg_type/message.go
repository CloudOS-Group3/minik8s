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
	// service topic
	ServiceTopic = "service"
	// node topic
	NodeTopic = "node"
	// DNS topic
	DNSTopic = "dns"
	// endpoint topic
	EndpointTopic = "endpoint"
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

type NodeMsg struct {
	Opt     string   `json:"opt"`
	NewNode api.Node `json:"new_node,omitempty"`
	OldNode api.Node `json:"old_node,omitempty"`
}

type DNSMsg struct {
	Opt    string  `json:"opt"`
	NewDNS api.DNS `json:"new_dns,omitempty"`
	OldDNS api.DNS `json:"old_dns,omitempty"`
}

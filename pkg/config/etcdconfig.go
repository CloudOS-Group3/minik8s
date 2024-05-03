package config

const (

	// usage: /registry/nodes/<node-name>
	EtcdNodePath = "/registry/nodes/"

	// usage: /registry/pods/<namespace>/<pod-name>
	EtcdPodPath = "/registry/pods"

	// usage: /registry/services/<namespace>/<svc-name>
	EtcdServicePath = "/registry/services/"

	// usgae: /registry/svclabels/<label-key>/<label-value>/<svc-uuid>
	EtcdServiceSelectorPath = "/registry/svclabels/"

	// usage: /registry/endpoints/namespace/<label-key>
	EndpointPath = "/registry/endpoints/"

	// usage: /registry/jobs/<namespace>/<job-name>
	EtcdJobPath = "/registry/jobs/"
)

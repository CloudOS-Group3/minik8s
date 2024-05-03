package config

// defines request urls for apiserver
const (
	TestURL = "/"

	NodesURL = "/api/v1/nodes"
	NodeURL  = "/api/v1/nodes/:name"

	PodsURL = "/api/v1/namespaces/:namespace/pods"
	PodURL  = "/api/v1/namespaces/:namespace/pods/:name"

	ServicesURL = "/api/v1/namespaces/:namespace/services"
	ServiceURL  = "/api/v1/namespaces/:namespace/services/:name"

	NamespacesURL = "/api/v1/namespaces"
	NamespaceURL  = "/api/v1/namespaces/:namespace"

	ConfigMapsURL = "/api/v1/namespaces/:namespace/configmaps"
	ConfigMapURL  = "/api/v1/namespaces/:namespace/configmaps/:name"

	PersistentVolumesURL = "/api/v1/persistentvolumes"
	PersistentVolumeURL  = "/api/v1/persistentvolumes/:name"

	PersistentVolumeClaimsURL = "/api/v1/namespaces/:namespace/persistentvolumeclaims"
	PersistentVolumeClaimURL  = "/api/v1/namespaces/:namespace/persistentvolumeclaims/:name"

	DeploymentsURL = "/apis/apps/v1/namespaces/:namespace/deployments"
	DeploymentURL  = "/apis/apps/v1/namespaces/:namespace/deployments/:name"

	ReplicaSetsURL = "/apis/apps/v1/namespaces/:namespace/replicasets"
	ReplicaSetURL  = "/apis/apps/v1/namespaces/:namespace/replicasets/:name"

	StatefulSetsURL = "/apis/apps/v1/namespaces/:namespace/statefulsets"
	StatefulSetURL  = "/apis/apps/v1/namespaces/:namespace/statefulsets/:name"

	JobsURL = "/apis/batch/v1/namespaces/:namespace/jobs"
	JobURL  = "/apis/batch/v1/namespaces/:namespace/jobs/:name"

	CronJobsURL = "/apis/batch/v1/namespaces/:namespace/cronjobs"
	CronJobURL  = "/apis/batch/v1/namespaces/:namespace/cronjobs/:name"

	HPAsURL = "/apis/v1/namespaces/:namespace/hpa"
	HPAURL = "/apis/v1/namespaces/:namespace/hpa/:name"

)

// const used to send and parse url
const (
	NamePlaceholder      = ":name"
	NamespacePlaceholder = ":namespace"
	NameParam            = "name"
	NamespaceParam       = "namespace"
)

const (
	JsonContent = "application/json"
)

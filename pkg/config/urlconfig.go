package config

// defines request urls for apiserver
const (
	TestURL = "/"

	// Nodes API
	NodesURL = "/api/v1/nodes"
	NodeURL  = "/api/v1/nodes/:name"

	// Pods API
	PodsURL = "/api/v1/namespaces/:namespace/pods"
	PodURL  = "/api/v1/namespaces/:namespace/pods/:name"

	// Services API
	ServicesAllURL = "/api/v1/services"
	ServicesURL    = "/api/v1/namespaces/:namespace/services"
	ServiceURL     = "/api/v1/namespaces/:namespace/services/:name"

	// Namespaces API
	NamespacesURL = "/api/v1/namespaces"
	NamespaceURL  = "/api/v1/namespaces/:namespace"

	// ConfigMaps API
	ConfigMapsURL = "/api/v1/namespaces/:namespace/configmaps"
	ConfigMapURL  = "/api/v1/namespaces/:namespace/configmaps/:name"

	// PersistentVolumes API
	PersistentVolumesURL = "/api/v1/persistentvolumes"
	PersistentVolumeURL  = "/api/v1/persistentvolumes/:name"

	// PersistentVolumeClaims API
	PersistentVolumeClaimsURL = "/api/v1/namespaces/:namespace/persistentvolumeclaims"
	PersistentVolumeClaimURL  = "/api/v1/namespaces/:namespace/persistentvolumeclaims/:name"

	// Deployments API
	DeploymentsURL = "/apis/apps/v1/namespaces/:namespace/deployments"
	DeploymentURL  = "/apis/apps/v1/namespaces/:namespace/deployments/:name"

	// ReplicaSets API
	ReplicaSetsURL = "/apis/apps/v1/namespaces/:namespace/replicasets"
	ReplicaSetURL  = "/apis/apps/v1/namespaces/:namespace/replicasets/:name"

	// StatefulSets API
	StatefulSetsURL = "/apis/apps/v1/namespaces/:namespace/statefulsets"
	StatefulSetURL  = "/apis/apps/v1/namespaces/:namespace/statefulsets/:name"

	// Jobs API
	JobsURL = "/apis/batch/v1/namespaces/:namespace/jobs"
	JobURL  = "/apis/batch/v1/namespaces/:namespace/jobs/:name"

	// CronJobs API
	CronJobsURL = "/apis/batch/v1/namespaces/:namespace/cronjobs"
	CronJobURL  = "/apis/batch/v1/namespaces/:namespace/cronjobs/:name"

	//Endpoints API
	EndpointURL = "/api/v1/namespaces/:namespace/endpoints/:label"

	//LabelIndex API
	LabelIndexURL = "/api/v1/labelIndex/:label"
)

// const used to send and parse url
const (
	NamePlaceholder      = ":name"
	NamespacePlaceholder = ":namespace"
	NameParam            = "name"
	NamespaceParam       = "namespace"
	LabelParam           = "label"
)

const (
	JsonContent = "application/json"
)

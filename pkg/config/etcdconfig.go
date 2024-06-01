package config

const (

	// usage: /registry/nodes/<node-name>
	EtcdNodePath = "/registry/nodes/"

	// usage: /registry/pods/<namespace>/<pod-name>
	EtcdPodPath = "/registry/pods/"

	// usage: /registry/deployments/<namespace>/<deployment-name>
	EtcdDeploymentPath = "/registry/deployment/"

	// usage: /registry/services/<namespace>/<svc-name>
	EtcdServicePath = "/registry/services/"

	// usgae: /registry/svclabels/<label-key>/<label-value>/<svc-uuid>
	EtcdServiceSelectorPath = "/registry/svclabels/"

	// usage: /registry/endpoints/namespace/<label-key>
	EndpointPath = "/registry/endpoints/"

	// usage: /registry/labelIndex/<label>
	LabelIndexPath = "/registry/labelIndex/"

	// usage: /registry/services/<namespace>/<svc-name>
	ServicePath = "/registry/services/"

	// usage: /registry/hpas/<namespace>/<hpa-name>
	EtcdHPAPath = "/registry/hpas/"

	// usage: /registry/dns/<name>
	EtcdDNSPath = "/registry/dns/"

	// usage: /registry/trigger/<function_namespace>/<function_name>
	EtcdTriggerPath = "/registry/trigger/"

	// usage: /registry/trigger/workflow/<wf_namespace>/<wf_name>
	EtcdTriggerWorkflowPath = "/registry/trigger/workflow/"

	// usage: /registry/result/trigger/<trigger_uuid>/
	TriggerResultPath = "/registry/result/trigger/"

	// usage: /trigger/<function_namespace>/<function_name>
	UserTriggerPath = "/trigger"

	// usage: /workflowtrigger/<function_namespace>/<function_name>
	UserTriggerWorkflowPath = "/workflowtrigger"

	// usage: /registry/function/<name>
	FunctionPath = "/registry/function/"

	// usage: /registry/workflow/<name>
	WorkflowPath = "/registry/workflow/"

	// usage: /registry/job/<job_id>
	EtcdJobPath = "/registry/jobs/"
)

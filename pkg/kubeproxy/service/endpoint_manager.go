package service

import (
	"minik8s/pkg/api"
	"minik8s/pkg/util"
)

func OnPodUpdate(pod *api.Pod, oldLabel map[string]string) *api.EndPoint {
	label := util.ConvertLabelToString(pod.Spec.NodeSelector)
	endpoint := GetEndpoint(label)
	if endpoint == nil {
		// add new endpoint
		endpoint = &api.EndPoint{
			NameSpace: pod.Metadata.NameSpace,
			PodName:   []string{pod.Metadata.Name},
		}
		// the label must be first appear
		// no service need to be updated
	} else {
		// update endpoint
		endpoint.PodName = append(endpoint.PodName, pod.Metadata.Name)
		// need to update service
		//for _, serviceName := range endpoint.ServiceName {
		//	// update service
		//}

	}
	if oldLabel != nil {
		oldLabelString := util.ConvertLabelToString(oldLabel)
		oldEndpoint := GetEndpoint(oldLabelString)
		if oldEndpoint != nil {
			// update old endpoint
			for i, name := range oldEndpoint.PodName {
				if name == pod.Metadata.Name {
					oldEndpoint.PodName = append(oldEndpoint.PodName[:i], oldEndpoint.PodName[i+1:]...)
					break
				}
			}
			// need to update service
			//for _, serviceName := range oldEndpoint.ServiceName {
			//	// update service
			//}

			// store the old endpoint

		}

	}
	// store the new endpoint
	return endpoint
}

func GetEndpoint(label string) *api.EndPoint {
	return nil
}

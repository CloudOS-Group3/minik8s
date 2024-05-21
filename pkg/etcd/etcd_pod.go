package etcd

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/util/log"
)

func (store Store) PutPod(value api.Pod) bool {
	// generate key
	key := ""
	if value.Metadata.NameSpace == "" {
		key = "/registry/pods/default/" + value.Metadata.Name
	} else {
		key = "/registry/pods/" + value.Metadata.NameSpace + "/" + value.Metadata.Name
	}

	// check whether the name is already exist
	res := store.GetEtcdPair(key)
	if len(res) != 0 {
		log.Info("Pod name %s already exists", key)
		return false
	}

	// marshal pod as a json
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Error("Error marshalling pod json %v", err)
		return false
	}

	// put json in etcd
	return store.PutEtcdPair(key, string(jsonValue))
}

func (store Store) GetPod(namespace, name string) (api.Pod, bool) {
	// generate key
	key := ""
	if namespace == "" {
		key = "/registry/pods/default/" + name
	} else {
		key = "/registry/pods/" + namespace + "/" + name
	}

	// get json in etcd
	res := store.GetEtcdPair(key)

	// if not exist
	if len(res) == 0 {
		log.Info("Pod %s not found", key)
		return api.Pod{}, false
	}

	// unmarshal json
	var pod api.Pod
	err := json.Unmarshal([]byte(res), &pod)
	if err != nil {
		log.Error("Error unmarshalling pod json %v", err)
		return api.Pod{}, false
	}
	return pod, true
}

func (store Store) UpdatePod(pod api.Pod) bool {

	name := pod.Metadata.Name
	namespace := pod.Metadata.NameSpace
	res, ok := store.GetPod(namespace, name)
	if !ok {
		log.Error("error get pod")
		return false
	}
	if res.Metadata.Name == "" {
		log.Error("%s pod doesn't exist", res.Metadata.Name)
		return false
	}

	jsonValue, err := json.Marshal(pod)
	if err != nil {
		log.Error("Error marshalling pod json %v", err)
		return false
	}

	key := ""
	if pod.Metadata.NameSpace == "" {
		key = "/registry/pods/default/" + pod.Metadata.Name
	} else {
		key = "/registry/pods/" + pod.Metadata.NameSpace + "/" + pod.Metadata.Name
	}

	return store.PutEtcdPair(key, string(jsonValue))
}

func (store Store) DeletePod(namespace, name string) bool {
	// generate key
	key := ""
	if namespace == "" {
		key = "/registry/pods/default/" + name
	} else {
		key = "/registry/pods/" + namespace + "/" + name
	}

	// check whether key is exist
	res := store.GetEtcdPair(key)
	if len(res) == 0 {
		log.Info("Pod %s not found", key)
		return false
	}

	return store.DeleteEtcdPair(key)
}

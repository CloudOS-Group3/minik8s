package stringutil

import (
	"fmt"
	"minik8s/pkg/etcd"
	"strings"
)

func EtcdResEntryToJSON(KVs []etcd.ResEntry) string {
	valueArray := []string{}
	for _, kv := range KVs {
		valueArray = append(valueArray, kv.Value)
	}
	jsonArray := strings.Join(valueArray, ",")
	jsonString := fmt.Sprint("[", jsonArray, "]")
	return jsonString
}
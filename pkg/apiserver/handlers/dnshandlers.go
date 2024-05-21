package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
)

func GetDNSs(context *gin.Context) {
	URL := config.EtcdDNSPath
	DNSList := etcdClient.PrefixGet(URL)
	jsonString := stringutil.EtcdResEntryToJSON(DNSList)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func AddDNS(context *gin.Context) {
	var newDNS api.DNS
	if err := context.ShouldBind(&newDNS); err != nil {
		log.Error("decode pod failed")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	EtcdPath := config.EtcdDNSPath + newDNS.Name
	res := etcdClient.GetEtcdPair(EtcdPath)
	if res != "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "duplicate DNS name",
		})
		return
	}
	jsonString, _ := json.Marshal(newDNS)
	etcdClient.PutEtcdPair(EtcdPath, string(jsonString))
	var message msg_type.DNSMsg
	message = msg_type.DNSMsg{
		Opt:    msg_type.Add,
		NewDNS: newDNS,
	}
	msgJson, _ := json.Marshal(message)
	publisher.Publish(msg_type.DNSTopic, string(msgJson))
}

func DeleteDNS(context *gin.Context) {
	name := context.Param(config.NameParam)
	EtcdPath := config.EtcdDNSPath + name
	oldDNS := etcdClient.GetEtcdPair(EtcdPath)
	if oldDNS == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "DNS not found",
		})
		return
	}
	etcdClient.DeleteEtcdPair(EtcdPath)
	var deletedDNS api.DNS
	var message msg_type.DNSMsg
	message = msg_type.DNSMsg{
		Opt:    msg_type.Delete,
		OldDNS: deletedDNS,
	}
	msgJson, _ := json.Marshal(message)
	publisher.Publish(msg_type.DNSTopic, string(msgJson))
}

func UpdateDNS(context *gin.Context) {
	var newDNS api.DNS
	if err := context.ShouldBind(&newDNS); err != nil {
		log.Error("decode pod failed")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	EtcdPath := config.EtcdDNSPath + newDNS.Name
	res := etcdClient.GetEtcdPair(EtcdPath)
	jsonString, _ := json.Marshal(newDNS)
	etcdClient.PutEtcdPair(EtcdPath, string(jsonString))
	var message msg_type.DNSMsg
	if res == "" {
		message = msg_type.DNSMsg{
			Opt:    msg_type.Add,
			NewDNS: newDNS,
		}
	} else {
		var oldDNS api.DNS
		_ = json.Unmarshal([]byte(res), &oldDNS)
		message = msg_type.DNSMsg{
			Opt:    msg_type.Update,
			OldDNS: oldDNS,
			NewDNS: newDNS,
		}
	}
	msgJson, _ := json.Marshal(message)
	publisher.Publish(msg_type.DNSTopic, string(msgJson))
}

func GetDNS(context *gin.Context) {
	name := context.Param(config.NameParam)
	EtcdPath := config.EtcdDNSPath + name
	res := etcdClient.GetEtcdPair(EtcdPath)
	if res == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "DNS not found",
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

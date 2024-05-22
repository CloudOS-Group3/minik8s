package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
)

func GetHPAs(context *gin.Context) {
	URL := config.EtcdHPAPath
	hpas := etcdClient.PrefixGet(URL)

	jsonString := stringutil.EtcdResEntryToJSON(hpas)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}
func AddHPA(context *gin.Context) {
	var newHPA api.HPA

	if err := context.ShouldBind(&newHPA); err != nil {
		log.Error("decode newHPA error")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newHPA.Metadata.UUID = uuid.NewString()

	byteArr, err := json.Marshal(newHPA)
	if err != nil {
		log.Error("marshal newHPA error")
		return
	}

	URL := config.EtcdHPAPath + "default/" + newHPA.Metadata.Name

	etcdClient.PutEtcdPair(URL, string(byteArr))
}
func GetHPA(context *gin.Context) {
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	if name == "" {
		log.Error("hpa name is empty")
		return
	}

	if namespace == "" {
		log.Error("hpa namespace is empty")
		return
	}

	key := config.EtcdHPAPath + namespace + "/" + name

	jsonString := etcdClient.GetEtcdPair(key)

	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func UpdateHPA(context *gin.Context) {

}
func DeleteHPA(context *gin.Context) {
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	if name == "" {
		log.Error("hpa name is empty")
		return
	}

	if namespace == "" {
		log.Error("hpa namespace is empty")
		return
	}

	key := config.EtcdHPAPath + namespace + "/" + name

	ok := etcdClient.DeleteEtcdPair(key)
	if !ok {
		log.Error("delete hpa error")
	}
}

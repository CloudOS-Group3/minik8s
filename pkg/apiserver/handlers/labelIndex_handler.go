package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/util"
	"minik8s/util/log"
	"net/http"
)

func GetLabelIndex(context *gin.Context) {
	label := context.Param(config.LabelParam)
	URL := config.LabelIndexPath + label
	res := etcdClient.GetEtcdPair(URL)
	log.Info("Get URL %s", URL)

	labelIndex := &api.LabelIndex{}
	if len(res) != 0 {
		err := json.Unmarshal([]byte(res), labelIndex)
		if err != nil {
			log.Error("Error unmarshalling labelIndex json %v", err)
			return
		}
	}
	byteArr, err := json.Marshal(labelIndex)
	if err != nil {
		log.Error("Error marshal labelIndex: %s", err.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}
func AddLabelIndex(context *gin.Context) {
	var labelIndex api.LabelIndex
	if err := context.ShouldBind(&labelIndex); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	labelIndexByteArr, err := json.Marshal(labelIndex)
	if err != nil {
		return
	}
	label := util.ConvertLabelToString(labelIndex.Labels)
	URL := config.LabelIndexPath + label
	etcdClient.PutEtcdPair(URL, string(labelIndexByteArr))
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func DeleteLabelIndex(context *gin.Context) {
	label := context.Param(config.LabelParam)
	URL := config.EndpointPath + label
	etcdClient.DeleteEtcdPair(URL)
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

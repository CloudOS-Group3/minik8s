package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/util"
	"net/http"
)

func GetLabelIndex(context *gin.Context) {
	label := context.Param(config.LabelParam)
	URL := config.LabelIndexPath + label
	labelIndex := etcdClient.PrefixGet(URL)

	context.JSON(http.StatusOK, gin.H{
		"data": labelIndex,
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

package handlers

import (
	// "encoding/json"
	// "minik8s/pkg/api"
	// "minik8s/pkg/config"
	// "minik8s/pkg/etcd"
	// "minik8s/pkg/kafka"
	// "minik8s/util/log"
	// "net/http"

	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDeployments(context *gin.Context) {
	log.Info("getting all deployments")

	URL := config.EtcdDeploymentPath
	deployments := etcdClient.PrefixGet(URL)

	log.Debug("all deployments are: %+v", deployments)

	jsonString := stringutil.EtcdResEntryToJSON(deployments)

	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})

}

func AddDeployment(context *gin.Context) {
	log.Info("adding a deployment")

	var newDeployment api.Deployment
	if err := context.ShouldBind(&newDeployment); err != nil {
		log.Error("decode deployment failed: %s", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"status": "wrong",
		})
		return
	}

	log.Debug("new deployment is: %+v", newDeployment)

	jsonString, err := json.Marshal(newDeployment)
	if err != nil {
		log.Error("json marshal error: %s", err.Error())
		return
	}

	URL := config.EtcdDeploymentPath + "default/" + newDeployment.Metadata.Name

	etcdClient.PutEtcdPair(URL, string(jsonString))
}

func GetDeployment(context *gin.Context) {
	log.Info("getting one deployment")

	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	if name == "" {
		log.Error("deployment name empty")
		return
	}

	if namespace == "" {
		log.Error("deployment namespace emtpy")
		return
	}

	key := config.EtcdDeploymentPath + namespace + "/" + name

	jsonString := etcdClient.GetEtcdPair(key)

	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func UpdateDeployment(context *gin.Context) {

}

func DeleteDeployment(context *gin.Context) {
	log.Info("deleting a deployment")

	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	if name == "" {
		log.Error("deployment name empty")
		return
	}

	if namespace == "" {
		log.Error("deployment namespace empty")
		return
	}

	key := config.EtcdDeploymentPath + namespace + "/" + name

	ok := etcdClient.DeleteEtcdPair(key)
	if !ok {
		log.Warn("delete deployment may have failed")
	}
}

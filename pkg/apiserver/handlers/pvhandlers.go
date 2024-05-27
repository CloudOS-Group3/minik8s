package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
	"os/exec"
)

func AddPV(context *gin.Context) {
	log.Info("Adding PV")

	var newPV api.PV
	if err := context.ShouldBindJSON(&newPV); err != nil {
		log.Error("error decoding pv")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := config.EtcdPVPath + newPV.Metadata.Name
	value, err := json.Marshal(newPV)
	if err != nil {
		log.Error("error encoding pv")
		return
	}
	etcdClient.PutEtcdPair(key, string(value))

	log.Debug("before executing 1")
	err = exec.Command("sh", "-c", fmt.Sprintf(`echo "%s *(rw,sync,no_root_squash)" >> /etc/exports`, newPV.Spec.NFS.Path)).Run()
	if err != nil {
		log.Error("error writing /etc/exports: %s", err.Error())
		return
	}
	log.Debug("before executing 2")
	err = exec.Command("sh", "-c", "exportfs", "-rv").Run()
	if err != nil {
		log.Error("error export fs: %s", err.Error())
		return
	}
	log.Debug("before executing 3")
	err = exec.Command("sh", "-c", "/etc/init.d/nfs-kernel-server", "restart").Run()
	if err != nil {
		log.Error("error restart init.d/nfs-kernel-server: %s", err.Error())
		return
	}

	log.Info("Successfully added PV")
}

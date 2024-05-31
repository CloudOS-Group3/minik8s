package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"regexp"
	"strings"
)

var UnitMap = map[string]int{
	"Ki": 1,
	"Mi": 2,
	"Gi": 3,
}

func AddPVC(context *gin.Context) {
	log.Info("received create pvc request")

	var newPVC api.PVC
	if err := context.ShouldBindJSON(&newPVC); err != nil {
		log.Error("error decode pvc from apiserver")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pvs := etcdClient.PrefixGet(config.EtcdPVPath)
	found := false
	for _, keyValuePair := range pvs {
		var pv api.PV
		if err := json.Unmarshal([]byte(keyValuePair.Key), &pv); err != nil {
			log.Error("error json unmarshalling pv")
			continue
		}
		amountRegExp := regexp.MustCompile(`^\d+`)
		amount := amountRegExp.FindString(pv.Spec.Capacity.Storage)
		requiredAmount := amountRegExp.FindString(newPVC.Spec.Resources.Requests.Storage)

		unitRegExp := regexp.MustCompile(`[a-zA-Z]+$`)
		unit := unitRegExp.FindString(pv.Spec.Capacity.Storage)
		requiredUnit := unitRegExp.FindString(newPVC.Spec.Resources.Requests.Storage)

		if UnitMap[unit] > UnitMap[requiredUnit] || (UnitMap[unit] == UnitMap[requiredUnit] && amount >= requiredAmount) {
			newPVC.Status.TargetPV = pv
			found = true
			break
		}
	}

	if !found {
		var newPV api.PV

		randomString := stringutil.GenerateRandomString(5)
		newPV.APIVersion = "v1"
		newPV.Kind = "PV"
		newPV.Metadata.UUID = uuid.NewString()
		newPV.Metadata.Name = "auto-created-pv-" + randomString
		newPV.Metadata.NameSpace = "default"
		newPV.Spec.Capacity.Storage = newPVC.Spec.Resources.Requests.Storage
		newPV.Spec.NFS.Path = config.NFSRootPath + randomString
		newPV.Spec.NFS.Server = config.NFSServerIP

		URL := config.GetUrlPrefix() + config.PersistentVolumesURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, newPV.Metadata.Name, -1)

		byteArr, err := json.Marshal(newPV)
		if err != nil {
			log.Error("error json marshalling pv")
			return
		}
		httputil.Post(URL, byteArr)

		newPVC.Status.TargetPV = newPV
	}

	key := config.ETcdPVCPath + newPVC.Metadata.Name
	value, err := json.Marshal(newPVC)
	if err != nil {
		log.Error("error json marshalling pvc")
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	etcdClient.PutEtcdPair(key, string(value))
}

func GetPVCs(context *gin.Context) {
	log.Info("received get pvcs request")

	URL := config.ETcdPVCPath
	log.Debug("before prefix get")
	pvcs := etcdClient.PrefixGet(URL)

	log.Debug("get all PVCs are: %+v", pvcs)

	jsonString := stringutil.EtcdResEntryToJSON(pvcs)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func GetPVC(context *gin.Context) {
	log.Info("getting pvc request")

	name := context.Param("name")
	URL := config.ETcdPVCPath + name
	jsonString := etcdClient.GetEtcdPair(URL)

	var pvc api.PVC
	if err := json.Unmarshal([]byte(jsonString), &pvc); err != nil {
		log.Error("error json unmarshalling pvc: %s", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	byteArr, err := json.Marshal(pvc)
	if err != nil {
		log.Error("error json marshalling pvc: %s", err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": string(byteArr)})
}

func DeletePVC(context *gin.Context) {
	log.Info("deleting pvc request")
	name := context.Param("name")
	URL := config.ETcdPVCPath + name

	etcdClient.DeleteEtcdPair(URL)
	log.Info("pvc deleted successfully")
}

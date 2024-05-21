package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	msg "minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/controller/controllers"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net"
	"net/http"
)

func GetService(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	name := context.Param(config.NameParam)
	URL := config.ServicePath + namespace + "/" + name
	svc := etcdClient.GetEtcdPair(URL)
	var service api.Service
	if len(svc) == 0 {
		log.Info("Service %s not found", name)
	} else {
		err := json.Unmarshal([]byte(svc), &service)
		if err != nil {
			log.Error("Error unmarshalling service json %v", err)
			return
		}
	}
	byteArr, err := json.Marshal(service)
	if err != nil {
		log.Error("Error marshal service: %s", err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func GetAllServices(context *gin.Context) {
	URL := config.ServicePath
	services := etcdClient.PrefixGet(URL)

	jsonString := stringutil.EtcdResEntryToJSON(services)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func GetServicesByNamespace(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	URL := config.ServicePath + namespace
	services := etcdClient.PrefixGet(URL)

	jsonString := stringutil.EtcdResEntryToJSON(services)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func AddService(context *gin.Context) {
	log.Info("AddService")
	var newService api.Service
	if err := context.ShouldBind(&newService); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// check if the service already exists
	oldService, _ := controllers.GetService(newService.Metadata.NameSpace, newService.Metadata.Name)

	// Allocate ClusterIP
	if oldService != nil && oldService.Status.ClusterIP != "" {
		newService.Status.ClusterIP = oldService.Status.ClusterIP
	} else {
		err := GenerateClusterIP(&newService)
		if err != nil {
			log.Fatal("Failed to generate ClusterIP: %v", err)
		}
	}

	serviceByteArray, err := json.Marshal(newService)
	if err != nil {
		log.Error("Failed to marshal service: %s", err.Error())
		return
	}
	log.Info("new service is: %+v", newService)
	URL := config.ServicePath + newService.Metadata.NameSpace + "/" + newService.Metadata.Name
	etcdClient.PutEtcdPair(URL, string(serviceByteArray))

	//construct message
	var message msg.ServiceMsg
	if oldService != nil {
		message = msg.ServiceMsg{
			Opt:        msg.Update,
			OldService: *oldService,
			NewService: newService,
		}
	} else {
		message = msg.ServiceMsg{
			Opt:        msg.Add,
			NewService: newService,
		}
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.ServiceTopic, string(msg_json))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func DeleteService(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	name := context.Param(config.NameParam)
	URL := config.ServicePath + namespace + "/" + name

	// check if the service already exists
	oldService, _ := controllers.GetService(namespace, name)
	if oldService == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	etcdClient.DeleteEtcdPair(URL)
	// delete ClusterIP
	allocatedIpsStr := etcdClient.GetEtcdPair(clusterIpEtcdPrefix)
	allocatedIPs := make(map[string]string)
	if len(allocatedIpsStr) != 0 {
		err := json.Unmarshal([]byte(allocatedIpsStr), &allocatedIPs)
		if err != nil {
			log.Fatal("Failed to unmarshal allocatedIPs: %s", err.Error())
		} else {
			delete(allocatedIPs, oldService.Status.ClusterIP)
			allocatedIpsByte, err := json.Marshal(allocatedIPs)
			if err != nil {
				log.Fatal("Failed to marshal allocatedIPs: %s", err.Error())
			} else {
				etcdClient.PutEtcdPair(clusterIpEtcdPrefix, string(allocatedIpsByte))
			}
		}
	}

	//construct message
	message := msg.ServiceMsg{
		Opt:        msg.Delete,
		OldService: *oldService,
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.ServiceTopic, string(msg_json))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// range: 10.96.0.0 - 10.96.255.255
const cidr = "10.96.0.0/16"
const clusterIpEtcdPrefix = "/registry/service/clusterIP"

//allocatedIPs map[string]string  // map IP -> service ns:name

func GenerateClusterIP(svc *api.Service) error {
	log.Info("GenerateClusterIP")
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	allocatedIpsStr := etcdClient.GetEtcdPair(clusterIpEtcdPrefix)
	allocatedIPs := make(map[string]string)
	if len(allocatedIpsStr) != 0 {
		err = json.Unmarshal([]byte(allocatedIpsStr), &allocatedIPs)
		if err != nil {
			return err
		}
	}

	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ipStr := ip.String()
		log.Info("ipStr: %s, %v", ipStr, allocatedIPs[ipStr])
		if allocatedIPs[ipStr] == "" {
			// find an available IP
			allocatedIPs[ipStr] = svc.Metadata.NameSpace + ":" + svc.Metadata.Name
			svc.Status.ClusterIP = ipStr
			// update etcd
			allocatedIpsByte, err := json.Marshal(allocatedIPs)
			if err != nil {
				return err
			}
			etcdClient.PutEtcdPair(clusterIpEtcdPrefix, string(allocatedIpsByte))
			return nil
		}
	}
	// no available IP
	return errors.New("no available IP addresses for ClusterIP allocation")

}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

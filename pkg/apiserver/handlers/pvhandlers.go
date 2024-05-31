package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"strings"
)

const (
	DefaultNFSServer = "192.168.3.6"
	SSHPort          = "22"
)

func execRemoteCommand(client *ssh.Client, cmd string) {
	session, err := client.NewSession()
	if err != nil {
		log.Error("Failed to create session: %v", err)
		return
	}
	defer session.Close()
	err = session.Run(cmd)
	if err != nil {
		log.Error("Failed to execute command: %v", err)
		return
	}
	log.Debug("Executed command successfully: %v", cmd)
}

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

	server := newPV.Spec.NFS.Server
	if server == "" {
		server = DefaultNFSServer
	}
	user := "root"
	password := "2024CloudOS"

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", server+":"+SSHPort, config)
	if err != nil {
		log.Error("error connecting to server")
		return
	}

	execRemoteCommand(client, fmt.Sprintf("mkdir -p %s", newPV.Spec.NFS.Path))
	execRemoteCommand(client, fmt.Sprintf("chmod 777 %s", newPV.Spec.NFS.Path))
	execRemoteCommand(client, fmt.Sprintf(`echo "%s *(rw,sync,no_root_squash)" >> /etc/exports`, newPV.Spec.NFS.Path))
	execRemoteCommand(client, "exportfs -rav")
	execRemoteCommand(client, "systemctl restart nfs-kernel-server")

	log.Info("Successfully added PV")
}

func GetPVs(context *gin.Context) {
	log.Info("received get pvs request")

	URL := config.EtcdPVPath
	log.Debug("before prefix get")
	pvs := etcdClient.PrefixGet(URL)

	log.Debug("get all PVs are: %+v", pvs)

	jsonString := stringutil.EtcdResEntryToJSON(pvs)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func GetPV(context *gin.Context) {
	log.Info("getting PV")
	name := context.Param("name")
	URL := config.EtcdPVPath + name

	jsonString := etcdClient.GetEtcdPair(URL)
	var pv api.PV
	if err := json.Unmarshal([]byte(jsonString), &pv); err != nil {
		log.Error("error decoding pv: %v", err.Error())
		return
	}

	byteArr, err := json.Marshal(pv)
	if err != nil {
		log.Error("error encoding pv: %v", err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": string(byteArr)})
	log.Info("Successfully retrieved PV")
}

func DeletePV(context *gin.Context) {
	log.Info("deleting PV")
	name := context.Param("name")
	URL := config.EtcdPVPath + name

	jsonString := etcdClient.GetEtcdPair(URL)
	var pv api.PV
	if err := json.Unmarshal([]byte(jsonString), &pv); err != nil {
		log.Error("error decoding pv: %v", err.Error())
		return
	}

	server := pv.Spec.NFS.Server
	if server == "" {
		server = DefaultNFSServer
	}
	user := "root"
	password := "2024CloudOS"

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", server+":"+SSHPort, config)
	if err != nil {
		log.Error("error connecting to server")
		return
	}
	execRemoteCommand(client, fmt.Sprintf("rm -rf %s", pv.Spec.NFS.Path))
	mountPathEscaped := strings.Replace(pv.Spec.NFS.Path, "/", "\\/", -1)
	execRemoteCommand(client, fmt.Sprintf(`sed -i "/%s *(rw,sync,no_root_squash)/d" /etc/exports`, mountPathEscaped))
	execRemoteCommand(client, "exportfs -rav")
	execRemoteCommand(client, "systemctl restart nfs-kernel-server")

	etcdClient.DeleteEtcdPair(URL)

	log.Info("Successfully deleted PV")
}

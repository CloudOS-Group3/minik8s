package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
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

	//session, err := client.NewSession()
	//if err != nil {
	//	log.Error("error creating session")
	//	return
	//}
	//
	//_, err = session.CombinedOutput(fmt.Sprintf("mkdir -p %s", newPV.Spec.NFS.Path))
	//if err != nil {
	//	log.Error("error executing command 1: %s", err.Error())
	//	return
	//}
	//_, err = session.CombinedOutput(fmt.Sprintf(`echo "%s *(rw,sync,no_root_squash)" >> /etc/exports`, newPV.Spec.NFS.Path))
	//if err != nil {
	//	log.Error("error executing command 2: %s", err.Error())
	//	return
	//}
	//_, err = session.CombinedOutput("exportfs -rav")
	//if err != nil {
	//	log.Error("error executing command 3: %s", err.Error())
	//	return
	//}
	//_, err = session.CombinedOutput("systemctl restart nfs-kernel-server")
	//if err != nil {
	//	log.Error("error executing command 4: %s", err.Error())
	//	return
	//}

	//err = exec.Command("sh", "-c", fmt.Sprintf(`echo "%s *(rw,sync,no_root_squash)" >> /etc/exports`, newPV.Spec.NFS.Path)).Run()
	//if err != nil {
	//	log.Error("error writing /etc/exports: %s", err.Error())
	//	return
	//}
	//err = exec.Command("sh", "-c", "exportfs", "-rv").Run()
	//if err != nil {
	//	log.Error("error export fs: %s", err.Error())
	//	return
	//}
	//err = exec.Command("sh", "-c", "/etc/init.d/nfs-kernel-server", "restart").Run()
	//if err != nil {
	//	log.Error("error restart init.d/nfs-kernel-server: %s", err.Error())
	//	return
	//}

	log.Info("Successfully added PV")
}

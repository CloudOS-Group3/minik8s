package config

import "strconv"

const (
	Local      = false
	localhost  = "localhost"
	Remotehost = "192.168.3.8" // IP of master
	Nodename   = "node1"
	port       = 6443
	protocol   = "http://"
)

const (
	EtcdDefaultPort = 9092
)

func GetUrlPrefix() string {
	if Local {
		return protocol + localhost + ":" + strconv.Itoa(port)
	} else {
		return protocol + Remotehost + ":" + strconv.Itoa(port)
	}
}

func GetHostAndPort() (string, int) {
	if Local {
		return localhost, port
	} else {
		return Remotehost, port
	}
}

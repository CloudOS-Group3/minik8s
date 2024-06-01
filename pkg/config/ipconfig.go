package config

import "strconv"

const (
	Local      = false
	localhost  = "localhost"
	Remotehost = "localhost" // IP of master
	port       = 6443
	protocol   = "http://"
)

var Nodename = ""

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

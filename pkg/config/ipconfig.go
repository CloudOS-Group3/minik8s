package config

import "strconv"

const (
	local      = true
	localhost  = "localhost"
	remotehost = "192.168.1.1"
	port       = 6443
	protocol   = "http://"
)

const (
	EtcdDefaultPort = 9092
)

func GetUrlPrefix() string {
	if local {
		return protocol + localhost + ":" + strconv.Itoa(port)
	} else {
		return protocol + remotehost + ":" + strconv.Itoa(port)
	}
}

func GetHostAndPort() (string, int) {
	if local {
		return localhost, port
	} else {
		return remotehost, port
	}
}

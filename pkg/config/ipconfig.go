package config

import "strconv"

const (
	local      = true
	localhost  = "localhost"
	Remotehost = "127.0.0.1"
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
		return protocol + Remotehost + ":" + strconv.Itoa(port)
	}
}

func GetHostAndPort() (string, int) {
	if local {
		return localhost, port
	} else {
		return Remotehost, port
	}
}

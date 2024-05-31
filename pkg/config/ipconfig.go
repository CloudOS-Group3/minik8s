package config

import "strconv"

const (
	Local      = true
	localhost  = "localhost"
	Remotehost = "127.0.0.1"
	port       = 6443
	protocol   = "http://"
)

const (
	EtcdDefaultPort = 9092
	NFSServerIP     = "192.168.3.6"
	NFSRootPath     = "/nfsroot/"
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

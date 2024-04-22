package config

const (
	local      = true
	localhost  = "localhost"
	remotehost = "192.168.1.1"
	port       = 6443
	protocol   = "http://"
)

func GetUrlPrefix() string {
	if local {
		return protocol + localhost + ":" + string(port)
	} else {
		return protocol + remotehost + ":" + string(port)
	}
}

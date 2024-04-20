package config

const (
	local      = true
	localhost  = "localhost"
	remotehost = "192.168.1.1"
	port       = 8080
)

func GetHostAndPort() (string, int) {
	if local {
		return localhost, port
	} else {
		return remotehost, port
	}
}
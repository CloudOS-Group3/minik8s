package host

import (
	"fmt"
	"minik8s/pkg/config"
	"os"
)

type KubeletHostManager struct {
	Hosts []string
}

func NewHostManager() *KubeletHostManager {
	return &KubeletHostManager{
		Hosts: make([]string, 0),
	}
}

func (m *KubeletHostManager) AddHost(host string) {
	m.Hosts = append(m.Hosts, host)
	m.WriteHost()
}

func (m *KubeletHostManager) RemoveHost(host string) {
	for i, h := range m.Hosts {
		if h == host {
			m.Hosts = append(m.Hosts[:i], m.Hosts[i+1:]...)
		}
	}
	m.WriteHost()
}

func (m *KubeletHostManager) WriteHost() {
	str := "127.0.0.1 localhost\n# The following lines are desirable for IPv6 capable hosts\n::1 ip6-localhost ip6-loopback\nfe00::0 ip6-localnet\nff00::0 ip6-mcastprefix\nff02::1 ip6-allnodes\nff02::2 ip6-allrouters\nff02::3 ip6-allhosts"
	for _, host := range m.Hosts {
		hostStr := fmt.Sprintf("%s %s\n", config.Remotehost, host)
		str = hostStr + str
	}
	os.WriteFile("/etc/hosts", []byte(str), 0644)
}

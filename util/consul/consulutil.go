package consul

import "github.com/hashicorp/consul/api"

var client *api.Client

func init() {
	config := api.DefaultConfig()
	config.Address = "192.168.3.8:8500"
	client, _ = api.NewClient(config)
}

func RegisterService(ID string, name string, addr string, port int) {

	registration := new(api.AgentServiceRegistration)
	registration.ID = ID
	registration.Name = name
	registration.Address = addr
	registration.Port = port
	_ = client.Agent().ServiceRegister(registration)
}

func DeRegisterService(ID string) {
	_ = client.Agent().ServiceDeregister(ID)
}

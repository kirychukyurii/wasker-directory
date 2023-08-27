package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ServiceRegistration struct {
	Id       string
	Service  string
	Host     string
	Port     int
	Protocol string
}

type ServiceConnection struct {
	Id      string
	Service string
	Host    string
	Port    int
}

type ServiceConnections []*ServiceConnection

func (c *ServiceDiscovery) GetByName(serviceName string) (ServiceConnections, error) {
	list, err := c.Client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, serviceName))
	if err != nil {
		return nil, err
	}

	result := make([]*ServiceConnection, 0, len(list))
	for _, v := range list {
		result = append(result, &ServiceConnection{
			Id:      v.ID,
			Service: v.Service,
			Host:    v.Address,
			Port:    v.Port,
		})
	}

	return result, nil
}

func (c *ServiceDiscovery) Register(register ServiceRegistration) error {
	registration := &api.AgentServiceRegistration{
		ID:      register.Id,
		Name:    register.Service,
		Port:    register.Port,
		Address: register.Host,
		Connect: &api.AgentServiceConnect{
			SidecarService: &api.AgentServiceRegistration{
				Proxy: &api.AgentServiceConnectProxyConfig{
					Config: map[string]interface{}{
						"protocol": register.Protocol,
					},
				},
			},
		},
	}

	return c.Client.Agent().ServiceRegister(registration)
}

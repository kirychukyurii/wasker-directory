package consul

import (
	"github.com/hashicorp/consul/api"

	"github.com/kirychukyurii/wasker-directory/internal/config"
	"github.com/kirychukyurii/wasker-directory/pkg/logger"
)

type ServiceDiscovery struct {
	Client *api.Client
}

func New(cfg config.Config, log logger.Logger) ServiceDiscovery {
	consulCfg := &api.Config{
		Address:    cfg.Consul.Addr(),
		Scheme:     cfg.Consul.Scheme,
		Datacenter: cfg.Consul.Datacenter,
	}

	consulClient, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("consul client")
	}

	return ServiceDiscovery{
		Client: consulClient,
	}
}

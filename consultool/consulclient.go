package consultool

import (
	"github.com/coffeehc/logger"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/hashicorp/consul/api"
)

func NewConsulClient(consulConfig *ConsulConfig) (*api.Client, base.Error) {
	if consulConfig == nil {
		consulConfig = &ConsulConfig{}
	}
	consulClient, err := api.NewClient(warpConsulConfig(consulConfig))
	if err != nil {
		logger.Error("创建 Consul Client 失败")
		return nil, base.NewError(base.ERROR_CODE_BASE_INIT_ERROR, err.Error())
	}
	return consulClient, nil
}
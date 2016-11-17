package consultool

import (
	"fmt"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/hashicorp/consul/api"
	"net"
	"strconv"
)

type consulServiceRegister struct {
	client    *api.Client
	serviceId string
	checkId   string
}

func NewConsulServiceRegister(consulClient *api.Client) (base.ServiceDiscoveryRegister, base.Error) {
	return &consulServiceRegister{
		client: consulClient,
	}, nil

}

func (this *consulServiceRegister) RegService(serviceInfo base.ServiceInfo, serviceAddr string) base.Error {
	if serviceAddr == "" {
		return base.NewError(-1, "serverAddr is nil")
	}
	this.serviceId = fmt.Sprintf("%s-%s", serviceInfo.GetServiceName(), serviceAddr)
	this.checkId = fmt.Sprintf("service:%s", this.serviceId)
	addr, port, err := net.SplitHostPort(serviceAddr)
	if err != nil {
		return base.NewError(-1, "serviceAddr is not a tcp addr")
	}
	p, _ := strconv.Atoi(port)
	registration := &api.AgentServiceRegistration{
		Name:              serviceInfo.GetServiceName(),
		Tags:              []string{serviceInfo.GetServiceTag()},
		Port:              p,
		Address:           addr,
		EnableTagOverride: true,
		Checks: api.AgentServiceChecks([]*api.AgentServiceCheck{
			{
				HTTP:     fmt.Sprintf("http://%s/health", serviceAddr),
				Interval: "10s",
			},
		}),
	}
	err = this.client.Agent().ServiceRegister(registration)
	if err != nil {
		logger.Error("注册服务失败:%s", err)
		return base.NewError(base.ERROR_CODE_BASE_SERVICE_REGISTER_ERROR, err.Error())
	}
	return nil
}

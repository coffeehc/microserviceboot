package server

import (
	"errors"

	"fmt"
	"net/http"

	"github.com/coffeehc/logger"
	"github.com/coffeehc/microserviceboot/common"
	"github.com/coffeehc/web"
)

type MicorService struct {
	server                  *web.Server
	service                 common.Service
	serviceDiscoveryRegedit ServiceDiscoveryRegister
}

func newMicorService(service common.Service, serviceDiscoveryRegedit ServiceDiscoveryRegister) (*MicorService, error) {
	serviceInfo := service.GetServiceInfo()
	if serviceInfo == nil {
		return nil, errors.New("没有指定 ServiceInfo")
	}
	webConfig := new(web.ServerConfig)
	webConfig.ServerAddr = fmt.Sprintf("%s:%d", common.GetLocalIp(), serviceInfo.GetServerPort())
	webConfig.DefaultTransport = web.Transport_Json
	switch serviceInfo.GetScheme() {
	//case common.RpcScheme_Http:
	case common.RpcScheme_https:
		webConfig.OpenTLS = true
		webConfig.CertFile, webConfig.KeyFile = serviceInfo.GetTLSCert()
	default:
		// nothing
	}
	logger.Info("ServiceName: %s", serviceInfo.GetServiceName())
	logger.Info("Version: %s", serviceInfo.GetVersion())
	logger.Info("Descriptor: %s", serviceInfo.GetDescriptor())
	return &MicorService{
		server:                  web.NewServer(webConfig),
		service:                 service,
		serviceDiscoveryRegedit: serviceDiscoveryRegedit,
	}, nil
}

func (this *MicorService) Start() error {
	serviceInfo := this.service.GetServiceInfo()
	logger.Info("MicorService start")
	err := this.regeditEndpoints()
	if err != nil {
		return err
	}
	if this.service.GetServiceInfo().GetDevModule() {
		logger.Debug("open dev module")
		apiDefineRquestHandler := buildApiDefineRquestHandler(serviceInfo)
		if apiDefineRquestHandler != nil {
			this.server.Regedit(fmt.Sprintf("/apidefine/%s.api", serviceInfo.GetServiceName()), web.GET, apiDefineRquestHandler)
		}
		this.server.AddFirstFilter("/*", web.SimpleAccessLogFilter)
	}
	err = this.server.Start()
	if err != nil {
		return err
	}
	if this.serviceDiscoveryRegedit != nil {
		err = this.serviceDiscoveryRegedit.RegService(serviceInfo, this.service.GetEndPoints())
		if err != nil {
			return err
		}
	}
	return nil
}

func buildApiDefineRquestHandler(serviceInfo common.ServiceInfo) web.RequestHandler {
	return func(request *http.Request, pathFragments map[string]string, reply web.Reply) {
		reply.With(serviceInfo.GetApiDefine()).As(web.Transport_Text)
	}
}

func (this *MicorService) regeditEndpoint(endPoint common.EndPoint) error {
	return this.server.Regedit(endPoint.Path, endPoint.Method, endPoint.HandlerFunc)
}

func (this *MicorService) regeditEndpoints() error {
	endPoints := this.service.GetEndPoints()
	if len(endPoints) == 0 {
		return errors.New("not regedit any endpoint")
	}
	for _, endPoint := range endPoints {
		err := this.regeditEndpoint(endPoint)
		if err != nil {
			return err
		}
	}
	return nil
}

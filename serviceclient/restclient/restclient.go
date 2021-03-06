package restclient

import (
	"context"
	"net/http"

	"github.com/coffeehc/commons/https/client"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coffeehc/microserviceboot/loadbalancer"
)

type restClient struct {
	options     *client.HTTPClientOptions
	transport   *http.Transport
	serviceInfo base.ServiceInfo
}

func (rc *restClient) getTransport() *http.Transport {
	return rc.transport
}

func newHttpClient(cxt context.Context, serviceInfo base.ServiceInfo, balancer loadbalancer.Balancer, defaultOption *client.HTTPClientOptions) *restClient {
	_clientOptions := &restClient{}
	_clientOptions.serviceInfo = serviceInfo
	transport := defaultOption.NewTransport((&_BalanceDialer{
		Timeout:   defaultOption.GetTimeout(),
		KeepAlive: defaultOption.GetDialerKeepAlive(),
		balancer:  balancer,
	}).DialContext)

	_clientOptions.transport = transport
	options := &client.HTTPClientOptions{
		Timeout:                        defaultOption.GetTimeout(),
		DialerTimeout:                  defaultOption.GetDialerTimeout(),
		DialerKeepAlive:                defaultOption.GetDialerKeepAlive(),
		TransportTLSHandshakeTimeout:   defaultOption.GetTransportTLSHandshakeTimeout(),
		TransportResponseHeaderTimeout: defaultOption.GetTransportResponseHeaderTimeout(),
		TransportIdleConnTimeout:       defaultOption.GetTransportIdleConnTimeout(),
		TransportMaxIdleConns:          defaultOption.GetTransportMaxIdleConns(),
		TransportMaxIdleConnsPerHost:   defaultOption.GetTransportMaxIdleConnsPerHost(),
		GlobalHeaderSetting:            defaultOption.GlobalHeaderSetting,
	}
	_clientOptions.options = options
	return _clientOptions
}

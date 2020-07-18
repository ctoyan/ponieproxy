package ponieproxy

import (
	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/filters"
	"github.com/elazarl/goproxy"
)

type PonieProxy struct {
	RequestFilters  []filters.RequestFilter
	ResponseFilters []filters.ResponseFilter
	Options         *config.Options
	ProxyInstance   *goproxy.ProxyHttpServer
}

func Init() *PonieProxy {
	setCA(caCert, caKey)
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	pp := new(PonieProxy)
	pp.ProxyInstance = proxy

	return pp
}

func (pp *PonieProxy) ApplyFilters() {
	for _, reqFilter := range pp.RequestFilters {
		pp.ProxyInstance.OnRequest(reqFilter.Conditions...).DoFunc(reqFilter.Handler)
	}

	for _, respFilter := range pp.ResponseFilters {
		pp.ProxyInstance.OnResponse(respFilter.Conditions...).DoFunc(respFilter.Handler)
	}
}

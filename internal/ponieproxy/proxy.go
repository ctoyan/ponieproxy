package ponieproxy

import (
	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/filters"
	"github.com/elazarl/goproxy"
)

type PonieProxy struct {
	Filters       *filters.Filters
	Options       *config.Options
	ProxyInstance *goproxy.ProxyHttpServer
}

func Init(options *config.Options) *PonieProxy {
	setCA(caCert, caKey)
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	pp := new(PonieProxy)
	pp.ProxyInstance = proxy
	pp.Filters = &filters.Filters{}

	pp.Filters.AddRequestFilters(options)
	pp.Filters.AddResponseFilters(options)

	for _, reqFilter := range pp.Filters.RequestFilters {
		pp.ProxyInstance.OnRequest(reqFilter.Conditions...).DoFunc(reqFilter.Handler)
	}

	for _, respFilter := range pp.Filters.ResponseFilters {
		pp.ProxyInstance.OnResponse(respFilter.Conditions...).DoFunc(respFilter.Handler)
	}

	return pp
}

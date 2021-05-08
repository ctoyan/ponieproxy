package ponieproxy

import (
	"io/ioutil"
	"log"

	"github.com/ctoyan/ponieproxy/filters"
	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/elazarl/goproxy"
)

type PonieProxy struct {
	RequestFilters  []filters.RequestFilter
	ResponseFilters []filters.ResponseFilter
	YAML            *config.YAML
	ProxyInstance   *goproxy.ProxyHttpServer
}

func Init() *PonieProxy {
	setCA(caCert, caKey)
	proxy := goproxy.NewProxyHttpServer()
	// Disable proxy logs
	proxy.Logger = log.New(ioutil.Discard, "", log.LstdFlags)

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

package filters

import (
	"github.com/elazarl/goproxy"
)

type requestFilter struct {
	Conditions []goproxy.ReqCondition
	Handler    goproxy.FuncReqHandler
}

type responseFilter struct {
	Conditions []goproxy.RespCondition
	Handler    goproxy.FuncRespHandler
}

type Filters struct {
	RequestFilters  []requestFilter
	ResponseFilters []responseFilter
}

func (f *Filters) addRespFilter(respFilter responseFilter) {
	f.ResponseFilters = append(f.ResponseFilters, respFilter)
}

func (f *Filters) addReqFilter(reqFilter requestFilter) {
	f.RequestFilters = append(f.RequestFilters, reqFilter)
}

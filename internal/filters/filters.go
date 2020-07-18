package filters

import (
	"github.com/elazarl/goproxy"
)

type RequestFilter struct {
	Conditions []goproxy.ReqCondition
	Handler    goproxy.FuncReqHandler
}

type ResponseFilter struct {
	Conditions []goproxy.RespCondition
	Handler    goproxy.FuncRespHandler
}

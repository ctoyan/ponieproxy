package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ctoyan/ponieproxy/filters"
	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/ponieproxy"
)

func main() {
	y := config.ParseYAML()

	pp := ponieproxy.Init()

	// Add your request filter functions here
	pp.RequestFilters = []filters.RequestFilter{
		filters.PopulateUserdata(y),
		filters.WriteReq(y),
		filters.HUNT(y),
		filters.SaveUrls(y),
		filters.DetectReqSecrets(y),
	}

	// Add your response filter functions here
	pp.ResponseFilters = []filters.ResponseFilter{
		filters.WriteResp(y),
		filters.SaveJs(y),
		filters.DetectRespSecrets(y),
	}

	// Apply all filters to the proxy
	pp.ApplyFilters()

	// Start the proxy
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", y.Host, y.Port), pp.ProxyInstance))
}

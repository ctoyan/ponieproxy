package main

import (
	"log"
	"net/http"

	"github.com/ctoyan/ponieproxy/filters"
	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/ponieproxy"
)

func main() {
	// Get your flags
	f := config.ParseFlags()

	// Return an instance of the ponieproxy
	pp := ponieproxy.Init()

	// Add your request filter functions here
	pp.RequestFilters = []filters.RequestFilter{
		filters.PopulateUserdata(f),
		filters.WriteReq(f),
		filters.HUNT(f),
		filters.SaveUrls(f),
		filters.SaveJs(f),
	}

	// Add your response filter functions here
	pp.ResponseFilters = []filters.ResponseFilter{
		filters.WriteResp(f),
	}

	// Apply all filters to the proxy
	pp.ApplyFilters()

	// Start the proxy
	log.Fatal(http.ListenAndServe(f.HostPort, pp.ProxyInstance))
}

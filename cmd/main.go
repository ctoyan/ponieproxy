package main

import (
	"log"
	"net/http"

	"github.com/ctoyan/ponieproxy/customFilters"
	"github.com/ctoyan/ponieproxy/internal/config"
	filters "github.com/ctoyan/ponieproxy/internal/filters"
	"github.com/ctoyan/ponieproxy/internal/ponieproxy"
)

func main() {
	// Get your flags
	f := config.ParseFlags()

	// Return an instance of the ponieproxy
	pp := ponieproxy.Init()

	// Add your request filter functions here
	pp.RequestFilters = []filters.RequestFilter{
		customFilters.WriteReq(f),
	}

	// Add your response filter functions here
	pp.ResponseFilters = []filters.ResponseFilter{
		customFilters.WriteResp(f),
	}

	// Apply all filters to the proxy
	pp.ApplyFilters()

	// Start the proxy
	log.Fatal(http.ListenAndServe(f.HostPort, pp.ProxyInstance))
}

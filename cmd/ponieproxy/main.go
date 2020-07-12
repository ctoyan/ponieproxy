package main

import (
	"log"
	"net/http"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/ponieproxy"
)

func main() {
	options := config.ParseOptions()
	ponieproxy := ponieproxy.Init(options)

	log.Fatal(http.ListenAndServe(options.HostPort, ponieproxy.ProxyInstance))
}

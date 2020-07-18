package config

import (
	"flag"
	"log"

	"github.com/ctoyan/ponieproxy/pkg/utils"
)

type Options struct {
	HostPort  string
	URL       string
	URLFile   string
	OutputDir string
}

func ParseFlags() *Options {
	o := &Options{}

	flag.StringVar(&o.HostPort, "h", ":8080", "Host and port. Default is :8080.")
	flag.StringVar(&o.URLFile, "u", "./urls.txt", "Path to a file, which contains a list of URL regexes to filter. Requires an existing file. Default is ./urls.txt")
	flag.StringVar(&o.OutputDir, "o", "./", "Path to a folder, which will contain uniquely named files with requests and responses.Every request and response have the same hash, but different extensions. Default is ./")

	flag.Parse()

	if !utils.FileExists(o.URLFile) {
		log.Fatalf("File %v doesn't exist", o.URLFile)
	}

	return o
}

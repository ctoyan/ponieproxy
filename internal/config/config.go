package config

import (
	"flag"
	"log"

	"github.com/ctoyan/ponieproxy/internal/utils"
)

type Flags struct {
	HostPort       string
	URL            string
	URLFile        string
	OutputDir      string
	SlackWebHook   string
	HuntOutputFile bool
	HuntExactMatch bool
}

func ParseFlags() *Flags {
	o := &Flags{}

	flag.StringVar(&o.HostPort, "h", ":8080", "Host and port")
	flag.StringVar(&o.URLFile, "u", "./urls.txt", "Path to a file, which contains a list of URL regexes to filter. Requires an existing file")
	flag.StringVar(&o.OutputDir, "o", "./", "Path to a folder, which will contain uniquely named files with requests and responses.Every request and response have the same hash, but different extensions")
	flag.StringVar(&o.SlackWebHook, "sw", "", "URL to slack webhook")
	flag.BoolVar(&o.HuntOutputFile, "ho", true, "Creates a checksumed file with the .hunt extension")
	flag.BoolVar(&o.HuntExactMatch, "hem", true, "Exact match for hunt params (case insensitive)")

	flag.Parse()

	if !utils.FileExists(o.URLFile) {
		log.Fatalf("File %v doesn't exist", o.URLFile)
	}

	return o
}

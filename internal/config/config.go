package config

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/ctoyan/ponieproxy/internal/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ExcludeReqContentTypes  []string `yaml:"excludeReqContentTypes"`
	IncludeReqContentTypes  []string `yaml:"includeReqContentTypes"`
	ExcludeRespContentTypes []string `yaml:"excludeRespContentTypes"`
	IncludeRespContentTypes []string `yaml:"includeRespContentTypes"`

	ExcludeReqFileTypes  []string `yaml:"excludeReqFileTypes"`
	IncludeReqFileTypes  []string `yaml:"includeReqFileTypes"`
	ExcludeRespFileTypes []string `yaml:"excludeRespFileTypes"`
	IncludeRespFileTypes []string `yaml:"includeRespFileTypes"`
}

type YAML struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Storage struct {
		Type string `yaml:"type"`
		DB   struct {
			Host string `yaml:"host"`
			User string `yaml:"user"`
			Pass string `yaml:"pass"`
		}
	}
	Settings struct {
		BaseOutputDir string   `yaml:"baseOutputDir"`
		InScope       []string `yaml:"inScope"`
		OutScope      []string `yaml:"outScope"`
		SlackHook     string   `yaml:"slackHook"`
	}

	Filters struct {
		Write struct {
			Active     bool   `yaml:"active"`
			ExactMatch bool   `yaml:"exactMatch"`
			Config     Config `yaml:"config"`
		}

		Hunt struct {
			Active         bool                `yaml:"active"`
			Config         Config              `yaml:"config"`
			ExactMatch     bool                `yaml:"exactMatch"`
			MatchingParams map[string][]string `yaml:"matchingParams"`
		}

		Urls struct {
			Active     bool   `yaml:"active"`
			OutputFile string `yaml:"outputFile"`
			Config     Config `yaml:"config"`
		}

		Secrets struct {
			Active    bool   `yaml:"active"`
			OutputDir string `yaml:"outputDir"`
			Config    Config `yaml:"config"`
		}

		Js struct {
			Active    bool   `yaml:"active"`
			OutputDir string `yaml:"outputDir"`
			Config    Config `yaml:"config"`
		}
	}
}

func ParseYAML() *YAML {
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "Config file path (e.g. ./config.yml)")
	flag.Parse()

	if !utils.FileExists(configFile) {
		log.Fatalf("File %v doesn't exist", configFile)
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading config.yml: %v ", err)
	}

	yConf := &YAML{}
	err = yaml.Unmarshal(yamlFile, yConf)
	if err != nil {
		log.Fatalf("error unmarshaling config.yml: %v", err)
	}

	return yConf
}

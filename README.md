# ponieproxy
Simple proxy which applies [filters](filters/README.md) (default or custom) to your requests and responses, while you browse a website.
It uses [goproxy](https://github.com/elazarl/goproxy).

It's useful to collect(saves js files, saves raw requests and resesponses, detects [HUNT](https://github.com/bugcrowd/HUNT) params and notifies you on slack, etc.) and manipulate data, while you browse a website and it does them using the [filters](filters/README.md) mentioned above.

## Install certificate
Create your own ca certificate and replace `ca.crt`, located in the root of this repository (or alternatively use the default `ca.crt`).
Then add it as a trusted certificate either in you browser or in your system, in order to be able to intercept TLS traffic.

## Configure browser
First of all, in order to use ponieproxy, you should set your browser to use ponieproxy as an HTTP proxy.

## Install
If you don't plan to write your own filters, you can just download it and run it like a normal tool.

Since it's written in Go and you can install it using:

```
go get -u github.com/ctoyan/ponieproxy/cmd/ponieproxy
```

Or you can [download a binary](https://github.com/ctoyan/ponieproxy/releases).

## Basic Usage

### With default filters
Create a `inscope.txt` file with regex for scoped URLs and run the proxy:

`ponieproxy -u URLS_FILE -o OUTPUT_DIR`

`Note:` The default filters add the regex lines between parens. For example - `(REGEX_ON_LINE_ONE)|(REGEX_ON_LINE_TWO)`

### With custom filters
Clone/Fork this repository. [Write your filters](filters/README.md). Then `cd` into the cloned repo in `cmd/ponieproxy` and run:

```
go install
ponieproxy -o OUTPUT_DIR -u URLS_FILE
```

## Arguments
```
Usage of ponieproxy:
  -c string
    	Config file path (e.g. ./config.yml) (default "./config.yml")
```

## What's different than Burp or ZAP?
With ponieproxy you can use the power of Go to write filters and manipulate the requests and responses, however you like. Depends on your skills and creativity. 

The other key difference is the way you use it - it's meant to be ran in the background and do whatever you program/set it to do and save everything in files, so you can analyze them later with other tools.

## Upcoming features/filters

- write filters with YAML, instead of Go (IN PROGRESS)
- add results to db instead of files
- reflected parameters detection
- find and replace in requests

# ponieproxy
Simple proxy which captures all requests and responses and saves them in uniquely named files.
It uses [goproxy](https://github.com/elazarl/goproxy).

It's useful for saving the traffic as text files and using it to find various data (like secrets, urls, endpoints), with the help of other bash tools.
It can also be used to apply the [HUNT Methodology](https://github.com/bugcrowd/HUNT) in a more bash friendly way.

## Install certificate
Add `ca.crt`, located in the root of this repository, as a trusted certificate either in you browser or in your system, in order to be able to intercept TLS traffic.

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
Create a `urls.txt` file with regex for scoped URLs and run the proxy:

`ponieproxy -u URLS_FILE -o OUTPUT_DIR`

`Note:` The default filters adds the regex lines between parens. For example - `(REGEX_ON_LINE_ONE)|(REGEX_ON_LINE_TWO)`

### With custom filters
Clone/Fork this repository. [Write your filters](filters/README.md). Then `cd` into the cloned repo in `cmd/ponieproxy` and run:

```
go install
ponieproxy -o OUTPUT_DIR -u URLS_FILE
```

## Arguments
```
-h string
    	Host and port. (default ":8080")
-u string
    	Path to a file, which contains a list of URL regexes to filter (to intercept all requests, use '.*'). Requires an existing file. (default "./urls.txt")
-hem
    	Exact match for hunt params (case insensitive). (default true)
-ho
    	Creates a checksumed file with the .hunt extension. (default true)
-o string
    	Path to a folder, which will contain uniquely named files with requests and responses.Every request and response have the same hash, but different extensions. (default "./")
-sw string
    	URL to slack webhook. No default
```

## Upcoming features/filters

- save all JS files (IN PROGRESS)
- reflected parameters detection
- find and replace in requests
- write filters with YAML, instead of Go

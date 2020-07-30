# ponieproxy
Simple proxy which captures all requests and responses and saves them in uniquely named files.
It uses [goproxy](https://github.com/elazarl/goproxy).

It's useful for saving the traffic as text files and using it to find various data (like secrets, urls, endpoints), with the help of other bash tools.
It can also be used to apply the [HUNT Methodology](https://github.com/bugcrowd/HUNT) in a more bash friendly way.

## Install certificate
Add `ca.crt`, located in the root of this repository, as a trusted certificate either in you browser of in your system, in order to be able to intercept TLS traffic.

## Configure browser
First of all, in order to use ponieproxy, you should set your browser to use ponieproxy as an HTTP proxy.

## Installation and usage if you want to use the DEFAULT FILTERS
### Install
If you don't plan to write your own filters, you can just download it and run it like a normal tool.

Since it's written in Go and you can install it using:

```
go get -u github.com/ctoyan/ponieproxy/cmd
```

Or you can [download a binary](https://github.com/ctoyan/ponieproxy/releases).

### Usage
This runs the proxy with default filters:

`ponieproxy -u ./urls.txt -o ./out`


## Installation and usage if you want to write CUSTOM FILTERS
### Install
You just need to clone this repository.

### Usage
`cd` into the cloned repo and run:
```
go run ./cmd/main.go -o OUTPUT_DIR -u URLS_FILE
```

## Arguments
```
-o string
    	Path to a folder, which will contain uniquely named files with requests and responses (default "./"). Every request and response have the same hash, but different extensions
-u string
    	Path to a file, which contains a list of URL regexes to intercept
```

`Note:` The default filters adds the regex lines between parens. For example - `(REGEX_ON_LINE_ONE)|(REGEX_ON_LINE_TWO)`

## Default Filters
You can check details in `customFilters/default.go` and these are the default ones currently:

- `WriteReq()` - writes uniquely hashed and unique requests for all matching regexes in urls.txt
- `WriteResp()` - writes uniquely hashed and unique responses for all matching regexes in urls.txt

## Writing your own filters
Since the ponieproxy is just a small wrapper over goproxy, a filter is a struct that combines a slice of goproxy conditions([Req](https://godoc.org/gopkg.in/elazarl/goproxy.v1#ReqCondition) and [Resp](https://godoc.org/gopkg.in/elazarl/goproxy.v1#RespCondition)) and a goproxy handler([Req](https://godoc.org/gopkg.in/elazarl/goproxy.v1#FuncReqHandler) and [Resp](https://godoc.org/gopkg.in/elazarl/goproxy.v1#FuncRespHandler)). Basically the conditions are applied to that handler and I've called it a filter.

There are some default filters added by me, which live in the `customFilters/default.go` file.

You can add or remove filters in the two arrays in the `main.go` file - `RequestFilters` and `ResponseFilters`.

Check out the default filters and you'll see how easy it is to write your own.

## Upcoming features/filters

- find [HUNT](https://github.com/bugcrowd/HUNT) params and send slack notifications (IN PROGRESS)
- add all requests paths to a file
- reflected parameters detection
- find and replace in requests
- write filters with YAML, instead of Go

# ponieproxy
Simple proxy which captures all requests and responses and saves them in uniquely named files.
It uses [goproxy](https://github.com/elazarl/goproxy).

It's useful for saving the traffic as text files and using it to find various data (like secrets, urls, endpoints), with the help of other bash tools.
It can also be used to apply the [HUNT Methodology](https://github.com/bugcrowd/HUNT) in a more bash friendly way.

## Install
The tool is written in Go and you can install it using:

```
go get -u github.com/ctoyan/ponieproxy
```

Or you can [download a binary](https://github.com/ctoyan/ponieproxy/releases).

### Install certificate
Add `ca.crt`, located in the root of this repository, as a trusted certificate either in you browser of in your system, in order to be able to intercept TLS traffic.

### Configure browser
First of all, in order to use ponieproxy, you should set your browser to use ponieproxy as an HTTP proxy.

## Basic Usage
The usage options are as follows:
```
-o string
    	Path to a folder, which will contain uniquely named files with requests and responses (default "./"). Every request and response have the same hash, but different extensions
-u string
    	Path to a file, which contains a list of URL regexes to intercept
```

`Note:` The regex for multiple urls is added between parens. For example - `(REGEX_ON_LINE_ONE)|(REGEX_ON_LINE_TWO)`

## Basic Example

`ponieproxy -u github.com -o ./out`

Outputs files sha1 summed files in the `./out` directory.

The content of a \*.req file is similar to:

```
POST /_private/browser/stats HTTP/1.1
Host: api.github.com
Accept: */*
Accept-Language: en-GB,en;q=0.5
Content-Length: 6453
Content-Type: text/plain;charset=UTF-8
Cookie: redacted
Origin: https://github.com
Referer: https://someurl
User-Agent: some user agent
Some-more-headers: here

{"stats":[{...some JSON here...}]}
```

The content of a \*.res file is similar to:

```
HTTP/1.1 200 OK
Content-Length: 0
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: ETag, Link, Location, Retry-After, X-GitHub-OTP, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, X-OAuth-Scopes, X-Accepted-OAuth-Scopes, X-Poll-Interval, X-GitHub-Media-Type, Deprecation, Sunset
Cache-Control: no-cache
Content-Security-Policy: default-src 'none'
Content-Type: text/plain
Referrer-Policy: origin-when-cross-origin, strict-origin-when-cross-origin
Server: GitHub.com
Status: 200 OK
Some-more-headers: here

{some JSON response here}
```

##Applying your own filters
Since ponieproxy uses goproxy behind the scenes, you can apply your own request and response filters using the goproxy `OnResponse` and `OnRequest` functions, along with conditions applied to them. Please check the [goproxy docs](https://godoc.org/gopkg.in/elazarl/goproxy.v1)

You can write your own filters when you go to `/internal/filters` and choose to write a request or response filter.

Then just use `f.addReqFilter` or `f.addRespFilter` and add the type of filter you want, which consists of Conditions and a Handler. Make sure to always return the requests and responses, so the proxy can forward them.

Take note that instead of just filtering the traffic, you can also modify the requests/responses simply by returning the modified request/response in the request/response filter handler.

# ponieproxy
Simple proxy which captures all requests and responses and saves them in uniquely named files.
It uses [goproxy](https://github.com/elazarl/goproxy).

It's useful for saving the traffic as text files and using it to find various data (like secrets, urls, endpoints), with the help of other bash tools.

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

## Usage
The usage options are as follows:
```
-o string
    	Path to a folder, which will contain uniquely named files with requests and responses (default "./")
-u string
    	Regex to match a single url for intercepting (ex. example.com/api)
-uL string
    	Path to a file, which contains a list of URL regexes to intercept
```

`Note:` The regex for multiple urls is added between parens. For example - `(REGEX_ON_LINE_ONE)|(REGEX_ON_LINE_TWO)`

## Example

`ponieproxy -u github.com -o ./out`

Outputs files sha1 summed files in the `./out` directory.

The content of a file is similar to:

```
-------------------------------------REQUEST-------------------------------------
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


-------------------------------------RESPONSE-------------------------------------
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

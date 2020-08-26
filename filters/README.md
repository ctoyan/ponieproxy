# Filters
Since the ponieproxy is just a small wrapper over [goproxy](https://github.com/elazarl/goproxy), a filter is a struct that combines a slice of goproxy conditions([Req](https://godoc.org/gopkg.in/elazarl/goproxy.v1#ReqCondition) and [Resp](https://godoc.org/gopkg.in/elazarl/goproxy.v1#RespCondition)) and a goproxy handler([Req](https://godoc.org/gopkg.in/elazarl/goproxy.v1#FuncReqHandler) and [Resp](https://godoc.org/gopkg.in/elazarl/goproxy.v1#FuncRespHandler)). Basically the conditions are applied to that handler and I've called it a filter.

You can add or remove filters in the two arrays in the `main.go` file - `RequestFilters` and `ResponseFilters`.

By default all filters are added (enabled). If you remove/comment them from `main.go`, they won't be applied.

## Basic Filters
The idea for those is to provide more basic functionality which most people would want.
You can check details in `filters/write.go` and `filters/populate.go`. Currently they are:

- `PopulateUserData()` - populates the ctx.UserData with some userful data, which is send across all captured requests/responses
- `WriteReq()` - writes uniquely hashed and unique requests for all matching regexes in inscope.txt
- `WriteResp()` - writes uniquely hashed and unique responses for all matching regexes in inscope.txt

## HUNT Filter
Ponieproxy applies a filter to ease the use of the [HUNT Methodology](https://github.com/bugcrowd/HUNT).

A valid question is - **What's the difference with the Burp and ZAP plugins that are already present?**

The answer is that you have a bit more control over the type of matching that it does. It always matches params case insensitively. The default matching style is exact (using `==`). 

If you set `-hem` to `false`, it will look for a substring within the param. Foe example if the filter is searching for `id`, it will positively match the following `userId`, `identification`, `ID`. With `-hem` set to the default `true`, it will only match `id`, `Id`, etc.

Ponieproxy only looks for matches in request query params (e.g. `?id=123&url=ssrf.com&nomatch=123`) and in JSON keys in the request body (e.g. `{id: 123, url: "ssrf.com"", nomatch: 123}`). In these two cases, ponieproxy will detect `id` and `url` and will write it in a `.hunt` file (same checksum name as the .req and .res files) and/or send a slack notification.

If you want an alert in slack, you should pass a slack webhook url to the `sw` option.

Currently the HUNT filter matches params for `IDOR`, `SQL Injection`, `SSRF`, `SSTI`, `LFI/RFI/Path Traversal`, `OSCI`, `Debug and Logic Parameters`, which are taken directly from the HUNT repo.

## Filter to Save Request URL
Appends to a file all unique, in-scope URLs, that you've requested. If you want to name the file, pass the filename to the `-su` flag.

## Filter to Save JS Files
Saves all `.js` files to their corresponding folder. Default save filder is `./js`. So for example a request for `https://somesite.com/path/to/file.js`, will save the file to `./js/path/to/file.js`

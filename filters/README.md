# Filters
Since the ponieproxy is just a small wrapper over goproxy, a filter is a struct that combines a slice of goproxy conditions([Req](https://godoc.org/gopkg.in/elazarl/goproxy.v1#ReqCondition) and [Resp](https://godoc.org/gopkg.in/elazarl/goproxy.v1#RespCondition)) and a goproxy handler([Req](https://godoc.org/gopkg.in/elazarl/goproxy.v1#FuncReqHandler) and [Resp](https://godoc.org/gopkg.in/elazarl/goproxy.v1#FuncRespHandler)). Basically the conditions are applied to that handler and I've called it a filter.

You can add or remove filters in the two arrays in the `main.go` file - `RequestFilters` and `ResponseFilters`.

Check out the default filters and you'll see how easy it is to write your own.

## Default Filters
The idea for those is to provide more basic functionality which most people would want.
You can check details in `filters/default.go`. Currently the default ones are:

- `PopulateUserData()` - Populates the ctx.UserData with some userful data, which is send across all captured requests/responses
- `WriteReq()` - writes uniquely hashed and unique requests for all matching regexes in urls.txt
- `WriteResp()` - writes uniquely hashed and unique responses for all matching regexes in urls.txt

## HUNT Filters
Ponieproxy applies filters to ease the use of the [HUNT Methodology](https://github.com/bugcrowd/HUNT).

A valid question is - What's the difference with the Burp and ZAP plugins that are already present?

The answer is that you have a bit more control over the type of matching that it does. It always matches params case insensitively. The default matching style is exact (using `==`). 

If you set `-hem` to `false`, it will look for a substring within the param. If the filter is searching for `id`, it will positively match the following `userId`, `identification`, `ID`.

Ponieproxy only looks for matches in request query params (e.g. `?id=123&url=ssrf.com&nomatch=123`) and in JSON keys in request body (e.g. `{id: 123, url: "ssrf.com"", nomatch: 123}`). In these two cases, ponieproxy will detect `id` and `url` and will write in a `.hunt` (checksum named) file and/or send a slack notification.

If you want an alert in slack, you should pass a slack webhook url to the `sw` option.

Currently there are filters integrated to match params for `IDOR`, `SQL Injection`, `SSRF`, `SSTI`, `LFI/RFI/Path Traversal`, `OSCI`, `Debug and Logic Parameters`, which are taken directly from the HUNT repo.

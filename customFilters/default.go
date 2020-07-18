package customFilters

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/filters"
	"github.com/ctoyan/ponieproxy/pkg/utils"
	"github.com/elazarl/goproxy"
)

/* Request filter
 * Write it to a uniquely named *.req file, in the output folder
 *
 * The only filter condition, wraps every line from your urls file
 * between braces and concatenates them, making the following regex:
 * (LINE_ONE)|(LINE_TWO)|(LINE_THREE), where LINE_N is a single line from your file.
 */
func WriteReq(f *config.Options) filters.RequestFilter {
	urlsList, err := utils.ReadLines(f.URLFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return filters.RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			reqBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("error reading reqBody: %v\n", err)
			}

			requestDump, err := httputil.DumpRequest(req, false)
			if err != nil {
				fmt.Printf("error on request dump: %v\n", err)
			}

			go utils.WriteUniqueFile(req.URL.Host, req.URL.Path, reqBody, f.OutputDir, requestDump, "req")

			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

			//pass reqBody to context for cache key construction
			ctx.UserData = reqBody
			return req, nil
		},
	}
}

/* Response filter
 * Write it to a uniquely named *.res file, in the output folder
 *
 * The only filter condition, wraps every line from your urls file
 * between braces and concatenates them, making the following regex:
 * (LINE_ONE)|(LINE_TWO)|(LINE_THREE), where LINE_N is a single line from your file.
 */
func WriteResp(f *config.Options) filters.ResponseFilter {
	urlsList, err := utils.ReadLines(f.URLFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return filters.ResponseFilter{
		Conditions: []goproxy.RespCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

			responseDump, err := httputil.DumpResponse(res, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			reqBody := ctx.UserData.([]byte)

			go utils.WriteUniqueFile(res.Request.URL.Host, res.Request.URL.Path, reqBody, f.OutputDir, responseDump, "res")

			return res
		},
	}
}

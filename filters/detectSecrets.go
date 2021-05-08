package filters

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/elazarl/goproxy"
)

/* Request filter
 * Detects secrets using regexs
 *
 * The only filter condition, wraps every line from your urls file
 * between braces and concatenates them, making the following regex:
 * (LINE_ONE)|(LINE_TWO)|(LINE_THREE), where LINE_N is a single line from your file.
 */
func DetectReqSecrets(y *config.YAML) RequestFilter {
	if !y.Filters.Secrets.Active {
		return RequestFilter{}
	}

	allSecrets := map[string]struct{}{}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			reqFileType(true, y.Filters.Secrets.Config.ExcludeReqFileTypes...),
			reqFileType(false, y.Filters.Secrets.Config.IncludeReqFileTypes...),
			reqContentType(true, y.Filters.Secrets.Config.ExcludeReqContentTypes...),
			reqContentType(false, y.Filters.Secrets.Config.IncludeReqContentTypes...),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			ud := ctx.UserData.(UserData)
			requestDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			outputDir := fmt.Sprintf("%v/%v", y.Settings.BaseOutputDir, y.Filters.Secrets.OutputDir)
			go detectSecrets(&allSecrets, requestDump, ud, outputDir)

			return req, nil
		},
	}
}

/* Response filter
 * Detects secrets using regexs
 *
 * The only filter condition, wraps every line from your urls file
 * between braces and concatenates them, making the following regex:
 * (LINE_ONE)|(LINE_TWO)|(LINE_THREE), where LINE_N is a single line from your file.
 */
func DetectRespSecrets(y *config.YAML) ResponseFilter {
	if !y.Filters.Secrets.Active {
		return ResponseFilter{}
	}

	allSecrets := map[string]struct{}{}

	return ResponseFilter{
		Conditions: []goproxy.RespCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			respFileType(true, y.Filters.Secrets.Config.ExcludeRespFileTypes...),
			respFileType(false, y.Filters.Secrets.Config.IncludeRespFileTypes...),
			respContentType(true, y.Filters.Secrets.Config.ExcludeRespContentTypes...),
			respContentType(true, y.Filters.Secrets.Config.IncludeRespContentTypes...),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			responseDump, err := httputil.DumpResponse(res, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			outputDir := fmt.Sprintf("%v/%v", y.Settings.BaseOutputDir, y.Filters.Secrets.OutputDir)
			go detectSecrets(&allSecrets, responseDump, ctx.UserData.(UserData), outputDir)

			return res
		},
	}
}

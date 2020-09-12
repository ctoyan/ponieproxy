package filters

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

/* Request filter
 * Detects secrets using regexs
 *
 * The only filter condition, wraps every line from your urls file
 * between braces and concatenates them, making the following regex:
 * (LINE_ONE)|(LINE_TWO)|(LINE_THREE), where LINE_N is a single line from your file.
 */
func DetectReqSecrets(f *config.Flags) RequestFilter {
	scopeUrls, err := utils.ReadLines(f.ScopeFile)
	allSecrets := map[string]struct{}{}
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(scopeUrls, ")|(")))),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			ud := ctx.UserData.(UserData)
			requestDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			go detectSecrets(&allSecrets, requestDump, ud, f.SavedSecretsDir)

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
func DetectRespSecrets(f *config.Flags) ResponseFilter {
	scopeUrls, err := utils.ReadLines(f.ScopeFile)
	allSecrets := map[string]struct{}{}
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return ResponseFilter{
		Conditions: []goproxy.RespCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(scopeUrls, ")|(")))),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			responseDump, err := httputil.DumpResponse(res, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			go detectSecrets(&allSecrets, responseDump, ctx.UserData.(UserData), f.SavedSecretsDir)

			return res
		},
	}
}

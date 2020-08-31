package filters

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

/* Request filter
 * Detect javascript files and save them in a folder
 */
func SaveJs(f *config.Flags) ResponseFilter {
	scopeUrls, err := utils.ReadLines(f.ScopeFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	// detect js file based on content-type response header
	hasJSContentType := goproxy.RespConditionFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		contentType := resp.Header.Get("Content-Type")
		ciContentType := strings.ToLower(contentType)

		return strings.Contains(ciContentType, "javascript") ||
			strings.Contains(ciContentType, "jscript") ||
			strings.Contains(ciContentType, "ecmascript")
	})

	return ResponseFilter{
		Conditions: []goproxy.RespCondition{
			hasJSContentType,
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(scopeUrls, ")|(")))),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			go func() {
				err := utils.DownloadFromURL(res.Request.URL, f.JsOutputDir)
				if err != nil {
					log.Printf("error downloading from url: %v", err)
				}
			}()

			return res
		},
	}
}

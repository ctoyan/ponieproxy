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
func SaveJs(f *config.Flags) RequestFilter {
	scopeUrls, err := utils.ReadLines(f.ScopeFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(scopeUrls, ")|(")))),
			goproxy.UrlMatches(regexp.MustCompile("^.+.js$")),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			go func() {
				err := utils.DownloadFromURL(req.URL, f.JsOutputDir)
				if err != nil {
					log.Printf("error downloading from url: %v", err)
				}
			}()

			return req, nil
		},
	}
}

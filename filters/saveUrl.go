package filters

import (
	"crypto/sha1"
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
 * Save every in scope url to a file containing a list of URLs
 */
func SaveUrls(f *config.Flags) RequestFilter {
	scopeUrls, err := utils.ReadLines(f.ScopeFile)
	savedUrls := map[[20]byte]struct{}{}
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(scopeUrls, ")|(")))),
			reqFileType(true, ".png", ".jpg", ".jpeg", ".woff", ".css", ".gif", ".js", ".ico"),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			currentUrl := req.URL.String()
			checksum := sha1.Sum([]byte(currentUrl))

			if _, ok := savedUrls[checksum]; !ok {
				go utils.AppendToFile(currentUrl, f.SavedUrlsFile)
				savedUrls[checksum] = struct{}{}
			}

			return req, nil
		},
	}
}

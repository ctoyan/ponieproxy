package filters

import (
	"crypto/sha1"
	"fmt"
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
func SaveUrls(y *config.YAML) RequestFilter {
	if !y.Filters.Urls.Active {
		return RequestFilter{}
	}

	savedUrls := map[[20]byte]struct{}{}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			reqFileType(true, y.Filters.Urls.Config.ExcludeReqFileTypes...),
			reqFileType(false, y.Filters.Urls.Config.IncludeReqFileTypes...),
			reqContentType(true, y.Filters.Urls.Config.ExcludeReqContentTypes...),
			reqContentType(false, y.Filters.Urls.Config.IncludeReqContentTypes...),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			currentUrl := req.URL.String()
			checksum := sha1.Sum([]byte(currentUrl))

			if _, ok := savedUrls[checksum]; !ok {
				outputFile := fmt.Sprintf("%v/%v", y.Settings.BaseOutputDir, y.Filters.Urls.OutputFile)
				fmt.Println(currentUrl)
				go utils.AppendToFile(currentUrl, outputFile)
				savedUrls[checksum] = struct{}{}
			}

			return req, nil
		},
	}
}

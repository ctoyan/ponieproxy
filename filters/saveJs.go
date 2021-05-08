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
func SaveJs(y *config.YAML) ResponseFilter {
	if !y.Filters.Js.Active {
		return ResponseFilter{}
	}

	return ResponseFilter{
		Conditions: []goproxy.RespCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			respFileType(true, y.Filters.Js.Config.ExcludeRespFileTypes...),
			respFileType(false, y.Filters.Js.Config.IncludeRespFileTypes...),
			respContentType(true, y.Filters.Js.Config.ExcludeRespContentTypes...),
			respContentType(false, y.Filters.Js.Config.IncludeRespContentTypes...),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			go func() {
				outputDir := fmt.Sprintf("%v/%v", y.Settings.BaseOutputDir, y.Filters.Js.OutputDir)
				err := utils.DownloadFromURL(res.Request.URL, outputDir)
				if err != nil {
					log.Printf("error downloading from url: %v", err)
				}
			}()

			return res
		},
	}
}

package filters

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

/* Request filter
 * Detect IDOR params
 *
 * Following the HUNT Methodology (https://github.com/bugcrowd/HUNT):
 * We're looking for exact or non-exact (case insensitive) matches for keywords
 */
func HUNT(y *config.YAML) RequestFilter {
	if !y.Filters.Hunt.Active {
		return RequestFilter{}
	}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			reqFileType(true, y.Filters.Hunt.Config.ExcludeReqFileTypes...),
			reqFileType(false, y.Filters.Hunt.Config.IncludeReqFileTypes...),
			reqContentType(true, y.Filters.Hunt.Config.ExcludeReqContentTypes...),
			reqContentType(false, y.Filters.Hunt.Config.IncludeReqContentTypes...),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			hunt := y.Filters.Hunt.MatchingParams
			ud := ctx.UserData.(UserData)

			// Get all keys from JSON
			var reqJsonKeys map[string]struct{}
			reqBodyJson, err := utils.UnmarshalReqBody(ud.ReqBody)
			if err == nil {
				reqJsonKeys = utils.CollectJsonKeys(reqBodyJson, map[string]struct{}{})
			}

			// Get all query params from request
			reqQueryParams := req.URL.Query()

			// Search for matches within HUNT list
			for huntKey, huntParams := range hunt {
				for _, param := range huntParams {
					if reqQueryParams != nil {
						go FindInQueryParams(huntKey, param, reqQueryParams, y, ud)
					}

					if reqJsonKeys != nil {
						go FindInJson(huntKey, param, reqJsonKeys, y, ud)
					}
					// handle non json requests (other mime types)
					// maybe search with string.contains all others?
				}
			}

			return req, nil
		},
	}
}

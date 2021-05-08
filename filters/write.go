package filters

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

/* Request filter
 * Write it to a uniquely named *.req file, in the output folder
 *
 * The only filter condition, wraps every line from your urls file
 * between braces and concatenates them, making the following regex:
 * (LINE_ONE)|(LINE_TWO)|(LINE_THREE), where LINE_N is a single line from your file.
 */
func WriteReq(y *config.YAML) RequestFilter {
	// if !y.Filters.WriteFiles.Active {
	// 	return RequestFilter{}
	// }

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			reqFileType(true, y.Filters.Write.Config.ExcludeReqFileTypes...),
			reqFileType(false, y.Filters.Write.Config.IncludeReqFileTypes...),
			reqContentType(true, y.Filters.Write.Config.ExcludeReqContentTypes...),
			reqContentType(false, y.Filters.Write.Config.IncludeReqContentTypes...),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			ud := ctx.UserData.(UserData)
			go utils.WriteUniqueFile(ud.Host, ud.FileChecksum, ud.ReqBody, y.Settings.BaseOutputDir, ud.ReqDump, ".req")

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
func WriteResp(y *config.YAML) ResponseFilter {
	// if !y.Filters.WriteFiles.Active {
	// 	return ResponseFilter{}
	// }

	return ResponseFilter{
		Conditions: []goproxy.RespCondition{
			goproxy.Not(goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.OutScope, ")|("))))),
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(y.Settings.InScope, ")|(")))),
			reqFileType(true, y.Filters.Write.Config.ExcludeReqFileTypes...),
			reqFileType(false, y.Filters.Write.Config.IncludeReqFileTypes...),
			reqContentType(true, y.Filters.Write.Config.ExcludeReqContentTypes...),
			reqContentType(false, y.Filters.Write.Config.IncludeReqContentTypes...),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			responseDump, err := httputil.DumpResponse(res, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			ud := ctx.UserData.(UserData)
			go utils.WriteUniqueFile(ud.Host, ud.FileChecksum, ud.ReqBody, y.Settings.BaseOutputDir, string(responseDump), ".res")

			return res
		},
	}
}

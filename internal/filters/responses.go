package filters

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

/*
 * For all requests matching the response conditions apply this handler
 */
func (f *Filters) AddResponseFilters(options *config.Options) {
	urlsList, err := utils.ReadLines(options.URLFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	/*
	 * Write all requests(gotten from the context) and responses to files
	 */
	f.addRespFilter(responseFilter{
		Conditions: []goproxy.RespCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))),
		},
		Handler: func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			requestDump, err := httputil.DumpRequest(res.Request, false)
			if err != nil {
				fmt.Printf("error on request dump: %v\n", err)
			}

			responseDump, err := httputil.DumpResponse(res, true)
			if err != nil {
				fmt.Printf("error on response dump: %v\n", err)
			}

			reqBody := string(ctx.UserData.([]byte))

			go func() {
				if options.OutputDir != "./" {
					os.MkdirAll(options.OutputDir, os.ModePerm)
				}

				cacheKey := fmt.Sprintf("%s%s%s", res.Request.URL.Host, res.Request.URL.Path, reqBody)
				hashedPair := sha1.Sum([]byte(cacheKey))

				filePath := fmt.Sprintf("%v/%x", options.OutputDir, hashedPair)
				reqFilePath := fmt.Sprintf("%v.req", filePath)
				resFilePath := fmt.Sprintf("%v.res", filePath)

				if !utils.FileExists(reqFilePath) {
					constructedReq := fmt.Sprintf(`%s %s`, requestDump, reqBody)
					err := utils.AppendToFile(constructedReq, reqFilePath)
					if err != nil {
						log.Fatalf("error writing to file: %v", err)
					}
				}

				if !utils.FileExists(resFilePath) {
					err := utils.AppendToFile(string(responseDump), resFilePath)
					if err != nil {
						log.Fatalf("error writing to file: %v", err)
					}
				}
			}()

			return res
		},
		/*Add more filters here*/
	})
}

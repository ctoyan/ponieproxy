package filters

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
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
 * Write various params to UserData.
 *
 * UserData is a part of the proxy context.
 * It is passed to every request and response.
 */
func PopulateUserdata(f *config.Flags) RequestFilter {
	urlsList, err := utils.ReadLines(f.URLFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			reqBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("error reading reqBody: %v\n", err)
			}

			requestDump, err := httputil.DumpRequest(req, false)
			if err != nil {
				fmt.Printf("error on request dump: %v\n", err)
			}

			checksum := sha1.Sum([]byte(fmt.Sprintf("%s%s%s", req.URL.Host, req.URL.Path, reqBody)))
			ctx.UserData = UserData{
				ReqBody:      string(reqBody),
				ReqDump:      string(requestDump),
				FileChecksum: hex.EncodeToString(checksum[:]),
				Host:         req.URL.Host,
				ReqURL:       req.URL.String(),
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
			return req, nil
		},
	}
}

package filters

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

func (f *Filters) AddRequestFilters(options *config.Options) {
	urlsList, err := utils.ReadLines(options.URLFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	/*
	 * Add the request body to the context, so it can be used later in responses
	 */
	f.addReqFilter(requestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("error reading body: %v\n", err)
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			ctx.UserData = body
			return req, nil
		},
		/*Add more filters here*/
	})
}

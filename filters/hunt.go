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
 * Detect IDOR params
 *
 * Following the HUNT Methodology (https://github.com/bugcrowd/HUNT):
 * We're looking for exact or non-exact (case insensitive) matches for keywords
 */
func HUNT(f *config.Flags) RequestFilter {
	urlsList, err := utils.ReadLines(f.URLFile)
	if err != nil {
		log.Fatalf("error reading lines from file: %v", err)
	}

	return RequestFilter{
		Conditions: []goproxy.ReqCondition{
			goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))),
		},
		Handler: func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			hunt := map[string][]string{}
			hunt["IDOR"] = []string{"account", "doc", "edit", "email", "group", "id", "key", "no", "number", "order", "profile", "report", "user"}
			hunt["SSRF"] = []string{"callback", "continue", "data", "dest", "dir", "domain", "feed", "host", "html", "navigation", "next", "open", "out", "page", "path", "port", "redirect", "reference", "return", "show", "site", "to", "uri", "url", "val", "validate", "view", "window"}
			hunt["OSCI"] = []string{"cli", "cmd", "daemon", "dir", "download", "execute", "ip", "log", "upload"}
			hunt["SQLI"] = []string{"column", "delete", "fetch", "field", "filter", "from", "id", "keyword", "name", "number", "order", "params", "process", "query", "report", "results", "role", "row", "search", "sel", "select", "sleep", "sort", "string", "table", "update", "user", "view", "where"}
			hunt["FIPT"] = []string{"doc", "document", "file", "folder", "path", "pdf", "pg", "php_path", "root", "style", "template"}
			hunt["SSTI"] = []string{"activity", "content", "id", "name", "preview", "redirect", "template", "view"}
			hunt["LOGIC"] = []string{"access", "adm", "admin", "alter", "cfg", "clone", "config", "create", "dbg", "debug", "delete", "disable", "edit", "enable", "exec", "execute", "grant", "load", "make", "modify", "rename", "reset", "root", "shell", "test", "toggle"}
			hunt["RCE"] = []string{"cmd", "exec", "command", "execute", "ping", "query", "jump", "code", "reg", "do", "func", "arg", "option", "load", "process", "step", "read", "function", "req", "feature", "exe", "module", "payload", "run", "print"}

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
						go FindInQueryParams(huntKey, param, reqQueryParams, f, ud)
					}

					if reqJsonKeys != nil {
						go FindInJson(huntKey, param, reqJsonKeys, f, ud)
					}
					// handle non json requests (other mime types)
					// maybe search with string.contains all others?
				}
			}

			return req, nil
		},
	}
}

package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	singleUrl := flag.String("u", "", "Regex to match a single url for intercepting (ex. example.com/api)")
	urlListFile := flag.String("uL", "", "Path to a file, which contains a list of URL regexes to intercept")
	outputDir := flag.String("o", "./", "Path to a folder, which will contain uniquely named files with requests and responses")

	flag.Parse()
	setCA(caCert, caKey)
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	var reqConds []goproxy.ReqCondition
	var respConds []goproxy.RespCondition

	if *urlListFile != "" {
		urlsList, err := readLines(*urlListFile)
		if err != nil {
			log.Fatalf("error reading lines from file: %v", err)
		}

		reqConds = append(reqConds, goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))))
		respConds = append(respConds, goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(urlsList, ")|(")))))
	}

	if *singleUrl != "" {
		reqConds = append(reqConds, goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("%v", *singleUrl))))
		respConds = append(respConds, goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("%v", *singleUrl))))
	}

	proxy.OnRequest(reqConds...).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("error reading body: %v\n", err)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		ctx.UserData = body
		return req, nil
	})

	proxy.OnResponse(respConds...).DoFunc(func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
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
			reqRespPair := fmt.Sprintf(`-------------------------------------REQUEST-------------------------------------
%s
%s


-------------------------------------RESPONSE-------------------------------------
%s`, requestDump, reqBody, responseDump)

			if *outputDir != "./" {
				os.MkdirAll(*outputDir, os.ModePerm)
			}

			hashedPair := sha1.Sum([]byte(reqRespPair))
			filePath := fmt.Sprintf("%v/%x", *outputDir, hashedPair)
			if !fileExists(filePath) {
				err := appendToFile(reqRespPair, filePath)
				if err != nil {
					log.Fatalf("error writing to file: %v", err)
				}
			}
		}()

		return res
	})

	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func appendToFile(data string, filePath string) error {
	if filePath != "" {
		f, err := os.OpenFile(filePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(data + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

package filters

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
	"github.com/elazarl/goproxy"
)

/**
 * Searches for a string in the JSON request body
 * Sends a slack notification
 */
func FindInJson(huntType string, huntParam string, reqJsonKeys map[string]struct{}, y *config.YAML, ud UserData) {
	var slackMsg strings.Builder
	var fileMsg strings.Builder

	for jsonKey, _ := range reqJsonKeys {
		forSlack := fmt.Sprintf("*%v* \nREQUEST JSON PARAM: `%v`\nREQ URL: %v \nFILE: `%v` \n", huntType, jsonKey, ud.ReqURL, ud.FileChecksum)
		forFile := fmt.Sprintf("%v \nREQUEST JSON PARAM: %v\nREQ URL: %v, \nFILE: %v \n\n", huntType, jsonKey, ud.ReqURL, ud.FileChecksum)

		constructMsg(jsonKey, huntParam, forSlack, forFile, &slackMsg, &fileMsg, y)
	}

	if fileMsg.String() != "" && y.Storage.Type == "files" {
		utils.WriteUniqueFile(ud.Host, ud.FileChecksum, "", y.Settings.BaseOutputDir, fileMsg.String(), ".hunt")
	}

	if slackMsg.String() != "" && y.Settings.SlackHook != "" {
		utils.SendSlackNotification(y.Settings.SlackHook, slackMsg.String())
	}
}

/**
 * Searches for a string in request query param
 * Sends a slack notification
 */
func FindInQueryParams(huntType string, huntParam string, reqQueryParams map[string][]string, y *config.YAML, ud UserData) {
	var slackMsg strings.Builder
	var fileMsg strings.Builder

	for queryParam := range reqQueryParams {
		forSlack := fmt.Sprintf("*%v* \nREQUEST QUERY PARAM: `%v`\nREQ URL: %v \nFILE: `%v` \n", huntType, queryParam, ud.ReqURL, ud.FileChecksum)
		forFile := fmt.Sprintf("%v \nREQUEST QUERY PARAM: %v \nREQ URL: %v \nFILE: %v \n\n", huntType, queryParam, ud.ReqURL, ud.FileChecksum)

		constructMsg(queryParam, huntParam, forSlack, forFile, &slackMsg, &fileMsg, y)
	}

	if fileMsg.String() != "" && y.Storage.Type == "files" {
		utils.WriteUniqueFile(ud.Host, ud.FileChecksum, "", y.Settings.BaseOutputDir, fileMsg.String(), ".hunt")
	}

	if slackMsg.String() != "" && y.Settings.SlackHook != "" {
		utils.SendSlackNotification(y.Settings.SlackHook, slackMsg.String())
	}
}

/*
 * Construct messages for slack and for the files based on user defined conditions
 */
func constructMsg(reqParam string, huntParam string, forSlack string, forFile string, slackMsg *strings.Builder, fileMsg *strings.Builder, y *config.YAML) {
	if y.Filters.Hunt.ExactMatch && strings.ToLower(reqParam) == strings.ToLower(huntParam) {
		if y.Storage.Type == "files" {
			fileMsg.WriteString(forFile)
		}

		if y.Settings.SlackHook != "" {
			slackMsg.WriteString(forSlack)
		}
	}

	if !y.Filters.Hunt.ExactMatch && strings.Contains(strings.ToLower(reqParam), strings.ToLower(huntParam)) {
		if y.Storage.Type == "files" {
			fileMsg.WriteString(forFile)
		}

		if y.Settings.SlackHook != "" {
			slackMsg.WriteString(forSlack)
		}
	}
}

/**
 * Detects all sorts of keys and secrets using regexs
 */
func detectSecrets(allSecrets *map[string]struct{}, dump []byte, ud UserData, saveSecretsDir string) {
	// taken from https://github.com/xyele/zile
	str := map[string]string{
		"slack_token":                   "(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})",
		"slack_webhook":                 "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
		"facebook_oauth":                "[f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].{0,30}['\"\\s][0-9a-f]{32}['\"\\s]",
		"twitter_oauth":                 "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].{0,30}['\"\\s][0-9a-zA-Z]{35,44}['\"\\s]",
		"twitter_access_token":          "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*[1-9][0-9]+-[0-9a-zA-Z]{40}",
		"heroku_api":                    "[h|H][e|E][r|R][o|O][k|K][u|U].{0,30}[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}",
		"mailgun_api":                   "key-[0-9a-zA-Z]{32}",
		"mailchamp_api":                 "[0-9a-f]{32}-us[0-9]{1,2}",
		"picatic_api":                   "sk_live_[0-9a-z]{32}",
		"google_oauth_id":               "[0-9(+-[0-9A-Za-z_]{32}.apps.googleusercontent.com",
		"google_api":                    "AIza[0-9A-Za-z-_]{35}",
		"google_captcha":                "6L[0-9A-Za-z-_]{38}",
		"google_oauth":                  "ya29\\.[0-9A-Za-z\\-_]+",
		"amazon_aws_access_key_id":      "AKIA[0-9A-Z]{16}",
		"amazon_mws_auth_token":         "amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
		"amazonaws_url":                 "(s3-|s3\\.)?(.*)\\.amazonaws\\.com",
		"facebook_access_token":         "EAACEdEose0cBA[0-9A-Za-z]+",
		"mailgun_api_key":               "key-[0-9a-zA-Z]{32}",
		"twilio_api_key":                "SK[0-9a-fA-F]{32}",
		"twilio_account_sid":            "AC[a-zA-Z0-9_\\-]{32}",
		"twilio_app_sid":                "AP[a-zA-Z0-9_\\-]{32}",
		"paypal_braintree_access_token": "access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}",
		"square_oauth_secret":           "sq0csp-[ 0-9A-Za-z\\-_]{43}",
		"square_access_token":           "sqOatp-[0-9A-Za-z\\-_]{22}",
		"stripe_standard_api":           "sk_live_[0-9a-zA-Z]{24}",
		"stripe_restricted_api":         "rk_live_[0-9a-zA-Z]{24}",
		"github_access_token":           "[a-zA-Z0-9_-]*:[a-zA-Z0-9_\\-]+@github\\.com*",
		"private_ssh_key":               "-----BEGIN PRIVATE KEY-----[a-zA-Z0-9\\S]{100,}-----END PRIVATE KEY-----",
		"private_rsa_key":               "-----BEGIN RSA PRIVATE KEY-----[a-zA-Z0-9\\S]{100,}-----END RSA PRIVATE KEY-----",
		"gpg_private_key_block":         "-----BEGIN PGP PRIVATE KEY BLOCK-----",
		"generic_api_key":               "[a|A][p|P][i|I][_]?[k|K][e|E][y|Y].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"generic_secret":                "[s|S][e|E][c|C][r|R][e|E][t|T].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"ip_address":                    "(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])",
		"urls":                          "(?:\"|'|`)(((?:[a-zA-Z]{1,10}:\\/\\/|\\/\\/)[^\"'\\/]{1,}\\.[a-zA-Z]{2,}[^\"']{0,})|((?:\\/|\\.\\.\\/|\\.\\/)[^\"'><,;| *()(%%$^\\/\\\\\\[\\]][^\"'><,;|()]{1,})|([a-zA-Z0-9_\\-\\/]{1,}\\/[a-zA-Z0-9_\\-\\/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-\\/]{1,}\\/[a-zA-Z0-9_\\-\\/\\$\\{\\}]{3,}(?:[\\?|#][^\"|']{0,}|))|([a-zA-Z0-9_\\-]{1,}\\.(?:php|asp|aspx|jsp|json|action|html|js|txt|xml)(?:[\\?|#][^\"|']{0,}|)))(?:\"|'|`)",
	}

	if !utils.FileExists(saveSecretsDir) {
		os.MkdirAll(saveSecretsDir, os.ModePerm)
	}

	for key, value := range str {
		re := regexp.MustCompile(value)
		matches := re.FindAllString(string(dump), -1)
		allSecretsForType := ""
		if len(matches) > 0 {
			for _, secretMatch := range matches {
				if _, exists := (*allSecrets)[secretMatch]; !exists {
					allSecretsForType += fmt.Sprintf("%v\n", secretMatch)
				}
				(*allSecrets)[secretMatch] = struct{}{}
			}
		}
		if allSecretsForType != "" {
			go utils.AppendToFile(allSecretsForType, fmt.Sprintf("%v/%v", saveSecretsDir, key))
		}
	}
}

/**
 * Block or allow responses which contain one of the passed file types.
 * If 'shouldBlock' param is false, it will allow all given file types.
 */
func respFileType(shouldBlock bool, fileTypes ...string) goproxy.RespConditionFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		if fileTypes == nil {
			return true
		}

		url := resp.Request.URL.Path
		for _, ft := range fileTypes {
			if strings.Contains(strings.ToLower(url), ft) {
				return !shouldBlock
			}
		}

		return shouldBlock
	}
}

/**
 * Block or allow requests which contain one of the passed file types.
 * If 'shouldBlock' param is false, it will allow all given file types.
 */
func reqFileType(shouldBlock bool, fileTypes ...string) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		if fileTypes == nil {
			return true
		}

		url := req.URL.Path
		for _, ft := range fileTypes {
			if strings.Contains(strings.ToLower(url), ft) {
				return !shouldBlock
			}
		}

		return shouldBlock
	}
}

/**
 * Block or allow responses which contain one of the passed content types.
 * If 'shouldBlock' param is false, it will apply a allow filter for the given file types.
 */
func respContentType(shouldBlock bool, contentTypes ...string) goproxy.RespConditionFunc {
	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
		if contentTypes == nil {
			return true
		}

		contentType := resp.Header.Get("Content-Type")
		for _, ct := range contentTypes {
			if strings.Contains(strings.ToLower(contentType), ct) {
				return !shouldBlock
			}
		}

		return shouldBlock
	}
}

/**
 * Block or allow requests which contain one of the passed content types.
 * If 'shouldBlock' param is false, it will apply a allow filter for the given file types.
 */
func reqContentType(shouldBlock bool, contentTypes ...string) goproxy.ReqConditionFunc {
	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
		if contentTypes == nil {
			return true
		}

		contentType := req.Header.Get("Content-Type")
		for _, ct := range contentTypes {
			if strings.Contains(strings.ToLower(contentType), ct) {
				return !shouldBlock
			}
		}

		return shouldBlock
	}
}

// /**
//  * For all responses - Include all regexes from inScope param, from config.
//  */
// func includeInRespScope(inScopeRegexes []string) goproxy.RespCondition {
// 	return goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(inScopeRegexes, ")|("))))
// }

// /**
//  * For all responses - Exclude all regexes from outScope param, from config.
//  */
// func excludeFromRespScope(outScopeRegexes []string) goproxy.RespConditionFunc {
// 	out := regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(outScopeRegexes, ")|(")))
// 	return func(resp *http.Response, ctx *goproxy.ProxyCtx) bool {
// 		return !(out.MatchString(resp.Request.URL.Path) || out.MatchString(resp.Request.URL.Host+resp.Request.URL.Path))
// 	}
// }

// /**
//  * For all requrests - Include all regexes from inScope param, from config.
//  */
// func includeInReqScope(inScopeRegexes []string) goproxy.ReqCondition {
// 	return goproxy.UrlMatches(regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(inScopeRegexes, ")|("))))
// }

// /**
//  * For all requests - Exclude all regexes from outScope param, from config.
//  */
// func excludeFromReqScope(outScopeRegexes []string) goproxy.ReqConditionFunc {
// 	out := regexp.MustCompile(fmt.Sprintf("(%v)", strings.Join(outScopeRegexes, ")|(")))
// 	return func(req *http.Request, ctx *goproxy.ProxyCtx) bool {
// 		return !(out.MatchString(req.URL.Path) || out.MatchString(req.URL.Host+req.URL.Path))
// 	}
// }

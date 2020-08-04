package filters

import (
	"fmt"
	"strings"

	"github.com/ctoyan/ponieproxy/internal/config"
	"github.com/ctoyan/ponieproxy/internal/utils"
)

/*
 * Searches for a string in the JSON request body
 * Sends a slack notification
 */
func FindInJson(huntType string, huntParam string, reqJsonKeys map[string]struct{}, flags *config.Flags, ud UserData) {
	var slackMsg strings.Builder
	var fileMsg strings.Builder

	for jsonKey, _ := range reqJsonKeys {
		forSlack := fmt.Sprintf("*%v* \nREQUEST JSON PARAM: `%v` \nFILE: `%v` \n", huntType, jsonKey, ud.FileChecksum)
		forFile := fmt.Sprintf("%v \nREQUEST JSON PARAM: %v \nFILE: %v \n\n", huntType, jsonKey, ud.FileChecksum)

		constructMsg(jsonKey, huntParam, forSlack, forFile, &slackMsg, &fileMsg, flags)
	}

	if fileMsg.String() != "" && flags.HuntOutputFile {
		utils.WriteUniqueFile(ud.Host, ud.FileChecksum, "", flags.OutputDir, fileMsg.String(), "hunt")
	}

	if slackMsg.String() != "" && flags.SlackWebHook != "" {
		utils.SendSlackNotification(flags.SlackWebHook, slackMsg.String())
	}
}

/*
 * Searches for a string in request query param
 * Sends a slack notification
 */
func FindInQueryParams(huntType string, huntParam string, reqQueryParams map[string][]string, flags *config.Flags, ud UserData) {
	var slackMsg strings.Builder
	var fileMsg strings.Builder

	for queryParam := range reqQueryParams {
		forSlack := fmt.Sprintf("*%v* \nREQUEST QUERY PARAM: `%v` \nFILE: `%v` \n", huntType, queryParam, ud.FileChecksum)
		forFile := fmt.Sprintf("%v \nREQUEST QUERY PARAM: %v \nFILE: %v \n\n", huntType, queryParam, ud.FileChecksum)

		constructMsg(queryParam, huntParam, forSlack, forFile, &slackMsg, &fileMsg, flags)
	}

	if fileMsg.String() != "" && flags.HuntOutputFile {
		utils.WriteUniqueFile(ud.Host, ud.FileChecksum, "", flags.OutputDir, fileMsg.String(), "hunt")
	}

	if slackMsg.String() != "" && flags.SlackWebHook != "" {
		utils.SendSlackNotification(flags.SlackWebHook, slackMsg.String())
	}
}

/*
 * Construct messages for slack and for the files based on user defined conditions
 */
func constructMsg(reqParam string, huntParam string, forSlack string, forFile string, slackMsg *strings.Builder, fileMsg *strings.Builder, flags *config.Flags) {
	if flags.HuntExactMatch && strings.ToLower(reqParam) == strings.ToLower(huntParam) {
		if flags.HuntOutputFile {
			fileMsg.WriteString(forFile)
		}

		if flags.SlackWebHook != "" {
			slackMsg.WriteString(forSlack)
		}
	}

	if !flags.HuntExactMatch && strings.Contains(strings.ToLower(reqParam), strings.ToLower(huntParam)) {
		if flags.HuntOutputFile {
			fileMsg.WriteString(forFile)
		}

		if flags.SlackWebHook != "" {
			slackMsg.WriteString(forSlack)
		}
	}
}

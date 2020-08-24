package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

/*
 * Unmarshal the request body
 */
func UnmarshalReqBody(reqBody string) (unmarshaled map[string]interface{}, err error) {
	unmarshaled = map[string]interface{}{}
	err = json.Unmarshal([]byte(reqBody), &unmarshaled)
	if err != nil {
		return nil, err
	}

	return unmarshaled, nil
}

/*
 * Recursively collects keys from JSON object
 * Returns a collection of unique keys as map[string]struct{}
 */
func CollectJsonKeys(json map[string]interface{}, reqJsonKeys map[string]struct{}) map[string]struct{} {
	for k, v := range json {
		if v == nil {
			return reqJsonKeys
		}
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		// if value is a slice of map[string]interface{} iterate over it again
		case reflect.Slice:
			s := reflect.ValueOf(v)
			for i := 0; i < s.Len(); i++ {
				el, ok := s.Index(i).Interface().(map[string]interface{})
				if ok {
					CollectJsonKeys(el, reqJsonKeys)
				}
			}
		// if not - all keys of other value types are added
		default:
			if _, ok := reqJsonKeys[k]; !ok {
				reqJsonKeys[k] = struct{}{}
			}
		}
	}

	return reqJsonKeys
}

/*
 * Creates a uniquely named file and write http dump to it.
 * The file name represents a hashed host, path and request body.
 * This means that any request and it's resposne are named with the same hash.
 * This makes it easy to go through and read them, when opened with "vim *"
 */
func WriteUniqueFile(host string, checksum string, body string, outputDir string, httpDump string, ext string) {
	folderPath := fmt.Sprintf("%v/%v", outputDir, host)
	filePath := fmt.Sprintf("%v/%v.%v", folderPath, checksum, ext)

	if !FileExists(folderPath) {
		os.MkdirAll(folderPath, os.ModePerm)
	}

	if !FileExists(filePath) {
		var constructed string
		if ext == "req" {
			constructed = fmt.Sprintf(`%v %v`, httpDump, body)
		}
		if ext == "res" || ext == "hunt" {
			constructed = fmt.Sprintf(`%v`, httpDump)
		}

		err := AppendToFile(constructed, filePath)
		if err != nil {
			log.Fatalf("error writing to file: %v", err)
		}
	}
}

/*
 * Takes file path and returns lines
 */
func ReadLines(path string) ([]string, error) {
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

/*
 * Takes data and writes it to a file
 */
func AppendToFile(data string, filePath string) error {
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

/*
 * Check if a file exists
 */
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

/*
 * Send a slack notification via webhook
 */
func SendSlackNotification(webhookUrl string, msg string) error {
	slackBody, _ := json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: msg,
	})

	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

/*
 * Download a file from a url into a folder
 */
func DownloadFromURL(url *url.URL, folder string) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileLocation := fmt.Sprintf("%v%v", folder, url.Path)

	if !FileExists(filepath.Dir(fileLocation)) {
		os.MkdirAll(filepath.Dir(fileLocation), 0770)
	}

	out, err := os.Create(fileLocation)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

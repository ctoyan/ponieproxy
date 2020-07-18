package utils

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"log"
	"os"
)

/*
 * Creates a uniquely named file, based on the host, path and req/res body.
 * The file is named, based on the hashed host, path and request body.
 * This means that any request and it's resposne are named with the same hash.
 * This makes it easy to go through and read them, when opened with "vim *"
 */
func WriteUniqueFile(host string, path string, body []byte, outputDir string, httpDump []byte, ext string) {
	if outputDir != "./" {
		os.MkdirAll(outputDir, os.ModePerm)
	}

	cacheKey := fmt.Sprintf("%s%s%s", host, path, body)
	hashed := sha1.Sum([]byte(cacheKey))

	filePath := fmt.Sprintf("%v/%x.%v", outputDir, hashed, ext)

	if !FileExists(filePath) {
		var constructed string
		if ext == "req" {
			constructed = fmt.Sprintf(`%s %s`, httpDump, body)
		}
		if ext == "res" {
			constructed = fmt.Sprintf(`%s`, httpDump)
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

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := regexp.MustCompile("@")
	seenBrowsers := make(map[string]bool)
	uniqueBrowsers := 0
	foundUsers := ""

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		lineCount++

		user := make(map[string]interface{})
		err := json.Unmarshal(line, &user)
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			continue
		}

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}

			if strings.Contains(browser, "Android") {
				isAndroid = true

				if !seenBrowsers[browser] {
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true

				if !seenBrowsers[browser] {
					seenBrowsers[browser] = true
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := r.ReplaceAllString(user["email"].(string), " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", lineCount-1, user["name"], email)
	}

	fmt.Fprintln(out, "found users:")
	fmt.Fprintln(out, foundUsers)
	fmt.Fprintln(out, "Total unique browsers", uniqueBrowsers)
}

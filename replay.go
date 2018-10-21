package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"./config"
)

func fireGetRequestToURL(url string) (bool, int, error) {
	response, err := http.Get(url)

	if err != nil {
		return false, 0, err
	}

	if response.StatusCode == 200 {
		return true, 200, nil
	}

	return false, response.StatusCode, errors.New("Unknown")
}

func getTimestampFromLogLine(logLine string) int64 {

	startIndex := strings.Index(logLine, "[")
	startIndex = startIndex + 1

	endIndex := strings.Index(logLine, "]")

	timeStamp := logLine[startIndex:endIndex]

	nginxLogDateFormat := "02/Jan/2006:15:04:05 -0700"

	parsedTime, err := time.Parse(nginxLogDateFormat, timeStamp)

	if err != nil {
		log.Fatalln(err)
	}

	parsedTimeInUTC := parsedTime.UTC()

	return (parsedTimeInUTC.Unix() * 1000)
}

func main() {

	config := config.InitializeConfig()

	if config == nil {
		os.Exit(1)
	}

	successPtr, err := os.OpenFile("succeeded.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer successPtr.Close()

	failurePtr, err := os.OpenFile("failed.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer failurePtr.Close()

	failureReqsPtr, err := os.OpenFile("reqs-failed.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer failureReqsPtr.Close()

	// Input
	file, err := os.Open(config.LogFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	total := 0
	succeeded := 0
	failed := 0

	for scanner.Scan() {

		s := scanner.Text()

		if config.RegexFilterEnabled {
			filterMatch := config.RegexFilter.MatchString(s)

			if !filterMatch {
				continue
			}
		}

		if config.RegexExcludeEnabled {
			excludeMatch := config.RegexExclude.MatchString(s)

			if excludeMatch {
				continue
			}
		}

		firstPosParams := strings.Index(s, "GET /track")
		firstPosParams = firstPosParams + 4

		lastPosParams := strings.Index(s, " HTTP/")

		if firstPosParams == -1 || lastPosParams == -1 {
			fmt.Println(s)
			continue
		}

		if lastPosParams <= firstPosParams {
			fmt.Println(s)
			continue
		}

		url := config.BaseURL + s[firstPosParams:lastPosParams]

		if config.IncludeTimeStamp {
			utcTimestamp := getTimestampFromLogLine(s)
			url = url + "&timestamp=" + strconv.FormatInt(utcTimestamp, 10)
		}

		success := true
		code := 200
		err = nil

		if config.DryRun {
			success, code, err = fireGetRequestToURL(url)
		} else {
			fmt.Println(url)
		}

		if success {
			succeeded++

			if _, err2 := successPtr.WriteString(s + "\n"); err2 != nil {
				fmt.Println(s)
				fmt.Println(err2)
			}
		} else {
			failed++

			if _, err2 := failurePtr.WriteString(s + "\n"); err2 != nil {
				fmt.Println(s)
				fmt.Println(err2)
			}

			requestFailureLog := ""

			requestFailureLog = strconv.Itoa(code) + ",\"" + err.Error() + "\"," + url + "\n"

			if _, err2 := failureReqsPtr.WriteString(requestFailureLog); err2 != nil {
				fmt.Println(requestFailureLog)
				fmt.Println(err2)
			}
		}

		time.Sleep(time.Duration(2) * time.Millisecond)

		total++
	}

	fmt.Println("succeeded: ", succeeded)
	fmt.Println("failed: ", failed)
	fmt.Println("total: ", total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

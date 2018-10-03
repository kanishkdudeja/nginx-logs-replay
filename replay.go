package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func readCommandLineParams() (string, string, string, string) {
	var baseURL string
	var logFilePath string
	var dryRun string
	var withTimestamp string

	flag.StringVar(&baseURL, "base-url", "uninitialized", "Denotes the host name to which requests will be replayed. Eg: https://website.com / 1.1.1.1")
	flag.StringVar(&logFilePath, "file", "uninitialized", "Denotes the path at which the log file is present. Eg: /var/log/nginx/access.log")
	flag.StringVar(&dryRun, "dry-run", "uninitialized", "Denotes whether it's a dry run or not")
	flag.StringVar(&withTimestamp, "with-timestamp", "uninitialized", "Denotes whether we need to send the UNIX timestamp along with the URL")

	flag.Parse()

	if baseURL == "uninitialized" {
		log.Fatalln("Please supply the baseURL (with http/https) as a parameter. Eg: ./replay --base-url=https://website.com")
	}

	if logFilePath == "uninitialized" {
		log.Fatalln("Please supply the path of the log file as a parameter. Eg: ./replay --file=/var/log/nginx/access.log")
	}

	if dryRun == "uninitialized" {
		log.Fatalln("Please supply the dry-run parameter. Pass as 'true' if you want the script to only print the URLs. Eg: ./replay --dry-run=true/false")
	}

	if dryRun != "true" && dryRun != "false" {
		log.Fatalln("The dry-run parameter can only have a value of true/false. Eg: ./replay --dry-run=true/false")
	}

	if withTimestamp != "uninitialized" && withTimestamp != "true" && withTimestamp != "false" {
		log.Fatalln("The with-timestamp parameter can only have a value of true/false. Eg: ./replay --with-timestamp=true/false")
	}

	return baseURL, logFilePath, dryRun, withTimestamp
}

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

	baseURL, logFilePath, dryRun, withTimestamp := readCommandLineParams()

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
	file, err := os.Open(logFilePath)
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

		url := baseURL + s[firstPosParams:lastPosParams]

		if withTimestamp == "true" {
			utcTimestamp := getTimestampFromLogLine(s)
			url = url + "&timestamp=" + strconv.FormatInt(utcTimestamp, 10)
		}

		success := true
		code := 200
		err = nil

		if dryRun == "false" {
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

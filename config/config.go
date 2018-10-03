// Package config provides configuration variables for the application
// Values of these variables depend on the environment in which the code is being run right now
package config

import (
	"errors"
	"flag"
)

type Config struct {
	DryRun           bool   `json:"dry_run"`
	BaseURL          string `json:"base_url"`
	LogFilePath      string `json:"log_file_path"`
	IncludeTimeStamp bool   `json:"include_time_stamp"`
}

// The InitializeConfig() function returns the command line flags and calls the config package to load up the configuration
func InitializeConfig() (*Config, error) {
	var dryRun string
	var baseURL string
	var logFilePath string
	var includeTimeStamp string

	flag.StringVar(&dryRun, "dry-run", "uninitialized", "Denotes whether it's a dry run or not")
	flag.StringVar(&baseURL, "base-url", "uninitialized", "Denotes the host name to which requests will be replayed. Eg: https://website.com / 1.1.1.1")
	flag.StringVar(&logFilePath, "file", "uninitialized", "Denotes the path at which the log file is present. Eg: /var/log/nginx/access.log")
	flag.StringVar(&includeTimeStamp, "include-timestamp", "uninitialized", "Denotes whether we need to send the UNIX timestamp along with the URL")

	flag.Parse()

	if baseURL == "uninitialized" {
		return nil, errors.New("Please supply the baseURL (with http/https) as a parameter. Eg: ./replay --base-url=https://website.com")
	}

	if logFilePath == "uninitialized" {
		return nil, errors.New("Please supply the path of the log file as a parameter. Eg: ./replay --file=/var/log/nginx/access.log")
	}

	if dryRun == "uninitialized" {
		return nil, errors.New("Please supply the dry-run parameter. Pass as 'true' if you want the script to only print the URLs. Eg: ./replay --dry-run=true/false")
	}

	if dryRun != "true" && dryRun != "false" {
		return nil, errors.New("The dry-run parameter can only have a value of true/false. Eg: ./replay --dry-run=true/false")
	}

	if includeTimeStamp != "uninitialized" && includeTimeStamp != "true" && includeTimeStamp != "false" {
		return nil, errors.New("The include-timestamp parameter can only have a value of true/false. Eg: ./replay --include-timestamp=true/false")
	}

	var configObj Config

	if dryRun == "true" {
		configObj.DryRun = true
	} else {
		configObj.DryRun = false
	}

	configObj.BaseURL = baseURL
	configObj.LogFilePath = logFilePath

	if includeTimeStamp == "true" {
		configObj.IncludeTimeStamp = true
	} else {
		configObj.IncludeTimeStamp = false
	}

	return &configObj, nil
}

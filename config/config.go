/*
Package config provides configuration variables
for the application.

Values of these variables depend on the environment in
which the code is being run right now
*/
package config

import (
	"errors"
	"flag"

	"../utils"
)

// Config Struct representing the config
type Config struct {
	DryRun           bool   `json:"dry_run"`
	BaseURL          string `json:"base_url"`
	LogFilePath      string `json:"log_file_path"`
	IncludeTimeStamp bool   `json:"include_time_stamp"`
}

// InitializeConfig Returns the command line flags
// and calls the config package to load up the configuration
func InitializeConfig() (*Config, error) {
	var dryRun bool
	var baseURL string
	var logFilePath string
	var includeTimeStamp bool

	flag.BoolVar(&dryRun, "dry-run", false, "Denotes whether it's a dry run or not")
	flag.StringVar(&baseURL, "base-url", "uninitialized", "Denotes the host name to which requests will be replayed. Eg: https://website.com / 1.1.1.1")
	flag.StringVar(&logFilePath, "log-file-path", "uninitialized", "Denotes the path at which the log file is present. Eg: /var/log/nginx/access.log")
	flag.BoolVar(&includeTimeStamp, "include-timestamp", false, "Denotes whether we need to send the UNIX timestamp along with the URL")

	flag.Parse()

	if baseURL == "uninitialized" {
		return nil, errors.New("Please supply the baseURL (with http/https) as a parameter. Eg: ./replay --base-url=https://website.com")
	}

	if logFilePath == "uninitialized" {
		return nil, errors.New("Please supply the path of the log file as a parameter. Eg: ./replay --log-file-path=/var/log/nginx/access.log")
	}

	err := utils.ValidateBaseURL(baseURL)

	if err != nil {
		return nil, err
	}

	var configObj Config

	configObj.DryRun = dryRun
	configObj.BaseURL = baseURL
	configObj.LogFilePath = logFilePath
	configObj.IncludeTimeStamp = includeTimeStamp

	return &configObj, nil
}

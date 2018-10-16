/*
Package config provides configuration variables
for the application.

Values of these variables depend on the environment in
which the code is being run right now
*/
package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"../utils"
)

// Config Struct representing the config
type Config struct {
	DryRun              bool
	BaseURL             string
	LogFilePath         string
	RegexFilterEnabled  bool
	RegexFilter         *regexp.Regexp
	RegexExcludeEnabled bool
	RegexExclude        *regexp.Regexp
	IncludeTimeStamp    bool
}

// InitializeConfig Returns the command line flags
// and calls the config package to load up the configuration
func InitializeConfig() *Config {
	var help bool
	var dryRun bool
	var baseURL string
	var logFilePath string
	var regexFilter string
	var regexExclude string
	var includeTimeStamp bool
	var regexpFilter, regexpExclude *regexp.Regexp

	flag.BoolVar(&help, "help", false, "Prints the usage string for the program")
	flag.BoolVar(&dryRun, "dry-run", false, "Denotes whether it's a dry run or not")
	flag.StringVar(&baseURL, "base-url", "", "Denotes the host name to which requests will be replayed. Eg: https://website.com / 1.1.1.1")
	flag.StringVar(&logFilePath, "log-file-path", "", "Denotes the path at which the log file is present. Eg: /var/log/nginx/access.log")
	flag.StringVar(&regexFilter, "regex-filter", "", "Denotes the Regex pattern to filter logs to replay. Eg: '/abc/'")
	flag.StringVar(&regexExclude, "regex-exclude", "", "Denotes the Regex pattern to filter logs to exclude replaying. Eg: '/abc/'")
	flag.BoolVar(&includeTimeStamp, "include-timestamp", false, "Denotes whether we need to send the UNIX timestamp along with the URL")

	flag.Parse()

	if len(os.Args) == 1 {
		printError("Please supply a configuration parameter for the program")
		flag.Usage()
		return nil
	}

	if help {
		flag.Usage()
		return nil
	}

	if baseURL == "" {
		printError("Please supply the baseURL (with http/https) as a parameter. Eg: ./replay --base-url=https://website.com")
		return nil
	}

	if logFilePath == "" {
		printError("Please supply the path of the log file as a parameter. Eg: ./replay --log-file-path=/var/log/nginx/access.log")
		return nil
	}

	err := utils.ValidateBaseURL(baseURL)

	if err != nil {
		printError(err.Error())
		return nil
	}

	if regexFilter != "" && regexExclude != "" {
		printError("You can only use one of the --regex-filter and --regex-exclude parameters at once.")
		return nil
	}

	if regexFilter != "" {
		regexpFilter, err = utils.CompileRegularExpression(regexFilter)

		if err != nil {
			printError("Encountered error in compiling regular expression passed in the --regex-filter parameter")
			printError(err.Error())
			return nil
		}
	}

	if regexExclude != "" {
		regexpExclude, err = utils.CompileRegularExpression(regexExclude)

		if err != nil {
			printError("Encountered error in compiling regular expression passed in the --regex-exclude parameter")
			printError(err.Error())
			return nil
		}
	}

	var configObj Config

	configObj.DryRun = dryRun
	configObj.BaseURL = baseURL
	configObj.LogFilePath = logFilePath

	if regexFilter != "" {
		configObj.RegexFilterEnabled = true
		configObj.RegexFilter = regexpFilter
	}

	if regexExclude != "" {
		configObj.RegexExcludeEnabled = true
		configObj.RegexExclude = regexpExclude
	}

	configObj.IncludeTimeStamp = includeTimeStamp

	return &configObj
}

func printError(message string) {
	fmt.Println(message)
}

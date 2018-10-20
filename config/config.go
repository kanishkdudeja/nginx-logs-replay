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
	"fmt"
	"os"
	"regexp"

	"../utils"
)

// Config Struct representing the config
type Config struct {
	Help                bool
	DryRun              bool
	BaseURL             string
	LogFilePath         string
	RegexFilterEnabled  bool
	RegexFilterString   string
	RegexFilter         *regexp.Regexp
	RegexExcludeEnabled bool
	RegexExcludeString  string
	RegexExclude        *regexp.Regexp
	IncludeTimeStamp    bool
}

var configObj Config

func (config *Config) parseAndSetFlags() {
	var help bool
	var dryRun bool
	var baseURL string
	var logFilePath string
	var regexFilter string
	var regexExclude string
	var includeTimeStamp bool

	flag.BoolVar(&help, "help", false, "Prints the usage string for the program")
	flag.BoolVar(&dryRun, "dry-run", false, "Denotes whether it's a dry run or not")
	flag.StringVar(&baseURL, "base-url", "", "Denotes the host name to which requests will be replayed. Eg: https://website.com / 1.1.1.1")
	flag.StringVar(&logFilePath, "log-file-path", "", "Denotes the path at which the log file is present. Eg: /var/log/nginx/access.log")
	flag.StringVar(&regexFilter, "regex-filter", "", "Denotes the Regex pattern to filter logs to replay. Eg: '/abc/'")
	flag.StringVar(&regexExclude, "regex-exclude", "", "Denotes the Regex pattern to filter logs to exclude replaying. Eg: '/abc/'")
	flag.BoolVar(&includeTimeStamp, "include-timestamp", false, "Denotes whether we need to send the UNIX timestamp along with the URL")

	flag.Parse()

	config.Help = help
	config.DryRun = dryRun
	config.BaseURL = baseURL
	config.LogFilePath = logFilePath
	config.RegexFilterString = regexFilter
	config.RegexExcludeString = regexExclude
	config.IncludeTimeStamp = includeTimeStamp
}

func (config *Config) validateConfig() error {
	if len(os.Args) == 1 {
		printError("Please supply a configuration parameter for the program")
		flag.Usage()
		return errors.New("no-config-parameters-supplied")
	}

	if config.Help {
		flag.Usage()
		return errors.New("help-required")
	}

	if config.BaseURL == "" {
		printError("Please supply the baseURL (with http/https) as a parameter. Eg: ./replay --base-url=https://website.com")
		return errors.New("base-url-not-supplied")
	}

	if config.LogFilePath == "" {
		printError("Please supply the path of the log file as a parameter. Eg: ./replay --log-file-path=/var/log/nginx/access.log")
		return errors.New("log-file-path-not-supplied")
	}

	err := utils.ValidateBaseURL(config.BaseURL)

	if err != nil {
		printError(err.Error())
		return err
	}

	if config.RegexFilterString != "" && config.RegexExcludeString != "" {
		printError("You can only use one of the --regex-filter and --regex-exclude parameters at once.")
		return errors.New("attempts-to-use-both-regex-options-at-once")
	}

	var regexpFilter, regexpExclude *regexp.Regexp

	if config.RegexFilterString != "" {
		regexpFilter, err = utils.CompileRegularExpression(config.RegexFilterString)

		if err != nil {
			printError("Encountered error in compiling regular expression passed in the --regex-filter parameter")
			printError(err.Error())
			return errors.New("error-in-compiling-regex-filter-expression")
		}

		config.RegexFilterEnabled = true
		config.RegexFilter = regexpFilter
	}

	if config.RegexExcludeString != "" {
		regexpExclude, err = utils.CompileRegularExpression(config.RegexExcludeString)

		if err != nil {
			printError("Encountered error in compiling regular expression passed in the --regex-exclude parameter")
			printError(err.Error())
			return errors.New("error-in-compiling-regex-exclude-expression")
		}

		config.RegexExcludeEnabled = true
		config.RegexExclude = regexpExclude
	}

	return nil
}

// InitializeConfig Returns the command line flags
// and calls the config package to load up the configuration
func InitializeConfig() *Config {

	configObj.parseAndSetFlags()
	err := configObj.validateConfig()

	if err != nil {
		return nil
	}

	return &configObj
}

func printError(message string) {
	fmt.Println(message)
}

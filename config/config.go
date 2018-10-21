/*
Package config provides functionality related to configuration
parameters for this program.

That includes parsing, validating and storing configuration parameters
supplied by the user at the time of running the program
*/
package config

import (
	"errors"
	"flag"
	"fmt"
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

var configuration Config

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

	if config.BaseURL == "" {
		return errors.New("Please supply the baseURL (with http/https) as a parameter. Eg: ./replay --base-url=https://website.com")
	}

	if config.LogFilePath == "" {
		return errors.New("Please supply the path of the log file as a parameter. Eg: ./replay --log-file-path=/var/log/nginx/access.log")
	}

	err := utils.ValidateBaseURL(config.BaseURL)

	if err != nil {
		return errors.New("Encountered error in validating --base-url parameter: " + err.Error())
	}

	if config.RegexFilterString != "" && config.RegexExcludeString != "" {
		return errors.New("You can only use one of the --regex-filter and --regex-exclude parameters at once")
	}

	var regexpFilter, regexpExclude *regexp.Regexp

	if config.RegexFilterString != "" {
		regexpFilter, err = utils.CompileRegularExpression(config.RegexFilterString)

		if err != nil {
			return errors.New("Encountered error in compiling regular expression passed in the --regex-filter parameter: " + err.Error())
		}

		config.RegexFilterEnabled = true
		config.RegexFilter = regexpFilter
	}

	if config.RegexExcludeString != "" {
		regexpExclude, err = utils.CompileRegularExpression(config.RegexExcludeString)

		if err != nil {
			return errors.New("Encountered error in compiling regular expression passed in the --regex-exclude parameter: " + err.Error())
		}

		config.RegexExcludeEnabled = true
		config.RegexExclude = regexpExclude
	}

	return nil
}

// InitializeConfig Returns the command line flags
// and calls the config package to load up the configuration
func InitializeConfig() *Config {

	configuration.parseAndSetFlags()

	// Print usage string if no configuration parameters are supplied
	if flag.NFlag() == 0 {
		printMessage("Please supply a configuration parameter for the program")
		flag.Usage()
		return nil
	}

	// Print usage string if user has supplied the --help configuration parameter
	if configuration.Help {
		flag.Usage()
		return nil
	}

	// Validate configuration
	err := configuration.validateConfig()

	// Show error if validation of configuration parameters failed
	if err != nil {
		printMessage(err.Error())
		return nil
	}

	return &configuration
}

func printMessage(message string) {
	fmt.Println(message)
}

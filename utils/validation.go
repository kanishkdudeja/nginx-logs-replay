package utils

import (
	"errors"
	"net/url"
	"regexp"
)

// ValidateBaseURL validates if the BaseURL (provided as a
// part of the configuration) is valid or not.
//
// To check if it's valid, it makes sure that it contains either http:// or https://
// It also check that the base URL should contain a valid hostname (domain name or IP address)
func ValidateBaseURL(baseURL string) error {
	url, err := url.Parse(baseURL)

	if err != nil {
		return err
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return errors.New("Base URL's scheme must be either HTTP/HTTPS")
	}

	if url.Host == "" {
		return errors.New("Base URL must have a hostname (can be either a domain name or an IP address)")
	}

	return nil
}

// CompileRegularExpression validates if the provided regular expression
// is valid or not.
func CompileRegularExpression(regularExpression string) (*regexp.Regexp, error) {
	regexp, err := regexp.Compile(regularExpression)

	return regexp, err
}

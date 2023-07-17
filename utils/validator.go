package utils

import (
	"html/template"
	"regexp"
)

// validator
func IsSafeFileName(name string) bool {
	return isMatchRegex("^[a-zA-Z0-9-_\\.]+$", name)
}

func IsPath(url string) bool {
	return isMatchRegex("[\\w-]+(/[\\w-./?%&=]*)?", url)
}

func IsChinesePhone(phone string) bool {
	return isMatchRegex("^1\\d{10}$", phone)
}

func isMatchRegex(pattern string, v string) bool {
	matched, err := regexp.MatchString(pattern, v)
	if err != nil {
		return false
	}

	return matched
}

// filter
func XSSFilter(input string) (output string) {
	return template.HTMLEscapeString(input)
}

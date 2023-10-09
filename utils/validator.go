package utils

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// number
	RePositiveInterger           = "^[1-9]\\d*$"
	RePositiveScientificNotation = "^(\\d+(.{0}|.\\d+))[Ee]{1}([\\+|-]?\\d+)$"
	RePositiveFloatPoint         = "^(?:[1-9][0-9]*\\.[0-9]+|0\\.(?!0+$)[0-9]+)$"

	// file&path
	ReURL      = "[\\w-]+(/[\\w-./?%&=]*)?"
	ReFileName = "^[a-zA-Z0-9-_\\.]+$"

	// phone
	ReChinesePhone = "^1\\d{10}$"

	// name
	ReUserName = "^[a-zA-Z0-9_-]+$"
)

// validator
func IsPositiveInterger(num string) bool {
	return isMatchRegex(RePositiveInterger, num)
}

func IsPositiveScientificNotation(num string) bool {
	return isMatchRegex(RePositiveScientificNotation, num)
}

func IsPositiveFloatPoint(num string) bool {
	return isMatchRegex(RePositiveFloatPoint, num)
}

func IsSafeFileName(name string) bool {
	return isMatchRegex(ReFileName, name)
}

func IsPath(url string) bool {
	return isMatchRegex(ReURL, url)
}

func IsChinesePhone(phone string) bool {
	return isMatchRegex(ReChinesePhone, phone)
}

func IsUserName(name string) bool {
	if length := StrLen(name); length > 20 || length < 3 {
		return false
	}

	return isMatchRegex(ReUserName, name)
}

func IsPictureName(pictureName string) bool {
	ext := filepath.Ext(pictureName)
	ext = strings.ToLower(ext)

	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	allowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	return allowed
}

func IsTxt(fileName string) bool {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)

	allowedExtensions := []string{".txt"}
	allowed := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	return allowed
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

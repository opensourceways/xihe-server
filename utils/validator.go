package utils

import (
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
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

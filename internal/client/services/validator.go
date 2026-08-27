package services

import (
	"regexp"
	"unicode"
)

func (s *Service) ValidateEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (s *Service) ValidateStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper, hasLower, hasDigit := false, false, false

	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
		}
		if unicode.IsLower(c) {
			hasLower = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

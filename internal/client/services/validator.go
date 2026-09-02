package services

import (
	"regexp"
	"unicode"
)

var emailRe = regexp.MustCompile(`^[\w.+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)

func (s *Service) ValidateEmailFormat(email string) bool {
	return emailRe.MatchString(email)
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

// Проверка на Luhn
func (s *Service) ValidatePAN(pan string) bool {
	if len(pan) < 13 || len(pan) > 19 {
		return false
	}

	sum := 0
	double := false

	for i := len(pan) - 1; i >= 0; i-- {
		if pan[i] < '0' || pan[i] > '9' {
			return false
		}

		digit := int(pan[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}

func (s *Service) ValidateCVC(cvc string) bool {
	if len(cvc) != 3 {
		return false
	}

	for i := range cvc {
		if cvc[i] < '0' || cvc[i] > '9' {
			return false
		}
	}

	return true
}

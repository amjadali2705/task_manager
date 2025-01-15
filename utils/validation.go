package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func ValidateDetails(name, email string, mobile int64, gender, password, confirmPassword string) error {

	//validate name
	if len(name) < 3 {
		return errors.New("name should be atleast 3 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name) {
		return errors.New("name should contain only alphabets and spaces")
	}

	//validate email
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return errors.New("invalid email")
	}

	//validate mobile
	mobilestr := strconv.FormatInt(mobile, 10)
	if len(mobilestr) != 10 || !regexp.MustCompile((`^\d{10}$`)).Match([]byte(mobilestr)) {
		return errors.New("mobile number should be 10 digits")
	}

	//validate gender
	gender = strings.ToLower(gender)
	if gender != "male" && gender != "female" && gender != "other" {
		return errors.New("gender must be 'male', 'female', or 'other'")
	}

	//validate password
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return errors.New("password should be atleast 8 characters")
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password should contain atleast one uppercase, one lowercase, one digit, and one special character")
	}

	//validate confirm password
	if password != confirmPassword {
		return errors.New("password and confirm password do not match")
	}

	return nil
}

func ValidateUser(name string, mobileno int64) error {
	// Validate name
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}

	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name) {
		return errors.New("name must contain only alphabets and spaces")
	}

	// Validate mobile
	mobileStr := strconv.FormatInt(mobileno, 10)
	if len(mobileStr) != 10 || !regexp.MustCompile(`^\d{10}$`).Match([]byte(mobileStr)) {
		return errors.New("mobile number must be 10 digits")
	}

	return nil
}

func ValidatePassword(password string) error {
	//validate password
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return errors.New("password should be atleast 8 characters")
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password should contain atleast one uppercase, one lowercase, one digit, and one special character")
	}
	return nil
}

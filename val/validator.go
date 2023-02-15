package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

// ValidateString Function named ValidateString responsible for validating if a string has a min and max length
func ValidateString(value string, min, max int) error {
	if len(value) < min || len(value) > max {
		return fmt.Errorf("must contain from %d-%d characters", min, max)
	}
	return nil
}

// ValidateUsername Function named ValidateUsername responsible for validating if a username is valid
func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	if !isValidUsername(username) {
		return fmt.Errorf("invalid username: must contain only lowercase letters, numbers, and underscores")
	}

	return nil
}

// ValidateFullName Function named ValidateFullName responsible for validating if a username is valid
func ValidateFullName(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	if !isValidFullName(username) {
		return fmt.Errorf("invalid full name: must contain only letters or spaces")
	}

	return nil
}

// ValidatePassword Function named ValidatePassword responsible for validating if a password is valid
func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

// ValidateEmail Function named ValidateEmail responsible for validating if an email is valid
func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 100); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	return nil
}

package utils

import (
    "regexp"
)

// CheckEmail checks if the email is valid
func CheckEmail(email string) bool {
    // Basic regex for email
    emailRegex := regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
    return emailRegex.MatchString(email)
}

// CheckPassword checks if password is strong enough
func CheckPassword(password string) bool {
    // At least 8 chars
    return len(password) >= 8
}



package utils

import (
    "regexp"
)

// ValidateEmail checks if the string is a valid email address.
func ValidateEmail(email string) bool {
      re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return re.MatchString(email)
}

// ValidateCode checks if the string is exactly 4 digits.
func ValidateCode(code string) bool {
    re := regexp.MustCompile(`^\d{4}$`)
    return re.MatchString(code)
}
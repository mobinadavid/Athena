package policies

import (
	"fmt"
	"regexp"
)

// PasswordPolicy represents the rules for a password policy.
type PasswordPolicy struct {
	MinLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireDigit     bool
	RequireSpecial   bool
}

// ValidatePassword checks if the provided password meets the password policy rules.
func (policy *PasswordPolicy) ValidatePassword(password string) error {
	// Check minimum length
	if len(password) < policy.MinLength {
		return fmt.Errorf("password must be at least %d characters long", policy.MinLength)
	}

	// Check uppercase requirement
	if policy.RequireUppercase && !containsUppercase(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check lowercase requirement
	if policy.RequireLowercase && !containsLowercase(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check digit requirement
	if policy.RequireDigit && !containsDigit(password) {
		return fmt.Errorf("password must contain at least one digit")
	}

	// Check special character requirement
	if policy.RequireSpecial && !containsSpecial(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func containsUppercase(s string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(s)
}

func containsLowercase(s string) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(s)
}

func containsDigit(s string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(s)
}

func containsSpecial(s string) bool {
	return regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(s)
}

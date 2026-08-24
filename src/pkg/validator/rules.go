package validator

import (
	"athena/src/pkg/i18n"
	"athena/src/pkg/policies"
	"fmt"
	"regexp"
	"strconv"
	"unicode/utf8"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func RegisterRules(val *validator.Validate, trans *ut.UniversalTranslator) {
	// Map of rule names to their corresponding validation functions
	ruleToFunc := map[string]validator.Func{
		"is-uuid":                        isValidUuid,
		"is-wallet-address":              isValidWalletAddress,
		"iranian-national-identity-code": iranianNationalCodeValidation,
		"iranian-mobile":                 iranianMobileValidation,
		"max-runes":                      validateMaxRunes,
		"is-strong-password":             isStrongPassword,
	}

	for ruleName, ruleFunc := range ruleToFunc {

		// Register the validation.
		_ = val.RegisterValidation(ruleName, ruleFunc)

		// Register validation messages as well.
		for _, lang := range i18n.Locales {
			translator, _ := trans.GetTranslator(lang)
			_ = val.RegisterTranslation(ruleName, translator, func(ut ut.Translator) error {
				return ut.Add(ruleName, i18n.Localize(lang, fmt.Sprintf("invalid-%s", ruleName)), false)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(fe.Tag(), fe.Field())
				return t
			})
		}
	}
}

// isValidUuid Custom validator function to validate UUID format
func isValidUuid(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}

// validate wallet-address
func isValidWalletAddress(fl validator.FieldLevel) bool {
	walletAddress := fl.Field().String()
	blockchainField := fl.Parent().FieldByName("Blockchain")
	if !blockchainField.IsValid() {
		return false
	}

	blockchain, ok := blockchainField.Interface().(string)
	if !ok {
		return false
	}

	var regex string
	switch blockchain {
	case "Bitcoin":
		regex = `^([13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[a-z0-9]{39,59})$`
	case "Ethereum", "Binance":
		regex = `^0x[a-fA-F0-9]{40}$`
	case "Tron":
		regex = `^T[a-zA-Z0-9]{33}$`
	default:
		return false
	}

	matched, _ := regexp.MatchString(regex, walletAddress)
	return matched
}

// iranianNationalCodeValidation validates the iranian national code.
func iranianNationalCodeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	matched, _ := regexp.MatchString(`^\d{8,10}$`, value)
	if !matched {
		return false
	}

	sequentialMatch, _ := regexp.MatchString(`^[0]{10}|[1]{10}|[2]{10}|[3]{10}|[4]{10}|[5]{10}|[6]{10}|[7]{10}|[8]{10}|[9]{10}$`, value)
	if sequentialMatch {
		return false
	}

	value = fmt.Sprintf("%010s", value)
	var sub int
	for i, char := range value[:9] {
		digit, _ := strconv.Atoi(string(char))
		sub += digit * (10 - i)
	}

	var control int
	if sub%11 < 2 {
		control = sub % 11
	} else {
		control = 11 - (sub % 11)
	}

	lastDigit, _ := strconv.Atoi(string(value[9]))
	return lastDigit == control
}

// iranianMobileValidation validates the iranian mobile number.
// valid example: 09123456789
func iranianMobileValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Must be exactly 11 chars and start with 09
	matched, err := regexp.MatchString(`^09[0-9]{9}$`, phone)
	if err != nil {
		return false
	}

	return matched
}

// validateMaxRunes checks the character (rune) count of a string.
func validateMaxRunes(fl validator.FieldLevel) bool {
	// Get the max value from the tag (e.g., "maxrunes=3000")
	param := fl.Param()
	maximum, err := strconv.Atoi(param)
	if err != nil {
		// Or panic, log, etc. if the tag is invalid
		return false
	}

	// Get the string value from the field
	str := fl.Field().String()

	// Count the runes (characters)
	count := utf8.RuneCountInString(str)

	// Check if the count is within the max
	return count <= maximum
}

func isStrongPassword(fl validator.FieldLevel) bool {
	passwordPolicy := &policies.PasswordPolicy{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireDigit:     true,
		RequireSpecial:   true,
	}

	policyCheck := passwordPolicy.ValidatePassword(fl.Field().String())

	return policyCheck == nil
}

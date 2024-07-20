package validator

import (
	"athena/src/pkg/i18n"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"regexp"
)

func RegisterRules(val *validator.Validate, trans *ut.UniversalTranslator) {
	// Map of rule names to their corresponding validation functions
	ruleToFunc := map[string]validator.Func{
		"is-uuid": isValidUuid,
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
func IsValidWalletAddress(blockchain string, walletAddress string) bool {

	var regex string

	switch blockchain {
	case "BTC":
		regex = `^([13][a-km-zA-HJ-NP-Z1-9]{25,34}|bc1[a-z0-9]{39,59})$`
	case "ETH", "BSC":
		regex = `^0x[a-fA-F0-9]{40}$`
	case "TRX":
		regex = `^T[a-zA-Z0-9]{33}$`

	default:
		return false
	}

	matched, _ := regexp.MatchString(regex, walletAddress)

	return matched
}

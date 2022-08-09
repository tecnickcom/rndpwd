// Package validator extends the parent validator package with custom checks.
package validator

import (
	"context"
	"regexp"

	vt "github.com/go-playground/validator/v10"
	val "github.com/nexmoinc/gosrvlib/pkg/validator"
)

const (
	// ValidCharset is a string containing the valid characters for a password.
	ValidCharset = "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
)

var (
	regexValidCharset = regexp.MustCompile("[^" + regexp.QuoteMeta(ValidCharset) + "]")
)

// Validator is the contract with the gosrvlib validator.
type Validator interface {
	ValidateStruct(interface{}) error
	ValidateStructCtx(context.Context, interface{}) error
}

// New instantiate a new Validator.
func New(fieldTagName string) (Validator, error) {
	customValidationTags := map[string]vt.FuncCtx{
		"rndcharset": validateRandomCharset(),
	}

	errorTemplates := map[string]string{
		"rndcharset": `{{.Namespace}} must contain only characters:` + ValidCharset,
	}

	//nolint:wrapcheck
	return val.New(
		val.WithFieldNameTag(fieldTagName),
		val.WithCustomValidationTags(val.CustomValidationTags()),
		val.WithCustomValidationTags(customValidationTags),
		val.WithErrorTemplates(val.ErrorTemplates()),
		val.WithErrorTemplates(errorTemplates),
	)
}

func validateRandomCharset() vt.FuncCtx {
	return func(ctx context.Context, fl vt.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			// empty fields are already checked by 'required'
			return true
		}

		return !regexValidCharset.MatchString(value)
	}
}

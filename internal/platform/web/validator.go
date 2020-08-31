package web

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidator() *Validator {
	validate := validator.New()

	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom rules
	validate.RegisterValidation("password", validatePassword)

	// Add custom translations
	translations := []struct {
		tag         string
		translation string
	}{
		{
			tag:         "required",
			translation: fmt.Sprintf("{0} is a required field"),
		},
		{
			tag:         "email",
			translation: fmt.Sprintf("{0} must be a valid email address"),
		},
		{
			tag:         "password",
			translation: fmt.Sprintf("{0} must greater than 5 characters and contain a capital letter, lower case letter, number, and special character"),
		},
	}

	for _, t := range translations {
		_ = validate.RegisterTranslation(t.tag, trans, register(t.tag, t.translation), translate)
	}

	return &Validator{
		validator:  validate,
		translator: trans,
	}
}

func register(tag string, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, true); err != nil {
			return
		}
		return
	}
}

func translate(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
	if err != nil {
		return fe.(error).Error()
	}
	return t
}

func validatePassword(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	var (
		length  = false
		upper   = false
		lower   = false
		number  = false
		special = false
	)

	if len(s) >= 6 {
		length = true
	}

	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			upper = true
		case unicode.IsLower(char):
			lower = true
		case unicode.IsNumber(char):
			number = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		}
	}

	return length && upper && lower && number && special
}

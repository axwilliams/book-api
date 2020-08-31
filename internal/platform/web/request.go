package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	if err := decoder.Decode(val); err != nil {
		return NewRequestError(fmt.Errorf("Unable to decode JSON: %w", err), http.StatusBadRequest)
	}

	v := NewValidator()

	if err := v.validator.Struct(val); err != nil {
		var fields []string

		vErrs := err.(validator.ValidationErrors)
		for _, vErr := range vErrs {
			field := fmt.Sprintf("%s", vErr.Translate(v.translator))
			fields = append(fields, field)
		}

		return &RequestError{
			Err:    ErrValidation,
			Status: http.StatusUnprocessableEntity,
			Fields: fields,
		}
	}

	return nil
}

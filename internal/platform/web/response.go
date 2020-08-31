package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func Message(k string, v interface{}) interface{} {
	return map[string]interface{}{k: v}
}

func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if statusCode == http.StatusNoContent || (statusCode == http.StatusOK && string(json) == "null") {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}

func RespondError(w http.ResponseWriter, err error) {
	re, ok := err.(*RequestError)
	if !ok {
		Respond(w, Message("message", ErrInternalServer.Error()), http.StatusInternalServerError)
		log.Printf("[error] %+v", err)
		return
	}

	if re.Err == ErrValidation {
		Respond(w, Validation{Message: "Validation failed", Errors: re.Fields}, re.Status)
		return
	}

	Respond(w, Message("message", re.Error()), re.Status)
}

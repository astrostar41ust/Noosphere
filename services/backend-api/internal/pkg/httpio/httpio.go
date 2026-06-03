package httpio

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, target interface{}, v *validator.Validate) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload format")
		return err
	}

	if err := v.Struct(target); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return err
	}

	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal serialization failure"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
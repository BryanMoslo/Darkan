package keywords

import (
	validation "darkan/internal/validation"

	"darkan/internal/response"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	keyword := Keyword{}
	err := json.NewDecoder(r.Body).Decode(&keyword)
	if err != nil {
		json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// Validate incoming data.
	validatior := validation.Validator{}
	validatior.Add(
		keyword.ValidateValue(),
		keyword.ValidateCallback(),
	)
	errors := validatior.Validate()
	if len(errors) > 0 {
		json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusBadRequest, "invalid data").WithData(map[string][]string{
			"errors": errors,
		}))
		return
	}

	keyword.Found = false
	keywordService := r.Context().Value("keywordService").(*service)

	err = keywordService.Create(&keyword)

	if isDuplicateKeyError(err) {
		slog.Error("keyword already exists for this keyword_id and source_url")
	}

	if err != nil && !isDuplicateKeyError(err) {
		slog.Error(fmt.Sprintf("Internal server error saving keyword: %s", err.Error()))
		json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusInternalServerError, "internal server error saving keyword"))
		return
	}

	response := response.SuccessResponse(http.StatusCreated, "keyword registered successfully").WithData(map[string]string{
		"keyword":      keyword.Value,
		"callback_url": keyword.CallbackURL,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

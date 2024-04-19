package keywords

import (
	"darkan/internal/response"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	keyword := Instance{}
	err := json.NewDecoder(r.Body).Decode(&keyword)
	if err != nil {
		json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	keywordService := r.Context().Value("keywordService").(*service)

	err = keywordService.Reload(&keyword)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		slog.Error(fmt.Sprintf("Internal server error reloading keyword: %s", err.Error()))
		json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusInternalServerError, "internal server error reloading keyword"))
		return
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		if err := keywordService.Create(&keyword); err != nil {
			slog.Error(fmt.Sprintf("Internal server error creating keyword: %s", err.Error()))
			json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusInternalServerError, "internal server error creating keyword"))
			return
		}
	}

	keyword.Found = false

	// Q: Should we trigger a bg search here? (still thinking on it)
	go keyword.Search(keywordService)

	response := response.SuccessResponse(http.StatusCreated, "keyword registered successfully").WithData(map[string]string{
		"keyword":      keyword.Value,
		"callback_url": keyword.CallbackURL,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

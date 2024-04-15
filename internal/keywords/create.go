package keywords

import (
	"darkan/internal/response"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	keyword := Instance{}
	err := json.NewDecoder(r.Body).Decode(&keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	keyword.Found = false
	keywords := r.Context().Value("keywordService").(*service)

	err = keywords.Create(&keyword)
	if err != nil {
		slog.Info(fmt.Sprintf("Internal server error saving keyword: %s", err.Error()))
		json.NewEncoder(w).Encode(response.ErrorResponse(http.StatusInternalServerError, "Internal server error saving keyword."))
		return
	}

	// Q: Should we trigger a bg search here? (still thinking on it)
	go keyword.search()

	response := response.SuccessResponse(http.StatusCreated, "Keyword registered successfully").WithData(map[string]string{
		"keyword":      keyword.Value,
		"callback_url": keyword.CallbackURL,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

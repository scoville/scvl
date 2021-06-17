package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/scoville/scvl/src/domain"
	"github.com/scoville/scvl/src/engine"
)

func (api *API) shortenHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := context.Get(r, "user").(*domain.User)
	if !ok {
		sendErrorJSON(w, http.StatusUnauthorized, "user not found")
		return
	}

	var req *engine.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	req.UserID = int(user.ID)

	page, err := api.engine.Shorten(req)
	if err != nil {
		sendErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{
		"url":  page.URL,
		"slug": page.Slug,
	}, noHeader)
}

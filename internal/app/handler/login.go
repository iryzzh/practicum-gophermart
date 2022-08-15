package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		basicResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	u, err := h.store.User().FindByLogin(req.Login)

	if err != nil || !u.ComparePassword(req.Password) {
		basicResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// очистка незашифрованного пароля
	u.Sanitize()

	session, err := h.sessionsStore.Get(r, h.cookieName)
	if err != nil {
		failResponse(w, err)
		return
	}

	session.Values["user_id"] = u.ID
	if err := h.sessionsStore.Save(r, w, session); err != nil {
		failResponse(w, err)
		return
	}
}

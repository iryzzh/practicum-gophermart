package handler

import (
	"context"
	"net/http"
)

func (h *Handler) AuthenticateHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionsStore.Get(r, h.cookieName)
		if err != nil {
			failResponse(w, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			basicResponse(w, http.StatusUnauthorized, "not authenticated")
			return
		}

		u, err := h.store.User().FindByID(id.(int))
		if err != nil {
			basicResponse(w, http.StatusUnauthorized, "not authenticated")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

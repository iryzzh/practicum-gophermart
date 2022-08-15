package handler

import (
	"encoding/json"
	"errors"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"net/http"
)

func (h *Handler) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	withdrawals, err := h.store.Balance().Withdrawals(r.Context().Value(ctxKeyUser).(*model.User).ID)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			basicResponse(w, http.StatusNoContent, err.Error())
			return
		}

		failResponse(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		failResponse(w, err)
	}
}

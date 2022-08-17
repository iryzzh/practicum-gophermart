package handler

import (
	"encoding/json"
	"io"
	"iryzzh/practicum-gophermart/internal/app/model"
	"net/http"
)

func (h *Handler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	balance, err := h.store.Balance().Get(r.Context().Value(ctxKeyUser).(*model.User).ID)
	if err != nil {
		failResponse(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		failResponse(w, err)
	}
}

func (h *Handler) PostWithdrawHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		failResponse(w, err)
		return
	}

	var m *model.Withdraw
	if err = json.Unmarshal(data, &m); err != nil {
		failResponse(w, err)
		return
	}

	if err = m.Validate(); err != nil {
		basicResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	m.UserID = r.Context().Value(ctxKeyUser).(*model.User).ID

	if err = h.store.Balance().Withdraw(m); err != nil {
		basicResponse(w, http.StatusPaymentRequired, err.Error())
		return
	}
}

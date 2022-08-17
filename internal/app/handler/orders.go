package handler

import (
	"encoding/json"
	"errors"
	"io"
	"iryzzh/practicum-gophermart/internal/app/model"
	"iryzzh/practicum-gophermart/internal/app/store"
	"log"
	"net/http"
)

func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	orders, err := h.store.Order().GetByUserID(r.Context().Value(ctxKeyUser).(*model.User).ID)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			basicResponse(w, http.StatusNoContent, "no content")
			return
		}

		failResponse(w, err)
		return
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		failResponse(w, err)
	}
}

// PostOrdersHandler Загрузка номера заказа
// Возможные коды ответа:
// 200 — номер заказа уже был загружен этим пользователем;
// 202 — новый номер заказа принят в обработку;
// 400 — неверный формат запроса;
// 401 — пользователь не аутентифицирован;
// 409 — номер заказа уже был загружен другим пользователем; (зачем это?)
// 422 — неверный формат номера заказа;
// 500 — внутренняя ошибка сервера.
func (h *Handler) PostOrdersHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		failResponse(w, err)
		return
	}

	if len(b) == 0 {
		basicResponse(w, http.StatusUnprocessableEntity, "empty body")
		return
	}

	o := &model.Order{
		Number: string(b),
		UserID: r.Context().Value(ctxKeyUser).(*model.User).ID,
	}

	if err := o.Validate(); err != nil {
		basicResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.store.Order().Create(o); err != nil {
		if errors.Is(err, store.ErrOrderAlreadyExists) {
			basicResponse(w, http.StatusOK, "")
			return
		}
		if errors.Is(err, store.ErrOrderConflict) {
			basicResponse(w, http.StatusConflict, err.Error())
			return
		}

		log.Println("other error:", err)
		basicResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

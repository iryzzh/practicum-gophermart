package handler

import (
	"encoding/json"
	"iryzzh/practicum-gophermart/internal/app/model"
	"net/http"
)

// RegisterHandler регистрирует пользователя
// 200 — пользователь успешно зарегистрирован и аутентифицирован
// 400 — неверный формат запроса;
// 409 — логин уже занят;
// 500 — внутренняя ошибка сервера.
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		basicResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	u := &model.User{
		Login:    req.Login,
		Password: req.Password,
	}

	// проверяем логин в бд
	if user, _ := h.store.User().FindByLogin(u.Login); user != nil {
		basicResponse(w, http.StatusConflict, "login already exists")
		return
	}

	if err := h.store.User().Create(u); err != nil {
		failResponse(w, err)
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

	w.WriteHeader(http.StatusOK)
}

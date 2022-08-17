package handler

import (
	"compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"io"
	"iryzzh/practicum-gophermart/internal/app/store"
	"net/http"
	"strings"
	"time"
)

const (
	ctxKeyUser ctxKey = iota
)

type ctxKey int8

type Handler struct {
	*chi.Mux
	store         store.Store
	sessionsStore *sessions.CookieStore
	cookieName    string
}

func New(store store.Store, sessionKey []byte) *Handler {
	s := &Handler{
		Mux:           chi.NewMux(),
		store:         store,
		sessionsStore: sessions.NewCookieStore(sessionKey),
		cookieName:    "_session_",
	}

	s.Use(middleware.RequestID)
	s.Use(middleware.RealIP)
	s.Use(middleware.Logger)
	s.Use(middleware.Recoverer)
	s.Use(middleware.Compress(5))

	s.Use(middleware.Timeout(5 * time.Second))

	s.Use(gzipMiddleware)

	/*
			POST /api/user/register — регистрация пользователя;
			POST /api/user/login — аутентификация пользователя;
			GET /api/orders/{number} — получение информации о расчёте начислений баллов лояльности.
		AUTH ONLY:
			POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
			GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
			GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
			POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
			GET /api/user/balance/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
			GET /api/user/withdrawals - Получение информации о выводе средств
	*/

	s.Post("/api/user/register", s.RegisterHandler)
	s.Post("/api/user/login", s.LoginHandler)
	s.Route("/api/user", func(r chi.Router) {
		r.Use(s.AuthenticateHandler)
		r.Get("/orders", s.GetOrdersHandler)
		r.Post("/orders", s.PostOrdersHandler)
		r.Get("/withdrawals", s.GetWithdrawalsHandler)

		// balance
		r.Route("/balance", func(r chi.Router) {
			r.Get("/", s.GetBalanceHandler)
			r.Post("/withdraw", s.PostWithdrawHandler)
		})
	})

	return s
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzr, r)
	})
}

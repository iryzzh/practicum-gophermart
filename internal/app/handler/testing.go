package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"iryzzh/practicum-gophermart/cmd/gophermart/config"
	"iryzzh/practicum-gophermart/internal/app/store"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func NewTestServer(t *testing.T, st store.Store) (*httptest.Server, error) {
	t.Helper()

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", cfg.RunAddress)
	if err != nil {
		return nil, err
	}

	handler := New(st, []byte(cfg.SessionKey))

	ts := httptest.NewUnstartedServer(handler)
	ts.Listener.Close()
	ts.Listener = l

	ts.Start()

	return ts, nil
}

func TestLogin(t *testing.T, server *httptest.Server, user, password string) (*cookiejar.Jar, error) {
	t.Helper()

	endpoint := "/api/user/login"
	payload := map[string]interface{}{"login": user, "password": password}

	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(payload); err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	resp, _ := TestRequest(t, "POST", server.URL+endpoint, b, jar)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return jar, nil
	}

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

func TestRequest(t *testing.T, method, path string, body io.Reader, jar *cookiejar.Jar) (*http.Response, string) {
	t.Helper()

	req, err := http.NewRequest(method, path, body)
	require.NoError(t, err)

	if jar == nil {
		jar, err = cookiejar.New(nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

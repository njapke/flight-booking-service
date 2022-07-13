package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecoverMiddleware(t *testing.T) {
	s := initService(t)
	mw := s.recoverMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	}))
	res := sendRequest(mw, "GET", "/", nil)
	require.Equal(t, http.StatusInternalServerError, res.Code)
}

func TestRecoverMiddlewareErrorType(t *testing.T) {
	s := initService(t)
	mw := s.recoverMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(fmt.Errorf("test"))
	}))
	res := sendRequest(mw, "GET", "/", nil)
	require.Equal(t, http.StatusInternalServerError, res.Code)
	var m map[string]string
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &m))
	require.Equal(t, "test", m["error"])
}

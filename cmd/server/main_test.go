package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rsync-backup/internal/types"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, strings.Contains(w.Body.String(), "Hello, world"))
}

func TestGetIP(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ip", nil)
	req.RemoteAddr = "192.168.191.133:12345"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, strings.Contains(w.Body.String(), "ip"))
}

func TestMakeServerConfig(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/serverconfig", nil)
	req.RemoteAddr = "192.168.191.133:12345"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostConfig(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()
	requestData := types.JsonConfigRequest{
		Data: types.JsonConfigData{
			Content: "test content",
		},
	}

	body, _ := json.Marshal(requestData)

	req, _ := http.NewRequest(http.MethodPost, "/config", bytes.NewBuffer(body))
	req.RemoteAddr = "192.168.191.133:12345"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetConfig(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodGet, "/config", nil)
	req.RemoteAddr = "192.168.191.133:12345"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

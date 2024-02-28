package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"rua.plus/gymo/controllers"
)

func TestRoot(t *testing.T) {
	router := InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	status := &controllers.RootStatus{
		Status: "ok",
	}
	resp, _ := json.Marshal(status)
	assert.Equal(t, resp, w.Body.Bytes())
}

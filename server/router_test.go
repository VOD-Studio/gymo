package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"rua.plus/gymo/controllers"
	"rua.plus/gymo/db"
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

func TestGetUser(t *testing.T) {
	mock := db.NewMockDB()
	router := InitRouter()

	var w *httptest.ResponseRecorder
	var req *http.Request
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/user/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	var rows *sqlmock.Rows
	rows = sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "xfy", "i@rua.plus")
	mock.ExpectQuery("^*$").WillReturnRows(rows)

	// user found
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/user/?email=i@rua.plus", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	rows = sqlmock.NewRows([]string{"id", "username", "email"})
	mock.ExpectQuery("^*$").WillReturnRows(rows)

	// user not found
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/user/?email=i@rua.plus", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 204, w.Code)
}

func TestRegister(t *testing.T) {
	mock := db.NewMockDB()
	router := InitRouter()

	var w *httptest.ResponseRecorder
	var req *http.Request

	// invalid request
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/v1/register/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)

	var rows *sqlmock.Rows
	rows = sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(1, "xfy", "i@rua.plus")
	mock.ExpectQuery("^*$").WillReturnRows(rows)

	var body []byte
	// user already exist
	body = []byte(`{"username": "xfy", "email": "i@rua.plus", "password": "passwd"}`)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(
		"POST",
		"/v1/register/",
		bytes.NewBuffer(body),
	)
	router.ServeHTTP(w, req)
	assert.Equal(t, 409, w.Code)
}

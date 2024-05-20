package main

import (
	"bytes"
	"flag"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	config "github.com/nyudlts/go-medialog/config"
)

var r *gin.Engine

func TestAPI(t *testing.T) {

	flag.Parse()

	//set the environment variables
	var err error
	env, err = config.GetEnvironment(configuration, environment)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	t.Run("Test get router", func(t *testing.T) {
		r, err = setupRouter()
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v", r)
	})

	t.Run("Test GET login route", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, err := http.NewRequestWithContext(c, "GET", "/users/login", nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("content-type"))
	})

	var sessionCookie string
	t.Run("Test POST login route", func(t *testing.T) {
		w := httptest.NewRecorder()

		var b bytes.Buffer
		w2 := multipart.NewWriter(&b)
		w2.WriteField("email", "admin@nyu.edu")
		w2.WriteField("password_1", "test")
		w2.Close()
		reader := bytes.NewReader(b.Bytes())
		req, err := http.NewRequestWithContext(c, "POST", "/users/authenticate", reader)
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", w2.FormDataContentType())
		r.ServeHTTP(w, req)
		header := w.Header().Get("Set-Cookie")
		sessionCookie = strings.Split(strings.Split(header, ";")[0], "=")[1]
		log.Println(sessionCookie)
		assert.Equal(t, http.StatusFound, w.Code)

	})
}

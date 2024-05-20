package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
	config "github.com/nyudlts/go-medialog/config"
)

func TestAPI(t *testing.T) {

	flag.Parse()

	//set the environment variables
	var err error
	env, err = config.GetEnvironment(configuration, environment)
	if err != nil {
		panic(err)
	}

	t.Run("Test get router", func(t *testing.T) {
		r, err := setupRouter()
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v", r)
	})

	t.Run("TestLoginRoute", func(t *testing.T) {
		router, err := setupRouter()
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/login", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("content-type"))
	})
}

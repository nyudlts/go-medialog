package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	router "github.com/nyudlts/go-medialog/router"
	"github.com/stretchr/testify/assert"
)

const APIROOT = "/api/v0"

func TestAPI(t *testing.T) {
	flag.Parse()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	env, err := router.GetEnvironment(configuration, environment)
	if err != nil {
		t.Error(err)
	}

	r, err = router.SetupRouter(env, true, false)
	if err != nil {
		t.Error(err)
	}

	t.Run("Test get API Root", func(t *testing.T) {
		req, err := http.NewRequestWithContext(c, "GET", APIROOT, nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
	})

}

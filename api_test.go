package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nyudlts/go-medialog/models"
	router "github.com/nyudlts/go-medialog/router"
	"github.com/stretchr/testify/assert"
)

const APIROOT = "/api/v0"

var token string

func TestAPI(t *testing.T) {
	flag.Parse()
	gin.SetMode(gin.TestMode)

	var err error
	env, err = router.GetEnvironment(configuration, environment)
	if err != nil {
		t.Error(err)
	}

	r, err = router.SetupRouter(env, true, false)
	if err != nil {
		t.Error(err)
	}
	t.Logf("[INFO] Running Go-Medialog %s", version)

	t.Run("test get API root", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		req, err := http.NewRequestWithContext(c, "GET", APIROOT, nil)
		if err != nil {
			t.Error(err)
		}
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}
		mlInfo := models.MedialogInfo{}
		if err := json.Unmarshal(body, &mlInfo); err != nil {
			t.Error(token)
		}
		t.Logf("[INFO] %v", mlInfo)
	})

	t.Run("test login to API", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/users/%s/login?password=%s", APIROOT, env.TestCreds.Username, env.TestCreds.Password)
		req, err := http.NewRequestWithContext(c, "POST", url, nil)
		if err != nil {
			t.Error(err)
		}

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
		bodyBytes, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		sessionToken := models.Token{}
		if err := json.Unmarshal(bodyBytes, &sessionToken); err != nil {
			t.Error(err)
		}
		token = sessionToken.Token
	})

	t.Run("test unauthorized access no medialog session header", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/repositories", APIROOT)
		req, err := http.NewRequestWithContext(c, "GET", url, nil)
		if err != nil {
			t.Error(err)
		}

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 401, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

	})

	t.Run("test get repositories", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/repositories", APIROOT)
		req, err := http.NewRequestWithContext(c, "GET", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

	})

}

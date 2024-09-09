package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

	var repoID uint
	t.Run("test create a repository", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/repositories", APIROOT)
		form := url.Values{}
		form.Set("title", "Test Repository")
		form.Add("slug", "Test")
		req, err := http.NewRequestWithContext(c, "POST", requestURL, strings.NewReader(form.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		repo := models.Repository{}
		if err := json.Unmarshal(body, &repo); err != nil {
			t.Error(err)
		}

		repoID = repo.ID

	})

	t.Run("test get all repositories", func(t *testing.T) {
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

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		repositories := []models.Repository{}
		if err := json.Unmarshal(body, &repositories); err != nil {
			t.Error(err)
		}
		assert.GreaterOrEqual(t, len(repositories), 1)
	})

	var repository models.Repository
	t.Run("test get a repository", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/repositories/%d", APIROOT, repoID)
		req, err := http.NewRequestWithContext(c, "GET", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		if err := json.Unmarshal(body, &repository); err != nil {
			t.Error(err)
		}
	})

	//resource functions
	var resource = models.Resource{}
	t.Run("test create a resource", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/resources", APIROOT)
		form := url.Values{}
		form.Set("title", "Test Resource")
		form.Add("collection_code", "test001")
		form.Add("repository_id", fmt.Sprintf("%d", repository.ID))
		form.Add("partner_code", "test")
		req, err := http.NewRequestWithContext(c, "POST", requestURL, strings.NewReader(form.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		if err := json.Unmarshal(body, &resource); err != nil {
			t.Error(err)
		}

	})

	t.Run("test get all resources", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/resources", APIROOT)
		req, err := http.NewRequestWithContext(c, "GET", requestURL, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		resources := []models.Resource{}
		if err := json.Unmarshal(body, &resources); err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(t, len(resources), 1)

	})

	t.Run("test get a resource", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/resources/%d", APIROOT, resource.ID)
		req, err := http.NewRequestWithContext(c, "GET", requestURL, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		resource := models.Resource{}
		if err := json.Unmarshal(body, &resource); err != nil {
			t.Error(err)
		}

	})

	//accession functions
	var accession = models.Accession{}
	t.Run("test create an accession", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/accessions", APIROOT)
		form := url.Values{}
		form.Add("accession_num", "testAccession001")
		form.Add("resource_id", fmt.Sprintf("%d", resource.ID))
		req, err := http.NewRequestWithContext(c, "POST", requestURL, strings.NewReader(form.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		if err := json.Unmarshal(body, &accession); err != nil {
			t.Error(err)
		}
	})

	t.Run("test get all accessions", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/accessions", APIROOT)
		req, err := http.NewRequestWithContext(c, "GET", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		accessions := []models.Accession{}
		if err := json.Unmarshal(body, &accessions); err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(t, len(accessions), 1)
	})

	t.Run("test get an accession", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/accessions/%d", APIROOT, accession.ID)
		req, err := http.NewRequestWithContext(c, "GET", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		if err := json.Unmarshal(body, &accession); err != nil {
			t.Error(err)
		}
	})

	//Entry functions
	var entry = models.Entry{}
	t.Run("test create an entry", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/entries", APIROOT)
		form := url.Values{}
		form.Add("media_id", "1")
		form.Add("mediatype", "mediatype_floppy_3_5")
		form.Add("stock_unit", "MB")
		form.Add("stock_size_num", "3.5")
		form.Add("resource_id", fmt.Sprintf("%d", resource.ID))
		form.Add("repository_id", fmt.Sprintf("%d", repository.ID))
		form.Add("accession_id", fmt.Sprintf("%d", accession.ID))
		req, err := http.NewRequestWithContext(c, "POST", requestURL, strings.NewReader(form.Encode()))
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		if err := json.Unmarshal(body, &entry); err != nil {
			t.Error(err)
		}

		assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", entry.ID.String())
	})

	t.Run("test get all entry ids", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/entries?all_ids=true", APIROOT)
		req, err := http.NewRequestWithContext(c, "GET", requestURL, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}

		entries := []string{}
		if err := json.Unmarshal(body, &entries); err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(t, len(entries), 1)

	})

	t.Run("test get an entry", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		requestURL := fmt.Sprintf("%s/entries/%s", APIROOT, entry.ID.String())
		req, err := http.NewRequestWithContext(c, "GET", requestURL, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
		body, err := io.ReadAll(recorder.Body)
		if err != nil {
			t.Error(err)
		}
		e := models.Entry{}
		if err := json.Unmarshal(body, &e); err != nil {
			t.Error(err)
		}
		assert.NotEqual(t, "00000000-0000-0000-0000-000000000000", e.ID.String())

	})

	t.Run("test update location of an entry", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		location := "sl_rsw_acm_born_digital"
		url := fmt.Sprintf("%s/entries/%s/update_location?location=%s", APIROOT, entry.ID, location)
		req, err := http.NewRequestWithContext(c, "PATCH", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
	})

	t.Run("test update location failure", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		location := "test"
		url := fmt.Sprintf("%s/entries/%s/update_location?location=%s", APIROOT, entry.ID, location)
		req, err := http.NewRequestWithContext(c, "PATCH", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 400, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
	})

	//delete functions

	t.Run("test delete an entry", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/entries/%s", APIROOT, entry.ID)
		req, err := http.NewRequestWithContext(c, "DELETE", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

	})

	t.Run("test delete an accession", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/accessions/%d", APIROOT, accession.ID)
		req, err := http.NewRequestWithContext(c, "DELETE", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

	})

	t.Run("test delete a resource", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/resources/%d", APIROOT, resource.ID)
		req, err := http.NewRequestWithContext(c, "DELETE", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))

	})

	t.Run("test delete a repository", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/repositories/%d", APIROOT, repoID)
		req, err := http.NewRequestWithContext(c, "DELETE", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
	})

	t.Run("test delete all sessions", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/delete_sessions", APIROOT)
		req, err := http.NewRequestWithContext(c, "DELETE", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
	})

	t.Run("test delete a token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		url := fmt.Sprintf("%s/logout", APIROOT)
		req, err := http.NewRequestWithContext(c, "DELETE", url, nil)
		if err != nil {
			t.Error(err)
		}
		req.Header.Add("X-Medialog-Token", token)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("content-type"))
	})

}

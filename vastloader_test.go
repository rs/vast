package vast_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KargoGlobal/vast"
	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
)

func TestLoadVAST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/vast_inline_linear.xml")
	}))
	vastXml, err := vast.LoadURI(server.URL+"/vast", 1000)
	assert.Nil(t, err)
	assert.Equal(t, "601364", vastXml.Ads[0].ID)
}

func TestLoadVASTSlowServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 1000)
		http.ServeFile(w, r, "testdata/vast_inline_linear.xml")
	}))
	vastXml, err := vast.LoadURI(server.URL+"/vast", 1)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	assert.Nil(t, vastXml)
}

func TestLoadInvalidVAST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is not VAST")
	}))
	vastXml, err := vast.LoadURI(server.URL+"/vast", 5000)
	assert.Contains(t, err.Error(), "EOF")
	assert.Nil(t, vastXml)
}

func TestLoadInvalidURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	vastXml, err := vast.LoadURI(server.URL+"/vast", 5000)
	assert.Contains(t, err.Error(), "EOF")
	assert.Nil(t, vastXml)
}

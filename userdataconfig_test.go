package userdataconfig

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var TEST_DATA = `#!/bin/bash
cd ..
cat pippo >>/echo.txt
#VARS
#{"name":"a test name","specialkey":"wonderful"}`

var BAD_DATA = `#!/bin/bash
cd ..
cat pippo >>/echo.txt
#VARS
#{ame":"a test name","specialkey":"wonderful"}`

var NOVARS_DATA = `#!/bin/bash
cd ..
cat pippo >>/echo.txt
`

func TestGetVars(t *testing.T) {
	rq := require.New(t)
	//build a mock backend service that will reply
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(TEST_DATA))
	}).Methods("GET")
	backend := httptest.NewServer(r)
	defer backend.Close()
	vars, err := GetVars(&(backend.URL))
	rq.NoError(err)
	rq.Equal("a test name", vars["name"])
	rq.Equal("wonderful", vars["specialkey"])
}

func TestGetVarsBadJSON(t *testing.T) {
	rq := require.New(t)
	//build a mock backend service that will reply
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(BAD_DATA))
	}).Methods("GET")
	backend := httptest.NewServer(r)
	defer backend.Close()
	_, err := GetVars(&(backend.URL))
	rq.Error(err)
}

func TestGetVarsNoVars(t *testing.T) {
	rq := require.New(t)
	//build a mock backend service that will reply
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(NOVARS_DATA))
	}).Methods("GET")
	backend := httptest.NewServer(r)
	defer backend.Close()
	_, err := GetVars(&(backend.URL))
	rq.Error(err)
}

func TestGetVarsBadUrl(t *testing.T) {
	rq := require.New(t)
	u := "http://127.0.0.1:2450/blal"
	_, err := GetVars(&u)
	rq.Error(err)
}

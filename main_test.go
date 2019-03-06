// Copyright 2019 Ivan Sorokin.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResult(t *testing.T) {
	body := []byte("42")
	hash := md5.Sum(body)
	result := NewResult("abc", body)
	if result.hash != hash {
		t.Errorf("Expected err to be %x, but received %x", hash, result.hash)
	}
}

func TestFetchBodyOK(t *testing.T) {
	mockBody := "42"
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(mockBody))
	}))
	defer srv.Close()
	body, err := FetchBody(srv.Client(), srv.URL)
	if err != nil {
		t.Errorf("Expected err to be nil, but received %s", err)
	}
	if mockBody != string(body) {
		t.Errorf("Expected body to be %s, but received %s", mockBody, body)
	}
}

func TestFetchBodyFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()
	_, err := FetchBody(srv.Client(), srv.URL)
	if err == nil {
		t.Errorf("Expected err to be %s, but received %s", fmt.Errorf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)), err)
	}
}

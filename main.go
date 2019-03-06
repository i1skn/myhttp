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
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Result represents successful fetch and hash
type Result struct {
	url  string
	hash [md5.Size]byte
}

// NewResult creates Result from url and body
func NewResult(url string, body []byte) Result {
	return Result{url, md5.Sum(body)}
}

// Error represents failed fetch and hash
type Error struct {
	url string
	err error
}

// FetchBody fetches response body from url
func FetchBody(client *http.Client, url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	resp, e := client.Get(url)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	return body, nil
}

func worker(in chan string, out chan Result, err chan Error) {
	for url := range in {
		body, e := FetchBody(&http.Client{}, url)
		if e != nil {
			err <- Error{url, e}
			continue
		}
		out <- NewResult(url, body)
	}
}

func main() {
	parallelRequests := flag.Int("parallel", 10, "Parallel requests limit")
	flag.Parse()
	urls := flag.Args()
	workersCount := *parallelRequests
	// It does not make sense to have more workers than urls
	if workersCount > len(urls) {
		workersCount = len(urls)
	}

	in := make(chan string)
	out := make(chan Result)
	err := make(chan Error)

	for i := 0; i < workersCount; i++ {
		go worker(in, out, err)
	}

	go func() {
		for _, url := range urls {
			in <- url
		}
	}()

	for i := 0; i < len(urls); i++ {
		select {
		case r := <-out:
			fmt.Printf("%s %x\n", r.url, r.hash)
		case e := <-err:
			fmt.Printf("%s %s\n", e.url, e.err)
		}
	}
}

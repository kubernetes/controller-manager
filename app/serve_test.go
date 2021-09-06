/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/tools/leaderelection"
	componentbaseconfig "k8s.io/component-base/config"
)

const (
	pprofPath = "/debug/pprof"
	logPath   = "/debug/flags"
)

func TestNewBaseHandler(t *testing.T) {
	// Setup any healthz checks we will want to use.
	var checks []healthz.HealthChecker
	var electionChecker *leaderelection.HealthzAdaptor = leaderelection.NewLeaderHealthzAdaptor(time.Second * 20)
	checks = append(checks, electionChecker)

	//Test1: EnableProfiling=true
	c := &componentbaseconfig.DebuggingConfiguration{
		EnableProfiling:           true,
		EnableContentionProfiling: true,
	}
	mux := NewBaseHandler(c, checks...)

	assert.NotContains(t, mux.ListedPaths(), pprofPath)
	assert.NotContains(t, mux.ListedPaths(), logPath)

	s := httptest.NewServer(mux)
	defer s.Close()
	resp, _ := http.Get(s.URL + pprofPath)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, _ = http.Get(s.URL + logPath)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//Test2: EnableProfiling=false
	c = &componentbaseconfig.DebuggingConfiguration{
		EnableProfiling:           false,
		EnableContentionProfiling: false,
	}
	mux = NewBaseHandler(c, checks...)

	assert.NotContains(t, mux.ListedPaths(), pprofPath)
	assert.NotContains(t, mux.ListedPaths(), logPath)

	s = httptest.NewServer(mux)
	defer s.Close()
	resp, _ = http.Get(s.URL + pprofPath)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp, _ = http.Get(s.URL + logPath)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

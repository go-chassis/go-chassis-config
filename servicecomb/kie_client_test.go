/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicecomb

import (
	"encoding/json"
	"fmt"
	"github.com/apache/servicecomb-kie/pkg/model"
	config "github.com/go-chassis/go-chassis-config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func init() {
}

//TestKieClient_NewKieClient for NewClient.
func TestKieClient_NewKieClient(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	_, err := NewClient(config.Options{App: "", ServerURI: "http://127.0.0.1:49800", Version: "1", ServiceName: "test", Endpoint: "http://127.0.0.1:49800"})
	assert.Equal(t, err, nil)
}

//TestKieClient_PullConfig for test PullConfig.
func TestKieClient_PullConfig(t *testing.T) {
	//v1/calculator/kie/kv?q=env:+servicename:calculator+version:0.0.1+app:
	helper := startHttpServer(":49800", "/v1/test/kie/kv/test")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	kieClient, err := NewClient(config.Options{App: "", ServerURI: "http://127.0.0.1:49800", Version: "1", ServiceName: "test", Endpoint: "http://127.0.0.1:49800"})
	resp, err := kieClient.PullConfig("test", "1", "", "", "test", "")
	fmt.Println(resp, err)
	//assert.Equal(t, resp.StatusCode, 404)
	assert.Equal(t, err.Error(), "can not find value")
	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

//TestKieClient_PullConfigs for test PullConfigs.
func TestKieClient_PullConfigs(t *testing.T) {
	//v1/calculator/kie/kv?q=env:+servicename:calculator+version:0.0.1+app:
	helper := startHttpServer(":49800", "/v1/calculator/kie/kv?q=version:0.0.1+app:+env:+servicename:calculator")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	kieClient, err := NewClient(config.Options{App: "", ServerURI: "http://127.0.0.1:49800", Version: "1", ServiceName: "test", Endpoint: "http://127.0.0.1:49800"})
	resp, err := kieClient.PullConfigs("test", "1", "", "")
	fmt.Println(resp, err)
	//assert.Equal(t, resp.StatusCode, 404)
	assert.Equal(t, err.Error(), "can not find value")
	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

//TestKieClient_PushConfigs for test PushConfigs.
func TestKieClient_PushConfigs(t *testing.T) {
	//v1/calculator/kie/kv?q=env:+servicename:calculator+version:0.0.1+app:
	helper := startHttpServer(":49800", "/")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	kieClient, err := NewClient(config.Options{App: "", ServerURI: "http://127.0.0.1:49800", Version: "1", ServiceName: "test", Endpoint: "http://127.0.0.1:49800"})
	data := make(map[string]interface{})
	data["test_info"] = "test_info"
	resp, err := kieClient.PushConfigs(data, "test", "1", "", "")
	fmt.Println(resp, err)
	//assert.Equal(t, resp.StatusCode, 404)
	assert.Equal(t, err.Error(), "json: cannot unmarshal array into Go value of type model.KVDoc")
	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

//TestKieClient_DeleteConfigs for test DeleteConfigs.
func TestKieClient_DeleteConfigs(t *testing.T) {
	//v1/calculator/kie/kv?q=env:+servicename:calculator+version:0.0.1+app:
	helper := startHttpServer(":49800", "/v1/calculator/kie/kv/?kvID=s")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	kieClient, err := NewClient(config.Options{App: "", ServerURI: "http://127.0.0.1:49800", Version: "1", ServiceName: "test", Endpoint: "http://127.0.0.1:49800"})
	data := []string{"1"}
	resp, err := kieClient.DeleteConfigsByKeys(data, "test", "1", "", "")
	fmt.Println(resp, err)
	//assert.Equal(t, resp.StatusCode, 404)
	assert.Equal(t, err, nil)
	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

//startHttpServer
func startHttpServer(port string, pattern string) *http.Server {
	helper := &http.Server{Addr: port}
	var result model.KVResponse
	var req []*model.KVResponse
	req = append(req, &result)
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(req)
		w.Write(body)
	})
	go func() {
		if err := helper.ListenAndServe(); err != nil {
			fmt.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	time.Sleep(time.Second * 1)
	return helper
}

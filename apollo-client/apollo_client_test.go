package apolloclient

import (
	"encoding/json"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestApolloClient_NewApolloClient(t *testing.T) {

}

func TestApolloClient_HTTPDo(t *testing.T) {
	keepAlive := map[string]interface{}{
		"timeout": "500",
	}
	helper := startHttpServer(":9876", "/test", keepAlive)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}

	apolloClient := &ApolloClient{}
	apolloClient.NewApolloClient()

	// Test existing API 's
	resp, err := apolloClient.HTTPDo("GET", "http://127.0.0.1:9876/test", nil, nil)
	assert.NotEqual(t, resp, nil)
	assert.Equal(t, err, nil)

	// Test Non-existing API's
	resp, err = apolloClient.HTTPDo("GET", "http://127.0.0.1:9876/testUN", nil, nil)
	assert.Equal(t, resp.StatusCode, 404)
	assert.Equal(t, err, nil)

	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

func TestApolloClient_PullConfig(t *testing.T) {
	configurations := map[string]interface{}{
		"timeout": "500",
	}
	configBody := map[string]interface{}{
		"appId":          "TestApp",
		"cluster":        "default",
		"namespaceName":  "application",
		"configurations": configurations,
		"releaseKey":     "20180327130726-1dc5027439679153",
	}

	helper := startHttpServer(":9875", "/configs/TestApp/Default/application", configBody)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}

	apolloClient := &ApolloClient{}
	apolloClient.NewApolloClient()
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "http://127.0.0.1:9875"
	config.GlobalDefinition.Cse.Config.Client.ApolloServiceName = "TestApp"
	config.GlobalDefinition.Cse.Config.Client.ClusterName = "Default"
	config.GlobalDefinition.Cse.Config.Client.ApolloNameSpace = "application"

	//Test existing Services
	configResponse, error := apolloClient.PullConfig("TestApp", "1.0", "SampleApp", "Default", "timeout", "")
	assert.NotEqual(t, configResponse, nil)
	assert.Equal(t, error, nil)

	//Test the non-existing Key
	configResponse, error = apolloClient.PullConfig("TestApp", "1.0", "SampleApp", "Default", "non-exsiting", "")
	assert.Contains(t, error.Error(), "No Key found")

	// Test the non-exsisting Service
	config.GlobalDefinition.Cse.Config.Client.ApolloServiceName = "Non-exsitingAppID"
	configResponse, error = apolloClient.PullConfig("TestApp", "1.0", "SampleApp", "Default", "non-exsiting", "")
	assert.Contains(t, error.Error(), "Bad Response")

	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}

}

func TestApolloClient_PullConfigs(t *testing.T) {
	configurations := map[string]interface{}{
		"timeout": "500",
	}
	configBody := map[string]interface{}{
		"appId":          "SampleApp",
		"cluster":        "default",
		"namespaceName":  "application",
		"configurations": configurations,
		"releaseKey":     "20180327130726-1dc5027439679153",
	}

	helper := startHttpServer(":9874", "/configs/SampleApp/Default/application", configBody)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}

	apolloClient := &ApolloClient{}
	apolloClient.NewApolloClient()
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "http://127.0.0.1:9874"
	config.GlobalDefinition.Cse.Config.Client.ApolloServiceName = "SampleApp"
	config.GlobalDefinition.Cse.Config.Client.ClusterName = "Default"
	config.GlobalDefinition.Cse.Config.Client.ApolloNameSpace = "application"

	//Test existing Services
	configResponse, error := apolloClient.PullConfigs("SampleApp", "1.0", "SampleApp", "Default")
	assert.NotEqual(t, configResponse, nil)
	assert.Equal(t, error, nil)

	//Test the non-existing Services
	config.GlobalDefinition.Cse.Config.Client.ApolloServiceName = "Non-exsitingAppID"
	configResponse, error = apolloClient.PullConfigs("SampleApp", "1.0", "SampleApp", "Default")
	assert.Contains(t, error.Error(), "Bad Response")

	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

func startHttpServer(port string, pattern string, responseBody map[string]interface{}) *http.Server {
	helper := &http.Server{Addr: port}
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		body, _ := json.Marshal(responseBody)
		w.Write(body)
	})

	go func() {
		if err := helper.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	return helper
}

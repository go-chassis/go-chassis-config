package memberdiscovery

import (
	"math/rand"
	"os"
	"strconv"
	"testing"

	"encoding/json"
	"github.com/ServiceComb/go-archaius"
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/http-client"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
)

type TestingSource struct {
}

func (*TestingSource) GetConfigCenters() []string {
	configserver := []string{`10.18.206.218:30103`}

	return configserver
}

func TestShuffle(t *testing.T) {
	t.Log("Testing Shuffle function for errors")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.APIVersion.Version = "v2"
	ccClient := new(MemDiscovery)
	memDiscovery := ccClient.NewConfiCenterInit(nil, "default", false)

	er := memDiscovery.Shuffle()

	assert.Error(t, er)
}

/*func TestGetConfigServerIsInitErr(t *testing.T) {
	t.Log("Testing GetConfigServer function for errors")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()

	memDiscovery := NewConfiCenterInit(nil, "default", false)

	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t,err)

	_, er := memDiscovery.GetConfigServer()

	assert.Error(t, er)
}*/

func TestRefreshMembersConfigAddNil(t *testing.T) {
	t.Log("Testing RefreshMembers function")
	ccClient := new(MemDiscovery)
	memDiscovery := ccClient.NewConfiCenterInit(nil, "default", false)

	er := memDiscovery.RefreshMembers()
	assert.NoError(t, er)

}

/*func TestInit(t *testing.T) {
	t.Log("Testing ConfigurationInit function with errors")
	memDiscovery := NewConfiCenterInit(nil, "default", false)

	//testSource := &TestingSource{}
	//configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(nil)
	assert.Error(t, err)
}*/

/*func TestInitConfig(t *testing.T) {
	t.Log("Testing ConfigurationInit function without errors")
	memDiscovery := NewConfiCenterInit(nil, "default", false)

	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)
}*/
/*func TestGetConfigServerAutoDiscovery(t *testing.T) {
	t.Log("Testing GetConfigServer function Auto discovery")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.Autodiscovery = true
	memDiscovery := NewConfiCenterInit(nil, "default", false)
	//testSource := &TestingSource{}
	//configCenters := testSource.GetConfigCenters()
	//err := memDiscovery.ConfigurationInit(configCenters)
	//assert.NoError(t,err)

	_, er := memDiscovery.GetConfigServer()

	assert.NoError(t, er)
}*/

/*func TestGetConfigServer(t *testing.T) {
	t.Log("Testing GetConfigServer without errors after initializing configurations")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()

	memDiscovery := NewConfiCenterInit(nil, "default", false)
	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)

	_, er := memDiscovery.GetConfigServer()

	assert.NoError(t, er)
}*/

/*func TestRefreshMembers(t *testing.T) {
	t.Log("Testing RefreshMembers without error after initializing configuration")
	//func1 := func() http.Header {
	//	return nil
	//}
	//auth.GenAuthHeaders = func1
	memDiscovery := NewConfiCenterInit(nil, "default", false)
	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()

	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)

	er := memDiscovery.RefreshMembers()
	assert.NoError(t, er)

}*/
func TestGetDefaultHeadersArrayHeader(t *testing.T) {
	t.Log("Testing RefreshMembers without error after initializing configuration")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.MicroserviceDefinition = &model.MicroserviceCfg{}
	config.MicroserviceDefinition.ServiceDescription.Environment = "dev"
	/*func1 := func() http.Header {
		var sl []string
		sl = append(sl, "1")
		sl = append(sl, "2")
		h1 := http.Header{"abc": sl, "def": sl}
		return h1
	}

	auth.GenAuthHeaders = func1*/

	_ = GetDefaultHeaders("tenantName")
}

/*func TestGetWorkingConfigCenterIP(t *testing.T) {
	t.Log("Testing GetWorkingConfigCenterIP function")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()

	memDiscovery := NewConfiCenterInit(nil, "default", false)
	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)
	var endpoint = []string{"127.0.0.1", "5.6.7.8"}

	_, er := memDiscovery.GetWorkingConfigCenterIP(endpoint)

	assert.NoError(t, er)
}*/

func TestGetDefaultHeaders(t *testing.T) {
	t.Log("Headers should contain environment")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.MicroserviceDefinition = &model.MicroserviceCfg{}
	config.MicroserviceDefinition.ServiceDescription.Environment = ""

	if config.MicroserviceDefinition == nil {
		config.MicroserviceDefinition = &model.MicroserviceCfg{}
	}
	h := GetDefaultHeaders("")
	assert.Equal(t, "", h.Get(HeaderEnvironment))

	e := strconv.Itoa(rand.Int())
	config.MicroserviceDefinition.ServiceDescription.Environment = e
	h = GetDefaultHeaders("")
	assert.Equal(t, e, h.Get(HeaderEnvironment))
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

func TestMemDiscovery_HTTPDo(t *testing.T) {
	keepAlive := map[string]interface{}{
		"timeout": "500",
	}
	helper := startHttpServer(":9876", "/test", keepAlive)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}

	ccClient := new(MemDiscovery)
	//ccClient := NewConfiCenterInit(nil, "default", false)
	options := &httpclient.URLClientOption{
		SSLEnabled: false,
		TLSConfig:  nil,
		Compressed: false,
		Verbose:    false,
	}
	ccClient.client, _ = httpclient.GetURLClient(options)

	// Test existing API 's
	resp, err := ccClient.HTTPDo("GET", "http://127.0.0.1:9876/test", nil, nil)
	assert.NotEqual(t, resp, nil)
	assert.Equal(t, err, nil)

	// Test Non-existing API's
	resp, err = ccClient.HTTPDo("GET", "http://127.0.0.1:9876/testUN", nil, nil)
	assert.Equal(t, resp.StatusCode, 404)
	assert.Equal(t, err, nil)

	// Shutdown the helper server gracefully
	if err := helper.Shutdown(nil); err != nil {
		panic(err)
	}
}

func TestInitConfigCenter(t *testing.T) {
	t.Log("Testing InitConfigCenter function")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()
	registry.Enable()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = ""
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithTenantEmpty(t *testing.T) {
	t.Log("Testing InitConfigCenter function with autodiscovery true and tenant name empty")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.Autodiscovery = true
	config.GlobalDefinition.Cse.Config.Client.TenantName = ""
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithEmptyURI(t *testing.T) {
	t.Log("Testing InitConfigCenter function with empty ServerURI")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = ""
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithEmptyMicroservice(t *testing.T) {
	t.Log("Testing InitConfigCenter function with empty microservice definition")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()

	config.MicroserviceDefinition = &model.MicroserviceCfg{}
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithEnableSSl(t *testing.T) {
	t.Log("Testing InitConfigCenter function without initializing any parameter")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()

	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithInvalidURI(t *testing.T) {
	t.Log("Testing InitConfigCenter function with Invalid URI")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "hdhhhd:njdj"
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithSSL(t *testing.T) {
	t.Log("Testing InitConfigCenter function with ServerURI https://127.0.0.1:8787")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "https://127.0.0.1:8787"
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithInvalidName(t *testing.T) {
	t.Log("Testing InitConfigCenter function with serverURI and microservice definition")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	name := model.MicServiceStruct{Name: "qwertyuiopasdfghjklgsgdfsgdgafdggsahhhhh"}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "https://127.0.0.1:8787"
	config.MicroserviceDefinition = &model.MicroserviceCfg{ServiceDescription: name}
	ccClient := new(MemDiscovery)
	err = ccClient.InitConfigCenter()
	assert.Error(t, err)
	t.Log("HEllo", err)
}

func TestEvent(t *testing.T) {
	t.Log("Testing EventListener function")
	factoryObj, _ := goarchaius.NewConfigFactory(nil)

	factoryObj.Init()

	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	eventValue := &core.Event{Key: "refreshMode", Value: 6}
	evt := EventListener{Name: "EventHandler", Factory: factoryObj}
	evt.Event(eventValue)
}

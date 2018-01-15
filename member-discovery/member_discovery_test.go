package memberdiscovery

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/stretchr/testify/assert"

	//	"net/http"
	"os"
	"testing"
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
	os.Setenv("CHASSIS_HOME", gopath+"/src/code.huawei.com/cse/go-chassis-examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.APIVersion.Version = "v2"
	memDiscovery := NewConfiCenterInit(nil, "default", false)

	er := memDiscovery.Shuffle()

	assert.Error(t, er)
}

func TestGetConfigServerIsInitErr(t *testing.T) {
	t.Log("Testing GetConfigServer function for errors")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/code.huawei.com/cse/go-chassis-examples/discovery/server/")
	config.Init()

	memDiscovery := NewConfiCenterInit(nil, "default", false)

	/*testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t,err)
	*/
	_, er := memDiscovery.GetConfigServer()

	assert.Error(t, er)
}

func TestRefreshMembersConfigAddNil(t *testing.T) {
	t.Log("Testing RefreshMembers function")
	memDiscovery := NewConfiCenterInit(nil, "default", false)

	er := memDiscovery.RefreshMembers()
	assert.NoError(t, er)

}

func TestInit(t *testing.T) {
	t.Log("Testing ConfigurationInit function with errors")
	memDiscovery := NewConfiCenterInit(nil, "default", false)

	//testSource := &TestingSource{}
	//configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(nil)
	assert.Error(t, err)
}

func TestInitConfig(t *testing.T) {
	t.Log("Testing ConfigurationInit function without errors")
	memDiscovery := NewConfiCenterInit(nil, "default", false)

	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)
}
func TestGetConfigServerAutoDiscovery(t *testing.T) {
	t.Log("Testing GetConfigServer function Auto discovery")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/code.huawei.com/cse/go-chassis-examples/discovery/server/")
	config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.Autodiscovery = true
	memDiscovery := NewConfiCenterInit(nil, "default", false)
	/*testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t,err)
	*/
	_, er := memDiscovery.GetConfigServer()

	assert.NoError(t, er)
}

func TestGetConfigServer(t *testing.T) {
	t.Log("Testing GetConfigServer without errors after initializing configurations")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/code.huawei.com/cse/go-chassis-examples/discovery/server/")
	config.Init()

	memDiscovery := NewConfiCenterInit(nil, "default", false)
	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)

	_, er := memDiscovery.GetConfigServer()

	assert.NoError(t, er)
}

func TestRefreshMembers(t *testing.T) {
	t.Log("Testing RefreshMembers without error after initializing configuration")
	/*func1 := func() http.Header {
		return nil
	}
	auth.GenAuthHeaders = func1*/
	memDiscovery := NewConfiCenterInit(nil, "default", false)
	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()

	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)

	er := memDiscovery.RefreshMembers()
	assert.NoError(t, er)

}
func TestGetDefaultHeadersArrayHeader(t *testing.T) {
	t.Log("Testing RefreshMembers without error after initializing configuration")
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
func TestGetWorkingConfigCenterIP(t *testing.T) {
	t.Log("Testing GetWorkingConfigCenterIP function")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/code.huawei.com/cse/go-chassis-examples/discovery/server/")
	config.Init()

	memDiscovery := NewConfiCenterInit(nil, "default", false)
	testSource := &TestingSource{}
	configCenters := testSource.GetConfigCenters()
	err := memDiscovery.ConfigurationInit(configCenters)
	assert.NoError(t, err)
	var endpoint = []string{"1.2.3.4", "5.6.7.8"}

	_, er := memDiscovery.GetWorkingConfigCenterIP(endpoint)

	assert.NoError(t, er)
}

package client

import (
	"crypto/tls"
	//"github.com/ServiceComb/go-chassis/core/config"

	"github.com/ServiceComb/go-cc-client/apollo-client"
	"github.com/ServiceComb/go-cc-client/member-discovery"
)

var configClientPlugins = make(map[string]func(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient)

//DefaultClient is config server's client
var DefaultClient ConfigClient

const (
	defaultConfigServer = "config_center"
)

//InstallConfigClientPlugin install a config client plugin
func InstallConfigClientPlugin(name string, f func(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient) {
	configClientPlugins[name] = f
}

//ConfigClient is the interface of config server client, it has basic func to interact with config server
type ConfigClient interface {
	//Init the Configuration for the Server
	Init()
	//PullConfigs pull all configs from remote
	PullConfigs(serviceName, version, app, env string) (map[string]interface{}, error)
	//PullConfig pull one config from remote
	PullConfig(serviceName, version, app, env, key, contentType string) (interface{}, error)
	//PullConfigsByDI pulls the configurations with customized DimensionInfo/Project
	PullConfigsByDI(dimensionInfo, diInfo string) (map[string]map[string]interface{}, error)
}

//Enable enable config server client
func Enable(clientType string) {
	switch clientType {
	case "apollo":
		InstallConfigClientPlugin("apollo", InitConfigApollo)
	case "config_center":
		InstallConfigClientPlugin("config_center", InitConfigCenterNew)
	default:
		InstallConfigClientPlugin("config_center", InitConfigCenterNew)
	}

	plugins := configClientPlugins[clientType]
	if plugins == nil {
		panic("Default Plugin not found")
	}
	var tlsConfig *tls.Config
	DefaultClient = plugins("", "", "", "", "", tlsConfig)

	//Initiaizing the Client
	DefaultClient.Init()
}

//InitConfigApollo initialize the Apollo Client
func InitConfigApollo(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient {
	apolloClient := &apolloclient.ApolloClient{}
	apolloClient.NewApolloClient()
	return apolloClient
}

//InitConfigCenterNew initialize the Config-Center Client
func InitConfigCenterNew(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient {
	configSourceClient := &memberdiscovery.ConfigSourceClient{}
	configSourceClient.Init()
	return configSourceClient
}

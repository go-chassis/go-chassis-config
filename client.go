package client

import (
	"crypto/tls"
	"log"
)

var configClientPlugins = make(map[string]func(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient)

//DefaultClient is config server's client
var DefaultClient ConfigClient

//InstallConfigClientPlugin install a config client plugin
func InstallConfigClientPlugin(name string, f func(endpoint, serviceName, app, env, version string, tlsConfig *tls.Config) ConfigClient) {
	configClientPlugins[name] = f
	log.Printf("Installed %s Plugin", name)
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
	plugins := configClientPlugins[clientType]
	if plugins == nil {
		panic("Default Plugin not found")
	}
	var tlsConfig *tls.Config
	DefaultClient = plugins("", "", "", "", "", tlsConfig)

	//Initializing the Client
	DefaultClient.Init()
	log.Printf("%s Plugin is enabled", clientType)
}

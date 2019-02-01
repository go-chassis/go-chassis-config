package ccclient

import (
	"errors"
	"fmt"

	"github.com/go-mesh/openlogging"
)

var configClientPlugins = make(map[string]func(options Options) ConfigClient)

//DefaultClient is config server's client
var DefaultClient ConfigClient

//InstallConfigClientPlugin install a config client plugin
func InstallConfigClientPlugin(name string, f func(options Options) ConfigClient) {
	configClientPlugins[name] = f
	openlogging.GetLogger().Infof("Installed %s Plugin", name)
}

//ConfigClient is the interface of config server client, it has basic func to interact with config server
type ConfigClient interface {
	//Init the Configuration for the Server
	//PullConfigs pull all configs from remote
	PullConfigs(serviceName, version, app, env string) (map[string]interface{}, error)
	//PullConfig pull one config from remote
	PullConfig(serviceName, version, app, env, key, contentType string) (interface{}, error)
	//PullConfigsByDI pulls the configurations with customized DimensionInfo/Project
	PullConfigsByDI(dimensionInfo, diInfo string) (map[string]map[string]interface{}, error)
	// PushConfigs push config to cc
	PushConfigs(data map[string]interface{}, dimensionInfo string) (map[string]interface{}, error)
	// DeleteConfigsByKeys delete config for cc by keys
	DeleteConfigsByKeys(keys []string, dimensionInfo string) (map[string]interface{}, error)
}

//Enable enable config server client
func NewClient(name string, options Options) (ConfigClient, error) {
	plugins := configClientPlugins[name]
	if plugins == nil {
		return nil, errors.New(fmt.Sprintf("plugin [%s] not found", name))
	}
	DefaultClient = plugins(options)

	openlogging.GetLogger().Infof("%s plugin is enabled", name)
	return DefaultClient, nil
}

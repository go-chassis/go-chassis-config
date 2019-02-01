/*
 * Copyright 2017 Huawei Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package memberdiscovery created on 2017/6/20.
package configcenter

import (
	"github.com/go-chassis/go-cc-client"
	"github.com/go-chassis/go-cc-client/serializers"
	"github.com/go-chassis/go-chassis/pkg/httpclient"
	"github.com/go-mesh/openlogging"

	"errors"
	"net/http"
	"os"
	"strings"
)

var (
	//HeaderTenantName is a variable of type string
	HeaderTenantName = "X-Tenant-Name"
	//ConfigMembersPath is a variable of type string
	ConfigMembersPath = ""
	//ConfigPath is a variable of type string
	ConfigPath = ""
	//ConfigRefreshPath is a variable of type string
	ConfigRefreshPath = ""
	//MemberDiscoveryService is a variable
	MemberDiscoveryService MemberDiscovery
	autoDiscoverable       = false
	apiVersionConfig       = ""
	environmentConfig      = ""
)

const (
	//StatusUP is a variable of type string
	StatusUP = "UP"
	//HeaderContentType is a variable of type string
	HeaderContentType = "Content-Type"
	//HeaderUserAgent is a variable of type string
	HeaderUserAgent = "User-Agent"
	//HeaderEnvironment specifies the environment of a service
	HeaderEnvironment        = "X-Environment"
	members                  = "/configuration/members"
	dimensionsInfo           = "dimensionsInfo"
	dynamicConfigAPI         = `/configuration/refresh/items`
	getConfigAPI             = `/configuration/items`
	defaultContentType       = "application/json"
	envProjectID             = "CSE_PROJECT_ID"
	packageInitError         = "package not initialize successfully"
	emptyConfigServerMembers = "empty config server member"
	emptyConfigServerConfig  = "empty config server passed"
	// Name of the Plugin
	Name = "config_center"
)

//Client is Client Implementation of ConfigClient
type Client struct {
	memDiscovery *MemDiscovery
}

//Instance is a struct
type Instance struct {
	Status      string   `json:"status"`
	ServiceName string   `json:"serviceName"`
	IsHTTPS     bool     `json:"isHttps"`
	EntryPoints []string `json:"endpoints"`
}

//Members is a struct
type Members struct {
	Instances []Instance `json:"instances"`
}

//NewConfigCenter is a function
func NewConfigCenter(options ccclient.Options) ccclient.ConfigClient {
	memDiscovery := new(MemDiscovery)
	//memDiscovery.Logger = logger
	memDiscovery.TLSConfig = options.TLSConfig
	memDiscovery.TenantName = options.TenantName
	memDiscovery.EnableSSL = options.EnableSSL
	var apiVersion string
	apiVersionConfig = options.APIVersion
	autoDiscoverable = options.AutoDiscovery
	environmentConfig = options.Env

	switch apiVersionConfig {
	case "v2":
		apiVersion = "v2"
	case "V2":
		apiVersion = "v2"
	case "v3":
		apiVersion = "v3"
	case "V3":
		apiVersion = "v3"
	default:
		apiVersion = "v3"
	}
	//Update the API Base Path based on the Version
	updateAPIPath(apiVersion)

	//Initiate RestClient from http-client package
	opts := &httpclient.URLClientOption{
		SSLEnabled: options.EnableSSL,
		TLSConfig:  options.TLSConfig,
		Compressed: false,
		Verbose:    false,
	}
	memDiscovery.client, _ = httpclient.GetURLClient(opts)
	ccclient := &Client{
		memDiscovery: memDiscovery,
	}

	configCenters := strings.Split(options.ServerURI, ",")
	cCenters := make([]string, 0)
	for _, value := range configCenters {
		value = strings.Replace(value, " ", "", -1)
		cCenters = append(cCenters, value)
	}
	memDiscovery.ConfigurationInit(cCenters)
	return ccclient
}

//Update the Base PATH and HEADERS Based on the version of ConfigCenter used.
func updateAPIPath(apiVersion string) {

	//Check for the env Name in Container to get Domain Name
	//Default value is  "default"
	projectID, isExsist := os.LookupEnv(envProjectID)
	if !isExsist {
		projectID = "default"
	}
	switch apiVersion {
	case "v3":
		ConfigMembersPath = "/v3/" + projectID + members
		ConfigPath = "/v3/" + projectID + getConfigAPI
		ConfigRefreshPath = "/v3/" + projectID + dynamicConfigAPI
	case "v2":
		ConfigMembersPath = "/members"
		ConfigPath = "/configuration/v2/items"
		ConfigRefreshPath = "/configuration/v2/refresh/items"
	default:
		ConfigMembersPath = "/v3/" + projectID + members
		ConfigPath = "/v3/" + projectID + getConfigAPI
		ConfigRefreshPath = "/v3/" + projectID + dynamicConfigAPI
	}
}

//GetDefaultHeaders gets default headers
func GetDefaultHeaders(tenantName string) http.Header {
	headers := http.Header{
		HeaderContentType: []string{"application/json"},
		HeaderUserAgent:   []string{"cse-configcenter-client/1.0.0"},
		HeaderTenantName:  []string{tenantName},
	}
	if environmentConfig != "" {
		headers.Set(HeaderEnvironment, environmentConfig)
	}

	return headers
}

// PullConfigs is the implementation of ConfigClient to pull all the configurations from Config-Server
func (cclient *Client) PullConfigs(serviceName, version, app, env string) (map[string]interface{}, error) {
	// serviceName is the dimensionInfo passed from ConfigClient (small hack)
	configurations, error := cclient.memDiscovery.pullConfigurationsFromServer(serviceName)
	if error != nil {
		return nil, error
	}
	return configurations, nil
}

// PullConfig is the implementation of ConfigClient to pull specific configurations from Config-Server
func (cclient *Client) PullConfig(serviceName, version, app, env, key, contentType string) (interface{}, error) {

	// serviceName is the dimensionInfo passed from ConfigClient (small hack)
	// TODO use the contentType to return the configurations
	configurations, error := cclient.memDiscovery.pullConfigurationsFromServer(serviceName)
	if error != nil {
		return nil, error
	}
	configurationsValue, ok := configurations[key]
	if !ok {
		openlogging.GetLogger().Error("Error in fetching the configurations for particular value,No Key found : " + key)
	}

	return configurationsValue, nil
}

// PullConfigsByDI pulls the configuration for custom DimensionInfo
func (cclient *Client) PullConfigsByDI(dimensionInfo, diInfo string) (map[string]map[string]interface{}, error) {
	// update dimensionInfo value
	type GetConfigAPI map[string]map[string]interface{}
	configAPIRes := make(GetConfigAPI)
	parsedDimensionInfo := strings.Replace(diInfo, "#", "%23", -1)
	restApi := ConfigPath + "?" + dimensionsInfo + "=" + parsedDimensionInfo
	err := cclient.memDiscovery.call(http.MethodGet, restApi, nil, nil, &configAPIRes)
	if err != nil {
		openlogging.GetLogger().Error("Pull config by DI failed:" + err.Error())
		return nil, err

	}
	return configAPIRes, nil
}

// PushConfigs push configs to ConfigSource cc , success will return { "Result": "Success" }
func (cclient *Client) PushConfigs(items map[string]interface{}, dimensionInfo string) (map[string]interface{}, error) {
	if len(items) == 0 {
		em := "data is empty , which data need to send cc"
		openlogging.GetLogger().Error(em)
		return nil, errors.New(em)
	}
	type CreateConfigApi struct {
		DimensionInfo string                 `json:"dimensionsInfo"`
		Items         map[string]interface{} `json:"items"`
	}
	configApi := CreateConfigApi{
		DimensionInfo: dimensionInfo,
		Items:         items,
	}

	return addDeleteConfig(cclient, configApi, http.MethodPost)
}

// DeleteConfigsByKeys
func (cclient *Client) DeleteConfigsByKeys(keys []string, dimensionInfo string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		em := "not key need to delete for cc, please check keys"
		openlogging.GetLogger().Error(em)
		return nil, errors.New(em)
	}
	type DeleteConfigApi struct {
		DimensionInfo string   `json:"dimensionsInfo"`
		Keys          []string `json:"keys"`
	}
	configApi := DeleteConfigApi{
		DimensionInfo: dimensionInfo,
		Keys:          keys,
	}

	return addDeleteConfig(cclient, configApi, http.MethodDelete)
}
func addDeleteConfig(cclient *Client, data interface{}, method string) (map[string]interface{}, error) {
	type ConfigAPI map[string]interface{}
	configAPIS := make(ConfigAPI)
	body, err := serializers.Encode(serializers.JsonEncoder, data)
	if err != nil {
		openlogging.GetLogger().Errorf("serializer data failed , err :", err.Error())
		return nil, err
	}
	err = cclient.memDiscovery.call(method, ConfigPath, nil, body, &configAPIS)
	if err != nil {
		return nil, err
	}
	return configAPIS, nil
}
func init() {
	ccclient.InstallConfigClientPlugin(Name, NewConfigCenter)
}

func isStatusSuccess(i int) bool {
	return i >= http.StatusOK && i < http.StatusBadRequest
}

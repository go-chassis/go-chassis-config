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
package memberdiscovery

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/ServiceComb/go-cc-client"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/http-client"
	"os"

	"github.com/ServiceComb/go-cc-client/serializers"
	"io/ioutil"
)

var (
	memDiscovery      *MemDiscovery
	//HeaderTenantName is a variable of type string
	HeaderTenantName  = "X-Tenant-Name"
	//ConfigMembersPath is a variable of type string
	ConfigMembersPath = ""
)

const (
	//StatusUP is a variable of type string
	StatusUP            = "UP"
	//HeaderContentType is a variable of type string
	HeaderContentType  = "Content-Type"
	//HeaderUserAgent is a variable of type string
	HeaderUserAgent    = "User-Agent"
	members             = "/configuration/members"
	defaultContentType  = "application/json"
)
//MemberDiscoveryService
var MemberDiscoveryService MemberDiscovery
//MemberDiscovery is a interface
type MemberDiscovery interface {
	ConfigurationInit([]string) error
	GetConfigServer() ([]string, error)
	RefreshMembers() error
	Shuffle() error
	GetWorkingConfigCenterIP([]string) ([]string, error)
}
//MemDiscovery is a struct
type MemDiscovery struct {
	ConfigServerAddresses []string
	//Logger                *log.Entry
	IsInit     bool
	TLSConfig  *tls.Config
	TenantName string
	EnableSSL  bool
	sync.RWMutex
	client *httpclient.URLClient
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
//NewConfiCenterInit is a function
func NewConfiCenterInit(tlsConfig *tls.Config, tenantName string, enableSSL bool) MemberDiscovery {
	if memDiscovery == nil {
		memDiscovery = new(MemDiscovery)
		//memDiscovery.Logger = logger
		memDiscovery.TLSConfig = tlsConfig
		memDiscovery.TenantName = tenantName
		memDiscovery.EnableSSL = enableSSL
		var apiVersion string

		switch config.GlobalDefinition.Cse.Config.Client.APIVersion.Version {
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
		options := &httpclient.URLClientOption{
			SSLEnabled: enableSSL,
			TLSConfig:  tlsConfig,
			Compressed: false,
			Verbose:    false,
		}
		memDiscovery.client, _ = httpclient.GetURLClient(options)
	}
	return memDiscovery
}

//HttpDo Use http-client package for rest communication
func (memDis *MemDiscovery) HTTPDo(method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	if len(headers) == 0 {
		headers = make(http.Header)
	}
	for k, v := range GetDefaultHeaders(memDis.TenantName) {
		headers[k] = v
	}
	return memDis.client.HttpDo(method, rawURL, headers, body)
}

//Update the Base PATH and HEADERS Based on the version of ConfigCenter used.
func updateAPIPath(apiVersion string) {

	//Check for the env Name in Container to get Domain Name
	//Default value is  "default"
	projectID, isExsist := os.LookupEnv(common.EnvProjectID)
	if !isExsist {
		projectID = "default"
	}
	switch apiVersion {
	case "v3":
		ConfigMembersPath = "/v3/" + projectID + members
		HeaderTenantName = "X-Tenant-Name"
	case "v2":
		ConfigMembersPath = "/members"
		HeaderTenantName = "X-Tenant-Name"
	default:
		ConfigMembersPath = "/v3/" + projectID + members
		HeaderTenantName = "X-Tenant-Name"
	}
}
//ConfigurationInit is a method for creating a configuration
func (memDis *MemDiscovery) ConfigurationInit(initConfigServer []string) error {
	if memDis.IsInit == true {
		return nil
	}

	if memDis.ConfigServerAddresses == nil {
		if initConfigServer == nil && len(initConfigServer) == 0 {
			err := errors.New(client.EmptyConfigServerConfig)
			lager.Logger.Error(client.EmptyConfigServerConfig, err)
			return err
		}

		memDis.ConfigServerAddresses = make([]string, 0)
		for _, server := range initConfigServer {
			memDis.ConfigServerAddresses = append(memDis.ConfigServerAddresses, server)
		}

		memDis.Shuffle()
	}

	memDis.IsInit = true
	return nil
}
//GetConfigServer is a method used for getting server configuration
func (memDis *MemDiscovery) GetConfigServer() ([]string, error) {
	if memDis.IsInit == false {
		err := errors.New(client.PackageInitError)
		lager.Logger.Error(client.PackageInitError, err)
		return nil, err
	}

	if len(memDis.ConfigServerAddresses) == 0 {
		err := errors.New(client.EmptyConfigServerMembers)
		lager.Logger.Error(client.EmptyConfigServerMembers, err)
		return nil, err
	}

	if config.GlobalDefinition.Cse.Config.Client.Autodiscovery {
		err := memDis.RefreshMembers()
		if err != nil {
			lager.Logger.Error("refresh member is failed", err)
			return nil, err
		}
	} else {
		tmpConfigAddrs := memDis.ConfigServerAddresses
		for key := range tmpConfigAddrs {
			if !strings.Contains(memDis.ConfigServerAddresses[key], "https") && memDis.EnableSSL {
				memDis.ConfigServerAddresses[key] = `https://` + memDis.ConfigServerAddresses[key]

			} else if !strings.Contains(memDis.ConfigServerAddresses[key], "http") {
				memDis.ConfigServerAddresses[key] = `http://` + memDis.ConfigServerAddresses[key]
			}
		}
	}

	err := memDis.Shuffle()
	if err != nil {
		lager.Logger.Error("member shuffle is failed", err)
		return nil, err
	}

	memDis.RLock()
	defer memDis.RUnlock()
	lager.Logger.Debugf("member server return %s", memDis.ConfigServerAddresses[0])
	return memDis.ConfigServerAddresses, nil
}
//RefreshMembers is a method
func (memDis *MemDiscovery) RefreshMembers() error {
	var (
		errorStatus bool
		errorInfo   error
		count       int
	)

	endpointMap := make(map[string]bool)

	if len(memDis.ConfigServerAddresses) == 0 {
		return nil
	}

	tmpConfigAddrs := memDis.ConfigServerAddresses
	confgCenterIP := len(tmpConfigAddrs)
	instances := new(Members)
	for _, host := range tmpConfigAddrs {
		errorStatus = false
		lager.Logger.Debugf("RefreshMembers hosts ", host)
		resp, err := memDis.HTTPDo("GET", host+ConfigMembersPath, nil, nil)
		if err != nil {
			errorStatus = true
			errorInfo = err
			count++
			if confgCenterIP > count {
				errorStatus = false
			}
			lager.Logger.Error("member request failed with error", err)
			continue
		}
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		contentType := resp.Header.Get("Content-Type")
		if len(contentType) > 0 && (len(defaultContentType) > 0 && !strings.Contains(contentType, defaultContentType)) {
			lager.Logger.Error("config source member request failed with error", errors.New("content type mis match"))
			continue
		}
		error := serializers.Decode(defaultContentType, body, &instances)
		if error != nil {
			lager.Logger.Error("config source member request failed with error", errors.New("error in decoding the request"))
			continue
		}
		for _, instance := range instances.Instances {
			if instance.Status != StatusUP {
				continue
			}
			for _, entryPoint := range instance.EntryPoints {
				endpointMap[entryPoint] = memDis.EnableSSL
			}
		}
	}
	if errorStatus {
		return errorInfo
	}

	memDis.Lock()
	// flush old config
	memDis.ConfigServerAddresses = make([]string, 0)
	var entryPoint string
	for endPoint, isHTTPSEnable := range endpointMap {
		parsedEndpoint := strings.Split(endPoint, `://`)
		if len(parsedEndpoint) != 2 {
			continue
		}
		if isHTTPSEnable {
			entryPoint = `https://` + parsedEndpoint[1]
		} else {
			entryPoint = `http://` + parsedEndpoint[1]
		}
		memDis.ConfigServerAddresses = append(memDis.ConfigServerAddresses, entryPoint)
	}
	memDis.Unlock()
	return nil
}
//GetDefaultHeaders gets default headers
func GetDefaultHeaders(tenantName string) http.Header {
	headers := http.Header{
		HeaderContentType: []string{"application/json"},
		HeaderUserAgent:   []string{"cse-configcenter-client/1.0.0"},
		HeaderTenantName:  []string{tenantName},
	}

	return headers
}
//Shuffle is a method to log error
func (memDis *MemDiscovery) Shuffle() error {
	if memDis.ConfigServerAddresses == nil || len(memDis.ConfigServerAddresses) == 0 {
		err := errors.New(client.EmptyConfigServerConfig)
		lager.Logger.Error(client.EmptyConfigServerConfig, err)
		return err
	}

	perm := rand.Perm(len(memDis.ConfigServerAddresses))

	memDis.Lock()
	defer memDis.Unlock()
	lager.Logger.Debugf("Before Suffled member %s ", memDis.ConfigServerAddresses)
	for i, v := range perm {
		lager.Logger.Debugf("shuffler %d %d", i, v)
		tmp := memDis.ConfigServerAddresses[v]
		memDis.ConfigServerAddresses[v] = memDis.ConfigServerAddresses[i]
		memDis.ConfigServerAddresses[i] = tmp
	}

	lager.Logger.Debugf("Suffled member %s", memDis.ConfigServerAddresses)
	return nil
}
//GetWorkingConfigCenterIP is a method which gets working configuration center IP
func (memDis *MemDiscovery) GetWorkingConfigCenterIP(entryPoint []string) ([]string, error) {
	instances := new(Members)
	ConfigServerAddresses := make([]string, 0)
	for _, server := range entryPoint {
		resp, err := memDis.HTTPDo("GET", server+ConfigMembersPath, nil, nil)
		if err != nil {
			lager.Logger.Error("config source member request failed with error", err)
			continue
		}
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		contentType := resp.Header.Get("Content-Type")
		if len(contentType) > 0 && (len(defaultContentType) > 0 && !strings.Contains(contentType, defaultContentType)) {
			lager.Logger.Error("config source member request failed with error", errors.New("content type mis match"))
			continue
		}
		error := serializers.Decode(defaultContentType, body, &instances)
		if error != nil {
			lager.Logger.Error("config source member request failed with error", errors.New("error in decoding the request"))
			continue
		}
		ConfigServerAddresses = append(ConfigServerAddresses, server)
	}
	return ConfigServerAddresses, nil
}

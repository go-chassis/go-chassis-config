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

package servicecombkie

import (
	"context"
	"errors"
	"fmt"
	sckieclient "github.com/apache/servicecomb-kie/client"
	"github.com/apache/servicecomb-kie/pkg/model"
	"github.com/go-chassis/go-archaius/sources/utils"
	"github.com/go-chassis/go-chassis-config"
	"github.com/go-mesh/openlogging"
)

// Client contains the implementation of Client
type Client struct {
	KieClient     *sckieclient.Client
	serviceName   string
	version       string
	URI           string
	EnableSSL     bool
	AutoDiscovery bool
	Namespace     string
	TenantName    string
}

const (
	//Name of the Plugin
	Name = "servicecomb-kie"
)

const (
	RetSucc   = 0
	RetDelErr = 1
)

// NewKieClient init's the necessary objects needed for seamless communication to Kie Server
func (kieClient *Client) NewKieClient() {
	defaultLabels := make(map[string]string)
	configInfo := sckieclient.Config{Endpoint: kieClient.URI, DefaultLabels: defaultLabels, VerifyPeer: kieClient.EnableSSL}
	var err error
	kieClient.KieClient, err = sckieclient.New(configInfo)
	if err != nil {
		openlogging.GetLogger().Error("KieClient Initialization Failed: " + err.Error())
	}
	openlogging.GetLogger().Debugf("KieClient Initialized successfully")
}

// PullConfigs is used for pull config from servicecomb-kie
func (kieClient *Client) PullConfigs(serviceName, version, app, env string) (map[string]interface{}, error) {
	openlogging.GetLogger().Debugf("KieClient begin PullConfigs")
	labels := map[string]string{"servicename": serviceName, "version": version, "app": app, "env": env}
	configsInfo := make(map[string]interface{})
	configurationsValue, err := kieClient.KieClient.SearchByLabels(context.TODO(), sckieclient.WithGetProject(serviceName), sckieclient.WithLabels(labels))
	if err != nil {
		openlogging.GetLogger().Errorf("Error in Querying the Response from Kie %s %s %s %s %s", err.Error(), serviceName, version, app, env)
		return nil, err
	}
	openlogging.GetLogger().Debugf("KieClient begin PullConfigs1")
	openlogging.GetLogger().Debugf("KieClient SearchByLabels. %s %s %s %s", serviceName, version, app, env)
	//Parse config result.
	for _, docRes := range configurationsValue {
		for _, docInfo := range docRes.Data {
			configsInfo[docInfo.Key] = docInfo.Value
			configDetail, err := utils.Convert2JavaProps(docInfo.Key, []byte(docInfo.Value))
			if err != nil {
				openlogging.GetLogger().Errorf("Error in Parse the Response from Kie %s %s %s %s %s ", err.Error(), serviceName, version, app, env)
			}
			for key, value := range configDetail {
				configsInfo[key] = value
			}
		}
	}
	return configsInfo, nil
}

// PullConfig get config by key and labels.
func (kieClient *Client) PullConfig(serviceName, version, app, env, key, contentType string) (interface{}, error) {
	fmt.Println("PullConfig")
	labels := map[string]string{"servicename": serviceName, "version": version, "app": app, "env": env}
	configurationsValue, err := kieClient.KieClient.Get(context.TODO(), key, sckieclient.WithGetProject(serviceName), sckieclient.WithLabels(labels))
	if err != nil {
		openlogging.GetLogger().Error("Error in Querying the Response from Kie: " + err.Error())
		return nil, err
	}
	for _, doc := range configurationsValue {
		for _, kvDoc := range doc.Data {
			if key == kvDoc.Key {
				openlogging.GetLogger().Debugf("The Key Value of : ", kvDoc.Value)
				return doc, nil
			}
		}
	}
	return nil, errors.New("can not find value")
}

//PullConfigsByDI not implemented
func (kieClient *Client) PullConfigsByDI(dimensionInfo string) (map[string]map[string]interface{}, error) {
	// TODO Return the configurations for customized Projects in Kie Configs
	fmt.Println(" PullConfigsByDI")
	return nil, errors.New("not implemented")
}

//PushConfigs put config in kie by key and labels.
func (kieClient *Client) PushConfigs(data map[string]interface{}, serviceName, version, app, env string) (map[string]interface{}, error) {
	fmt.Println("PushConfigs")
	var configReq model.KVDoc
	labels := map[string]string{"servicename": serviceName, "version": version, "app": app, "env": env}
	configResult := make(map[string]interface{})
	for key, configValue := range data {
		configReq.Key = key
		configReq.Value = configValue.(string)
		configReq.Labels = labels
		configurationsValue, err := kieClient.KieClient.Put(context.TODO(), configReq, sckieclient.WithProject(serviceName))
		if err != nil {
			openlogging.GetLogger().Error("Error in PushConfigs to Kie: " + err.Error())
			return nil, err
		}
		openlogging.GetLogger().Debugf("The Key Value of : ", configurationsValue)
		configResult[configurationsValue.Key] = configurationsValue.Value
	}
	return configResult, nil
}

//DeleteConfigsByKeys use keyId for delete
func (kieClient *Client) DeleteConfigsByKeys(keys []string, serviceName, version, app, env string) (map[string]interface{}, error) {
	fmt.Println("DeleteConfigsByKeys")
	result := make(map[string]interface{})
	for _, keyId := range keys {
		err := kieClient.KieClient.Delete(context.TODO(), keyId, "", sckieclient.WithProject(serviceName))
		if err != nil {
			openlogging.GetLogger().Errorf("Error in Delete from Kie. %s " + err.Error())
			result[keyId] = RetDelErr
		} else {
			openlogging.GetLogger().Debugf("Delete The KeyId:%s", keyId)
			result[keyId] = RetSucc
		}
	}
	return result, nil
}

//Watch not implemented because kie not support.
func (kieClient *Client) Watch(f func(map[string]interface{}), errHandler func(err error)) error {
	return nil
}

//Options.
func (kieClient *Client) Options() config.Options {
	fmt.Println("Options")
	optionInfo := config.Options{
		ServiceName:   kieClient.serviceName,
		Version:       kieClient.version,
		Endpoint:      kieClient.URI,
		EnableSSL:     kieClient.EnableSSL,
		AutoDiscovery: kieClient.AutoDiscovery,
		Namespace:     kieClient.Namespace,
		TenantName:    kieClient.TenantName,
	}
	return optionInfo
}

//InitConfigKie initialize the Kie Client
func InitConfigKie(options config.Options) (config.Client, error) {
	fmt.Println("InitConfigKie")
	kieClient := &Client{
		serviceName:   options.ServiceName,
		version:       options.Version,
		URI:           options.ServerURI,
		EnableSSL:     options.EnableSSL,
		AutoDiscovery: options.AutoDiscovery,
		Namespace:     options.Namespace,
		TenantName:    options.TenantName,
	}
	kieClient.NewKieClient()
	return kieClient, nil
}

func init() {
	config.InstallConfigClientPlugin(Name, InitConfigKie)
	//fmt.Println("init config client plugin:%s", Name)
}

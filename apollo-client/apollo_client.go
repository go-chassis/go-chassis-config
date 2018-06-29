package apolloclient

import (
	"errors"
	"github.com/ServiceComb/go-cc-client/serializers"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/http-client"
	"io/ioutil"
	"net/http"
	"strings"
)

// ApolloClient contains the implementation of ConfigClient
type ApolloClient struct {
	name   string
	client *httpclient.URLClient
}

const apolloServerAPI = ":ServerURL/configs/:appID/:clusterName/:nameSpace"
const defaultContentType = "application/json"

// NewApolloClient init's the necessary objects needed for seamless communication to apollo Server
func (apolloClient *ApolloClient) NewApolloClient() {
	options := &httpclient.URLClientOption{
		SSLEnabled: false,
		TLSConfig:  nil, //TODO Analyse the TLS configuration of Apollo Server
		Compressed: false,
		Verbose:    false,
	}
	var err error
	apolloClient.client, err = httpclient.GetURLClient(options)
	if err != nil {
		lager.Logger.Error("ApolloClient Initialization Failed", err)
	}
	lager.Logger.Debugf("ApolloClient Initialized successfully")
}

// Init will initialize the needed parameters
func (apolloClient *ApolloClient) Init() {
	lager.Logger.Debugf("ApolloClient Initialized successfully")
}

// HTTPDo Use http-client package for rest communication
func (apolloClient *ApolloClient) HTTPDo(method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	return apolloClient.client.HttpDo(method, rawURL, headers, body)
}

// PullConfigs is the implementation of ConfigClient and pulls all the configuration for a given serviceName
func (apolloClient *ApolloClient) PullConfigs(serviceName, version, app, env string) (map[string]interface{}, error) {
	/*
		1. Compose the URL
		2. Make a Http Request to Apollo Server
		3. Unmarshal the response
		4. Return back the configuration/error
		Note: Currently the input to this function in not used, need to check it's feasibility of using it, as the serviceName/version can be different in Apollo
	*/

	// Compose the URL
	pullConfigurationURL := composeURL()

	// Make a Http Request to Apollo Server
	resp, err := apolloClient.HTTPDo("GET", pullConfigurationURL, nil, nil)
	if err != nil {
		lager.Logger.Error("Error in Querying the Response from Apollo", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		lager.Logger.Error("Bad Response : ", errors.New("Response from Apollo Server "+resp.Status))
		return nil, errors.New("Bad Response from Apollo Server " + resp.Status)
	}
	/*
		Sample Response from Apollo Server
		{
			"appId": "SampleApp",
			"cluster": "default",
			"namespaceName": "application",
			"configurations": {
				"timeout": "500"
			},
			"releaseKey": "20180327130726-1dc5027439679153"
		}
	*/

	//Unmarshal the response
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	var configurations map[string]interface{}
	error := serializers.Decode(defaultContentType, body, &configurations)
	if error != nil {
		lager.Logger.Error("Error in Unmarshalling the Response from Apollo", error)
		return nil, error
	}

	lager.Logger.Debugf("The Marshaled response of the body is : ", configurations["configurations"])

	return configurations, nil
}

// PullConfig is the implementation of the ConfigClient
func (apolloClient *ApolloClient) PullConfig(serviceName, version, app, env, key, contentType string) (interface{}, error) {
	/*
		1. Compose the URL
		2. Make a Http Request to Apollo Server
		3. Unmarshal the response
		4. Get the particular key/value
		4. Return back the value/error
		//TODO Use the contentType to send the response
	*/

	// Compose the URL
	pullConfigurationURL := composeURL()

	// Make a Http Request to Apollo Server
	resp, err := apolloClient.HTTPDo("GET", pullConfigurationURL, nil, nil)
	if err != nil {
		lager.Logger.Error("Error in Querying the Response from Apollo", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		lager.Logger.Error("Bad Response : ", errors.New("Response from Apollo Server "+resp.Status))
		return nil, errors.New("Bad Response from Apollo Server " + resp.Status)
	}

	//Unmarshal the response
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	var configurations map[string]interface{}
	error := serializers.Decode(defaultContentType, body, &configurations)
	if error != nil {
		lager.Logger.Error("Error in Unmarshalling the Response from Apollo", error)
		return nil, err
	}

	//Find the particular Key
	configList := configurations["configurations"]
	configurationsValue := ""
	isFound := false

	for configKey, configValue := range configList.(map[string]interface{}) {
		if configKey == key {
			configurationsValue = configValue.(string)
			isFound = true
		}
	}

	if !isFound {
		lager.Logger.Error("Error in fetching the configurations for particular value", errors.New("No Key found : "+key))
		return nil, errors.New("No Key found : " + key)
	}
	lager.Logger.Debugf("The Key Value of : ", configurationsValue)
	return configurationsValue, nil
}

// composeURL composes the URL based on the configurations given in chassis.yaml
func composeURL() string {
	pullConfigurationURL := strings.Replace(apolloServerAPI, ":ServerURL", config.GlobalDefinition.Cse.Config.Client.ServerURI, 1)
	pullConfigurationURL = strings.Replace(pullConfigurationURL, ":appID", config.GlobalDefinition.Cse.Config.Client.ApolloServiceName, 1)
	pullConfigurationURL = strings.Replace(pullConfigurationURL, ":clusterName", config.GlobalDefinition.Cse.Config.Client.ClusterName, 1)
	pullConfigurationURL = strings.Replace(pullConfigurationURL, ":nameSpace", config.GlobalDefinition.Cse.Config.Client.ApolloNameSpace, 1)
	return pullConfigurationURL
}

//PullConfigsByDI returns the configuration for additional Projects in Apollo
func (apolloClient *ApolloClient) PullConfigsByDI(dimensionInfo, diInfo string) (map[string]map[string]interface{}, error) {
	// TODO Return the configurations for customized Projects in Apollo Configs
	return nil, nil
}

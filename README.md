### go-cc-client
[![Build Status](https://travis-ci.org/go-chassis/go-cc-client.svg?branch=master)](https://travis-ci.org/go-chassis/go-cc-client)  

Supported config center:

- ctrip apollo https://github.com/ctripcorp/apollo
- huawei cloud CSE config center https://www.huaweicloud.com/product/cse.html


# Example
Get a client of config center

1. import the config client you want to use 
``go
_ "github.com/go-chassis/go-cc-client/configcenter"
``

2. New a client 
``go
c, err := ccclient.NewClient("config_center", ccclient.Options{
		ServerURI: "http://127.0.0.1:30200",
	})
``

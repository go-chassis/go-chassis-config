### go-cc-client
[![Build Status](https://travis-ci.org/go-chassis/go-cc-client.svg?branch=master)](https://travis-ci.org/go-chassis/go-cc-client)  
config center client can pull and push configs in distributed configuration 
management service

Supported distributed configuration management service:

| name       | import                                         |description    |
|----------|----------|:-------------:|
|config_name |github.com/go-chassis/go-cc-client/configcenter |huawei cloud CSE config center https://www.huaweicloud.com/product/cse.html |
|apollo      |github.com/go-chassis/go-cc-client/apollo       |ctrip apollo https://github.com/ctripcorp/apollo |

# Example
Get a client of config center

1. import the config client you want to use 
```go
import _ "github.com/go-chassis/go-cc-client/configcenter"
```

2. Create a client 
```go
c, err := ccclient.NewClient("config_center", ccclient.Options{
		ServerURI: "http://127.0.0.1:30200",
	})
````

# Use huawei cloud 
```go
import (
	"github.com/huaweicse/auth"
	"github.com/go-chassis/foundation/httpclient"
	_ "github.com/go-chassis/go-cc-client/configcenter"
)

func main() {
	var err error
	httpclient.SignRequest,err =auth.GetShaAKSKSignFunc("your ak", "your sk", "")
	if err!=nil{
        //handle err
	}
	ccclient.NewClient("config_center",ccclient.Options{
		ServerURI:"the address of CSE endpoint",
	})
}

```
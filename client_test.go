package ccclient_test

import (
	"github.com/go-chassis/go-cc-client"
	_ "github.com/go-chassis/go-cc-client/configcenter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnable(t *testing.T) {
	c, err := ccclient.NewClient("config_center", ccclient.Options{
		ServerURI: "http://127.0.0.1:30100",
	})
	assert.NoError(t, err)
	_, err = c.PullConfigs("service", "app", "1.0", "")
	assert.Error(t, err)
}

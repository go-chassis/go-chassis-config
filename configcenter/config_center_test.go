package configcenter_test

import (
	"github.com/go-chassis/go-chassis-config"
	"github.com/go-chassis/go-chassis-config/configcenter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfigCenter(t *testing.T) {
	c, err := configcenter.NewConfigCenter(config.Options{
		ServerURI: "http://",
		Labels:    map[string]string{"app": "default"}})
	assert.NoError(t, err)
	assert.Equal(t, "default", c.Options().Labels["app"])
}

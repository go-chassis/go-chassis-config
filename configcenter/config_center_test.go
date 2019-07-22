package configcenter_test

import (
	"github.com/go-chassis/go-chassis-config"
	"github.com/go-chassis/go-chassis-config/configcenter"
	"github.com/go-chassis/paas-lager"
	"github.com/go-mesh/openlogging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	log.Init(log.Config{
		LoggerLevel:   "DEBUG",
		EnableRsyslog: false,
		LogFormatText: false,
		Writers:       []string{"stdout"},
	})
	l := log.NewLogger("test")
	openlogging.SetLogger(l)
}

func TestNewConfigCenter(t *testing.T) {
	c, _ := configcenter.NewConfigCenter(config.Options{App: "default"})
	assert.Equal(t, "default", c.Options().App)
}

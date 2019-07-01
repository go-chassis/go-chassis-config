package configcenter_test

import (
	"github.com/go-chassis/go-chassis-config/pkg/configcenter"
	"github.com/go-chassis/paas-lager"
	"github.com/go-mesh/openlogging"
	"testing"
)

func init() {
	log.Init(log.Config{
		LoggerLevel:   "DEBUG",
		EnableRsyslog: false,
		LogFormatText: true,
		Writers:       []string{"stdout"},
	})
	l := log.NewLogger("test")
	openlogging.SetLogger(l)
}
func TestNew(t *testing.T) {
	configcenter.New(configcenter.Options{})
}

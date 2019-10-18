package configcenter_test

import (
	"github.com/go-chassis/go-chassis-config/pkg/configcenter"
	"testing"
)

func TestNew(t *testing.T) {
	configcenter.New(configcenter.Options{})
}

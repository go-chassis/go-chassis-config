package client_test

import (
	"github.com/go-chassis/go-cc-client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnable(t *testing.T) {
	err := client.Enable("")
	assert.Error(t, err)
}

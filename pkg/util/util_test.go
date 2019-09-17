package util_test

import (
	"github.com/go-chassis/go-chassis-config/pkg/util"
	"testing"
)

func TestMap2String(t *testing.T) {
	m := make(map[string]string)
	m["s"] = "a"
	m["c"] = "c"
	m["d"] = "b"
	t.Log(util.Map2String(m))
}

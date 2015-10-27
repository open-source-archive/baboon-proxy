package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowLTMVirtualServer(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	object := ShowLTMVirtualServer(url)
	assert.NotNil(t, object)
}

func TestShowLTMVirtualServerName(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	vserver := "vs_www.zalando.de_80"
	object := ShowLTMVirtualServerName(url, vserver)
	assert.NotNil(t, object)
}

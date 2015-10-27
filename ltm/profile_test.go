package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowLTMProfile(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	vserver := "vs_www.zalando.de_80"
	object := ShowLTMProfile(url, vserver)
	assert.NotNil(t, object)
}

package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowLTMFWRules(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	vserver := "vs_www.zalando.de_80"
	object := ShowLTMFWRules(url, vserver)
	assert.NotNil(t, object)
}

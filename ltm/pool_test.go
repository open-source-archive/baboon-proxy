package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowLTMPools(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	object := client.ShowLTMPools(url)
	assert.NotNil(t, object)
}

func TestShowLTMPool(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	pool := "itr-http"
	object := client.ShowLTMPool(url, pool)
	assert.NotNil(t, object)
}

func TestShowLTMPoolMember(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/ltm/"
	pool := "itr-http"
	object := client.ShowLTMPoolMember(url, pool)
	assert.NotNil(t, object)
}

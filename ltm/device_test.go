package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const checkMark = "\u2713"
const ballotX = "\u2717"

func TestShowLTMDevice(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/cm/device/"
	object := ShowLTMDevice(url)
	assert.NotNil(t, object)
}

func TestShowLTMDeviceName(t *testing.T) {
	url := "https://itr-ltm01/mgmt/tm/cm/device/"
	host := "itr-ltm01"
	object := ShowLTMDeviceName(url, host)
	t.Log("\tShould return specific information about loadbalancer device", checkMark)
	assert.NotNil(t, object)
}

package ltm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShowLTMDataGroup(t *testing.T) {
	host := "https://itr-ltm01/mgmt/tm/ltm/"
	direction := "external"
	object := ShowLTMDataGroup(host, direction)
	assert.NotNil(t, object)
}

func TestShowLTMDataGroupName(t *testing.T) {
	host := "https://itr-ltm01/mgmt/tm/ltm/"
	direction := "external"
	datagroupname := "aws_allowed"
	object := ShowLTMDataGroupName(host, direction, datagroupname)
	assert.NotNil(t, object)
}

func TestPostLTMDataGroupName(t *testing.T) {
	host := "https://itr-ltm01/mgmt/tm/ltm/"
	direction := "external"
	ccg := CreateLTMDatagroup{Name: "aws_allowed", Type: "ip"}
	object, _ := PostLTMDatagroup(host, direction, &ccg)
	//assert.Nil(t, err)
	assert.NotNil(t, object)
}

// Needs to be implemented
//func TestPostLTMDataGroupNameItem(t *testing.T) {
//}

func TestDeleteLTMDataGroup(t *testing.T) {
	host := "https://itr-ltm01/mgmt/tm/ltm/"
	direction := "external"
	ddg := RemoveLTMDatagroup{Name: "aws_allowed"}
	object, err := DeleteLTMDatagroup(host, direction, ddg.Name)
	assert.Nil(t, err)
	assert.NotNil(t, object)
}

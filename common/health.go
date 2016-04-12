package common

import (
	"github.com/alouca/gosnmp"
)

type HealthDevices []HealthDevice

type HealthDevice struct {
	Name         string `json:"name"`
	ManagementIP string `json:"managementIP"`
	Status       string `json:"status"`
	MajorVersion string `json:"majorVersion"`
	MinorVersion string `json:"minorVersion"`
}

type HealthDeviceChannel struct {
	Status       *gosnmp.SnmpPacket
	MajorVersion *gosnmp.SnmpPacket
	MinorVersion *gosnmp.SnmpPacket
}

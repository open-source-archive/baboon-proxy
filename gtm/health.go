package gtm

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

type HealthPools []HealthPool

type HealthPool struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type HealthWideIPs []HealthWideIP

type HealthWideIP struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type HealthPoolChannel struct {
	Name   []gosnmp.SnmpPDU
	Status []gosnmp.SnmpPDU
}

type HealthWideIPChannel struct {
	Name   []gosnmp.SnmpPDU
	Status []gosnmp.SnmpPDU
}

type HealthGSLBServers struct {
	Name   string             `json:"name"`
	Server []HealthGSLBServer `json:"gslbServer"`
}

type HealthGSLBServer struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type HealthGSLBChannel struct {
	Name   []gosnmp.SnmpPDU
	Status []gosnmp.SnmpPDU
}

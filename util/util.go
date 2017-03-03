package util

import (
	"bytes"
	"net"
	"strings"

	"github.com/zalando-techmonkeys/baboon-proxy/common"
)

// ReplaceCommon takes common partition and replaces with an empty string
func ReplaceCommon(s string) string {
	const partition string = "/Common/"
	const emptystring string = ""
	return strings.Replace(s, partition, emptystring, -1)
}

// ReplaceColon splits IP port by colon and return IP only
func ReplaceColon(s string) string {
	const colon string = ":"
	value := strings.Split(s, colon)
	ip, _ := value[0], value[1]
	return ip
}

// ReplaceLTMUritoDeviceURI takes local taffic manager path
// and replaces it with device path
func ReplaceLTMUritoDeviceURI(s string) string {
	return strings.Replace(s, common.LtmURI, common.DeviceURI, -1)
}

// ReplaceLTMUritoAddressListURI takes local taffic manager path
// and replaces it with address list path
func ReplaceLTMUritoAddressListURI(s string) string {
	return strings.Replace(s, common.LtmURI, common.AddressListURI, -1)
}

// ReplaceGTMWipUritoGTMPoolURI replace wideips with pools
func ReplaceGTMWipUritoGTMPoolURI(s string) string {
	return strings.Replace(s, "wideips", "pools/", -1)
}

// VerifyIPv4Scope verifies if ip is in scope
// this snippet checks if a ip which should be blocked
// is not in any whitelist range
func VerifyIPv4Scope(sourceIP, rangeIP string) bool {
	if !(strings.Contains(rangeIP, "-")) {
		if sourceIP == rangeIP {
			return true
		}
		return false
	}
	r := strings.Split(rangeIP, "-")
	fromIP := net.ParseIP(r[0])
	toIP := net.ParseIP(r[1])
	s := net.ParseIP(sourceIP)

	if s.To4() == nil {
		return false
	}
	if bytes.Compare(s, fromIP) >= 0 && bytes.Compare(s, toIP) <= 0 {
		return true
	}
	return false
}

func CheckDeviceStatus(signal int) string {
	var status string
	switch signal {
	case 1:
		status = "Offline"
	case 2:
		status = "ForcedOffline"
	case 3:
		status = "Standby"
	case 4:
		status = "Active"
	default:
		status = "Unknown"
	}
	return status
}

func CheckPoolStatus(signal string) string {
	var status string
	switch signal {
	case "Available":
		status = "Available"
	case "No enabled pool members available":
		status = "No enabled pool members available"
	default:
		status = "Unknown"
	}
	return status
}

func CheckWideIPStatus(signal string) string {
	var status string
	switch signal {
	case "Available":
		status = "Available"
	case "No enabled pools available":
		status = "No enabled pools available"
	default:
		status = "Unknown"
	}
	return status
}

func CheckGSLBServerStatus(signal int) string {
	var status string
	switch signal {
	case 1:
		status = "Available"
	case 2:
		status = "Unavailable"
	case 3:
		status = "No enabled Virtual Server available"
	case 4:
		status = "Unknown"
	case 5:
		status = "Unlicensed"
	}
	return status
}

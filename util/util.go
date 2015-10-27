package util

import (
	"bytes"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net"
	"strings"
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

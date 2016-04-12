package snmp

import (
	"fmt"
	"github.com/alouca/gosnmp"
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"github.com/zalando-techmonkeys/baboon-proxy/gtm"
	"net"
	"time"
)

var (
	conf = config.LoadConfig()
	s    = conf.SNMP
)

func Register(ip string) (*gosnmp.GoSNMP, error) {
	switch s.Version {
	case "2c":
		return gosnmp.NewGoSNMP(ip, s.Community, gosnmp.Version2c, int64(s.TimeOut))
	default:
		return nil, fmt.Errorf("SNMP Version %s is not supported", s.Version)
	}
}

func GetServer(ip string, cs chan common.HealthDeviceChannel) {
	globalOid := s.OIDs["common"].(map[interface{}]interface{})
	c, err := Register(ip)
	checkErr(err)
	versionOid := globalOid["version"]
	hotfixOid := globalOid["hotfix"]
	deviceStatusOid := globalOid["devicestatus"]
	version, err := c.Get(versionOid.(string))
	checkErr(err)
	hotfix, err := c.Get(hotfixOid.(string))
	checkErr(err)
	status, err := c.Get(deviceStatusOid.(string))
	checkErr(err)
	cs <- common.HealthDeviceChannel{status, version, hotfix}
}

func GetGTMPool(ip string, cs chan gtm.HealthPoolChannel) {
	gtmOid := s.OIDs["gtm"].(map[interface{}]interface{})
	c, err := Register(ip)
	checkErr(err)
	poolNameOid := gtmOid["pool"]
	poolStatusOid := gtmOid["poolstatus"]
	poolName, err := c.Walk(poolNameOid.(string))
	checkErr(err)
	poolStatus, err := c.Walk(poolStatusOid.(string))
	checkErr(err)
	cs <- gtm.HealthPoolChannel{poolName, poolStatus}
}

func GetGTMWideIP(ip string, cs chan gtm.HealthWideIPChannel) {
	gtmOid := s.OIDs["gtm"].(map[interface{}]interface{})
	c, err := Register(ip)
	checkErr(err)
	wideIPNameOid := gtmOid["wideip"]
	wideIPStatusOid := gtmOid["wideipstatus"]
	wideIPName, err := c.Walk(wideIPNameOid.(string))
	checkErr(err)
	wideIPStatus, err := c.Walk(wideIPStatusOid.(string))
	checkErr(err)
	cs <- gtm.HealthWideIPChannel{wideIPName, wideIPStatus}
}

func GetGTMGSLBServer(ip string, cs chan gtm.HealthGSLBChannel) {
	gtmOid := s.OIDs["gtm"].(map[interface{}]interface{})
	c, err := Register(ip)
	checkErr(err)
	serverNameOid := gtmOid["server"]
	serverStatusOid := gtmOid["serverstatus"]
	serverName, err := c.Walk(serverNameOid.(string))
	checkErr(err)
	serverStatus, err := c.Walk(serverStatusOid.(string))
	checkErr(err)
	cs <- gtm.HealthGSLBChannel{serverName, serverStatus}
}

func CheckConnection(trafficManager map[string]string, firstMatchTrafficManager string) string {
	for _, ip := range trafficManager {
		_, err := net.DialTimeout("udp4", fmt.Sprintf("%s:%s", ip, s.Port), time.Duration(s.TimeOut)*time.Second)
		if err != nil {
			glog.Errorf("Could not establish connection to %s, reason: %s", ip, err.Error())
			continue
		}
		firstMatchTrafficManager = ip
		break
	}
	return firstMatchTrafficManager
}

func checkErr(err error) {
	if err != nil {
		glog.Errorf("%s", err.Error())
	}
}

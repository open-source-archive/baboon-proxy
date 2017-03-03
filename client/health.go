package client

import (
	"github.com/gin-gonic/gin"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/gtm"
	"github.com/zalando-techmonkeys/baboon-proxy/snmp"
	"github.com/zalando-techmonkeys/baboon-proxy/util"
	"net/http"
)

var (
	ITMServer = conf.ITMMgmtIP
	LTMServer = conf.LTMMgmtIP
	GTMServer = conf.GTMMgmtIP
)

func HealthServer(c *gin.Context) {
	var (
		devices        common.HealthDevices
		trafficManager map[string]string
	)
	tm := c.Params.ByName("trafficmanager")
	switch tm {
	case "itm":
		trafficManager = ITMServer
	case "gtm":
		trafficManager = GTMServer
	case "ltm":
		trafficManager = LTMServer
	default:
		respondWithStatus(http.StatusNotFound, tm, nil, "TrafficManager not found", conf.Documentation["healthdocumentationuri"], c)
		return
	}
	devChan := make(chan common.HealthDeviceChannel)
	for name, ip := range trafficManager {
		go snmp.GetServer(ip, devChan)
		dev := <-devChan
		for i, v := range dev.MajorVersion.Variables {
			devStatus := util.CheckDeviceStatus(dev.Status.Variables[i].Value.(int))
			devices = append(devices,
				common.HealthDevice{
					Name:         name,
					ManagementIP: ip,
					Status:       devStatus,
					MajorVersion: v.Value.(string),
					MinorVersion: dev.MinorVersion.Variables[i].Value.(string),
				})
		}
	}
	respondWithStatus(http.StatusOK, "", devices, "", conf.Documentation["healthdocumentationuri"], c)
}

func HealthPools(c *gin.Context) {
	var (
		pools                    gtm.HealthPools
		trafficManager           map[string]string
		firstMatchTrafficManager string
	)
	tm := c.Params.ByName("trafficmanager")
	switch tm {
	case "itm":
		trafficManager = ITMServer
	case "gtm":
		trafficManager = GTMServer
	default:
		respondWithStatus(http.StatusNotFound, tm, nil, "TrafficManager not found", conf.Documentation["healthdocumentationuri"], c)
		return
	}
	poolChan := make(chan gtm.HealthPoolChannel)
	firstMatchTrafficManager = snmp.CheckConnection(trafficManager, firstMatchTrafficManager)
	go snmp.GetGTMPool(firstMatchTrafficManager, poolChan)
	p := <-poolChan
	for i, _ := range p.Name {
		status := util.CheckPoolStatus(p.Status[i].Value.(string))
		pools = append(pools,
			gtm.HealthPool{
				Name:   p.Name[i].Value.(string),
				Status: status})
	}
	respondWithStatus(http.StatusOK, "", pools, "", conf.Documentation["healthdocumentationuri"], c)
}

func HealthWideIPs(c *gin.Context) {
	var (
		wideIPs                  gtm.HealthWideIPs
		trafficManager           map[string]string
		firstMatchTrafficManager string
	)
	tm := c.Params.ByName("trafficmanager")
	switch tm {
	case "itm":
		trafficManager = ITMServer
	case "gtm":
		trafficManager = GTMServer
	default:
		respondWithStatus(http.StatusNotFound, tm, nil, "TrafficManager not found", conf.Documentation["healthdocumentationuri"], c)
		return
	}
	wideIPChan := make(chan gtm.HealthWideIPChannel)
	firstMatchTrafficManager = snmp.CheckConnection(trafficManager, firstMatchTrafficManager)
	go snmp.GetGTMWideIP(firstMatchTrafficManager, wideIPChan)
	wch := <-wideIPChan
	for i, _ := range wch.Name {
		status := util.CheckWideIPStatus(wch.Status[i].Value.(string))
		wideIPs = append(wideIPs,
			gtm.HealthWideIP{
				Name:   wch.Name[i].Value.(string),
				Status: status})
	}
	respondWithStatus(http.StatusOK, "", wideIPs, "", conf.Documentation["healthdocumentationuri"], c)
}

func HealthGSLB(c *gin.Context) {
	var (
		gslbServers    []gtm.HealthGSLBServers
		server         []gtm.HealthGSLBServer
		trafficManager map[string]string
	)
	tm := c.Params.ByName("trafficmanager")
	switch tm {
	case "itm":
		trafficManager = ITMServer
	case "gtm":
		trafficManager = GTMServer
	default:
		respondWithStatus(http.StatusNotFound, tm, nil, "TrafficManager not found", conf.Documentation["healthdocumentationuri"], c)
		return
	}
	gslbChan := make(chan gtm.HealthGSLBChannel)
	for name, ip := range trafficManager {
		go snmp.GetGTMGSLBServer(ip, gslbChan)
		gch := <-gslbChan
		for i, _ := range gch.Name {
			status := util.CheckGSLBServerStatus(gch.Status[i].Value.(int))
			server = append(server, gtm.HealthGSLBServer{
				Name:   gch.Name[i].Value.(string),
				Status: status,
			})
		}
		gslbServers = append(gslbServers,
			gtm.HealthGSLBServers{
				Name:   name,
				Server: server,
			})
	}
	respondWithStatus(http.StatusOK, "", gslbServers, "", conf.Documentation["healthdocumentationuri"], c)
}

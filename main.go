package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/client"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"github.com/zalando-techmonkeys/gin-oauth2"
	"os"
	"strconv"
)

var (
	port       *int
	sslenabled *bool
	gtmenabled *bool
	ltmenabled *bool
)

func usage() {
	fmt.Fprint(os.Stderr, "usage: baboon-proxy -port=80 -ssl-enabled=false -ltm-enabled=true -gtm-enabled-true -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	port = flag.Int("port", 80, "Default Port")
	sslenabled = flag.Bool("ssl-enabled", false, "enable SSL")
	ltmenabled = flag.Bool("ltm-enabled", false, "enable LTM")
	gtmenabled = flag.Bool("gtm-enabled", false, "enable GTM")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	app := gin.New()
	conf := config.LoadConfig()

	// Use Logger, Cross-Origin Ressource and GZIP compression middleware
	app.Use(client.LoggerMiddleware())
	app.Use(client.CORSMiddleware())
	app.Use(gzip.Gzip(gzip.DefaultCompression))

	rootusers, emergencyusers, OAuth2Endpoint, err := config.LoadAuthConf(conf)
	if err != nil {
		glog.Errorf("Could not load configuration. Reason: %s", err.Message)
		panic("Could not load configuration for Baboon. Exiting.")
	}
	glog.Infof("%+v", rootusers)
	if *gtmenabled {

		publicGTM := app.Group("/api/gtms/:trafficmanager")
		{
			publicGTM.GET("/pools", client.GTMPoolList)
			publicGTM.GET("/pools/:pool", client.GTMPoolNameList)
			publicGTM.GET("/pools/:pool/members", client.GTMPoolMemberList)
			publicGTM.GET("/wideips", client.GTMWipList)
			publicGTM.GET("/wideips/:wideip", client.GTMWipNameList)
		}
		privateGTM := app.Group("/api/gtms/:trafficmanager")
		privateGTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, rootusers))
		{
			privateGTM.POST("/pools", client.GTMPoolPost)
			privateGTM.POST("/wideips", client.GTMWideipPost)
			privateGTM.DELETE("/pools", client.GTMPoolDelete)
			privateGTM.DELETE("/wideips", client.GTMWipDelete)
		}
	}
	if *ltmenabled {
		publicLTM := app.Group("/api/ltms/:lbpair")
		{
			publicLTM.GET("/pools", client.LTMPoolList)
			publicLTM.GET("/pools/:poolname", client.LTMPoolNameList)
			publicLTM.GET("/pools/:poolname/members", client.LTMPoolMemberList)
			publicLTM.GET("/devices", client.LTMDeviceList)
			publicLTM.GET("/devices/:devicename", client.LTMDeviceNameList)
			publicLTM.GET("/virtuals", client.LTMVirtualServerList)
			publicLTM.GET("/virtuals/:virtual", client.LTMVirtualServerNameList)
			publicLTM.GET("/virtuals/:virtual/rules", client.LTMFWRuleList)
			publicLTM.GET("/virtuals/:virtual/profiles", client.LTMProfileList)
			publicLTM.GET("/datagroups", client.LTMDataGroupList)
			publicLTM.GET("/datagroups/:datagroupname", client.LTMDataGroupNameList)
			publicLTM.GET("/blockips", client.LTMAddressList)
		}
		privateLTM := app.Group("/api/ltms/:lbpair")
		privateLTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, rootusers))
		{
			privateLTM.POST("/pools", client.LTMPoolPost)
			privateLTM.POST("/virtuals", client.LTMVirtualServerPost)
			privateLTM.POST("/pools/:poolname/members", client.LTMPoolMemberPost)
			privateLTM.POST("/datagroups", client.LTMDataGroupPost)
			privateLTM.PUT("/pools", client.LTMPoolPut)
			privateLTM.PUT("/pools/:poolname/members", client.LTMPoolMemberPut)
			privateLTM.DELETE("/pools", client.LTMPoolDelete)
			privateLTM.DELETE("/pools/:poolname/members", client.LTMPoolMemberDelete)
			privateLTM.DELETE("/datagroups", client.LTMDataGroupDelete)
			privateLTM.PUT("/datagroups/:datagroupname", client.LTMDataGroupItemPut)
			privateLTM.PATCH("/datagroups/:datagroupname", client.LTMDataGroupItemPatch)
			//To do: privateLTM.DELETE("/datagroups/:direction/:datagroupname", client.LTMDataGroupItemDelete)
		}
		emergencyLTM := app.Group("/api/ltms/:lbpair")
		emergencyLTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, emergencyusers))
		{
			emergencyLTM.PATCH("/blockips", client.LTMBlockIPPatch)
			emergencyLTM.PATCH("/whiteips", client.LTMWhiteIPPatch)
			emergencyLTM.DELETE("/blockips", client.LTMRemoveBlockIPPatch)
		}
	}
	switch {
	case *sslenabled:
		run := app.RunTLS(fmt.Sprintf(":%s", strconv.Itoa(*port)),
			conf.Security["certFile"],
			conf.Security["keyFile"])
		if run != nil {
			fmt.Println("Could not start web server,", run.Error())
		}
	default:
		run := app.Run(fmt.Sprintf(":%s", strconv.Itoa(*port)))
		if run != nil {
			fmt.Println("Could not start web server,", run.Error())
		}
	}
}

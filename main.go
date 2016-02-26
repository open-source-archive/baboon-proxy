package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/client"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"github.com/zalando-techmonkeys/gin-oauth2"
)

var (
	port       *int
	sslenabled *bool
	gtmenabled *bool
	ltmenabled *bool
	// BuildTime for Debugging
	BuildTime = "No BuildTime Provided"
	// GitHash for Debugging
	GitHash = "No GitHash Provided"
)

func usage() {
	fmt.Fprint(os.Stderr, fmt.Sprintf("Build Time: %s\nGit Commit Hash: %s\n\nUsage: ./baboon-proxy \n\t-port=80 \n\t-ssl-enabled=false \n\t-ltm-enabled=false \n\t-gtm-enabled=false \n\t-stderrthreshold=[INFO|WARN|FATAL] \n\t-log_dir=[string]\n\nExplanation:\n", BuildTime, GitHash))
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	port = flag.Int("port", 80, "Default Port")
	sslenabled = flag.Bool("ssl-enabled", false, "enable SSL")
	ltmenabled = flag.Bool("ltm-enabled", false, "enable LTM feature")
	gtmenabled = flag.Bool("gtm-enabled", false, "enable GTM feature")
	flag.Usage = usage
	flag.Parse()
}

func main() {
	app := gin.New()
	var conf *config.Config

	// Version
	version := client.Version{Build: BuildTime, Hash: GitHash}
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
	app.GET("/api/version", version.BaboonVersion)
	if *gtmenabled {

		publicGTM := app.Group("/api/gtms/:trafficmanager")
		{
			publicGTM.GET("/pools", client.GTMPoolList)
			publicGTM.GET("/pools/:pool", client.GTMPoolNameList)
			publicGTM.GET("/pools/:pool/members", client.GTMPoolMemberList)
			publicGTM.GET("/wideips", client.GTMWipList)
			publicGTM.GET("/wideips/:wideip", client.GTMWipNameList)
			publicGTM.GET("/irules", client.GTMIRuleList)
			publicGTM.GET("/irules/:irule", client.GTMIRuleNameList)
		}
		privateGTM := app.Group("/api/gtms/:trafficmanager")
		privateGTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, rootusers))
		{
			privateGTM.POST("/pools", client.GTMPoolPost)
			privateGTM.POST("/pools/:pool/members", client.GTMPoolMemberPost)
			privateGTM.POST("/wideips", client.GTMWideipPost)
			privateGTM.DELETE("/pools/:pool", client.GTMPoolDelete)
			privateGTM.DELETE("/wideips/:wideip", client.GTMWipDelete)
			privateGTM.DELETE("/pools/:pool/members", client.GTMPoolMemberDelete)
		}
		emergencyGTM := app.Group("/api/gtms/:trafficmanager")
		emergencyGTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, emergencyusers))
		{
			emergencyGTM.PUT("/pools/:pool/members", client.GTMPoolMemberStatusPut)
			emergencyGTM.PUT("/pools", client.GTMPoolStatusPut)
		}
	}
	if *ltmenabled {
		publicLTM := app.Group("/api/ltms/:lbpair")
		{
			publicLTM.GET("/pools", client.LTMPoolList)
			publicLTM.GET("/pools/:pool", client.LTMPoolNameList)
			publicLTM.GET("/pools/:pool/members", client.LTMPoolMemberList)
			publicLTM.GET("/devices", client.LTMDeviceList)
			publicLTM.GET("/devices/:devicename", client.LTMDeviceNameList)
			publicLTM.GET("/virtuals", client.LTMVirtualServerList)
			publicLTM.GET("/virtuals/:virtual", client.LTMVirtualServerNameList)
			publicLTM.GET("/virtuals/:virtual/rules", client.LTMFWRuleList)
			publicLTM.GET("/virtuals/:virtual/profiles", client.LTMProfileList)
			publicLTM.GET("/datagroups", client.LTMDataGroupList)
			publicLTM.GET("/datagroups/:datagroupname", client.LTMDataGroupNameList)
			publicLTM.GET("/blacklist", client.LTMBlackAddressList)
			publicLTM.GET("/whitelist", client.LTMWhiteAddressList)
			publicLTM.GET("/irules", client.LTMIRuleList)
			publicLTM.GET("/irules/:irule", client.LTMIRuleNameList)
		}
		privateLTM := app.Group("/api/ltms/:lbpair")
		privateLTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, rootusers))
		{
			privateLTM.POST("/pools", client.LTMPoolPost)
			privateLTM.POST("/virtuals", client.LTMVirtualServerPost)
			privateLTM.POST("/pools/:pool/members", client.LTMPoolMemberPost)
			privateLTM.POST("/datagroups", client.LTMDataGroupPost)
			privateLTM.PUT("/pools", client.LTMPoolPut)
			privateLTM.PUT("/pools/:pool/members", client.LTMPoolMemberPut)
			privateLTM.DELETE("/pools/:pool", client.LTMPoolDelete)
			privateLTM.DELETE("/pools/:pool/members", client.LTMPoolMemberDelete)
			privateLTM.DELETE("/datagroups", client.LTMDataGroupDelete)
			privateLTM.PUT("/datagroups/:datagroupname", client.LTMDataGroupItemPut)
			privateLTM.PATCH("/datagroups/:datagroupname", client.LTMDataGroupItemPatch)
			//To do: privateLTM.DELETE("/datagroups/:direction/:datagroupname", client.LTMDataGroupItemDelete)
		}
		emergencyLTM := app.Group("/api/ltms/:lbpair")
		emergencyLTM.Use(ginoauth2.Auth(ginoauth2.UidCheck, *OAuth2Endpoint, emergencyusers))
		{
			emergencyLTM.PATCH("/blacklist", client.LTMBlockIPPatch)
			emergencyLTM.PATCH("/whitelist", client.LTMWhiteIPPatch)
			emergencyLTM.DELETE("/whitelist", client.LTMRemoveWhiteIPPatch)
			emergencyLTM.DELETE("/blacklist", client.LTMRemoveBlockIPPatch)
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

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
	"github.com/zalando-techmonkeys/baboon-proxy/gtm"
	"github.com/zalando-techmonkeys/baboon-proxy/ltm"
	"github.com/zalando-techmonkeys/baboon-proxy/util"
)

// Version struct provides information of current build
type Version struct {
	Build string `json:"buildTime"`
	Hash  string `json:"gitHash"`
}

// ErrorResponse provides detailed information regarding errors
// including a link to baboon documentation
type ErrorResponse struct {
	Type   string `json:"type"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

var (
	returnerror errors.Error
	conf        = config.LoadConfig()
)

// BaboonVersion shows current githash and build time for debugging
func (v Version) BaboonVersion(c *gin.Context) {
	respondWithStatus(http.StatusOK, "", v, "", "", c)
}

// GTMWipDelete delete wide ip
func GTMWipDelete(c *gin.Context) {
	wideip := c.Params.ByName("wideip")
	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, wideip, nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	res, err := gtm.DeleteGTMWip(f5url, wideip)
	if err != nil {
		respondWithStatus(err.Status, wideip, nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "WideIP deleted", wideip,
		returnerror.ErrorMessage(), conf.Documentation["gtmwideipdocumentationuri"], c)
}

// GTMWipList show all wide ips
func GTMWipList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	res, gtmwiplist, err := gtm.ShowGTMWips(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	poolsURI := util.ReplaceGTMWipUritoGTMPoolURI(c.Request.RequestURI)
	for _, wip := range gtmwiplist.Items {
		for i, pools := range wip.Pools {
			u := new(url.URL)
			u.Scheme = common.Protocol
			u.Path = path.Join(c.Request.Host, poolsURI, pools.Name)
			wip.Pools[i].PoolsReference = u.String()
		}
	}
	respondWithStatus(res.Status, "", gtmwiplist, returnerror.ErrorMessage(), conf.Documentation["gtmwideipdocumentationuri"], c)
}

// GTMWipNameList show a specific wide ip
func GTMWipNameList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	wideip := c.Params.ByName("wideip")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, wideip, nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	res, gtmwipnamelist, err := gtm.ShowGTMWip(f5url, wideip)
	if err != nil {
		respondWithStatus(err.Status, wideip, nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	for i, pool := range gtmwipnamelist.Pools {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, "/api/gtms", tm, "/pools/", pool.Name)
		gtmwipnamelist.Pools[i].PoolsReference = u.String()
	}
	respondWithStatus(res.Status, "", gtmwipnamelist, returnerror.ErrorMessage(), conf.Documentation["gtmwideipdocumentationuri"], c)
}

// LTMPoolList show local traffic manager pools
func LTMPoolList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	f5url, err := ltm.Loadbalancer(lbpair, conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	res, poollist, err := ltm.ShowLTMPools(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumenationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	for i, v := range poollist.Items {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", v.Name, common.MembersURI)
		poollist.Items[i].MembersReference = u.String()
	}
	respondWithStatus(res.Status, "", poollist, returnerror.ErrorMessage(), conf.Documentation["ltmpooldocumentationuri"], c)
}

// GTMPoolList show global traffic manager pools
func GTMPoolList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	res, gtmpoollist, err := gtm.ShowGTMPools(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	for i, v := range gtmpoollist.Items {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", v.Name, common.MembersURI)
		gtmpoollist.Items[i].MembersReference = u.String()
	}
	respondWithStatus(res.Status, "", gtmpoollist, returnerror.ErrorMessage(), conf.Documentation["gtmpooldocumentationuri"], c)
}

// GTMIRuleList show global traffic manager iRules
func GTMIRuleList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmiruledocumentationuri"], c)
		return
	}
	res, gtmirulelist, err := gtm.ShowGTMIRules(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmiruledocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", gtmirulelist, returnerror.ErrorMessage(), conf.Documentation["gtmiruledocumentationuri"], c)
}

// GTMIRuleNameList show specific global traffic manager iRule
func GTMIRuleNameList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	irule := c.Params.ByName("irule")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, irule, nil, err.Message, conf.Documentation["gtmiruledocumentationuri"], c)
		return
	}
	res, gtmirulenamelist, err := gtm.ShowGTMIRule(f5url, irule)
	if err != nil {
		respondWithStatus(err.Status, irule, nil, err.Message, conf.Documentation["gtmiruledocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", gtmirulenamelist, returnerror.ErrorMessage(), conf.Documentation["gtmiruledocumentationuri"], c)
}

// LTMPoolNameList show specific local traffic manager pool
func LTMPoolNameList(c *gin.Context) {
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	pool := c.Params.ByName("pool")
	res, poollist, err := ltm.ShowLTMPool(f5url, pool)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	u := new(url.URL)
	u.Scheme = common.Protocol
	u.Path = path.Join(c.Request.Host, c.Request.RequestURI, common.MembersURI)
	poollist.MembersReference = u.String()
	respondWithStatus(res.Status, "", poollist, returnerror.ErrorMessage(), conf.Documentation["ltmpooldocumentationuri"], c)
}

// LTMIRuleNameList show specific iRule
func LTMIRuleNameList(c *gin.Context) {
	irule := c.Params.ByName("irule")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, irule, nil, err.Message, conf.Documentation["ltmiruledocumentationuri"], c)
		return
	}
	res, irulenamelist, err := ltm.ShowLTMIRule(f5url, irule)
	if err != nil {
		respondWithStatus(err.Status, irule, nil, err.Message, conf.Documentation["ltmiruledocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", irulenamelist, returnerror.ErrorMessage(), conf.Documentation["ltmiruledocumentationuri"], c)
}

// LTMIRuleList show all iRules
func LTMIRuleList(c *gin.Context) {
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmiruledocumentationuri"], c)
		return
	}
	res, irulelist, err := ltm.ShowLTMIRules(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmiruledocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", irulelist, returnerror.ErrorMessage(), conf.Documentation["ltmiruledocumentationuri"], c)
}

// GTMPoolNameList show specific global traffic manager pool
func GTMPoolNameList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	res, gtmpoollist, err := gtm.ShowGTMPool(f5url, pool)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	u := new(url.URL)
	u.Scheme = common.Protocol
	u.Path = path.Join(c.Request.Host, c.Request.RequestURI, common.MembersURI)
	gtmpoollist.MembersReference = u.String()
	respondWithStatus(res.Status, "", gtmpoollist, returnerror.ErrorMessage(), conf.Documentation["gtmpooldocumentationuri"], c)
}

// GTMPoolMemberList show global traffic manager members in a specific pool
func GTMPoolMemberList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	res, poolmemberlist, err := gtm.ShowGTMPoolMembers(f5url, pool)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", poolmemberlist, returnerror.ErrorMessage(), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
}

// LTMPoolMemberList show local traffic manager members in a specific pool
func LTMPoolMemberList(c *gin.Context) {
	pool := c.Params.ByName("pool")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	res, poolmemberlist, err := ltm.ShowLTMPoolMember(f5url, pool)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", poolmemberlist, returnerror.ErrorMessage(), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMDeviceList show local traffic manager devices
func LTMDeviceList(c *gin.Context) {
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdevicedocumentationuri"], c)
		return
	}
	res, devicelist, err := ltm.ShowLTMDevice(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdevicedocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", devicelist, returnerror.ErrorMessage(), conf.Documentation["ltmdevicedocumentationuri"], c)
}

// LTMDeviceNameList show local traffic manager specific device
func LTMDeviceNameList(c *gin.Context) {
	device := c.Params.ByName("devicename")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, device, nil, err.Message, conf.Documentation["ltmdevicedocumentationuri"], c)
		return
	}
	res, devicenamelist, err := ltm.ShowLTMDeviceName(device, f5url, conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, device, nil, err.Message, conf.Documentation["ltmdevicedocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", devicenamelist, returnerror.ErrorMessage(), conf.Documentation["ltmdevicedocumentationuri"], c)
}

// LTMVirtualServerList show local traffic manager virtual servers
func LTMVirtualServerList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	f5url, err := ltm.Loadbalancer(lbpair, conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	res, virtualserverlist, err := ltm.ShowLTMVirtualServer(f5url)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	for i, v := range virtualserverlist.Items {
		u1 := new(url.URL)
		u1.Scheme = common.Protocol
		u1.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", util.ReplaceCommon(v.Name), common.ProfilesURI)
		u2 := new(url.URL)
		u2.Scheme = common.Protocol
		u2.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", util.ReplaceCommon(v.Name), common.FwURI)
		virtualserverlist.Items[i].ProfilesReference = u1.String()
		virtualserverlist.Items[i].FwRulesReference = u2.String()
		if len(v.Pool) > 0 {
			u := new(url.URL)
			u.Scheme = common.Protocol
			u.Path = path.Join(c.Request.Host, "/api/ltms/", lbpair, common.PoolsURI, util.ReplaceCommon(v.Pool))
			virtualserverlist.Items[i].PoolsReference = u.String()
		}
	}
	respondWithStatus(res.Status, "", virtualserverlist, returnerror.ErrorMessage(), conf.Documentation["ltmvirtualdocumentationuri"], c)
}

// LTMVirtualServerNameList show local traffic manager specific virtual server
func LTMVirtualServerNameList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	vservername := c.Params.ByName("virtual")
	f5url, err := ltm.Loadbalancer(lbpair, conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, vservername, nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	res, virtualservernamelist, err := ltm.ShowLTMVirtualServerName(f5url, vservername)
	if err != nil {
		respondWithStatus(err.Status, vservername, nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	u1 := new(url.URL)
	u1.Scheme = common.Protocol
	u1.Path = path.Join(c.Request.Host, c.Request.RequestURI, common.ProfilesURI)
	u2 := new(url.URL)
	u2.Scheme = common.Protocol
	u2.Path = path.Join(c.Request.Host, c.Request.RequestURI, common.FwURI)
	virtualservernamelist.ProfilesReference = u1.String()
	virtualservernamelist.FwRulesReference = u2.String()
	if len(virtualservernamelist.Pool) > 0 {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, "/api/ltms/", lbpair, common.PoolsURI, util.ReplaceCommon(virtualservernamelist.Pool))
		virtualservernamelist.PoolsReference = u.String()
	}
	respondWithStatus(res.Status, "", virtualservernamelist, returnerror.ErrorMessage(), conf.Documentation["ltmvirtualdocumentationuri"], c)
}

// LTMProfileList show local traffic manager profiles of a specific virtual server
func LTMProfileList(c *gin.Context) {
	vservername := c.Params.ByName("virtual")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, vservername, nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	res, profilelist, err := ltm.ShowLTMProfile(f5url, vservername)
	if err != nil {
		respondWithStatus(err.Status, vservername, nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", profilelist, returnerror.ErrorMessage(), conf.Documentation["ltmvirtualdocumentationuri"], c)
}

// LTMFWRuleList show local traffic manager iRules of a specific virtual server
func LTMFWRuleList(c *gin.Context) {
	vservername := c.Params.ByName("virtual")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, vservername, nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	res, rulelist, err := ltm.ShowLTMFWRules(f5url, vservername)
	if err != nil {
		respondWithStatus(err.Status, vservername, nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", rulelist, returnerror.ErrorMessage(), conf.Documentation["ltmvirtualdocumentationuri"], c)
}

// LTMDataGroupList show local traffic manager internal data groups
func LTMDataGroupList(c *gin.Context) {
	direction := common.InternalDataGroup
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	res, datagrouplist, err := ltm.ShowLTMDataGroup(f5url, direction)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", datagrouplist, returnerror.ErrorMessage(), conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupNameList show local traffic manager internal specific data group
func LTMDataGroupNameList(c *gin.Context) {
	direction := common.InternalDataGroup
	datagroupname := c.Params.ByName("datagroupname")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, datagroupname, nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	res, datagroupnamelist, err := ltm.ShowLTMDataGroupName(f5url, direction, datagroupname)
	if err != nil {
		respondWithStatus(err.Status, datagroupname, nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", datagroupnamelist, returnerror.ErrorMessage(), conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMVirtualServerPost create virtual server
func LTMVirtualServerPost(c *gin.Context) {
	var vservercreate ltm.CreateVirtualServer
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}

	if err := c.BindJSON(&vservercreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create virtual server",
			fmt.Sprintf("%s", err), conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	res, err := ltm.PostLTMVirtualServer(f5url, &vservercreate)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmvirtualdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Virtual server added", vservercreate.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmvirtualdocumentationuri"], c)
}

// LTMSSLKeyPost install a new ssl key on local traffic manager
func LTMSSLKeyPost(c *gin.Context) {
	var sslkeycreate ltm.CreateSSLKey
	sslkeycreate.Command = "install"
	//lbpair := c.Params.ByName("lbpair")
	//f5url := DeviceActive(lbpair)
	c.BindJSON(&sslkeycreate)
	res, _ := ltm.PostLTMSSLKey(common.CryptoURL, &sslkeycreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	switch res.Status {
	case http.StatusOK:
		{
			u := new(url.URL)
			u.Scheme = common.Protocol
			u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", sslkeycreate.Name)
			c.Set("message", "SSL key installed "+sslkeycreate.Name)
			c.Header("location", u.String())
			c.JSON(http.StatusCreated, gin.H{"message": "SSL key installed " + sslkeycreate.Name})
		}
	default:
		c.Set("message", returnerror.ErrorMessage())
		c.Header("Content-Type", "application/problem+json")
		c.JSON(res.Status, gin.H{"status": res.Status, "title": returnerror.ErrorMessage()})
	}
}

// LTMPoolPost create a new pool with members and a monitoring check on a local traffic manager
func LTMPoolPost(c *gin.Context) {
	var poolcreate ltm.CreatePool

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	if err := c.BindJSON(&poolcreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create pool",
			fmt.Sprintf("%s", err), conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	res, err := ltm.PostLTMPool(f5url, &poolcreate)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Pool created", poolcreate.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmpooldocumentationuri"], c)
}

// GTMPoolPost create a new wide IP pool with members and a monitoring check on a global traffic manager
func GTMPoolPost(c *gin.Context) {
	var poolcreate gtm.CreatePool

	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolcreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create pool",
			fmt.Sprintf("%s", err), conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	res, err := gtm.PostGTMPool(f5url, &poolcreate)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Pool created", poolcreate.Name,
		returnerror.ErrorMessage(), conf.Documentation["gtmpooldocumentationuri"], c)
}

// GTMPoolMemberPost adds additional LTM virtual server on a global traffic manager pool
func GTMPoolMemberPost(c *gin.Context) {
	var poolmember gtm.CreatePoolMember
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolmember); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Add pool member",
			fmt.Sprintf("%s", err), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	res, err := gtm.PostGTMPoolMember(f5url, pool, &poolmember)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Poolmember added", poolmember.Name,
		returnerror.ErrorMessage(), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
}

// GTMWideipPost create new wide IP on a global traffic manager
func GTMWideipPost(c *gin.Context) {
	var wideipcreate gtm.CreateWip

	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&wideipcreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create wideip",
			fmt.Sprintf("%s", err), conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	res, err := gtm.PostGTMWip(f5url, &wideipcreate)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["gtmwideipdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "WideIP created", wideipcreate.Name,
		returnerror.ErrorMessage(), conf.Documentation["gtmwideipdocumentationuri"], c)
}

// LTMPoolMemberPost add new members to a specific pool on a local traffic manager
func LTMPoolMemberPost(c *gin.Context) {
	var poolmembercreate ltm.CreatePoolMember

	pool := c.Params.ByName("pool")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}

	if err := c.BindJSON(&poolmembercreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create pool member",
			fmt.Sprintf("%s", err), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	res, err := ltm.PostLTMPoolMember(f5url, pool, &poolmembercreate)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)

	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Poolmember added", poolmembercreate.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMPoolPut modify pool (old ones will be deleted) or change monitoring
func LTMPoolPut(c *gin.Context) {
	var poolmodify ltm.ModifyPool

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}

	if err := c.BindJSON(&poolmodify); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Modify pool",
			fmt.Sprintf("%s", err), conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	res, err := ltm.PutLTMPool(f5url, poolmodify.Name, &poolmodify)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)

	respondWithStatus(res.Status, "Pool modified", poolmodify.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmpooldocumentationuri"], c)
}

// LTMPoolDelete delete a pool on a local traffic manager
func LTMPoolDelete(c *gin.Context) {
	pool := c.Params.ByName("pool")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	res, err := ltm.DeleteLTMPool(f5url, pool)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Pool deleted", pool,
		returnerror.ErrorMessage(), conf.Documentation["ltmpooldocumentationuri"], c)
}

// GTMPoolDelete delete a pool on a global traffic manager
func GTMPoolDelete(c *gin.Context) {
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	res, err := gtm.DeleteGTMPool(f5url, pool)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Pool deleted", pool,
		returnerror.ErrorMessage(), conf.Documentation["gtmpooldocumentationuri"], c)
}

// GTMPoolMemberDelete delete specific pool members on a global traffic manager
func GTMPoolMemberDelete(c *gin.Context) {
	var poolmemberdelete gtm.RemovePoolMember
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolmemberdelete); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Delete pool member",
			fmt.Sprintf("%s", err), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	res, err := gtm.DeleteGTMPoolMember(f5url, pool, &poolmemberdelete)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Poolmember deleted", poolmemberdelete.Name,
		returnerror.ErrorMessage(), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
}

// GTMPoolMemberStatusPut modify pool member status on a global traffic manager (enabled, disabled)
func GTMPoolMemberStatusPut(c *gin.Context) {
	var poolmemberstatus gtm.ModifyPoolMemberStatus
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolmemberstatus); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Modify pool member status",
			fmt.Sprintf("%s", err), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	res, err := gtm.PutGTMPoolMemberStatus(f5url, pool, &poolmemberstatus)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["gtmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Poolmember modified", poolmemberstatus.Name,
		returnerror.ErrorMessage(), conf.Documentation["gtmpoolmemberdocumentationuri"], c)
}

// GTMPoolStatusPut modify pool member status on a global traffic manager (enabled, disabled)
func GTMPoolStatusPut(c *gin.Context) {
	var poolstatus gtm.ModifyPoolStatus
	f5url, err := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	if err != nil {
		respondWithStatus(err.Status, poolstatus.Name, nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolstatus); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Modify pool status",
			fmt.Sprintf("%s", err), conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	res, err := gtm.PutGTMPoolStatus(f5url, &poolstatus)
	if err != nil {
		respondWithStatus(err.Status, poolstatus.Name, nil, err.Message, conf.Documentation["gtmpooldocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Pool modified", poolstatus.Name,
		returnerror.ErrorMessage(), conf.Documentation["gtmpooldocumentationuri"], c)
}

// LTMPoolMemberPut modify pool members on a local traffic manager (enabled, disabled, force-offline)
func LTMPoolMemberPut(c *gin.Context) {
	var poolmembermodify ltm.ModifyPoolMember
	pool := c.Params.ByName("pool")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolmembermodify); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Modify pool member",
			fmt.Sprintf("%s", err), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	res, err := ltm.PutLTMPoolMember(f5url, pool, poolmembermodify.Name, poolmembermodify.Status)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Poolmember modified", poolmembermodify.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMPoolMemberDelete delete specific pool members on a local traffic manager
func LTMPoolMemberDelete(c *gin.Context) {
	var poolmemberdelete ltm.RemovePoolMember
	pool := c.Params.ByName("pool")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&poolmemberdelete); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Delete pool member",
			fmt.Sprintf("%s", err), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	res, err := ltm.DeleteLTMPoolMember(f5url, pool, poolmemberdelete.Name)
	if err != nil {
		respondWithStatus(err.Status, pool, nil, err.Message, conf.Documentation["ltmpoolmemberdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Poolmember deleted", poolmemberdelete.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMDataGroupPost add new internal datagroup on a local traffic manager
func LTMDataGroupPost(c *gin.Context) {
	var datagroupcreate ltm.CreateDataGroup
	direction := common.InternalDataGroup
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&datagroupcreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create datagroup item",
			fmt.Sprintf("%s", err), conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	res, err := ltm.PostLTMDataGroup(f5url, direction, &datagroupcreate)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Datagroup added", datagroupcreate.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupDelete delete a internal datagroup on a local traffic manager
func LTMDataGroupDelete(c *gin.Context) {
	var datagroupdelete ltm.RemoveDataGroup
	direction := common.InternalDataGroup
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	if err := c.BindJSON(&datagroupdelete); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Delete datagroup item",
			fmt.Sprintf("%s", err), conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	res, err := ltm.DeleteLTMDataGroup(f5url, direction, datagroupdelete.Name)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "Datagroup deleted", datagroupdelete.Name,
		returnerror.ErrorMessage(), conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupItemPut remove all existing records in a datagroup and add new records (ip or string)
// on a local traffic manager
func LTMDataGroupItemPut(c *gin.Context) {
	var datagroupitemcreate ltm.CreateDataGroupItem
	direction := common.InternalDataGroup
	datagroupname := c.Params.ByName("datagroupname")

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, datagroupname, nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}

	if err := c.BindJSON(&datagroupitemcreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create datagroup item",
			fmt.Sprintf("%s", err), conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	res, err := ltm.PatchLTMDataGroupItem(f5url, direction, datagroupname, &datagroupitemcreate)
	if err != nil {
		respondWithStatus(err.Status, datagroupname, nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Datagroup item added in", datagroupname,
		returnerror.ErrorMessage(), conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupItemPatch add an item to a datagroup (ip or string) on a local traffic manager
func LTMDataGroupItemPatch(c *gin.Context) {
	var datagroupitemcreate ltm.CreateDataGroupItem
	direction := common.InternalDataGroup
	datagroupname := c.Params.ByName("datagroupname")

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, datagroupname, nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}

	if err := c.BindJSON(&datagroupitemcreate); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Create datagroup item",
			fmt.Sprintf("%s", err), conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	res, err := ltm.PatchLTMDataGroupItem(f5url, direction, datagroupname, &datagroupitemcreate)
	if err != nil {
		respondWithStatus(err.Status, datagroupname, nil, err.Message, conf.Documentation["ltmdatagroupdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "Datagroup item added in", datagroupname,
		returnerror.ErrorMessage(), conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMBlackAddressList show local traffic blacklisted ip addresses
func LTMBlackAddressList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	f5url, err := ltm.Loadbalancer(lbpair, conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	f5url = util.ReplaceLTMUritoAddressListURI(f5url)
	res, addresslist, err := ltm.ShowLTMAddressList(f5url, common.BlackList)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", addresslist, returnerror.ErrorMessage(), conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMWhiteAddressList show local traffic whitelisted ip addresses
func LTMWhiteAddressList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	f5url, err := ltm.Loadbalancer(lbpair, conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	f5url = util.ReplaceLTMUritoAddressListURI(f5url)
	res, addresslist, err := ltm.ShowLTMAddressList(f5url, common.WhiteList)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "", addresslist, returnerror.ErrorMessage(), conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMBlockIPPatch add ips which will be blocked
func LTMBlockIPPatch(c *gin.Context) {
	var blockips ltm.CreateAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	f5url = util.ReplaceLTMUritoAddressListURI(f5url)
	if err := c.BindJSON(&blockips); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Block IPs",
			fmt.Sprintf("%s", err), conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	res, err := ltm.PatchLTMBlockAddresses(f5url, &blockips)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "IP(s) blocked successfully", blockips,
		returnerror.ErrorMessage(), conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMWhiteIPPatch add ips which will be whitelisted
func LTMWhiteIPPatch(c *gin.Context) {
	var whiteips ltm.CreateAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	f5url = util.ReplaceLTMUritoAddressListURI(f5url)
	if err := c.BindJSON(&whiteips); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "White IPs",
			fmt.Sprintf("%s", err), conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	res, err := ltm.PatchLTMWhiteAddresses(f5url, &whiteips)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == http.StatusOK {
		res.Status = http.StatusCreated
	}
	respondWithStatus(res.Status, "IP(s) whitelisted successfully", whiteips,
		returnerror.ErrorMessage(), conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMRemoveBlockIPPatch remove ips which are currently blocked
func LTMRemoveBlockIPPatch(c *gin.Context) {
	var unblockips ltm.DeleteAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	f5url = util.ReplaceLTMUritoAddressListURI(f5url)
	if err := c.BindJSON(&unblockips); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Unblock IPs",
			fmt.Sprintf("%s", err), conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	res, err := ltm.DeleteLTMBlockAddresses(f5url, &unblockips, common.BlackList)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "IP(s) removed successfully from blacklist", fmt.Sprintf("%+v", unblockips),
		returnerror.ErrorMessage(), conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMRemoveWhiteIPPatch remove ips which are currently whitelisted
func LTMRemoveWhiteIPPatch(c *gin.Context) {
	var removewhiteips ltm.DeleteAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), conf.Ltmdevicenames)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	f5url = util.ReplaceLTMUritoAddressListURI(f5url)
	if err := c.BindJSON(&removewhiteips); err != nil {
		respondWithStatus(http.StatusBadRequest, "Invalid JSON data", "Unblock IPs",
			fmt.Sprintf("%s", err), conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	res, err := ltm.DeleteLTMBlockAddresses(f5url, &removewhiteips, common.WhiteList)
	if err != nil {
		respondWithStatus(err.Status, "", nil, err.Message, conf.Documentation["ltmaddresslistdocumentationuri"], c)
		return
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	respondWithStatus(res.Status, "IP(s) removed successfully from whitelist", fmt.Sprintf("%+v", removewhiteips),
		returnerror.ErrorMessage(), conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LoggerMiddleware log user activity
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set time to measure request average
		t := time.Now()
		source := c.ClientIP()
		// Request
		c.Next()

		// Measure how long the request takes
		latency := time.Since(t)

		// Get post message, dont use c.MustGet() otherwise if empty it would panic
		value, ok := c.Get("message")
		if !ok {
			value = ""
		}
		message := value.(string)

		// Important things to log
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		switch {
		case status >= 400 && status <= 499:
			{
				glog.Warningf("[GIN] | %3d | %12v | %s | %s | %s | %s",
					status, latency, source, method, path, message)
			}
		case status == 500:
			{
				glog.Errorf("[GIN] | %3d | %12v | %s | %s | %s | %s",
					status, latency, source, method, path, message)
			}
		case status == 200:
			{
				glog.Infof("[GIN] | %3d | %12v | %s | %s | %s",
					status, latency, source, method, path)
			}
		default:
			glog.Infof("[GIN] | %3d | %12v | %s | %s | %s | %s",
				status, latency, source, method, path, message)
		}
	}
}

// CORSMiddleware handle Access-Control-Allow-Origin and allowed methods
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}

// respondWithStatus return JSON response
func respondWithStatus(code int, message, name interface{}, e, documentation string, c *gin.Context) {
	switch code {
	case http.StatusOK:
		{
			if message != "" {
				c.Set("message", fmt.Sprintf("%s %+v", message, name))
				c.JSON(code, gin.H{"message": fmt.Sprintf("%s %+v", message, name)})
			} else {
				c.Set("message", message)
				c.JSON(code, gin.H{"message": name})
			}
		}
	case http.StatusCreated:
		c.Set("message", fmt.Sprintf("%s %+v", message, name))
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", fmt.Sprintf("%#v", name))
		c.Header("location", u.String())
		c.JSON(code, gin.H{"message": fmt.Sprintf("%s %+v", message, name)})
	case http.StatusBadRequest:
		c.Set("message", e)
		c.JSON(code, ErrorResponse{Type: documentation, Status: code,
			Title: "Bad Request", Detail: e})
	case http.StatusNotFound:
		c.Set("message", e)
		c.JSON(code, ErrorResponse{Type: documentation, Status: code,
			Title: "Not Found", Detail: e})
	case http.StatusConflict:
		c.Set("message", e)
		c.Header("Content-Type", "application/problem+json")
		c.JSON(code, ErrorResponse{Type: documentation, Status: code,
			Title: "Conflict", Detail: e})
	default:
		c.Set("message", e)
		c.Header("Content-Type", "application/problem+json")
		c.JSON(code, ErrorResponse{Type: documentation, Status: code,
			Title: e, Detail: e})
	}
}

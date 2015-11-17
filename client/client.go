package client

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/gtm"
	"github.com/zalando-techmonkeys/baboon-proxy/ltm"
	"github.com/zalando-techmonkeys/baboon-proxy/util"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

var returnerror ltm.ErrorLTM

// GTMWipDelete delete wide ip
func GTMWipDelete(c *gin.Context) {
	var wipdelete gtm.RemoveWip

	f5url, _ := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	c.Bind(&wipdelete)
	res, err := gtm.DeleteGTMWip(f5url, wipdelete.Name)
	if err != nil {
		glog.Errorf("%s", err)
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Wide IP deleted", wipdelete.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["gtmwideipdocumentationuri"], c)
}

// GTMWipList show all wide ips
func GTMWipList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	}
	gtmwiplist, err := gtm.ShowGTMWips(f5url)
	if err != nil {
		glog.Errorf("%s", err)
	}
	poolsURI := util.ReplaceGTMWipUritoGTMPoolURI(c.Request.RequestURI)
	for _, wip := range gtmwiplist.Items {
		for i, pools := range wip.Pools {
			u := new(url.URL)
			u.Scheme = common.Protocol
			u.Path = path.Join(c.Request.Host, poolsURI, pools.Name)
			wip.Pools[i].PoolsReference = u.String()
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": gtmwiplist})
}

// GTMWipNameList show a specific wide ip
func GTMWipNameList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	wideip := c.Params.ByName("wideip")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	}
	gtmwipnamelist, err := gtm.ShowGTMWip(f5url, wideip)
	if err != nil {
		glog.Errorf("%s", err)
	}
	poolsURI := util.ReplaceGTMWipUritoGTMPoolURI(c.Request.RequestURI)
	for i, pool := range gtmwipnamelist.Pools {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, poolsURI, pool.Name)
		gtmwipnamelist.Pools[i].PoolsReference = u.String()
	}
	c.JSON(http.StatusOK, gin.H{"message": gtmwipnamelist})
}

// LTMPoolList show local traffic manager pools
func LTMPoolList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	glog.Infof("%v", common.Conf.Ltmdevicenames)
	f5url, err := ltm.Loadbalancer(lbpair, common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	poollist, err := ltm.ShowLTMPools(f5url)
	if err != nil {
		glog.Errorf("%s", err)
	}
	for i, v := range poollist.Items {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", v.Name, common.MembersURI)
		poollist.Items[i].MembersReference = u.String()
	}
	c.JSON(http.StatusOK, gin.H{"message": poollist})
}

// GTMPoolList show global traffic manager pools
func GTMPoolList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	}
	gtmpoollist, err := gtm.ShowGTMPools(f5url)
	if err != nil {
		glog.Errorf("%s", err)
	}
	for i, v := range gtmpoollist.Items {
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", v.Name, common.MembersURI)
		gtmpoollist.Items[i].MembersReference = u.String()
	}
	c.JSON(http.StatusOK, gin.H{"message": gtmpoollist})
}

// LTMPoolNameList show specific local traffic manager pool
func LTMPoolNameList(c *gin.Context) {
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	poolname := c.Params.ByName("poolname")
	poolnamelist, err := ltm.ShowLTMPool(f5url, poolname)
	if err != nil {
		glog.Errorf("%s", err)
	}
	u := new(url.URL)
	u.Scheme = common.Protocol
	u.Path = path.Join(c.Request.Host, c.Request.RequestURI, common.MembersURI)
	poolnamelist.MembersReference = u.String()
	c.JSON(http.StatusOK, gin.H{"message": poolnamelist})
}

// GTMPoolNameList show specific global traffic manager pool
func GTMPoolNameList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	}
	gtmpoolnamelist, err := gtm.ShowGTMPool(f5url, pool)
	if err != nil {
		glog.Errorf("%s", err)
	}
	u := new(url.URL)
	u.Scheme = common.Protocol
	u.Path = path.Join(c.Request.Host, c.Request.RequestURI, common.MembersURI)
	gtmpoolnamelist.MembersReference = u.String()
	c.JSON(http.StatusOK, gin.H{"message": gtmpoolnamelist})
}

// GTMPoolMemberList show global traffic manager members in a specific pool
func GTMPoolMemberList(c *gin.Context) {
	tm := c.Params.ByName("trafficmanager")
	pool := c.Params.ByName("pool")
	f5url, err := gtm.Trafficmanager(tm)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	}
	poolmemberlist, err := gtm.ShowGTMPoolMembers(f5url, pool)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": poolmemberlist})
}

// LTMPoolMemberList show local traffic manager members in a specific pool
func LTMPoolMemberList(c *gin.Context) {
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	poolname := c.Params.ByName("poolname")
	poolmemberlist, err := ltm.ShowLTMPoolMember(f5url, poolname)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": poolmemberlist})
}

// LTMDeviceList show local traffic manager devices
func LTMDeviceList(c *gin.Context) {
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	devicelist, err := ltm.ShowLTMDevice(f5url)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": devicelist})
}

// LTMDeviceNameList show local traffic manager specific device
func LTMDeviceNameList(c *gin.Context) {
	device := c.Params.ByName("devicename")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	devicenamelist, err := ltm.ShowLTMDeviceName(device, f5url, common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": devicenamelist})
}

// LTMVirtualServerList show local traffic manager virtual servers
func LTMVirtualServerList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	f5url, err := ltm.Loadbalancer(lbpair, common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	virtualserverlist, err := ltm.ShowLTMVirtualServer(f5url)
	if err != nil {
		glog.Errorf("%s", err)
	}
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
	c.JSON(http.StatusOK, gin.H{"message": virtualserverlist})
}

// LTMVirtualServerNameList show local traffic manager specific virtual server
func LTMVirtualServerNameList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	vservername := c.Params.ByName("virtual")
	f5url, err := ltm.Loadbalancer(lbpair, common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	virtualservernamelist, err := ltm.ShowLTMVirtualServerName(f5url, vservername)
	if err != nil {
		glog.Errorf("%s", err)
	}
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
	c.JSON(http.StatusOK, gin.H{"message": virtualservernamelist})
}

// LTMProfileList show local traffic manager profiles of a specific virtual server
func LTMProfileList(c *gin.Context) {
	vservername := c.Params.ByName("virtual")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	profilelist, err := ltm.ShowLTMProfile(f5url, vservername)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": profilelist})
}

// LTMFWRuleList show local traffic manager iRules of a specific virtual server
func LTMFWRuleList(c *gin.Context) {
	vservername := c.Params.ByName("virtual")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	rulelist := ltm.ShowLTMFWRules(f5url, vservername)
	c.JSON(http.StatusOK, gin.H{"message": rulelist})
}

// LTMDataGroupList show local traffic manager internal data groups
func LTMDataGroupList(c *gin.Context) {
	direction := common.InternalDataGroup
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	datagrouplist, err := ltm.ShowLTMDataGroup(f5url, direction)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": datagrouplist})
}

// LTMDataGroupNameList show local traffic manager internal specific data group
func LTMDataGroupNameList(c *gin.Context) {
	direction := common.InternalDataGroup
	datagroupname := c.Params.ByName("datagroupname")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	datagroupnamelist, err := ltm.ShowLTMDataGroupName(f5url, direction, datagroupname)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": datagroupnamelist})
}

// LTMVirtualServerPost create virtual server
func LTMVirtualServerPost(c *gin.Context) {
	var vservercreate ltm.CreateVirtualServer
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}

	c.Bind(&vservercreate)
	res, err := ltm.PostLTMVirtualServer(f5url, &vservercreate)
	if err != nil {
		glog.Errorf("%s", err)
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Virtual server added", vservercreate.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmvirtualdocumentationuri"], c)
}

// LTMSSLKeyPost install a new ssl key on local traffic manager
func LTMSSLKeyPost(c *gin.Context) {
	var sslkeycreate ltm.CreateSSLKey
	sslkeycreate.Command = "install"
	//lbpair := c.Params.ByName("lbpair")
	//f5url := DeviceActive(lbpair)
	c.Bind(&sslkeycreate)
	res, _ := ltm.PostLTMSSLKey(common.CryptoURL, &sslkeycreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	switch res.Status {
	case 200:
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

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.Bind(&poolcreate)
	res, _ := ltm.PostLTMPool(f5url, &poolcreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool created", poolcreate.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmpooldocumentationuri"], c)
}

// GTMPoolPost create a new wide IP pool with members and a monitoring check on a global traffic manager
func GTMPoolPost(c *gin.Context) {
	var poolcreate gtm.CreatePool

	f5url, _ := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	c.Bind(&poolcreate)
	res, _ := gtm.PostGTMPool(f5url, &poolcreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool created", poolcreate.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["gtmpooldocumentationuri"], c)
}

// GTMWideipPost create new wide IP on a global traffic manager
func GTMWideipPost(c *gin.Context) {
	var wideipcreate gtm.CreateWip

	f5url, _ := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	c.Bind(&wideipcreate)
	res, _ := gtm.PostGTMWip(f5url, &wideipcreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "WideIP created", wideipcreate.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["gtmwideipdocumentationuri"], c)
}

// LTMPoolMemberPost add new members to a specific pool on a local traffic manager
func LTMPoolMemberPost(c *gin.Context) {
	var poolmembercreate ltm.CreatePoolMember

	poolname := c.Params.ByName("poolname")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}

	c.Bind(&poolmembercreate)
	res, _ := ltm.PostLTMPoolMember(f5url, poolname, &poolmembercreate)
	json.Unmarshal([]byte(res.Body), &returnerror)

	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool member added", poolmembercreate.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMPoolPut modify pool (old ones will be deleted) or change monitoring
func LTMPoolPut(c *gin.Context) {
	var poolmodify ltm.ModifyPool

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}

	c.Bind(&poolmodify)
	res, _ := ltm.PutLTMPool(f5url, poolmodify.Name, &poolmodify)
	json.Unmarshal([]byte(res.Body), &returnerror)

	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool modified", poolmodify.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmpooldocumentationuri"], c)
}

// LTMPoolDelete delete a pool on a local traffic manager
func LTMPoolDelete(c *gin.Context) {
	var pooldelete ltm.RemovePool
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.Bind(&pooldelete)
	res, _ := ltm.DeleteLTMPool(f5url, pooldelete.Name)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool deleted", pooldelete.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmpooldocumentationuri"], c)
}

// GTMPoolDelete delete a pool on a global traffic manager
func GTMPoolDelete(c *gin.Context) {
	var pooldelete gtm.RemovePool
	f5url, _ := gtm.Trafficmanager(c.Params.ByName("trafficmanager"))
	c.Bind(&pooldelete)
	res, _ := gtm.DeleteGTMPool(f5url, pooldelete.Name)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool deleted", pooldelete.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["gtmpooldocumentationuri"], c)
}

// LTMPoolMemberPut modify pool members on a local traffic manager (enabled, disabled, force-offline)
func LTMPoolMemberPut(c *gin.Context) {
	var poolmembermodify ltm.ModifyPoolMember
	poolname := c.Params.ByName("poolname")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.Bind(&poolmembermodify)
	res, _ := ltm.PutLTMPoolMember(f5url, poolname, poolmembermodify.Name, poolmembermodify.Status)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool member modified", poolmembermodify.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMPoolMemberDelete delete specific pool members on a local traffic manager
func LTMPoolMemberDelete(c *gin.Context) {
	var poolmemberdelete ltm.RemovePoolMember
	poolname := c.Params.ByName("poolname")
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.Bind(&poolmemberdelete)
	res, _ := ltm.DeleteLTMPoolMember(f5url, poolname, poolmemberdelete.Name)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Pool member deleted", poolmemberdelete.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmpoolmemberdocumentationuri"], c)
}

// LTMDataGroupPost add new internal datagroup on a local traffic manager
func LTMDataGroupPost(c *gin.Context) {
	var datagroupcreate ltm.CreateDataGroup
	direction := common.InternalDataGroup
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.Bind(&datagroupcreate)
	res, _ := ltm.PostLTMDataGroup(f5url, direction, &datagroupcreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Datagroup added", datagroupcreate.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupDelete delete a internal datagroup on a local traffic manager
func LTMDataGroupDelete(c *gin.Context) {
	var datagroupdelete ltm.RemoveDataGroup
	direction := common.InternalDataGroup
	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.Bind(&datagroupdelete)
	res, _ := ltm.DeleteLTMDataGroup(f5url, direction, datagroupdelete.Name)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Datagroup deleted", datagroupdelete.Name,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupItemPut remove all existing records in a datagroup and add new records (ip or string)
// on a local traffic manager
func LTMDataGroupItemPut(c *gin.Context) {
	var datagroupitemcreate ltm.CreateDataGroupItem
	direction := common.InternalDataGroup
	datagroupname := c.Params.ByName("datagroupname")

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}

	c.Bind(&datagroupitemcreate)
	res, _ := ltm.PatchLTMDataGroupItem(f5url, direction, datagroupname, &datagroupitemcreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Datagroup item added in", datagroupname,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMDataGroupItemPatch add an item to a datagroup (ip or string) on a local traffic manager
func LTMDataGroupItemPatch(c *gin.Context) {
	var datagroupitemcreate ltm.CreateDataGroupItem
	direction := common.InternalDataGroup
	datagroupname := c.Params.ByName("datagroupname")

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}

	c.Bind(&datagroupitemcreate)
	res, _ := ltm.PatchLTMDataGroupItem(f5url, direction, datagroupname, &datagroupitemcreate)
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "Datagroup item added in", datagroupname,
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmdatagroupdocumentationuri"], c)
}

// LTMAddressList show local traffic blocked ip addresses
func LTMAddressList(c *gin.Context) {
	lbpair := c.Params.ByName("lbpair")
	f5url, err := ltm.Loadbalancer(lbpair, common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	f5url = strings.Replace(f5url, common.LtmURI, "", -1)
	addresslist, err := ltm.ShowLTMAddressList(f5url, common.BlackList)
	if err != nil {
		glog.Errorf("%s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": addresslist})
}

// LTMBlockIPPatch add ips which will be blocked
func LTMBlockIPPatch(c *gin.Context) {
	var blockips ltm.CreateAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	f5url = strings.Replace(f5url, common.LtmURI, "", -1)
	c.Bind(&blockips)
	res, err := ltm.PatchLTMBlockAddresses(f5url, &blockips)
	if err != nil {
		glog.Errorf("%s", err)
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "IP(s) blocked successfully", fmt.Sprintf("%+v", blockips),
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMWhiteIPPatch add ips which will be whitelisted
// Not yet implemented
func LTMWhiteIPPatch(c *gin.Context) {
	var whiteips ltm.CreateAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	f5url = strings.Replace(f5url, common.LtmURI, "", -1)
	c.Bind(&whiteips)
	res, err := ltm.PatchLTMWhiteAddresses(f5url, &whiteips)
	if err != nil {
		glog.Errorf("%s", err)
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "IP(s) whitelisted successfully", fmt.Sprintf("%+v", whiteips),
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmaddresslistdocumentationuri"], c)
}

// LTMRemoveBlockIPPatch remove ips which are currently blocked
func LTMRemoveBlockIPPatch(c *gin.Context) {
	var unblockips ltm.DeleteAddresses

	f5url, err := ltm.Loadbalancer(c.Params.ByName("lbpair"), common.Conf.Ltmdevicenames)
	if err != nil {
		glog.Errorf("%s", err)
	}
	f5url = strings.Replace(f5url, common.LtmURI, "", -1)
	c.Bind(&unblockips)
	res, err := ltm.DeleteLTMBlockAddresses(f5url, &unblockips)
	if err != nil {
		glog.Errorf("%s", err)
	}
	json.Unmarshal([]byte(res.Body), &returnerror)
	if res.Status == 200 {
		res.Status = 201
	}
	respondWithStatus(res.Status, "IP(s) removed successfully", fmt.Sprintf("%+v", unblockips),
		returnerror.ErrorMessage(), common.Conf.Documentation["ltmaddresslistdocumentationuri"], c)
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
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// respondWithStatus return JSON response
func respondWithStatus(code int, message, name, err, documentation string, c *gin.Context) {
	switch code {
	case 200:
		{
			c.Set("message", fmt.Sprintf("%s %s", message, name))
			c.JSON(code, gin.H{"message": fmt.Sprintf("%s %s", message, name)})
		}
	case 201:
		c.Set("message", fmt.Sprintf("%s %s", message, name))
		u := new(url.URL)
		u.Scheme = common.Protocol
		u.Path = path.Join(c.Request.Host, c.Request.RequestURI, "/", name)
		c.Header("location", u.String())
		c.JSON(code, gin.H{"message": fmt.Sprintf("%s %s", message, name)})
	case 400:
		c.Set("message", err)
		c.Header("Content-Type", "application/problem+json")
		c.JSON(code, gin.H{"type": documentation, "status": code,
			"title": "Invalid JSON data", "detail": err})
	case 409:
		c.Set("message", err)
		c.Header("Content-Type", "application/problem+json")
		c.JSON(code, gin.H{"type": documentation, "status": code,
			"title": "Conflict", "detail": err})
	default:
		c.Set("message", err)
		c.Header("Content-Type", "application/problem+json")
		c.JSON(code, gin.H{"type": documentation, "status": code,
			"title": err, "detail": err})
	}
}

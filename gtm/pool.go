package gtm

import (
	"fmt"
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net/url"
	"path"
)

// Pools struct provides information about multiple pools
type Pools struct {
	Kind  string `json:"kind"`
	Items []struct {
		Pool
	} `json:"items"`
}

// Pool struct provides information about a specific pool
type Pool struct {
	Kind                      string `json:"kind"`
	Name                      string `json:"name"`
	Partition                 string `json:"partition"`
	Fullpath                  string `json:"fullPath"`
	Generation                int    `json:"generation"`
	Alternatemode             string `json:"alternateMode"`
	Dynamicratio              string `json:"dynamicRatio"`
	Enabled                   bool   `json:"enabled"`
	Fallbackipv4              string `json:"fallbackIpv4"`
	Fallbackipv6              string `json:"fallbackIpv6"`
	Fallbackmode              string `json:"fallbackMode"`
	Limitmaxbps               int    `json:"limitMaxBps"`
	Limitmaxbpsstatus         string `json:"limitMaxBpsStatus"`
	Limitmaxconnections       int    `json:"limitMaxConnections"`
	Limitmaxconnectionsstatus string `json:"limitMaxConnectionsStatus"`
	Limitmaxpps               int    `json:"limitMaxPps"`
	Limitmaxppsstatus         string `json:"limitMaxPpsStatus"`
	Loadbalancingmode         string `json:"loadBalancingMode"`
	Manualresume              string `json:"manualResume"`
	Maxaddressreturned        int    `json:"maxAddressReturned"`
	Monitor                   string `json:"monitor"`
	Qoshitratio               int    `json:"qosHitRatio"`
	Qoshops                   int    `json:"qosHops"`
	Qoskilobytessecond        int    `json:"qosKilobytesSecond"`
	Qoslcs                    int    `json:"qosLcs"`
	Qospacketrate             int    `json:"qosPacketRate"`
	Qosrtt                    int    `json:"qosRtt"`
	Qostopology               int    `json:"qosTopology"`
	Qosvscapacity             int    `json:"qosVsCapacity"`
	Qosvsscore                int    `json:"qosVsScore"`
	TTL                       int    `json:"ttl"`
	Verifymemberavailability  string `json:"verifyMemberAvailability"`
	MembersReference          string `json:"membersReference"`
}

// CreatePool struct to create a pool
type CreatePool struct {
	Name    string `json:"name" binding:"required"`
	Members []struct {
		Name         string `json:"name" binding:"required"`
		Loadbalancer string `json:"loadbalancer,omitempty"`
		Partition    string `json:"partition,omitempty"`
		Subpath      string `json:"subPath,omitempty"`
		Fullpath     string `json:"fullPath,omitempty"`
	} `json:"members" binding:"required"`
	Monitor string `json:"monitor",binding:"required"`
}

// RemovePool struct to delete a pool
type RemovePool struct {
	Name string `json:"name" binding:"required"`
}

// PoolMembers struct provides information about multiple members in one pool
type PoolMembers struct {
	Kind  string `json:"kind"`
	Items []struct {
		PoolMember
	} `json:"items"`
}

// PoolMember struct provides information about a specific member in one pool
type PoolMember struct {
	Kind                      string `json:"kind"`
	Name                      string `json:"name"`
	Partition                 string `json:"partition"`
	Subpath                   string `json:"subPath"`
	Fullpath                  string `json:"fullPath"`
	Generation                int    `json:"generation"`
	Enabled                   bool   `json:"enabled"`
	Limitmaxbps               int    `json:"limitMaxBps"`
	Limitmaxbpsstatus         string `json:"limitMaxBpsStatus"`
	Limitmaxconnections       int    `json:"limitMaxConnections"`
	Limitmaxconnectionsstatus string `json:"limitMaxConnectionsStatus"`
	Limitmaxpps               int    `json:"limitMaxPps"`
	Limitmaxppsstatus         string `json:"limitMaxPpsStatus"`
	Monitor                   string `json:"monitor"`
	Order                     int    `json:"order"`
	Ratio                     int    `json:"ratio"`
	Dependson                 []struct {
		Name      string `json:"name"`
		Partition string `json:"partition"`
		Subpath   string `json:"subPath"`
	} `json:"dependsOn"`
}

// ShowGTMPools shows all declared pools on gtm
func ShowGTMPools(host string) *Pools {
	gtmpools := new(Pools)
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	backend.Request(common.GET, u.String(), &gtmpools)
	return gtmpools
}

// ShowGTMPool shows specific declared pool on gtm
func ShowGTMPool(host, pool string) *Pool {
	// Declaration GTM Pool
	gtmpool := new(Pool)
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri, "/~Common~"+pool)
	backend.Request(common.GET, u.String(), &gtmpool)
	return gtmpool
}

//ShowGTMPoolMembers shows members on a specific pool
func ShowGTMPoolMembers(host, pool string) *PoolMembers {
	// Declaration GTM Pool Member
	gtmpoolmembers := new(PoolMembers)
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri, "/~Common~"+pool, "/members")
	backend.Request(common.GET, u.String(), &gtmpoolmembers)
	return gtmpoolmembers
}

//PostGTMPool creates a new pool on a trafficmanager
func PostGTMPool(host string, json *CreatePool) (*backend.Response, error) {
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)

	for i := range json.Members {
		json.Members[i].Partition = "Common"
		json.Members[i].Subpath = json.Members[i].Loadbalancer + ":/Common"
		json.Members[i].Fullpath = "/Common/" + json.Members[i].Loadbalancer + ":/Common/" + json.Members[i].Name
		json.Members[i].Loadbalancer = ""
	}
	r, err := backend.Request(common.POST, u.String(), &json)
	fmt.Println(json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteGTMPool deletes a pool on a trafficmanager
func DeleteGTMPool(host, pool string) (*backend.Response, error) {
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	u.Path = path.Join(u.Path, "/~Common~"+pool)
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

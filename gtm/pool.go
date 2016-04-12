package gtm

import (
	"fmt"
	"net/url"
	"path"

	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
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
	Canonicialname            string `json:"canonicalName"`
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
	MembersReference          string `json:"memReference"`
}

// CreatePool struct to create a pool
type CreatePool struct {
	Name          string             `json:"name"`
	Members       []CreatePoolMember `json:"members"`
	Monitor       string             `json:"monitor"`
	CanonicalName string             `json:"canonicalName,omitempty"`
}

// CreatePoolMember struct to create a member in a pool
type CreatePoolMember struct {
	Name         string `json:"name"`
	Loadbalancer string `json:"loadbalancer,omitempty"`
	Partition    string `json:"partition,omitempty"`
	Subpath      string `json:"subPath,omitempty"`
	Fullpath     string `json:"fullPath,omitempty"`
}

// RemovePoolMember struct to delete a member in a pool
type RemovePoolMember struct {
	Name         string `json:"name"`
	Loadbalancer string `json:"loadbalancer"`
}

// ModifyPoolMemberStatus struct to modify pool member status on gtm
type ModifyPoolMemberStatus struct {
	Name         string `json:"name"`
	Loadbalancer string `json:"loadbalancer"`
	Status       bool   `json:"status"`
}

// ModifyPoolStatus struct to modify a pool status client-side
type ModifyPoolStatus struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

// EnablePoolMemberStatus enables a gtm pool member client-side
type EnablePoolMemberStatus struct {
	Enabled bool `json:"enabled"`
}

// DisablePoolMemberStatus disables a gtm pool member client-side
type DisablePoolMemberStatus struct {
	Disabled bool `json:"disabled"`
}

// EnablePoolStatus enables a gtm pool client-side
type EnablePoolStatus struct {
	Enabled bool `json:"enabled"`
}

// DisablePoolStatus disables a gtm pool client-side
type DisablePoolStatus struct {
	Disabled bool `json:"disabled"`
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
func ShowGTMPools(host string) (*backend.Response, *Pools, *errors.Error) {
	gtmpools := new(Pools)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	res, err := backend.Request(common.GET, u.String(), &gtmpools)
	if err != nil {
		return nil, nil, err
	}
	return res, gtmpools, nil
}

// ShowGTMPool shows specific declared pool on gtm
func ShowGTMPool(host, pool string) (*backend.Response, *Pool, *errors.Error) {
	// Declaration GTM Pool
	gtmpool := new(Pool)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri, fmt.Sprintf("/~%s~%s", gtmPartition, pool))
	res, err := backend.Request(common.GET, u.String(), &gtmpool)
	if err != nil {
		return nil, nil, err
	}
	return res, gtmpool, nil
}

//ShowGTMPoolMembers shows members on a specific pool
func ShowGTMPoolMembers(host, pool string) (*backend.Response, *PoolMembers, *errors.Error) {
	// Declaration GTM Pool Member
	gtmpoolmembers := new(PoolMembers)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri, fmt.Sprintf("/~%s~%s/members", gtmPartition, pool))
	res, err := backend.Request(common.GET, u.String(), &gtmpoolmembers)
	if err != nil {
		return nil, nil, err
	}
	return res, gtmpoolmembers, nil
}

//PostGTMPool creates a new pool on a trafficmanager
func PostGTMPool(host string, json *CreatePool) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	for i := range json.Members {
		json.Members[i].Partition = gtmPartition
		json.Members[i].Subpath = fmt.Sprintf("%s:/%s", json.Members[i].Loadbalancer, gtmPartition)
		json.Members[i].Fullpath = fmt.Sprintf("/%s/%s:/%s/%s", gtmPartition, json.Members[i].Loadbalancer, ltmPartition, json.Members[i].Name)
		json.Members[i].Loadbalancer = ""
	}
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//PostGTMPoolMember adds new members to an existing pool on a trafficmanager
func PostGTMPoolMember(host string, pool string, json *CreatePoolMember) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri, pool, common.MembersURI)

	json.Partition = gtmPartition
	json.Subpath = fmt.Sprintf("%s:/%s", json.Loadbalancer, gtmPartition)
	json.Fullpath = fmt.Sprintf("/%s/%s:/%s/%s", gtmPartition, json.Loadbalancer, ltmPartition, json.Name)
	json.Loadbalancer = ""
	r, err := backend.Request(common.POST, u.String(), &json)
	fmt.Println(json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteGTMPool deletes a pool on a trafficmanager
func DeleteGTMPool(host, pool string) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	u.Path = path.Join(u.Path, fmt.Sprintf("/~%s~%s", gtmPartition, pool))
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil

}

// DeleteGTMPoolMember delete pool member on a trafficmanager
func DeleteGTMPoolMember(host, pool string, poolmember *RemovePoolMember) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	u.Path = path.Join(u.Path, fmt.Sprintf("/%s/members/~%s~%s:~%s~%s", pool, gtmPartition,
		poolmember.Loadbalancer, ltmPartition, poolmember.Name))
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PutGTMPoolMemberStatus modify status of wideip pool member
func PutGTMPoolMemberStatus(host, pool string, poolmember *ModifyPoolMemberStatus) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	u.Path = path.Join(u.Path, fmt.Sprintf("/~%s~%s/members/~%s~%s:~%s~%s", gtmPartition, pool, gtmPartition,
		poolmember.Loadbalancer, ltmPartition, poolmember.Name))
	var memberstatus interface{}
	switch poolmember.Status {
	case true:
		{
			memberstatus = EnablePoolMemberStatus{Enabled: true}
		}
	case false:
		{
			memberstatus = DisablePoolMemberStatus{Disabled: true}
		}
	}

	r, err := backend.Request(common.PUT, u.String(), &memberstatus)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PutGTMPoolStatus modify status of wideip pool
func PutGTMPoolStatus(host string, poolmodify *ModifyPoolStatus) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmpoolsuri)
	u.Path = path.Join(u.Path, fmt.Sprintf("/~%s~%s", gtmPartition, poolmodify.Name))
	var poolstatus interface{}
	switch poolmodify.Status {
	case true:
		{
			poolstatus = EnablePoolStatus{Enabled: true}
		}
	case false:
		{
			poolstatus = DisablePoolStatus{Disabled: true}
		}
	}

	r, err := backend.Request(common.PUT, u.String(), &poolstatus)
	if err != nil {
		return nil, err
	}
	return r, nil
}

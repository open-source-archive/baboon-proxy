package ltm

import (
	"fmt"
	"net/url"
	"path"

	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
)

/*
{"kind":"tm:ltm:pool:poolcollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/pool?ver=11.5.1","items":[{}]}
*/

// Pools struct provides information
// about multiple pools
type Pools struct {
	Kind  string `json:"kind"`
	Items []struct {
		Pool
	} `json:"items"`
}

// Pool struct provides information
// about a specific pool
type Pool struct {
	Kind                   string `json:"kind"`
	Name                   string `json:"name"`
	Partition              string `json:"partition"`
	FullPath               string `json:"fullPath"`
	Generation             int    `json:"generation"`
	AllowNat               string `json:"allowNat"`
	AllowSnat              string `json:"allowSnat"`
	Description            string `json:"description"`
	IgnorePersistedWeight  string `json:"ignorePersistedWeight"`
	IPTosToClient          string `json:"ipTosToClient"`
	IPTosToServer          string `json:"ipTosToServer"`
	LinkQosToClient        string `json:"linkQosToClient"`
	LinkQosToServer        string `json:"linkQosToServer"`
	LoadBalancingMode      string `json:"loadBalancingMode"`
	MinActiveMembers       int    `json:"minActiveMembers"`
	MinUpMembers           int    `json:"minUpMembers"`
	MinUpMembersAction     string `json:"minUpMembersAction"`
	MinUpMembersChecking   string `json:"minUpMembersChecking"`
	Monitor                string `json:"monitor"`
	QueueDepthLimit        int    `json:"queueDepthLimit"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit"`
	QueueTimeLimit         int    `json:"queueTimeLimit"`
	ReselectTries          int    `json:"reselectTries"`
	SlowRampTime           int    `json:"slowRampTime"`
	MembersReference       string `json:"memReference"`
}

// PoolMembers struct provides information
// about multiple members in one pool
type PoolMembers struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind            string `json:"kind"`
		Name            string `json:"name"`
		Partition       string `json:"partition"`
		FullPath        string `json:"fullPath"`
		Generation      int    `json:"generation"`
		Address         string `json:"address"`
		ConnectionLimit int    `json:"connectionLimit"`
		DynamicRatio    int    `json:"dynamicRatio"`
		InheritProfile  string `json:"inheritProfile"`
		Logging         string `json:"logging"`
		Monitor         string `json:"monitor"`
		PriorityGroup   int    `json:"priorityGroup"`
		RateLimit       string `json:"rateLimit"`
		Ratio           int    `json:"ratio"`
		Session         string `json:"session"`
		State           string `json:"state"`
	} `json:"items"`
}

// CreatePool contain fields to create a pool
type CreatePool struct {
	Name      string `json:"name" binding:"required"`
	Partition string `json:"partition,omitempty"`
	Members   []struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"members,omitempty"`
	Monitor string `json:"monitor" binding:"required"`
}

// ModifyPool contain fields to modify a pool
type ModifyPool struct {
	Name      string `json:"name" binding:"required"`
	Partition string `json:"partition,omitempty"`
	Members   []struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"members,omitempty"`
	Monitor string `json:"monitor,omitempty"`
}

// CreatePoolMember contain fields to add a member in a pool
type CreatePoolMember struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// ModifyPoolMember contain fields to modify a member in a pool
type ModifyPoolMember struct {
	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"required"`
}

// ModifyPoolMemberStatus contain fields to modify a member status in a pool
type ModifyPoolMemberStatus struct {
	State   string `json:"state"`
	Session string `json:"session"`
}

// RemovePoolMember contain fields to delete a member in a pool
type RemovePoolMember struct {
	Name string `json:"name" binding:"required"`
}

// RemovePool contain fields to delete a pool
type RemovePool struct {
	Name string `json:"name" binding:"required"`
}

// ShowLTMPools show all declared pools
func ShowLTMPools(host string) (*backend.Response, *Pools, *errors.Error) {
	// Declaration LTM Pools
	ltmpools := new(Pools)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, "pool")
	res, err := backend.Request(common.GET, u.String(), &ltmpools)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmpools, nil
}

// ShowLTMPool show specific pool
func ShowLTMPool(host, pool string) (*backend.Response, *Pool, *errors.Error) {
	// Declaration LTM Pool
	ltmpool := new(Pool)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s", ltmPartition, pool))
	res, err := backend.Request(common.GET, u.String(), &ltmpool)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmpool, nil
}

// ShowLTMPoolMember show members of a specific pool
func ShowLTMPoolMember(host, pool string) (*backend.Response, *PoolMembers, *errors.Error) {
	// Declaration LTM Pool Member
	ltmpoolmember := new(PoolMembers)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s/members", ltmPartition, pool))
	poolmemberURL := u.String()
	res, err := backend.Request(common.GET, poolmemberURL, &ltmpoolmember)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmpoolmember, nil
}

// PostLTMPool create a new pool
func PostLTMPool(host string, json *CreatePool) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, "pool")
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PutLTMPool modfiy pool like health check
func PutLTMPool(host, poolname string, json *ModifyPool) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s", ltmPartition, poolname))
	r, err := backend.Request(common.PUT, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteLTMPool delete a pool
func DeleteLTMPool(host, pool string) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s", ltmPartition, pool))
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PostLTMPoolMember add member to a pool
func PostLTMPoolMember(host, poolname string, json *CreatePoolMember) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s/members", ltmPartition, poolname))
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PutLTMPoolMember modify status of pool member
func PutLTMPoolMember(host, poolname, member, status string) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s/members/~%s~%s", ltmPartition, poolname, ltmPartition, member))
	memberstatus := ModifyPoolMemberStatus{}
	switch status {
	case "enabled":
		{
			memberstatus = ModifyPoolMemberStatus{State: "user-up", Session: "user-enabled"}
		}
	case "disabled":
		{
			memberstatus = ModifyPoolMemberStatus{State: "user-up", Session: "user-disabled"}
		}
	case "offline":
		{
			memberstatus = ModifyPoolMemberStatus{State: "user-down", Session: "user-disabled"}
		}
	}

	r, err := backend.Request(common.PUT, u.String(), &memberstatus)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteLTMPoolMember delete pool member
func DeleteLTMPoolMember(host, poolname, poolmember string) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("pool/~%s~%s/members/~%s~%s", ltmPartition, poolname, ltmPartition, poolmember))
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

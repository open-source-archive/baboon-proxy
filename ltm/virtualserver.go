package ltm

import (
	"fmt"
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net/url"
	"path"
)

/*
{"kind":"tm:ltm:virtual:virtualcollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/virtual?ver=11.5.1","items":[{}]}
*/

// VirtualServers struct contain fields
// of multiple virtual server
type VirtualServers struct {
	Kind  string `json:"kind"`
	Items []struct {
		VirtualServer
	} `json:"items"`
}

// VirtualServer struct contain fields
// of specific virtual server
type VirtualServer struct {
	Kind                     string `json:"kind"`
	Name                     string `json:"name"`
	Partition                string `json:"partition"`
	FullPath                 string `json:"fullPath"`
	Generation               int    `json:"generation"`
	AddressStatus            string `json:"addressStatus"`
	AutoLasthop              string `json:"autoLasthop"`
	CmpEnabled               string `json:"cmpEnabled"`
	ConnectionLimit          int    `json:"connectionLimit"`
	Destination              string `json:"destination"`
	Enabled                  bool   `json:"enabled"`
	GtmScore                 int    `json:"gtmScore"`
	IPForward                bool   `json:"ipForward"`
	IPProtocol               string `json:"ipProtocol"`
	Mask                     string `json:"mask"`
	Mirror                   string `json:"mirror"`
	MobileAppTunnel          string `json:"mobileAppTunnel"`
	Nat64                    string `json:"nat64"`
	Pool                     string `json:"pool"`
	RateLimit                string `json:"rateLimit"`
	RateLimitDstMask         int    `json:"rateLimitDstMask"`
	RateLimitMode            string `json:"rateLimitMode"`
	RateLimitSrcMask         int    `json:"rateLimitSrcMask"`
	Source                   string `json:"source"`
	SourceAddressTranslation struct {
		Type string `json:"type"`
	} `json:"sourceAddressTranslation"`
	SourcePort        string   `json:"sourcePort"`
	SynCookieStatus   string   `json:"synCookieStatus"`
	TranslateAddress  string   `json:"translateAddress"`
	TranslatePort     string   `json:"translatePort"`
	VlansEnabled      bool     `json:"vlansEnabled"`
	VsIndex           int      `json:"vsIndex"`
	Rules             []string `json:"rules"`
	Vlans             []string `json:"vlans"`
	PoolsReference    string   `json:"poolsReference"`
	ProfilesReference string   `json:"profilesReference"`
	FwRulesReference  string   `json:"fwReference"`
	Persist           []struct {
		Name      string `json:"name"`
		Partition string `json:"partition"`
		TmDefault string `json:"tmDefault"`
	} `json:"persist"`
}

// CreateVirtualServer contain fields to create a virtual server
type CreateVirtualServer struct {
	Name        string   `json:"name" binding:"required"`
	Partition   string   `json:"partition,omitempty"`
	Destination string   `json:"destination" binding:"required"`
	Mask        string   `json:"mask" binding:"required"`
	IPForward   bool     `json:"ipForward,omitempty"`
	IPProtocol  string   `json:"ipProtocol" binding:"required"`
	Pool        string   `json:"pool" binding:"required"`
	Rules       []string `json:"rules,omitempty"`
	Vlans       []string `json:"vlans,omitempty"`
	Persist     []struct {
		Name string `json:"name,omitempty"`
	} `json:"persist,omitempty"`
	Profiles []struct {
		Name string `json:"name" binding:"required"`
	} `json:"profiles" binding:"required"`
	SourceAddressTranslation struct {
		Type string `json:"type,omitempty"`
	} `json:"sourceAddressTranslation,omitempty"`
}

// ShowLTMVirtualServer show all virtual server
func ShowLTMVirtualServer(host string) *VirtualServers {
	// Declaration LTM virtual server
	ltmvirtualserver := new(VirtualServers)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, "virtual")
	backend.Request(common.GET, u.String(), ltmvirtualserver)
	return ltmvirtualserver
}

// ShowLTMVirtualServerName show specific virtual server
func ShowLTMVirtualServerName(host, vserver string) *VirtualServer {
	// Declaration LTM virtual server name
	ltmvirtualservername := new(VirtualServer)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, fmt.Sprintf("virtual/~%s~%s", ltmPartition, vserver))
	backend.Request(common.GET, u.String(), &ltmvirtualservername)
	return ltmvirtualservername
}

// PostLTMVirtualServer create a new virtual server
func PostLTMVirtualServer(host string, json *CreateVirtualServer) (*backend.Response, error) {
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, "virtual")
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

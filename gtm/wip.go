package gtm

import (
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net/url"
	"path"
)

// Wips struct provides information about multiple wide ips
type Wips struct {
	Kind  string `json:"kind"`
	Items []struct {
		Wip
	} `json:"items"`
}

// Wip struct provides information about a specific wide ip
type Wip struct {
	Kind                string   `json:"kind"`
	Name                string   `json:"name"`
	Partition           string   `json:"partition"`
	Fullpath            string   `json:"fullPath"`
	Generation          int      `json:"generation"`
	Enabled             bool     `json:"enabled"`
	Ipv6Noerrornegttl   int      `json:"ipv6NoErrorNegTtl"`
	Ipv6Noerrorresponse string   `json:"ipv6NoErrorResponse"`
	Persistcidripv4     int      `json:"persistCidrIpv4"`
	Persistcidripv6     int      `json:"persistCidrIpv6"`
	Persistence         string   `json:"persistence"`
	Poollbmode          string   `json:"poolLbMode"`
	Ttlpersistence      int      `json:"ttlPersistence"`
	Rules               []string `json:"rules"`
	Pools               []struct {
		Name           string `json:"name"`
		Partition      string `json:"partition"`
		Order          int    `json:"order"`
		Ratio          int    `json:"ratio"`
		PoolsReference string `json:"poolsReference"`
	} `json:"pools"`
}

// CreateWip struct to add a wide ip
type CreateWip struct {
	Name  string `json:"name" binding:"required"`
	Pools []struct {
		Name string `json:"name"`
	} `json:"pools" binding:"required"`
	Poollbmode string `json:"poolLbMode"`
}

// RemoveWip struct to delete a wide ip
type RemoveWip struct {
	Name string `json:"name" binding:"required"`
}

// ShowGTMWips lists all wide ips on a trafficmanager
func ShowGTMWips(host string) *Wips {
	gtmwips := new(Wips)
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmwipsuri)
	backend.Request(common.GET, u.String(), &gtmwips)
	return gtmwips
}

// ShowGTMWip list a specific wide ip on a trafficmanager
func ShowGTMWip(host, wideip string) *Wip {
	gtmwip := new(Wip)
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmwipsuri, "/", wideip)
	backend.Request(common.GET, u.String(), &gtmwip)
	return gtmwip
}

// PostGTMWip creates a new wide ip on a trafficmanager
func PostGTMWip(host string, json *CreateWip) (*backend.Response, error) {
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmwipsuri)
	if !(len(json.Poollbmode) > 0) {
		json.Poollbmode = "global-availability"
	}

	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteGTMWip deletes a pool on a trafficmanager
func DeleteGTMWip(host, wideip string) (*backend.Response, error) {
	u, _ := url.Parse(host)
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmwipsuri)
	u.Path = path.Join(u.Path, "/~Common~"+wideip)
	r, err := backend.Request(common.DELETE, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

package ltm

import (
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net/url"
	"path"
)

/*
{"kind":"tm:ltm:virtual:fw-rules:fw-rulescollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/virtual/profiles?ver=11.5.1","items":[{}]}
*/

// IRules struct provides information
// about rules on a virtual server
type IRules struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind       string `json:"kind"`
		Name       string `json:"name"`
		Partition  string `json:"partition"`
		Fullpath   string `json:"fullPath"`
		Generation int    `json:"generation"`
		Context    string `json:"context"`
	} `json:"items"`
}

// ShowLTMFWRules shows firewall profile
func ShowLTMFWRules(host, vserver string) *IRules {
	fwrules := new(IRules)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, "virtual/~Common~"+vserver, "/fw-rules")
	backend.Request(common.GET, u.String(), &fwrules)
	return fwrules
}

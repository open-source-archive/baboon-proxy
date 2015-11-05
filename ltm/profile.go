package ltm

import (
	"fmt"
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net/url"
	"path"
)

/*
{"kind":"tm:ltm:virtual:profiles:profilescollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/virtual/profiles?ver=11.5.1","items":[{}]}
*/

// Profiles struct provides information
// about multiple profiles on a virtual server
type Profiles struct {
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

// ShowLTMProfile show profiles from a virtual server
func ShowLTMProfile(host, vserver string) *Profiles {
	// Declaration LTM Profile
	ltmprofile := new(Profiles)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, fmt.Sprintf("virtual/~%s~%s/profiles", ltmPartition, vserver))
	backend.Request(common.GET, u.String(), &ltmprofile)
	return ltmprofile
}

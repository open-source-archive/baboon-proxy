package ltm

import (
	"fmt"
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
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
func ShowLTMProfile(host, vserver string) (*backend.Response, *Profiles, *errors.Error) {
	// Declaration LTM Profile
	ltmprofile := new(Profiles)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("virtual/~%s~%s/profiles", ltmPartition, vserver))
	res, err := backend.Request(common.GET, u.String(), &ltmprofile)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmprofile, nil
}

package gtm

import (
	"fmt"
	"net/url"
	"path"

	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
)

// IRules struct provides information
// about rules on a wide ip
type IRules struct {
	Kind  string  `json:"kind"`
	Items []IRule `json:"items"`
}

// IRule struct provides information
// about a specific rule on a wide ip
type IRule struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Partition    string `json:"partition"`
	Fullpath     string `json:"fullPath"`
	Generation   int    `json:"generation"`
	Apianonymous string `json:"apiAnonymous"`
}

// ShowGTMIRules shows iRules
func ShowGTMIRules(host string) (*IRules, error) {
	iRs := new(IRules)
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmirulesuri)
	_, err = backend.Request(common.GET, u.String(), &iRs)
	if err != nil {
		return nil, err
	}
	return iRs, nil
}

// ShowGTMIRule shows a specific iRule
func ShowGTMIRule(host, iRuleName string) (*IRule, error) {
	iR := new(IRule)
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	u.Scheme = common.Protocol
	u.Path = path.Join(u.Path, common.Gtmirulesuri, fmt.Sprintf("/%s", iRuleName))
	_, err = backend.Request(common.GET, u.String(), &iR)
	if err != nil {
		return nil, err
	}
	return iR, nil
}

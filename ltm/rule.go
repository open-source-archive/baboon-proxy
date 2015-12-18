package ltm

import (
	"fmt"
	"net/url"
	"path"

	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
)

/*
{"kind":"tm:ltm:virtual:fw-rules:fw-rulescollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/virtual/profiles?ver=11.5.1","items":[{}]}
*/

// IRules struct provides information
// about rules on a virtual server
type IRules struct {
	Kind  string  `json:"kind"`
	Items []IRule `json:"items"`
}

// IRule struct provides information
// about a specific rule on a virtual server
type IRule struct {
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Partition    string `json:"partition"`
	Fullpath     string `json:"fullPath"`
	Generation   int    `json:"generation"`
	Apianonymous string `json:"apiAnonymous"`
}

// FirewallRules on a virtual server
type FirewallRules struct {
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
func ShowLTMFWRules(host, vserver string) *FirewallRules {
	fwrules := new(FirewallRules)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, fmt.Sprintf("virtual/~%s~%s/fw-rules", ltmPartition, vserver))
	backend.Request(common.GET, u.String(), &fwrules)
	return fwrules
}

// ShowLTMIRules shows iRules
func ShowLTMIRules(host string) (*IRules, error) {
	iRs := new(IRules)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, "rule")
	_, err := backend.Request(common.GET, u.String(), &iRs)
	if err != nil {
		return nil, err
	}
	return iRs, nil
}

// ShowLTMIRule shows a specific iRule
func ShowLTMIRule(host, iRuleName string) (*IRule, error) {
	iR := new(IRule)
	u, _ := url.Parse(host)
	u.Path = path.Join(u.Path, fmt.Sprintf("rule/~%s~%s", ltmPartition, iRuleName))
	_, err := backend.Request(common.GET, u.String(), &iR)
	if err != nil {
		return nil, err
	}
	return iR, nil
}

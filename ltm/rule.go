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
func ShowLTMFWRules(host, vserver string) (*backend.Response, *FirewallRules, *errors.Error) {
	fwrules := new(FirewallRules)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("virtual/~%s~%s/fw-rules", ltmPartition, vserver))
	res, err := backend.Request(common.GET, u.String(), &fwrules)
	if err != nil {
		return nil, nil, err
	}
	return res, fwrules, nil
}

// ShowLTMIRules shows iRules
func ShowLTMIRules(host string) (*backend.Response, *IRules, *errors.Error) {
	iRs := new(IRules)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, "rule")
	res, err := backend.Request(common.GET, u.String(), &iRs)
	if err != nil {
		return nil, nil, err
	}
	return res, iRs, nil
}

// ShowLTMIRule shows a specific iRule
func ShowLTMIRule(host, iRuleName string) (*backend.Response, *IRule, *errors.Error) {
	iR := new(IRule)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, fmt.Sprintf("rule/~%s~%s", ltmPartition, iRuleName))
	res, err := backend.Request(common.GET, u.String(), &iR)
	if err != nil {
		return nil, nil, err
	}
	return res, iR, nil
}

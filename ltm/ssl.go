package ltm

import (
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
	"net/url"
	"path"
)

/*
{"kind":"tm:ltm:pool:poolcollectionstate","selfLink":"https://localhost/mgmt/tm/ltm/pool?ver=11.5.1","items":[{}]}
*/

// CreateSSLKey contain fields to create a ssl key
type CreateSSLKey struct {
	Command string `json:"command" binding:"required"`
	Name    string `json:"name" binding:"required"`
	File    string `json:"from-local-file" binding:"required"`
}

// CreateSSLCert contain fields to create a ssl certificate
type CreateSSLCert struct {
	Command string `json:"command" binding:"required"`
	Name    string `json:"name" binding:"required"`
	File    string `json:"from-local-file" binding:"required"`
}

// CreateSSLProfile contain fields to create a ssl profile
type CreateSSLProfile struct {
	Name string `json:"name" binding:"required"`
}

// PostLTMSSLKey create a new ssl key
func PostLTMSSLKey(host string, json interface{}) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, "key")
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PostLTMSSLCert create a new ssl certificate
func PostLTMSSLCert(host string, json interface{}) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, "cert")
	r, err := backend.Request(common.POST, u.String(), &json)
	if err != nil {
		return nil, err
	}
	return r, nil
}

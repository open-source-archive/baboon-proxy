package backend

import (
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
)

var credentials *Credentials

// Credentials contain fields
// which are necessary to use
// F5 API
type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// InitCredentials load config file
// which contains local admin user and password
// from F5 device
func InitCredentials() {
	switch {
	case common.Conf.Backend["f5user"] == "":
		glog.Fatalf("Could not get F5 user from config file")
	case common.Conf.Backend["f5password"] == "":
		glog.Fatalf("Could not get F5 password from config file")
	default:
		credentials = &Credentials{common.Conf.Backend["f5user"], common.Conf.Backend["f5password"]}
	}

}

package backend

import (
	"crypto/tls"
	"github.com/jmcvetta/napping"
	"net/http"
	"net/url"
)

var (
	sess    napping.Session
	cln     http.Client
	tr      http.Transport
	tlsconf tls.Config
	headers http.Header
)

// InitSession initiliaze HTTPS session
func InitSession() {

	// insecure ssl cert, not trusted, skip it
	tr = http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cln = http.Client{Transport: &tr}
	headers = make(http.Header)
	// get a session with your credentials
	sess = napping.Session{
		Client:   &cln,
		Userinfo: url.UserPassword(credentials.User, credentials.Password),
		Header:   &headers,
	}
}

package backend

import (
	"encoding/json"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"log"
	"os"
)

var conf *config.Config

// Credentials contain fields
// which are necessary to use
// F5 API
type Credentials struct {
	User     string `json:"user"`
	Password string `json:"pass"`
}

// InitCredentials load config file
// which contains local admin user and password
// from F5 device
func InitCredentials() {
	conf = config.LoadConfig()

	credentials = Credentials{}

	apiFile := conf.Security["apiFile"]
	file, err := os.Open(apiFile)
	if err != nil {
		log.Fatalf("File not found %s", apiFile)
	}
	config := json.NewDecoder(file)
	err = config.Decode(&credentials)
	if err != nil {
		log.Fatalf("Cannot decode configuration: ", err)
	}
}

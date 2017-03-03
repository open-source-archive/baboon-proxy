//handles the configuration of the applications. Yaml files are mapped with the struct

package config

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/viper"
	"github.com/zalando-techmonkeys/gin-oauth2"
	"golang.org/x/oauth2"
)

// User contain fields for
// OAUTH2 implementation
type User struct {
	Username string
	Fullname string
	Role     string
	Group    string
}

// Config contain fields for config yaml file
type Config struct {
	Endpoints            map[string]string
	AllowedUsers         []User
	Security             map[string]string
	Documentation        map[string]string
	Ltmdevicenames       map[string]string
	Externalgtmlisteners map[string]string
	Internalgtmlisteners map[string]string
	Backend              map[string]string
	Partition            map[string]string
	LTMMgmtIP            map[string]string
	GTMMgmtIP            map[string]string
	ITMMgmtIP            map[string]string
	SNMP                 SNMPBase
}

// SNMPBase contain fields to check backend health
type SNMPBase struct {
	Community string
	Version   string
	TimeOut   int
	Port      string
	OIDs      map[string]interface{}
}

// Error contain fields of config error handling
//created a struct just for future usage
type Error struct {
	Message string
}

// Error return config error
// viper can not parse yaml file
func (e *Error) Error() string {
	return fmt.Sprintf(e.Message)
}

// InitConfig initiliaze configuration file
func InitConfig(filename string) (*Config, *Error) {
	viper.SetConfigType("YAML")
	f, err := os.Open(filename)
	if err != nil {
		return nil, &Error{"could not read configuration files."}
	}
	err = viper.ReadConfig(f)
	if err != nil {
		return nil, &Error{"configuration format is not correct."}
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		glog.Errorf("Cannot read configuration. Reason: %s", err)
		return nil, &Error{"cannot read configuration, something must be wrong."}
	}

	return &config, nil
}

// LoadAuthConf extract necessary oAuth2 information
func LoadAuthConf(config *Config) ([]ginoauth2.AccessTuple, []ginoauth2.AccessTuple, *oauth2.Endpoint, *Error) {
	var rootUsers = []ginoauth2.AccessTuple{}
	var emergencyUsers = []ginoauth2.AccessTuple{}
	var endpoint = oauth2.Endpoint{}
	for _, user := range config.AllowedUsers {
		username := user.Username
		fullname := user.Fullname
		role := user.Role
		if username == "" || fullname == "" || role == "" {
			return nil, nil, nil, &Error{"configuration is invalid. TokenRUL or AuthURL are missing"}
		}
		u := ginoauth2.AccessTuple{role, username, fullname}
		if user.Group == "root" {
			rootUsers = append(rootUsers, u)
		}
		if user.Group == "emergency" || user.Group == "root" {
			emergencyUsers = append(emergencyUsers, u)
		}
	}
	authURL := config.Endpoints["AuthURL"]
	tokenURL := config.Endpoints["TokenURL"]
	if authURL == "" || tokenURL == "" {
		return nil, nil, nil, &Error{"configuration is invalid. TokenURL or AuthURL are missing"}
	}
	endpoint = oauth2.Endpoint{authURL, tokenURL}
	return rootUsers, emergencyUsers, &endpoint, nil
}

// LoadConfig initiliaze config file
func LoadConfig() *Config {
	var err *Error
	conf, err := InitConfig("config.yaml")
	if err != nil {
		glog.Errorf("Cannot load configuration. Reason: %s", err.Message)
		panic("Cannot load configuration for Baboon. Exiting.")
	}
	return conf
}

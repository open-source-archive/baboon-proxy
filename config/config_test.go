package config

import "testing"
import "fmt"

//tests that if the config file is missing, the software recognizes the problem
func TestMissingFile(t *testing.T) {
	_, err := InitConfig("filenotfound.yaml")
	if err == nil {
		t.Fatalf("Func not returing error but file is missing")
	}
}

//tests that errors in the config are found
func TestWrongConfig(t *testing.T) {
	conf, err := InitConfig("wrongconfig.yaml")
	_, _, _, err = LoadAuthConf(conf)
	if err == nil {
		t.Fatalf("Config file not correct, but code is not returning error")
	}
}

//tests that a semantically correct configuration works
func TestOKConfig(t *testing.T) {
	conf, _ := InitConfig("example.yaml")
	users, _, _, err := LoadAuthConf(conf)
	if err != nil {
		t.Fatalf("Error loading configuration. Either configuration is incorrect or we have a bug")
	}

	for _, user := range users {
		if user.Realm == "" || user.Uid == "" && user.Cn == "" {
			fmt.Println(user)
			t.Fatalf("Error with user")
		}
	}

	//TODO: test all the remaining parts
}

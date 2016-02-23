package util

import (
	"testing"
)

type IPScope struct {
	Range string
	Match bool
}

func TestReplaceCommon(t *testing.T) {
	is := "/Common/vs_test_80"
	expected := "vs_test_80"

	got := ReplaceCommon(is)
	if got != expected {
		t.Error("Should replace /Common/")
	}
}

func TestReplaceColon(t *testing.T) {
	is := "192.168.0.1:80"
	expect := "192.168.0.1"

	got := ReplaceColon(is)

	if expect != got {
		t.Errorf("Expect: %s, but got %s", expect, got)
	} else {
		t.Logf("Expect : %s, and got %s", expect, got)
	}
}

func TestReplaceLTMUritoDeviceUri(t *testing.T) {
	is := "https://f5.com/mgmt/tm/ltm"
	expect := "https://f5.com/mgmt/tm/cm/device"

	got := ReplaceLTMUritoDeviceURI(is)
	if expect != got {
		t.Errorf("Expect: %s, but got %s", expect, got)
	} else {
		t.Logf("Expect : %s, and got %s", expect, got)
	}
}

func TestReplaceLTMUritoAddressListURI(t *testing.T) {
	is := "https://f5.com/mgmt/tm/ltm"
	expect := "https://f5.com/mgmt/tm/security/firewall/address-list/"

	got := ReplaceLTMUritoAddressListURI(is)
	if expect != got {
		t.Errorf("Expect: %s, but got %s", expect, got)
	} else {
		t.Logf("Expect : %s, and got %s", expect, got)
	}
}

func TestReplaceGTMWipUritoGTMPoolURI(t *testing.T) {
	is := "f5.com/mgmt/gtm/wideips"
	expect := "f5.com/mgmt/gtm/pools/"

	got := ReplaceGTMWipUritoGTMPoolURI(is)

	if expect != got {
		t.Errorf("Expect: %s, but got %s", expect, got)
	} else {
		t.Logf("Expect : %s, and got %s", expect, got)
	}
}

func TestVerifyIPv4Scope(t *testing.T) {
	testScope := map[string]string{
		"192.168.0.1":    "192.168.0.1",
		"192.168.0.3":    "192.168.0.2",
		"192.168.0.7777": "192.168.0.0-192.168.255.255",
		"192.168.0.230":  "192.168.0.0-192.168.255.255",
		"192.168.0.4":    "10.0.0.0-10.255.255.255",
	}
	expect := map[string]bool{
		"192.168.0.1":    true,
		"192.168.0.3":    false,
		"192.168.0.7777": false,
		"192.168.0.230":  true,
		"192.168.0.4":    false,
	}
	var got bool
	for k, v := range testScope {
		got = VerifyIPv4Scope(k, v)
		if got != expect[k] {
			t.Errorf("Expect %v but got %v - %s %s", expect[k], got, k, v)
		} else {
			if expect[k] {
				t.Logf("%s in scope of %s", k, v)
			} else {
				t.Logf("%s not in scope of %s", k, v)
			}
		}
	}
}

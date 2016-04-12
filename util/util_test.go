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

func TestCheckDeviceStatus(t *testing.T) {
	var signalTest = []struct {
		signal int
		expect string
	}{
		{1, "Offline"},
		{2, "ForcedOffline"},
		{3, "Standby"},
		{4, "Active"},
		{5, "Unknown"},
	}
	for _, tt := range signalTest {
		actual := CheckDeviceStatus(tt.signal)
		if actual != tt.expect {
			t.Errorf("CheckDeviceStatus(%d): expected %s, actual %s", tt.signal, tt.expect, actual)
		}
	}
}

func TestCheckPoolStatus(t *testing.T) {
	var signalTest = []struct {
		signal string
		expect string
	}{
		{"Available", "Available"},
		{"No enabled pool members available", "No enabled pool members available"},
		{"Unknown", "Unknown"},
	}
	for _, tt := range signalTest {
		actual := CheckPoolStatus(tt.signal)
		if actual != tt.expect {
			t.Errorf("CheckPoolStatus(%s): expected %s, actual %s", tt.signal, tt.expect, actual)
		}
	}
}

func TestCheckWideIPStatus(t *testing.T) {
	var signalTest = []struct {
		signal string
		expect string
	}{
		{"Available", "Available"},
		{"No enabled pools available", "No enabled pools available"},
		{"Unknown", "Unknown"},
	}
	for _, tt := range signalTest {
		actual := CheckWideIPStatus(tt.signal)
		if actual != tt.expect {
			t.Errorf("CheckWideIPStatus(%s): expected %s, actual %s", tt.signal, tt.expect, actual)
		}
	}
}

func TestCheckGSLBServerStatus(t *testing.T) {
	var signalTest = []struct {
		signal int
		expect string
	}{
		{1, "Available"},
		{2, "Unavailable"},
		{3, "No enabled Virtual Server available"},
		{4, "Unknown"},
		{5, "Unlicensed"},
	}
	for _, tt := range signalTest {
		actual := CheckGSLBServerStatus(tt.signal)
		if actual != tt.expect {
			t.Errorf("CheckGSLBServerStatus(%d): expected %s, actual %s", tt.signal, tt.expect, actual)
		}
	}
}

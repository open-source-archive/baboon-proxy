package ltm

import (
	"github.com/golang/glog"
	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
	"github.com/zalando-techmonkeys/baboon-proxy/util"
	"net/url"
	"path"
)

// DeleteAddresses include fields
// to delete IP addresses from IP blacklist
type DeleteAddresses struct {
	Addresses []struct {
		Name string `json:"name" binding:"required"`
	} `json:"addresses" binding:"required"`
}

// CreateAddresses include fields
// to block IP addresses
type CreateAddresses struct {
	Addresses []struct {
		Name string `json:"name" binding:"required"`
	} `json:"addresses" binding:"required"`
}

// AddressList showing all
// firewall address list fields
type AddressList struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Fullpath   string `json:"fullPath"`
	Generation int    `json:"generation"`
	Addresses  []struct {
		Name string `json:"name" binding:"required"`
	} `json:"addresses" binding:"required"`
}

// ShowLTMAddressList returns a specific address list on LB
func ShowLTMAddressList(host, address string) (*backend.Response, *AddressList, *errors.Error) {
	addresslist := new(AddressList)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.AddressList, address)
	res, err := backend.Request(common.GET, u.String(), addresslist)
	if err != nil {
		return nil, nil, err
	}
	return res, addresslist, nil
}

// ShowLTMAddressListName returns a specific address list on LB
func ShowLTMAddressListName(host, address string) (*backend.Response, *AddressList, *errors.Error) {
	addresslist := new(AddressList)
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.AddressList, address)
	res, err := backend.Request(common.GET, u.String(), addresslist)
	if err != nil {
		return nil, nil, err
	}
	return res, addresslist, nil
}

// PatchLTMBlockAddresses to block IPs
// Due not overwriting address list
// its necessary to get all entries first
func PatchLTMBlockAddresses(host string, blockIP *CreateAddresses) (*backend.Response, *errors.Error) {

	_, whiteIP, err := ShowLTMAddressListName(host, common.WhiteList)
	if err != nil {
		return nil, err
	}
	_, blackIP, err := ShowLTMAddressListName(host, common.BlackList)
	if err != nil {
		return nil, err
	}
	var whitelistedIP CreateAddresses

	// IPs which you want to block should be first checked
	// IPs in whitelist IP scope are not added to blacklist

	for _, w := range whiteIP.Addresses {
		for _, b := range blockIP.Addresses {
			if util.VerifyIPv4Scope(b.Name, w.Name) {
				whitelistedIP.Addresses = append(whitelistedIP.Addresses, b)
			}
		}
	}
	// first slices are compared to order loops
	// if blocking ips are in white list it skips IP
	// append(...) is doing a delete item of the slice
	// see https://github.com/golang/go/wiki/SliceTricks

	glog.Infof("IPs to block: %+v", blockIP.Addresses)
	if len(blockIP.Addresses) < len(whitelistedIP.Addresses) {
		for i, b := range blockIP.Addresses {
			for _, w := range whitelistedIP.Addresses {
				if b == w {
					blockIP.Addresses = append(blockIP.Addresses[:i], blockIP.Addresses[i+1:]...)
					break
				}
			}
		}
	} else {
		for _, w := range whitelistedIP.Addresses {
			for i, b := range blockIP.Addresses {
				if w == b {
					blockIP.Addresses = append(blockIP.Addresses[:i], blockIP.Addresses[i+1:]...)
					break
				}
			}
		}
	}

	glog.Infof("Blocked IPs %+v", blockIP.Addresses)

	// add already blacklisted IPs to our new black IPs

	for i := range blackIP.Addresses {
		blockIP.Addresses = append(blockIP.Addresses, blackIP.Addresses[i])
	}

	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.AddressList, common.BlackList)

	r, err := backend.Request(common.PATCH, u.String(), &blockIP)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteLTMBlockAddresses remove IPs from blacklist
func DeleteLTMBlockAddresses(host string, deleteIP *DeleteAddresses) (*backend.Response, *errors.Error) {
	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.AddressList, common.BlackList)

	_, blackIP, err := ShowLTMAddressListName(host, common.BlackList)
	if err != nil {
		return nil, err
	}

	if len(blackIP.Addresses) < len(deleteIP.Addresses) {
		for i := range blackIP.Addresses {
			for j := range deleteIP.Addresses {
				if deleteIP.Addresses[j] == blackIP.Addresses[i] {
					blackIP.Addresses = append(blackIP.Addresses[:i], blackIP.Addresses[i+1:]...)
					break
				}
			}
		}
	} else {
		for i := range deleteIP.Addresses {
			for j := range blackIP.Addresses {
				if deleteIP.Addresses[i] == blackIP.Addresses[j] {
					blackIP.Addresses = append(blackIP.Addresses[:j], blackIP.Addresses[j+1:]...)
					break
				}
			}
		}
	}

	newblacklist := DeleteAddresses{Addresses: blackIP.Addresses}
	r, err := backend.Request(common.PATCH, u.String(), &newblacklist)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// PatchLTMWhiteAddresses to whitelist ips
// Due not overwriting address list
// its necessary to get all entries first
func PatchLTMWhiteAddresses(host string, whiteIP *CreateAddresses) (*backend.Response, *errors.Error) {

	_, white, err := ShowLTMAddressListName(host, common.WhiteList)
	if err != nil {
		return nil, err
	}

	// first slices are compared to order loops
	// if white list has already white ips it will skip
	// append(...) is doing a delete item of the slice
	// see https://github.com/golang/go/wiki/SliceTricks

	if len(white.Addresses) < len(whiteIP.Addresses) {
		for i := range white.Addresses {
			for j := range whiteIP.Addresses {
				if whiteIP.Addresses[j] == white.Addresses[i] {
					whiteIP.Addresses = append(whiteIP.Addresses[:j], whiteIP.Addresses[j+1:]...)
					break
				}
			}
		}
	} else {
		for i := range whiteIP.Addresses {
			for j := range white.Addresses {
				if whiteIP.Addresses[i] == white.Addresses[j] {
					whiteIP.Addresses = append(whiteIP.Addresses[:i], whiteIP.Addresses[i+1:]...)
					break
				}
			}
		}
	}

	for i := range white.Addresses {
		whiteIP.Addresses = append(whiteIP.Addresses, white.Addresses[i])
	}

	u, errParse := url.Parse(host)
	if errParse != nil {
		return nil, &errors.ErrorCodeBadRequestParse
	}
	u.Path = path.Join(u.Path, common.AddressList, common.WhiteList)

	r, err := backend.Request(common.PATCH, u.String(), &whiteIP)
	if err != nil {
		return nil, err
	}
	return r, nil
}

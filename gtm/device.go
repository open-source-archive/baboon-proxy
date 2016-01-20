package gtm

import (
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
	"net"
	"time"
)

// Trafficmanager matches internal (ITM)
// or external (GTM)
// Check if ITM/GTM is available
// Required to create WIPs and Pools
func Trafficmanager(cluster string) (string, *errors.Error) {
	var (
		dnsserver         string
		seconds           = 2
		tm                map[string]string
		internalListeners = common.Conf.Internalgtmlisteners
		externalListeners = common.Conf.Externalgtmlisteners
	)

	switch cluster {
	case "itm":
		{
			tm = internalListeners
		}
	case "gtm":
		{
			tm = externalListeners
		}
	default:
		{
			return "", &errors.ErrorCodeNotFoundPattern
		}
	}
	timeOut := time.Duration(seconds) * time.Second
	for name, ipPort := range tm {
		_, err := net.DialTimeout("tcp", ipPort, timeOut)
		if err != nil {
			continue
		}
		dnsserver = name
		break
	}
	return dnsserver, nil
}

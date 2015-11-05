package gtm

import (
	"fmt"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"net"
	"time"
)

// Trafficmanager matches internal (ITM)
// or external (GTM)
// Check if ITM/GTM is available
// Required to create WIPs and Pools
func Trafficmanager(cluster string) (string, error) {
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
			return "", fmt.Errorf("Pattern %s not found, should be itm or gtm", cluster)
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

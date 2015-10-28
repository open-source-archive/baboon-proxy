package gtm

import (
	"fmt"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"net"
	"time"
)

var conf *config.Config

func init() {
	conf = config.LoadConfig()
}

// Trafficmanager matches internal (ITM)
// or external (GTM)
// Check if ITM/GTM is available
// Required to create WIPs and Pools
func Trafficmanager(cluster string) (string, error) {
	var dnsserver string
	var seconds = 2
	var tm map[string]string

	switch cluster {
	case "itm":
		{
			tm = conf.Internalgtmdevicenames
		}
	case "gtm":
		{
			tm = conf.Externalgtmdevicenames
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

package common

//Constants for HTTP methods
const (
	GET = iota + 1
	POST
	PUT
	PATCH
	DELETE
)

// Protocol Scheme
const Protocol = "https"

// LtmURI Local Traffic Manager Path
const LtmURI = "/mgmt/tm/ltm"

// DeviceURI Local Traffic Manager Device Path
const DeviceURI = "/mgmt/tm/cm/device"

// Gtmpoolsuri Global Traffic Manager Pool Path
const Gtmpoolsuri = "/mgmt/tm/gtm/pool"

// Gtmwipsuri Global Traffic Manager WideIP Path
const Gtmwipsuri = "/mgmt/tm/gtm/wideip"

// Gtmirulesuri Global Traffic Manager IRules Path
const Gtmirulesuri = "/mgmt/tm/gtm/rule"

// Dg Local Traffic Manager DataGroup Path
const Dg = "data-group/"

// CryptoURL Local Traffic Manager Crypto Path
const CryptoURL = "/mgmt/tm/sys/crypto/"

// MembersURI Members Path
const MembersURI = "/members"

// PoolsURI Pools Path
const PoolsURI = "/pools/"

// FwURI iRules Path
const FwURI = "/rules"

// ProfilesURI Virtual Server Profile Path
const ProfilesURI = "/profiles"

// InternalDataGroup Local Traffic Manager Internal DataGroup Path
const InternalDataGroup = "internal"

// ExternalDataGroup Local Traffic Manager External DataGroup Path
const ExternalDataGroup = "external"

// AddressList Local Traffic Manager Network Firewall Address List Path
const AddressListURI = "/mgmt/tm/security/firewall/address-list/"

// BlackList Local Traffic Manager BlackList Address List
// for future planning ipv6 addresses
const BlackList = "black_listed_ipv4_addr"

// WhiteList Local Traffic Manager WhiteList Address List
const WhiteList = "white_listed_ipv4_addr"

// types.go

/*
Package ssllabs These are the types used by SSLLabs/Qualys

This is for API v3
*/
package ssllabs

import (
	"encoding/json"
)

// LabsError is for whatever error we get from SSLLabs
type LabsError struct {
	Field   string
	Message string
}

// LabsErrorResponse is a set of errors
type LabsErrorResponse struct {
	ResponseErrors []LabsError `json:"errors"`
}

// Error() implements the interface
func (e LabsErrorResponse) Error() string {
	msg, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(msg)
}

// Info describes the current SSLLabs engine used
type Info struct {
	EngineVersion        string `json:"engineVersion"`
	CriteriaVersion      string `json:"criteriaVersion"`
	MaxAssessments       int    `json:"maxAssessments"`
	CurrentAssessments   int    `json:"currentAssessments"`
	NewAssessmentCoolOff int64  `json:"newAssessmentCoolOff"`
	Messages             []string
}

// Host is a one-site report
type Host struct {
	Host            string
	Port            int
	Protocol        string
	IsPublic        bool `json:"isPublic"`
	Status          string
	StatusMessage   string   `json:"statusMessage"`
	StartTime       int64    `json:"startTime"`
	TestTime        int64    `json:"testTime"`
	EngineVersion   string   `json:"engineVersion"`
	CriteriaVersion string   `json:"criteriaVersion"`
	CacheExpiryTime int64    `json:"cacheExpiryTime"`
	CertHostnames   []string `json:"certHostnames"`
	Endpoints       []Endpoint
	Certs           []Cert `json:"certs,omitempty"`
}

// Endpoint is an Endpoint (IPv4, IPv6)
type Endpoint struct {
	IPAddress            string `json:"ipAddress"`
	ServerName           string `json:"serverName"`
	StatusMessage        string `json:"statusMessage"`
	StatusDetails        string `json:"statusDetails"`
	StatusDetailsMessage string `json:"statusDetailsMessage"`
	Grade                string
	GradeTrustIgnored    string `json:"gradeTrustIgnored"`
	FutureGrade          string `json:"futureGrade"`
	HasWarnings          bool   `json:"hasWarnings"`
	IsExceptional        bool   `json:"isExceptional"`
	Progress             int
	Duration             int
	Eta                  int
	Delegation           int
	Details              EndpointDetails `json:"details,omitempty"`
}

// EndpointDetails gives the details of a given Endpoint
type EndpointDetails struct {
	HostStartTime                  int64              `json:"hostStartTime"`
	CertChains                     []CertificateChain `json:"certChains"`
	Protocols                      []Protocol
	Suites                         []ProtocolSuites
	NoSniSuites                    ProtocolSuites `json:"noSniSuites"`
	NamedGroups                    NamedGroups    `json:"namedGroups"`
	ServerSignature                string         `json:"serverSignature"`
	PrefixDelegation               bool           `json:"prefixDelegation"`
	NonPrefixDelegation            bool           `json:"nonPrefixDelegation"`
	VulnBeast                      bool           `json:"vulnBeast"`
	RenegSupport                   int            `json:"renegSupport"`
	SessionResumption              int            `json:"sessionResumption"`
	CompressionMethods             int            `json:"compressionMethods"`
	SupportsNpn                    bool           `json:"supportsNpn"`
	NpnProcotols                   string         `json:"npnProtocols"`
	SupportsAlpn                   bool           `json:"supportsAlpn"`
	AlpnProtocols                  string
	SessionTickets                 int    `json:"sessionTickets"`
	OcspStapling                   bool   `json:"ocspStapling"`
	StaplingRevocationStatus       int    `json:"staplingRevocationStatus"`
	StaplingRevocationErrorMessage string `json:"staplingRevocationErrorMessage"`
	SniRequired                    bool   `json:"sniRequired"`
	HTTPStatusCode                 int    `json:"httpStatusCode"`
	HTTPForwarding                 string `json:"httpForwarding"`
	SupportsRC4                    bool   `json:"supportsRc4"`
	RC4WithModern                  bool   `json:"rc4WithModern"`
	RC4Only                        bool   `json:"rc4Only"`
	ForwardSecrecy                 int    `json:"forwardSecrecy"`
	ProtocolIntolerance            int    `json:"protocolIntolerance"`
	MiscIntolerance                int    `json:"miscIntolerance"`
	Sims                           SimDetails
	Heartbleed                     bool
	Heartbeat                      bool
	OpenSSLCcs                     int `json:"openSslCcs"`
	OpenSSLLuckyMinus20            int `json:"openSSLLuckyMinus20"`
	Ticketbleed                    int
	Bleichenbacher                 int
	Poodle                         bool
	PoodleTLS                      int  `json:"poodleTLS"`
	FallbackScsv                   bool `json:"fallbackScsv"`
	Freak                          bool
	HasSct                         int      `json:"hasSct"`
	DhPrimes                       []string `json:"dhPrimes"`
	DhUsesKnownPrimes              int      `json:"dhUsesKnownPrimes"`
	DhYsReuse                      bool     `json:"dhYsReuse"`
	EcdhParameterReuse             bool     `json:"ecdhParameterReuse"`
	Logjam                         bool
	ChaCha20Preference             bool
	HstsPolicy                     HstsPolicy        `json:"hstsPolicy"`
	HstsPreloads                   []HstsPreload     `json:"hstsPreloads"`
	HpkpPolicy                     HpkpPolicy        `json:"hpkpPolicy"`
	HpkpRoPolicy                   HpkpPolicy        `json:"hpkpRoPolicy"`
	StaticPkpPolicy                SPkpPolicy        `json:"staticPkpPolicy"`
	HTTPTransactions               []HTTPTransaction `json:"httpTransactions"`
	DrownHosts                     []DrownHost       `json:"drownHosts"`
	DrownErrors                    bool              `json:"drownErrors"`
	DrownVulnerable                bool              `json:"drownVulnerable"`
}

// CertificateChain is the list of certificates
type CertificateChain struct {
	ID         string
	CertIds    []string    `json:"certIds"`
	Trustpaths []TrustPath `json:"trustpaths"`
	Issues     int
	NoSni      bool `json:"noSni"`
}

// TrustPath defines the path of trust in cert chain
type TrustPath struct {
	CertIds       []string `json:"certIds"`
	Trust         []Trust  `json:"trust"`
	IsPinned      bool     `json:"isPinned"`
	MatchedPins   int      `json:"matchedPins"`
	UnMatchedPins int      `json:"unMatchedPins"`
}

// Trust identifies the cert store for trust
type Trust struct {
	RootStore         string `json:"rootStore"`
	IsTrusted         bool   `json:"isTrusted"`
	TrustErrorMessage string `json:"trustErrorMessage"`
}

// Protocol describes the HTTP protocols
type Protocol struct {
	ID               int `json:"id"`
	Name             string
	Version          string
	V2SuitesDisabled bool `json:"v2SuitesDisabled"`
	Q                int
}

// ProtocolSuites is a set of protocols
type ProtocolSuites struct {
	Protocol   int
	List       []Suite
	Preference bool
}

func (ls *ProtocolSuites) len() int {
	return len(ls.List)
}

// Suite describes a single protocol
type Suite struct {
	ID             int `json:"id"`
	Name           string
	CipherStrength int    `json:"cipherStrength"`
	KxType         string `json:"kxType"`
	KxStrength     int    `json:"kxStrength"`
	DHP            int    `json:"dhP"`
	DHG            int    `json:"dhG"`
	DHYs           int    `json:"dhYs"`
	NamedGroupBits int    `json:"namedGroupBits"`
	NamedGroupID   int    `json:"namedGroupId"`
	NamedGroudName string `json:"namedGroupName"`
	Q              int
}

// NamedGroups is for groups
type NamedGroups struct {
	List       []NamedGroup
	Preference bool
}

// NamedGroup is a group
type NamedGroup struct {
	ID   int
	Name string
	Bits int
}

// SimDetails are the result of simulation
type SimDetails struct {
	Results []Simulation
}

// Simulation describes the simulation of a given client
type Simulation struct {
	Client         SimClient
	ErrorCode      int    `json:"errorCode"`
	ErrorMessage   string `json:"errorMessage"`
	Attempts       int
	CertChainID    string `json:"certChainId"`
	ProtocolID     int    `json:"protocolId"`
	SuiteID        int    `json:"suiteId"`
	SuiteName      string `json:"suiteName"`
	KxType         string `json:"kxType"`
	KxStrength     int    `json:"kxStrength"`
	DhBits         int    `json:"dhBits"`
	DHP            int    `json:"dhP"`
	DHG            int    `json:"dhG"`
	DHYs           int    `json:"dhYs"`
	NamedGroupBits int    `json:"namedGroupBits"`
	NamedGroupID   int    `json:"namedGroupId"`
	NamedGroupName string `json:"namedGroupName"`
	AlertType      int    `json:"alertType"`
	AlertCode      int    `json:"alertCode"`
	KeyAlg         string `json:"keyAlg"`
	KeySize        int    `json:"keySize"`
	SigAlg         string `json:"sigAlg"`
}

// SimClient is a simulated client
type SimClient struct {
	ID          int `json:"id"`
	Name        string
	Platform    string
	Version     string
	IsReference bool `json:"isReference"`
}

// HstsPolicy describes the HSTS policy
type HstsPolicy struct {
	LongMaxAge        int64 `json:"LONG_MAX_AGE"`
	Header            string
	Status            string
	Error             string
	MaxAge            int64 `json:"maxAge"`
	IncludeSubDomains bool  `json:"includeSubDomains"`
	Preload           bool
	Directives        map[string]string
}

// HstsPreload is for HSTS preloading
type HstsPreload struct {
	Source     string
	HostName   string `json:"hostName"`
	Status     string
	Error      string
	SourceTime int64 `json:"sourceTime"`
}

// HpkpPolicy describes the HPKP policy
type HpkpPolicy struct {
	Header            string
	Status            string
	Error             string
	MaxAge            int64 `json:"maxAge"`
	IncludeSubDomains bool  `json:"includeSubDomains"`
	ReportURI         string
	Pins              []HpkpPin
	MatchedPins       []HpkpPin `json:"matchedPins"`
	Directives        []HpkpDirective
}

// SPkpPolicy descries the Static PkpPolicy
type SPkpPolicy struct {
	Status               string   `json:"status"`
	Error                string   `json:"error"`
	IncludeSubDomains    bool     `json:"includeSubDomains"`
	ReportURI            string   `json:"reportUri"`
	Pins                 []string `json:"pins"`
	MatchedPins          []string `json:"matchedPins"`
	ForbiddenPins        []string `json:"forbiddenPins"`
	MatchedForbiddenPins []string `json:"matchedForbiddenPins"`
}

// HTTPTransaction gives the entire request/response
type HTTPTransaction struct {
	RequestURL        string       `json:"requestUrl"`
	StatusCode        int          `json:"statusCode"`
	RequestLine       string       `json:"requestLine"`
	RequestHeaders    []string     `json:"requestHeaders"`
	ResponseLine      string       `json:"responseLine"`
	ResponseRawHeader []string     `json:"responseRawHeader"`
	ResponseHeader    []HTTPHeader `json:"responseHeader"`
	FragileServer     bool         `json:"fragileServer"`
}

// HTTPHeader is obvious
type HTTPHeader struct {
	Name  string
	Value string
}

// DrownHost describes a potentially Drown-weak site
type DrownHost struct {
	IP      string `json:"ip"`
	Export  bool
	Port    int
	Special bool
	SSLv2   bool `json:"sslv2"`
	Status  string
}

// Cert describes an X.509 certificate
type Cert struct {
	ID                     string
	Subject                string
	SerialNumber           string    `json:"serialNumber"`
	CommonNames            []string  `json:"commonNames"`
	AltNames               []string  `json:"altNames"`
	NotBefore              int64     `json:"notBefore"`
	NotAfter               int64     `json:"notAfter"`
	IssuerSubject          string    `json:"issuerSubject"`
	SigAlg                 string    `json:"sigAlg"`
	RevocationInfo         int       `json:"revocationInfo"`
	CrlURIs                []string  `json:"crlURIs"`
	OcspURIs               []string  `json:"ocspURIs"`
	RevocationStatus       int       `json:"revocationStatus"`
	CrlRevocationStatus    int       `json:"crlRevocationStatus"`
	OcspRevocationStatus   int       `json:"ocspRevocationStatus"`
	DNSCaa                 bool      `json:"dnsCaa"`
	CaaPolicy              CaaPolicy `json:"caaPolicy"`
	MustStaple             bool      `json:"mustStaple"`
	Sgc                    int
	ValidationType         string `json:"validationType"`
	Issues                 int
	Sct                    bool
	SHA1Hash               string `json:"sha1Hash"`
	SHA256Hash             string `json:"sha256Hash"`
	PinSHA256              string `json:"pinSha256"`
	KeyAlg                 string `json:"keyAlg"`
	KeySize                int    `json:"keySize"`
	KeyStrength            int    `json:"keyStrength"`
	KeyKnownDebianInsecure bool   `json:"keyKnownDebianInsecure"`
	Raw                    string `json:"raw"`
}

// CaaPolicy is the policy around CAA usage
type CaaPolicy struct {
	PolicyHostname string      `json:"policyHostname"`
	CaaRecords     []CaaRecord `json:"caaRecords"`
}

// CaaRecord describe the DNS CAA record content
type CaaRecord struct {
	Tag   string
	Value string
	Flags int
}

// HpkpPin is for pinned keys
type HpkpPin struct {
	HashFunction string `json:"hashFunction"`
	Value        string
}

// HpkpDirective is related to HPKP handling
type HpkpDirective struct {
	Name  string
	Value string
}

// Hosts is a shortcut to all Host
type Hosts []Host

// LabsResults are all the result of a run w/ 1 or more sites
type LabsResults struct {
	reports   []Host
	responses []string
}

// StatusCodes describes all possible status code & translations
type StatusCodes struct {
	StatusDetails map[string]string `json:"statusDetails"`
}

package puppetdb

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// PuppetDB stores informations used to connect to a PuppetDB
type PuppetDB struct {
	options *Options
	client  *http.Client
}

// Options contains the options used to connect to a PuppetDB
type Options struct {
	URL        string
	CertPath   string
	CACertPath string
	KeyPath    string
	SSLVerify  bool
}

// Node is a structure returned by a PuppetDB
type Node struct {
	Certname                string `json:"certname"`
	Deactivated             string `json:"deactivated"`
	LatestReportStatus      string `json:"latest_report_status"`
	ReportEnvironment       string `json:"report_environment"`
	ReportTimestamp         string `json:"report_timestamp"`
	LatestReportHash        string `json:"latest_report_hash"`
	FactsEnvironment        string `json:"facts_environment"`
	CachedCatalogStatus     string `json:"cached_catalog_status"`
	LatestReportNoop        bool   `json:"latest_report_noop"`
	Expired                 string `json:"expired"`
	LatestReportNoopPending bool   `json:"latest_report_noop_pending"`
	CatalogTimestamp        string `json:"catalog_timestamp"`
	FactsTimestamp          string `json:"facts_timestamp"`
}

// ReportMetric is a structure returned by a PuppetDB
type ReportMetric struct {
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Category string  `json:"category"`
}

// NewClient creates a new PuppetDB client
func NewClient(options *Options) (p *PuppetDB, err error) {
	var transport *http.Transport

	puppetdbURL, err := url.Parse(options.URL)
	if err != nil {
		err = fmt.Errorf("failed to parse PuppetDB URL: %v", err)
		return
	}

	if puppetdbURL.Scheme != "http" && puppetdbURL.Scheme != "https" {
		err = fmt.Errorf("%s is not a valid http scheme", puppetdbURL.Scheme)
		return
	}

	if puppetdbURL.Scheme == "https" {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(options.CertPath, options.KeyPath)
		if err != nil {
			err = fmt.Errorf("failed to load keypair: %s", err)
			return nil, err
		}

		// Load CA cert
		caCert, err := os.ReadFile(options.CACertPath)
		if err != nil {
			err = fmt.Errorf("failed to load ca certificate: %s", err)
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: !options.SSLVerify,
		}
		// BuildNameToCertificate is deprecated; leave nil to let Go select the first compatible certificate.
		transport = &http.Transport{TLSClientConfig: tlsConfig}
	} else {
		transport = &http.Transport{}
	}

	p = &PuppetDB{
		client:  &http.Client{Transport: transport},
		options: options,
	}
	return
}

// Nodes returns the list of nodes
func (p *PuppetDB) Nodes() (nodes []Node, err error) {
	// Use the full PuppetDB query endpoint
	err = p.get("/pdb/query/v4/nodes", "[\"or\", [\"=\", [\"node\", \"active\"], false], [\"=\", [\"node\", \"active\"], true]]", &nodes)
	if err != nil {
		err = fmt.Errorf("failed to get nodes: %s", err)
		return
	}
	return
}

// ReportMetrics returns the list of reportMetrics
func (p *PuppetDB) ReportMetrics(reportHash string) (reportMetrics []ReportMetric, err error) {
	err = p.get(fmt.Sprintf("/pdb/query/v4/reports/%s/metrics", reportHash), "", &reportMetrics)
	if err != nil {
		err = fmt.Errorf("failed to get reports: %s", err)
		return
	}
	return
}

// GetRaw performs a GET against the given endpoint and returns the raw response body.
// Endpoint should be a path like "/status/v1/services" or "/metrics/v2/list".
func (p *PuppetDB) GetRaw(endpoint string, query string) (body []byte, err error) {
	base := strings.TrimRight(p.options.URL, "/")
	var myurl string
	if strings.HasPrefix(endpoint, "/") {
		myurl = fmt.Sprintf("%s%s", base, endpoint)
	} else {
		myurl = fmt.Sprintf("%s/%s", base, endpoint)
	}
	if query != "" {
		myurl = fmt.Sprintf("%s?query=%s", myurl, url.QueryEscape(query))
	}
	req, err := http.NewRequest("GET", myurl, strings.NewReader(""))
	if err != nil {
		err = fmt.Errorf("failed to build request: %s", err)
		return
	}
	resp, err := p.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to call API: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response: %s", err)
		return
	}
	return
}

// MetricsList returns the raw JSON body from /metrics/v2/list
func (p *PuppetDB) MetricsList() (body []byte, err error) {
	return p.GetRaw("/metrics/v2/list", "")
}

// Metrics returns the raw JSON body from /metrics/v2
func (p *PuppetDB) Metrics() (body []byte, err error) {
	return p.GetRaw("/metrics/v2", "")
}

// Reports queries /pdb/query/v4/reports with an optional PuppetDB query string
func (p *PuppetDB) Reports(query string) (body []byte, err error) {
	return p.GetRaw("/pdb/query/v4/reports", query)
}

// MetricsV2Response models a subset of the /metrics/v2 response.
type MetricsV2Response struct {
	Request   map[string]interface{} `json:"request"`
	Value     map[string]interface{} `json:"value"`
	Timestamp int64                  `json:"timestamp"`
	Status    int                    `json:"status"`
}

// MetricsV2 fetches and parses /metrics/v2 into a MetricsV2Response
func (p *PuppetDB) MetricsV2() (MetricsV2Response, error) {
	var resp MetricsV2Response
	body, err := p.GetRaw("/metrics/v2", "")
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return resp, fmt.Errorf("failed to unmarshal metrics v2: %s", err)
	}
	return resp, nil
}

func (p *PuppetDB) get(endpoint string, query string, object interface{}) (err error) {
	// Build URL by appending the provided endpoint to the base URL.
	// The caller should pass endpoint paths such as:
	//   "/status/v1/services"
	//   "/metrics/v2/list"
	//   "/metrics/v2"
	//   "/pdb/query/v4/nodes"
	//   "/pdb/query/v4/reports"
	base := strings.TrimRight(p.options.URL, "/")
	var myurl string
	// Ensure endpoint is appended with a single '/'
	if strings.HasPrefix(endpoint, "/") {
		myurl = fmt.Sprintf("%s%s", base, endpoint)
	} else {
		myurl = fmt.Sprintf("%s/%s", base, endpoint)
	}
	if query != "" {
		myurl = fmt.Sprintf("%s?query=%s", myurl, url.QueryEscape(query))
	}
	req, err := http.NewRequest("GET", myurl, strings.NewReader(""))
	if err != nil {
		err = fmt.Errorf("failed to build request: %s", err)
		return
	}
	resp, err := p.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to call API: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response: %s", err)
		return
	}
	err = json.Unmarshal(body, object)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal: %s", err)
		return
	}
	return
}

// Post performs a POST request to the PuppetDB API
func (p *PuppetDB) Post(endpoint string, contentType string, body []byte) (resp *http.Response, err error) {
	base := strings.TrimRight(p.options.URL, "/")
	var myurl string
	if strings.HasPrefix(endpoint, "/") {
		myurl = fmt.Sprintf("%s%s", base, endpoint)
	} else {
		myurl = fmt.Sprintf("%s/%s", base, endpoint)
	}

	req, err := http.NewRequest("POST", myurl, strings.NewReader(string(body)))
	if err != nil {
		err = fmt.Errorf("failed to build POST request: %s", err)
		return
	}

	req.Header.Set("Content-Type", contentType)
	resp, err = p.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to call POST API: %s", err)
		return
	}
	return
}

// ServiceInfo represents the JSON structure returned by /status/v1/services for each service key.
type ServiceInfo struct {
	ServiceVersion       string                 `json:"service_version"`
	ServiceStatusVersion int                    `json:"service_status_version"`
	DetailLevel          string                 `json:"detail_level"`
	State                string                 `json:"state"`
	Status               map[string]interface{} `json:"status"`
	ActiveAlerts         []interface{}          `json:"active_alerts"`
}

// Services returns a typed map of services from /status/v1/services
func (p *PuppetDB) Services() (map[string]ServiceInfo, error) {
	body, err := p.GetRaw("/status/v1/services", "")
	if err != nil {
		return nil, err
	}
	var services map[string]ServiceInfo
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %s", err)
	}
	return services, nil
}

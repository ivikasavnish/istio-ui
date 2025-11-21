package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"` // admin, user
	CreatedAt time.Time `json:"created_at"`
}

// Snapshot represents a traffic configuration snapshot
type Snapshot struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Config      string    `json:"config"` // JSON serialized config
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}

// ScheduledAction represents a scheduled Istio action
type ScheduledAction struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	ActionType  string    `json:"action_type"` // traffic_shift, apply_mtls, remove_fault, etc.
	CronExpr    string    `json:"cron_expr"`
	Config      string    `json:"config"` // JSON serialized action config
	Enabled     bool      `json:"enabled"`
	NextRun     time.Time `json:"next_run"`
	LastRun     time.Time `json:"last_run,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// ServiceCache represents cached service information
type ServiceCache struct {
	ID          int       `json:"id"`
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Subsets     string    `json:"subsets"` // JSON array
	LastSeen    time.Time `json:"last_seen"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// VirtualServiceSpec represents a VirtualService configuration
type VirtualServiceSpec struct {
	Name      string                    `json:"name"`
	Namespace string                    `json:"namespace"`
	Hosts     []string                  `json:"hosts"`
	Gateways  []string                  `json:"gateways,omitempty"`
	HTTP      []HTTPRoute               `json:"http,omitempty"`
}

// HTTPRoute represents an HTTP route
type HTTPRoute struct {
	Name    string        `json:"name,omitempty"`
	Match   []HTTPMatch   `json:"match,omitempty"`
	Route   []HTTPDestination `json:"route"`
	Timeout string        `json:"timeout,omitempty"`
	Retries *HTTPRetry    `json:"retries,omitempty"`
	Fault   *HTTPFault    `json:"fault,omitempty"`
	Mirror  *Destination  `json:"mirror,omitempty"`
}

// HTTPMatch represents HTTP match conditions
type HTTPMatch struct {
	URI     *StringMatch          `json:"uri,omitempty"`
	Headers map[string]StringMatch `json:"headers,omitempty"`
}

// StringMatch represents string matching
type StringMatch struct {
	Exact  string `json:"exact,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

// HTTPDestination represents a destination with weight
type HTTPDestination struct {
	Destination Destination `json:"destination"`
	Weight      int         `json:"weight"`
}

// Destination represents a service destination
type Destination struct {
	Host   string `json:"host"`
	Subset string `json:"subset,omitempty"`
	Port   *Port  `json:"port,omitempty"`
}

// Port represents a port specification
type Port struct {
	Number int `json:"number"`
}

// HTTPRetry represents retry configuration
type HTTPRetry struct {
	Attempts      int    `json:"attempts"`
	PerTryTimeout string `json:"perTryTimeout"`
}

// HTTPFault represents fault injection
type HTTPFault struct {
	Delay *FaultDelay `json:"delay,omitempty"`
	Abort *FaultAbort `json:"abort,omitempty"`
}

// FaultDelay represents delay fault injection
type FaultDelay struct {
	Percentage *Percentage `json:"percentage,omitempty"`
	FixedDelay string      `json:"fixedDelay"`
}

// FaultAbort represents abort fault injection
type FaultAbort struct {
	Percentage *Percentage `json:"percentage,omitempty"`
	HTTPStatus int         `json:"httpStatus,omitempty"`
	GRPCStatus string      `json:"grpcStatus,omitempty"`
}

// Percentage represents a percentage value
type Percentage struct {
	Value float64 `json:"value"`
}

// DestinationRuleSpec represents a DestinationRule configuration
type DestinationRuleSpec struct {
	Name              string              `json:"name"`
	Namespace         string              `json:"namespace"`
	Host              string              `json:"host"`
	TrafficPolicy     *TrafficPolicy      `json:"trafficPolicy,omitempty"`
	Subsets           []Subset            `json:"subsets,omitempty"`
}

// TrafficPolicy represents traffic policy settings
type TrafficPolicy struct {
	ConnectionPool  *ConnectionPool  `json:"connectionPool,omitempty"`
	LoadBalancer    *LoadBalancer    `json:"loadBalancer,omitempty"`
	OutlierDetection *OutlierDetection `json:"outlierDetection,omitempty"`
	TLS             *TLSSettings     `json:"tls,omitempty"`
}

// ConnectionPool represents connection pool settings
type ConnectionPool struct {
	TCP  *TCPSettings  `json:"tcp,omitempty"`
	HTTP *HTTPSettings `json:"http,omitempty"`
}

// TCPSettings represents TCP settings
type TCPSettings struct {
	MaxConnections int `json:"maxConnections,omitempty"`
	ConnectTimeout string `json:"connectTimeout,omitempty"`
}

// HTTPSettings represents HTTP settings
type HTTPSettings struct {
	HTTP1MaxPendingRequests  int `json:"http1MaxPendingRequests,omitempty"`
	HTTP2MaxRequests         int `json:"http2MaxRequests,omitempty"`
	MaxRequestsPerConnection int `json:"maxRequestsPerConnection,omitempty"`
}

// LoadBalancer represents load balancer settings
type LoadBalancer struct {
	Simple string `json:"simple,omitempty"` // ROUND_ROBIN, LEAST_CONN, RANDOM
}

// OutlierDetection represents circuit breaker settings
type OutlierDetection struct {
	ConsecutiveErrors  int    `json:"consecutiveErrors,omitempty"`
	Interval           string `json:"interval,omitempty"`
	BaseEjectionTime   string `json:"baseEjectionTime,omitempty"`
	MaxEjectionPercent int    `json:"maxEjectionPercent,omitempty"`
}

// TLSSettings represents TLS settings
type TLSSettings struct {
	Mode string `json:"mode"` // DISABLE, SIMPLE, MUTUAL, ISTIO_MUTUAL
}

// Subset represents a service subset
type Subset struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

// GatewaySpec represents a Gateway configuration
type GatewaySpec struct {
	Name      string         `json:"name"`
	Namespace string         `json:"namespace"`
	Selector  map[string]string `json:"selector"`
	Servers   []Server       `json:"servers"`
}

// Server represents a gateway server
type Server struct {
	Port  ServerPort `json:"port"`
	Hosts []string   `json:"hosts"`
	TLS   *ServerTLS `json:"tls,omitempty"`
}

// ServerPort represents a server port
type ServerPort struct {
	Number   int    `json:"number"`
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
}

// ServerTLS represents server TLS configuration
type ServerTLS struct {
	Mode           string `json:"mode"` // SIMPLE, MUTUAL, PASSTHROUGH
	ServerCertificate string `json:"serverCertificate,omitempty"`
	PrivateKey     string `json:"privateKey,omitempty"`
	CACertificates string `json:"caCertificates,omitempty"`
}

// PeerAuthenticationSpec represents PeerAuthentication configuration
type PeerAuthenticationSpec struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Selector  map[string]string `json:"selector,omitempty"`
	MtlsMode  string            `json:"mtls_mode"` // STRICT, PERMISSIVE, DISABLE
}

// AuthorizationPolicySpec represents AuthorizationPolicy configuration
type AuthorizationPolicySpec struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Selector  map[string]string `json:"selector,omitempty"`
	Action    string            `json:"action"` // ALLOW, DENY
	Rules     []AuthRule        `json:"rules,omitempty"`
}

// AuthRule represents an authorization rule
type AuthRule struct {
	From []RuleFrom `json:"from,omitempty"`
	To   []RuleTo   `json:"to,omitempty"`
	When []Condition `json:"when,omitempty"`
}

// RuleFrom represents source of a rule
type RuleFrom struct {
	Source *Source `json:"source,omitempty"`
}

// Source represents a traffic source
type Source struct {
	Principals []string `json:"principals,omitempty"`
	Namespaces []string `json:"namespaces,omitempty"`
}

// RuleTo represents destination of a rule
type RuleTo struct {
	Operation *Operation `json:"operation,omitempty"`
}

// Operation represents an operation
type Operation struct {
	Methods []string `json:"methods,omitempty"`
	Paths   []string `json:"paths,omitempty"`
}

// Condition represents a condition
type Condition struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

// TrafficWeightUpdate represents a traffic weight update request
type TrafficWeightUpdate struct {
	Service   string          `json:"service"`
	Namespace string          `json:"namespace"`
	Weights   map[string]int  `json:"weights"` // version -> weight
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

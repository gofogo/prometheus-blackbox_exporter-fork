// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prober

// MetricDef describes a Prometheus metric for documentation purposes.
type MetricDef struct {
	// Name is the full metric name.
	Name string
	// Type is the Prometheus metric type (gauge, counter, gaugevec, …).
	Type string
	// Prober lists which probe type(s) expose this metric.
	// "all" = every prober, "-" = global/exporter-level metric.
	Prober string
	// Labels is the ordered list of label names, empty for label-less metrics.
	Labels []string
	// Help is the metric description exposed in /metrics.
	Help string
}

// MetricRegistry is the list of metrics registered for documentation.
// Append to this slice via registerMetricDef; do not write to it directly.
var MetricRegistry []MetricDef

func registerMetricDef(d MetricDef) {
	MetricRegistry = append(MetricRegistry, d)
}

func init() {
	// Shared DNS resolution metrics — defined as GaugeOpts in utils.go and
	// registered by chooseProtocol, which is called by every prober that
	// resolves a hostname (all except unix).
	registerMetricDef(MetricDef{
		Name:   probeDNSLookupTimeSecondsOpts.Name,
		Type:   "gauge",
		Prober: "dns, grpc, http, icmp, tcp, unix, websocket",
		Help:   probeDNSLookupTimeSecondsOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeIPProtocolGaugeOpts.Name,
		Type:   "gauge",
		Prober: "all",
		Help:   probeIPProtocolGaugeOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeIPAddrHashOpts.Name,
		Type:   "gauge",
		Prober: "all",
		Help:   probeIPAddrHashOpts.Help,
	})

	// Core probe metrics — defined as GaugeOpts in handler.go, present for
	// every prober regardless of the probe type.
	registerMetricDef(MetricDef{
		Name:   probeSuccessGaugeOpts.Name,
		Type:   "gauge",
		Prober: "all",
		Help:   probeSuccessGaugeOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeDurationGaugeOpts.Name,
		Type:   "gauge",
		Prober: "all",
		Help:   probeDurationGaugeOpts.Help,
	})

	// Shared TLS/SSL metrics — defined as GaugeOpts in prober.go and reused
	// by multiple probers. Names and help strings are sourced directly from
	// those opts so that changes propagate to generated docs automatically.

	registerMetricDef(MetricDef{
		Name:   sslEarliestCertExpiryGaugeOpts.Name,
		Type:   "gauge",
		Prober: "http, grpc, tcp, unix",
		Help:   sslEarliestCertExpiryGaugeOpts.Help,
	})

	registerMetricDef(MetricDef{
		Name:   sslChainExpiryInTimeStampGaugeOpts.Name,
		Type:   "gauge",
		Prober: "http, grpc, tcp, unix",
		Help:   sslChainExpiryInTimeStampGaugeOpts.Help,
	})

	registerMetricDef(MetricDef{
		Name:   probeTLSInfoGaugeOpts.Name,
		Type:   "gaugevec",
		Prober: "http, grpc, tcp, unix",
		Labels: []string{"version"},
		Help:   probeTLSInfoGaugeOpts.Help,
	})

	registerMetricDef(MetricDef{
		Name:   probeTLSCipherGaugeOpts.Name,
		Type:   "gaugevec",
		Prober: "http",
		Labels: []string{"cipher"},
		Help:   probeTLSCipherGaugeOpts.Help,
	})

	// probe_dns_serial is only registered when the DNS query type is SOA and a
	// successful response is received, so it cannot be captured by calling
	// ProbeDNS with a dummy target.
	registerMetricDef(MetricDef{
		Name:   "probe_dns_serial",
		Type:   "gauge",
		Prober: "dns",
		Help:   "Returns the serial number of the zone",
	})

	// probe_failed_due_to_cel is only registered when the module has CEL
	// expressions configured (FailIfBodyJsonMatchesCEL / FailIfBodyJsonNotMatchesCEL).
	registerMetricDef(MetricDef{
		Name:   "probe_failed_due_to_cel",
		Type:   "gauge",
		Prober: "http",
		Help:   "Indicates if probe failed due to CEL expression not matching",
	})

	// probe_http_last_modified_timestamp_seconds is only registered when the
	// HTTP response includes a parseable Last-Modified header.
	registerMetricDef(MetricDef{
		Name:   "probe_http_last_modified_timestamp_seconds",
		Type:   "gauge",
		Prober: "http",
		Help:   "Returns the Last-Modified HTTP response header in unixtime",
	})

	// probe_websocket_duration_seconds is registered early in ProbeWebsocket
	// but no phase label values are seeded before chooseProtocol fails on a
	// dummy target, making the label invisible to Gather().
	registerMetricDef(MetricDef{
		Name:   "probe_websocket_duration_seconds",
		Type:   "gaugevec",
		Prober: "websocket",
		Labels: []string{"phase"},
		Help:   "Duration of websocket request by phase",
	})

	// probe_failed_due_to_bytes is registered by probeQueryResponses (used by
	// ProbeTCP and ProbeUnix) whenever a query-response step is configured.
	registerMetricDef(MetricDef{
		Name:   "probe_failed_due_to_bytes",
		Type:   "gauge",
		Prober: "tcp, unix",
		Help:   "Indicates if probe failed due to bytes",
	})

	// probe_expect_info is registered by probeQueryResponses only when a
	// query-response entry has labels configured. Its label names are
	// user-defined in the module config (QueryResponse.Labels), so they
	// cannot be enumerated statically.
	registerMetricDef(MetricDef{
		Name:   "probe_expect_info",
		Type:   "gaugevec",
		Prober: "tcp, unix",
		Labels: []string{"<user-defined>"},
		Help:   "Explicit content matched",
	})

	// probe_icmp_reply_hop_limit is only registered when an ICMP reply is
	// received and the hop limit (TTL) is retrievable from the control message.
	registerMetricDef(MetricDef{
		Name:   "probe_icmp_reply_hop_limit",
		Type:   "gauge",
		Prober: "icmp",
		Help:   "Replied packet hop limit (TTL for ipv4)",
	})

	// probe_grpc_healthcheck_response is registered early in ProbeGRPC but its
	// label values (serving_status) are only populated after a live gRPC call,
	// so Gather() returns the family with no series and the label is invisible.
	registerMetricDef(MetricDef{
		Name:   "probe_grpc_healthcheck_response",
		Type:   "gaugevec",
		Prober: "grpc",
		Labels: []string{"serving_status"},
		Help:   "Response HealthCheck response",
	})

	// probe_ssl_last_chain_info is a GaugeVec whose labels are only populated
	// when a TLS handshake succeeds, making it unreachable via a dummy call.
	registerMetricDef(MetricDef{
		Name:   "probe_ssl_last_chain_info",
		Type:   "gaugevec",
		Prober: "grpc, http, tcp, unix",
		Labels: []string{"fingerprint_sha256", "issuer", "serialnumber", "subject", "subjectalternative"},
		Help:   "Contains SSL leaf certificate information",
	})
}

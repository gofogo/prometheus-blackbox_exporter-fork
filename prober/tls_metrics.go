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

import "github.com/prometheus/client_golang/prometheus"

// TLS opts are package-level so they can be referenced by grpc.go, http.go,
// query_response.go (tcp/unix) which all share the same TLS metrics.
var (
	sslEarliestCertExpiryGaugeOpts = prometheus.GaugeOpts{
		Name: "probe_ssl_earliest_cert_expiry",
		Help: "Returns last SSL chain expiry in unixtime",
	}
	sslChainExpiryInTimeStampGaugeOpts = prometheus.GaugeOpts{
		Name: "probe_ssl_last_chain_expiry_timestamp_seconds",
		Help: "Returns last SSL chain expiry in timestamp",
	}
	probeTLSInfoGaugeOpts = prometheus.GaugeOpts{
		Name: "probe_tls_version_info",
		Help: "Returns the TLS version used or NaN when unknown",
	}
	probeTLSCipherGaugeOpts = prometheus.GaugeOpts{
		Name: "probe_tls_cipher_info",
		Help: "Returns the TLS cipher negotiated during handshake",
	}
	probeSSLLastChainInfoOpts = prometheus.GaugeOpts{
		Name: "probe_ssl_last_chain_info",
		Help: "Contains SSL leaf certificate information",
	}
)

func init() {
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
	registerMetricDef(MetricDef{
		Name:   probeSSLLastChainInfoOpts.Name,
		Type:   "gaugevec",
		Prober: "grpc, http, tcp, unix",
		Labels: []string{"fingerprint_sha256", "issuer", "serialnumber", "subject", "subjectalternative"},
		Help:   probeSSLLastChainInfoOpts.Help,
	})
}

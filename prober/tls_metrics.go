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

// TLS specs are package-level so grpc.go, http.go, and query_response.go
// (tcp/unix) all share the same single definition.
var (
	SSLEarliestCertExpiryGaugeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_ssl_earliest_cert_expiry", Help: "Returns last SSL chain expiry in unixtime"},
		Prober: "http, grpc, tcp, unix",
	}
	SSLChainExpiryInTimeStampGaugeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_ssl_last_chain_expiry_timestamp_seconds", Help: "Returns last SSL chain expiry in timestamp"},
		Prober: "http, grpc, tcp, unix",
	}
	ProbeTLSInfoGaugeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_tls_version_info", Help: "Returns the TLS version used or NaN when unknown"},
		Labels: []string{"version"},
		Prober: "http, grpc, tcp, unix",
	}
	ProbeTLSCipherGaugeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_tls_cipher_info", Help: "Returns the TLS cipher negotiated during handshake"},
		Labels: []string{"cipher"},
		Prober: "http",
	}
	ProbeSSLLastChainInfoSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_ssl_last_chain_info", Help: "Contains SSL leaf certificate information"},
		Labels: []string{"fingerprint_sha256", "subject", "issuer", "subjectalternative", "serialnumber"},
		Prober: "grpc, http, tcp, unix",
	}
)

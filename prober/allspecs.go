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

// AllSpecs is the canonical list of all MetricSpecs exposed by the exporter's
// probe endpoint. The documentation generator reads this slice directly.
// Adding a new metric: define a MetricSpec var in the relevant *_metrics.go
// file, then add it here.
var AllSpecs = []MetricSpec{
	// All probers
	ProbeSuccessGaugeSpec,
	ProbeDurationGaugeSpec,
	ProbeDNSLookupTimeSecondsSpec,
	ProbeIPProtocolGaugeSpec,
	ProbeIPAddrHashSpec,

	// DNS
	ProbeDNSDurationGaugeVecSpec,
	ProbeDNSAnswerRRSSpec,
	ProbeDNSAuthorityRRSSpec,
	ProbeDNSAdditionalRRSSpec,
	ProbeDNSQuerySucceededSpec,
	ProbeDNSSerialSpec,

	// gRPC
	ProbeGRPCDurationGaugeVecSpec,
	ProbeGRPCSSLSpec,
	ProbeGRPCStatusCodeSpec,
	ProbeGRPCHealthCheckResponseSpec,

	// HTTP
	ProbeHTTPDurationGaugeVecSpec,
	ProbeHTTPContentLengthSpec,
	ProbeHTTPUncompressedBodyLengthSpec,
	ProbeHTTPRedirectsSpec,
	ProbeHTTPSSLSpec,
	ProbeHTTPStatusCodeSpec,
	ProbeHTTPVersionSpec,
	ProbeFailedDueToCELSpec,
	ProbeHTTPLastModifiedSpec,

	// ICMP
	ProbeICMPDurationGaugeVecSpec,
	ProbeICMPReplyHopLimitSpec,

	// TLS (shared)
	SSLEarliestCertExpiryGaugeSpec,
	SSLChainExpiryInTimeStampGaugeSpec,
	ProbeTLSInfoGaugeSpec,
	ProbeTLSCipherGaugeSpec,
	ProbeSSLLastChainInfoSpec,

	// TCP / Unix query-response
	ProbeFailedDueToRegexSpec,
	ProbeFailedDueToBytesSpec,
	ProbeExpectInfoSpec,

	// WebSocket
	ProbeWebsocketStatusCodeSpec,
	ProbeWebsocketConnectionUpgradedSpec,
	ProbeWebsocketFailedDueToRegexSpec,
	ProbeWebsocketDurationGaugeVecSpec,
}

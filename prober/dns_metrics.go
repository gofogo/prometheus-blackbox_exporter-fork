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

var (
	probeDNSDurationGaugeVecOpts = prometheus.GaugeOpts{
		Name: "probe_dns_duration_seconds",
		Help: "Duration of DNS request by phase",
	}
	probeDNSAnswerRRSOpts = prometheus.GaugeOpts{
		Name: "probe_dns_answer_rrs",
		Help: "Returns number of entries in the answer resource record list",
	}
	probeDNSAuthorityRRSOpts = prometheus.GaugeOpts{
		Name: "probe_dns_authority_rrs",
		Help: "Returns number of entries in the authority resource record list",
	}
	probeDNSAdditionalRRSOpts = prometheus.GaugeOpts{
		Name: "probe_dns_additional_rrs",
		Help: "Returns number of entries in the additional resource record list",
	}
	probeDNSQuerySucceededOpts = prometheus.GaugeOpts{
		Name: "probe_dns_query_succeeded",
		Help: "Displays whether or not the query was executed successfully",
	}
	probeDNSSerialOpts = prometheus.GaugeOpts{
		Name: "probe_dns_serial",
		Help: "Returns the serial number of the zone",
	}
)

func init() {
	registerMetricDef(MetricDef{
		Name:   probeDNSDurationGaugeVecOpts.Name,
		Type:   "gaugevec",
		Prober: "dns",
		Labels: []string{"phase"},
		Help:   probeDNSDurationGaugeVecOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeDNSAnswerRRSOpts.Name,
		Type:   "gauge",
		Prober: "dns",
		Help:   probeDNSAnswerRRSOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeDNSAuthorityRRSOpts.Name,
		Type:   "gauge",
		Prober: "dns",
		Help:   probeDNSAuthorityRRSOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeDNSAdditionalRRSOpts.Name,
		Type:   "gauge",
		Prober: "dns",
		Help:   probeDNSAdditionalRRSOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeDNSQuerySucceededOpts.Name,
		Type:   "gauge",
		Prober: "dns",
		Help:   probeDNSQuerySucceededOpts.Help,
	})
	// probe_dns_serial is only registered when the query type is SOA and a
	// successful response is received.
	registerMetricDef(MetricDef{
		Name:   probeDNSSerialOpts.Name,
		Type:   "gauge",
		Prober: "dns",
		Help:   probeDNSSerialOpts.Help,
	})
}

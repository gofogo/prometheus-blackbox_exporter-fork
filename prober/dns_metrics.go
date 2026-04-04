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
	ProbeDNSDurationGaugeVecSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_duration_seconds", Help: "Duration of DNS request by phase"},
		Labels: []string{"phase"},
		Prober: "dns",
	}
	ProbeDNSAnswerRRSSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_answer_rrs", Help: "Returns number of entries in the answer resource record list"},
		Prober: "dns",
	}
	ProbeDNSAuthorityRRSSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_authority_rrs", Help: "Returns number of entries in the authority resource record list"},
		Prober: "dns",
	}
	ProbeDNSAdditionalRRSSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_additional_rrs", Help: "Returns number of entries in the additional resource record list"},
		Prober: "dns",
	}
	ProbeDNSQuerySucceededSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_query_succeeded", Help: "Displays whether or not the query was executed successfully"},
		Prober: "dns",
	}
	ProbeDNSSerialSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_serial", Help: "Returns the serial number of the zone"},
		Prober: "dns",
	}
)

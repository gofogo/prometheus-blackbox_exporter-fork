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
	// probeFailedDueToRegexOpts is shared by http.go and query_response.go
	// (used by ProbeTCP and ProbeUnix).
	probeFailedDueToRegexOpts = prometheus.GaugeOpts{
		Name: "probe_failed_due_to_regex",
		Help: "Indicates if probe failed due to regex",
	}
	probeFailedDueToBytesOpts = prometheus.GaugeOpts{
		Name: "probe_failed_due_to_bytes",
		Help: "Indicates if probe failed due to bytes",
	}
	probeExpectInfoOpts = prometheus.GaugeOpts{
		Name: "probe_expect_info",
		Help: "Explicit content matched",
	}
)

func init() {
	registerMetricDef(MetricDef{
		Name:   probeFailedDueToRegexOpts.Name,
		Type:   "gauge",
		Prober: "http, tcp, unix",
		Help:   probeFailedDueToRegexOpts.Help,
	})
	// probe_failed_due_to_bytes is registered when a query-response step is configured.
	registerMetricDef(MetricDef{
		Name:   probeFailedDueToBytesOpts.Name,
		Type:   "gauge",
		Prober: "tcp, unix",
		Help:   probeFailedDueToBytesOpts.Help,
	})
	// probe_expect_info is registered when a query-response entry has labels
	// configured. Its label names are user-defined in the module config.
	registerMetricDef(MetricDef{
		Name:   probeExpectInfoOpts.Name,
		Type:   "gaugevec",
		Prober: "tcp, unix",
		Labels: []string{"<user-defined>"},
		Help:   probeExpectInfoOpts.Help,
	})
}

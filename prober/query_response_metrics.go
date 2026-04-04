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
	// ProbeFailedDueToRegexSpec is shared by http.go and query_response.go
	// (used by ProbeTCP and ProbeUnix).
	ProbeFailedDueToRegexSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_failed_due_to_regex", Help: "Indicates if probe failed due to regex"},
		Prober: "http, tcp, unix",
	}
	ProbeFailedDueToBytesSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_failed_due_to_bytes", Help: "Indicates if probe failed due to bytes"},
		Prober: "tcp, unix",
	}
	// ProbeExpectInfoSpec uses dynamic (user-defined) labels at runtime.
	// Labels here is a documentation-only placeholder; production code calls
	// prometheus.NewGaugeVec(ProbeExpectInfoSpec.Opts, dynamicNames) directly.
	ProbeExpectInfoSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_expect_info", Help: "Explicit content matched"},
		Labels: []string{"<user-defined>"},
		Prober: "tcp, unix",
	}
)

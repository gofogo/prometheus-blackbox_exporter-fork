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

// MetricSpec is the single source of truth for a metric's name, help, variable
// labels, metric type, and which prober(s) expose it. Both production code and
// the documentation generator use this struct — labels are never duplicated.
type MetricSpec struct {
	// Opts carries Name, Help (and optionally Namespace/Subsystem).
	Opts prometheus.GaugeOpts
	// Labels is the ordered list of variable label names.
	// nil / empty means a plain Gauge; non-empty means a GaugeVec.
	Labels []string
	// Prober lists which probe type(s) expose this metric.
	// "all" = every prober, "-" = global/exporter-level metric.
	Prober string
}

// Type returns "gaugevec" when Labels are present, "gauge" otherwise.
func (s MetricSpec) Type() string {
	if len(s.Labels) > 0 {
		return "gaugevec"
	}
	return "gauge"
}

// NewGauge creates a Gauge from the spec. Use for specs with no Labels.
func (s MetricSpec) NewGauge() prometheus.Gauge {
	return prometheus.NewGauge(s.Opts)
}

// NewGaugeVec creates a GaugeVec from the spec using the spec's Labels.
// Use for specs with Labels defined.
func (s MetricSpec) NewGaugeVec() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(s.Opts, s.Labels)
}

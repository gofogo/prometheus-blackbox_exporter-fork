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
	// Type is the Prometheus metric type: gauge, gaugevec, counter, …
	Type string
	// Prober lists which probe type(s) expose this metric.
	// "all" = every prober, "-" = global/exporter-level metric.
	Prober string
	// Labels is the ordered list of label names; empty for label-less metrics.
	Labels []string
	// Help is the metric description exposed in /metrics.
	Help string
}

// MetricRegistry is the canonical list of all metrics exposed by the exporter.
// Each prober file registers its metrics via registerMetricDef in an init().
var MetricRegistry []MetricDef

func registerMetricDef(d MetricDef) {
	MetricRegistry = append(MetricRegistry, d)
}

func init() {
	// Config / exporter-level metrics — defined in config/config.go.
	registerMetricDef(MetricDef{
		Name:   "blackbox_exporter_config_last_reload_successful",
		Type:   "gauge",
		Prober: "-",
		Help:   "Blackbox exporter config loaded successfully.",
	})
	registerMetricDef(MetricDef{
		Name:   "blackbox_exporter_config_last_reload_success_timestamp_seconds",
		Type:   "gauge",
		Prober: "-",
		Help:   "Timestamp of the last successful configuration reload.",
	})
}

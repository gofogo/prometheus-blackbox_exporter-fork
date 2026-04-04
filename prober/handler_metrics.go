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
	probeSuccessGaugeOpts = prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Displays whether or not the probe was a success",
	}
	probeDurationGaugeOpts = prometheus.GaugeOpts{
		Name: "probe_duration_seconds",
		Help: "Returns how long the probe took to complete in seconds",
	}
)

func init() {
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
}

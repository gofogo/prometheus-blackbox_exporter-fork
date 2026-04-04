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
	probeHTTPDurationGaugeVecOpts = prometheus.GaugeOpts{
		Name: "probe_http_duration_seconds",
		Help: "Duration of http request by phase, summed over all redirects",
	}
	probeHTTPContentLengthOpts = prometheus.GaugeOpts{
		Name: "probe_http_content_length",
		Help: "Length of http content response",
	}
	probeHTTPUncompressedBodyLengthOpts = prometheus.GaugeOpts{
		Name: "probe_http_uncompressed_body_length",
		Help: "Length of uncompressed response body",
	}
	probeHTTPRedirectsOpts = prometheus.GaugeOpts{
		Name: "probe_http_redirects",
		Help: "The number of redirects",
	}
	probeHTTPSSLOpts = prometheus.GaugeOpts{
		Name: "probe_http_ssl",
		Help: "Indicates if SSL was used for the final redirect",
	}
	probeHTTPStatusCodeOpts = prometheus.GaugeOpts{
		Name: "probe_http_status_code",
		Help: "Response HTTP status code",
	}
	probeHTTPVersionOpts = prometheus.GaugeOpts{
		Name: "probe_http_version",
		Help: "Returns the version of HTTP of the probe response",
	}
	probeFailedDueToCELOpts = prometheus.GaugeOpts{
		Name: "probe_failed_due_to_cel",
		Help: "Indicates if probe failed due to CEL expression not matching",
	}
	probeHTTPLastModifiedOpts = prometheus.GaugeOpts{
		Name: "probe_http_last_modified_timestamp_seconds",
		Help: "Returns the Last-Modified HTTP response header in unixtime",
	}
)

func init() {
	registerMetricDef(MetricDef{
		Name:   probeHTTPDurationGaugeVecOpts.Name,
		Type:   "gaugevec",
		Prober: "http",
		Labels: []string{"phase"},
		Help:   probeHTTPDurationGaugeVecOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeHTTPContentLengthOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPContentLengthOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeHTTPUncompressedBodyLengthOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPUncompressedBodyLengthOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeHTTPRedirectsOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPRedirectsOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeHTTPSSLOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPSSLOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeHTTPStatusCodeOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPStatusCodeOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeHTTPVersionOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPVersionOpts.Help,
	})
	// probe_failed_due_to_cel is only registered when the module has CEL expressions configured.
	registerMetricDef(MetricDef{
		Name:   probeFailedDueToCELOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeFailedDueToCELOpts.Help,
	})
	// probe_http_last_modified_timestamp_seconds is only registered when the
	// HTTP response includes a parseable Last-Modified header.
	registerMetricDef(MetricDef{
		Name:   probeHTTPLastModifiedOpts.Name,
		Type:   "gauge",
		Prober: "http",
		Help:   probeHTTPLastModifiedOpts.Help,
	})
}

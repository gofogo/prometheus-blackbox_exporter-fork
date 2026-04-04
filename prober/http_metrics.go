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
	ProbeHTTPDurationGaugeVecSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_duration_seconds", Help: "Duration of http request by phase, summed over all redirects"},
		Labels: []string{"phase"},
		Prober: "http",
	}
	ProbeHTTPContentLengthSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_content_length", Help: "Length of http content response"},
		Prober: "http",
	}
	ProbeHTTPUncompressedBodyLengthSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_uncompressed_body_length", Help: "Length of uncompressed response body"},
		Prober: "http",
	}
	ProbeHTTPRedirectsSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_redirects", Help: "The number of redirects"},
		Prober: "http",
	}
	ProbeHTTPSSLSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_ssl", Help: "Indicates if SSL was used for the final redirect"},
		Prober: "http",
	}
	ProbeHTTPStatusCodeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_status_code", Help: "Response HTTP status code"},
		Prober: "http",
	}
	ProbeHTTPVersionSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_version", Help: "Returns the version of HTTP of the probe response"},
		Prober: "http",
	}
	ProbeFailedDueToCELSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_failed_due_to_cel", Help: "Indicates if probe failed due to CEL expression not matching"},
		Prober: "http",
	}
	ProbeHTTPLastModifiedSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_http_last_modified_timestamp_seconds", Help: "Returns the Last-Modified HTTP response header in unixtime"},
		Prober: "http",
	}
)

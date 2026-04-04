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
	probeGRPCDurationGaugeVecOpts = prometheus.GaugeOpts{
		Name: "probe_grpc_duration_seconds",
		Help: "Duration of gRPC request by phase",
	}
	probeGRPCSSLOpts = prometheus.GaugeOpts{
		Name: "probe_grpc_ssl",
		Help: "Indicates if SSL was used for the connection",
	}
	probeGRPCStatusCodeOpts = prometheus.GaugeOpts{
		Name: "probe_grpc_status_code",
		Help: "Response gRPC status code",
	}
	probeGRPCHealthCheckResponseOpts = prometheus.GaugeOpts{
		Name: "probe_grpc_healthcheck_response",
		Help: "Response HealthCheck response",
	}
)

func init() {
	registerMetricDef(MetricDef{
		Name:   probeGRPCDurationGaugeVecOpts.Name,
		Type:   "gaugevec",
		Prober: "grpc",
		Labels: []string{"phase"},
		Help:   probeGRPCDurationGaugeVecOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeGRPCSSLOpts.Name,
		Type:   "gauge",
		Prober: "grpc",
		Help:   probeGRPCSSLOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeGRPCStatusCodeOpts.Name,
		Type:   "gauge",
		Prober: "grpc",
		Help:   probeGRPCStatusCodeOpts.Help,
	})
	registerMetricDef(MetricDef{
		Name:   probeGRPCHealthCheckResponseOpts.Name,
		Type:   "gaugevec",
		Prober: "grpc",
		Labels: []string{"serving_status"},
		Help:   probeGRPCHealthCheckResponseOpts.Help,
	})
}

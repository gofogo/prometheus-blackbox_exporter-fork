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
	ProbeWebsocketStatusCodeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_websocket_status_code", Help: "Response HTTP status code"},
		Prober: "websocket",
	}
	ProbeWebsocketConnectionUpgradedSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_websocket_connection_upgraded", Help: "Indicates if the websocket connection was successfully upgraded"},
		Prober: "websocket",
	}
	ProbeWebsocketFailedDueToRegexSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_websocket_failed_due_to_regex", Help: "Indicates if probe failed due to regex"},
		Prober: "websocket",
	}
	ProbeWebsocketDurationGaugeVecSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_websocket_duration_seconds", Help: "Duration of websocket request by phase"},
		Labels: []string{"phase"},
		Prober: "websocket",
	}
)

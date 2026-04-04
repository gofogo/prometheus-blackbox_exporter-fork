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
	ProbeDNSLookupTimeSecondsSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_dns_lookup_time_seconds", Help: "Returns the time taken for probe dns lookup in seconds"},
		Prober: "dns, grpc, http, icmp, tcp, unix, websocket",
	}
	ProbeIPProtocolGaugeSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_ip_protocol", Help: "Specifies whether probe ip protocol is IP4 or IP6"},
		Prober: "all",
	}
	ProbeIPAddrHashSpec = MetricSpec{
		Opts:   prometheus.GaugeOpts{Name: "probe_ip_addr_hash", Help: "Specifies the hash of IP address. It's useful to detect if the IP address changes."},
		Prober: "all",
	}
)

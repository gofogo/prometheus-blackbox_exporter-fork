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

// main generates docs/monitoring/metrics.md.
// Run via: go run ./internal/gen/docs/metrics
// Or via:  make generate-metrics-documentation
package main

import (
	"context"
	"embed"
	"io"
	"log"
	"log/slog"
	"sort"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/blackbox_exporter/internal/gen/docs/render"
	"github.com/prometheus/blackbox_exporter/prober"
)

//go:embed templates
var templateFS embed.FS

// MetricDef describes a single Prometheus metric for documentation.
type MetricDef struct {
	Name   string
	Type   string
	Prober string
	Labels []string
	Help   string
}

// proberOverrides maps metric names to a Prober value that overrides what
// gatherFromRegistry would infer. Used for metrics shared across probers that
// cannot be expressed via a single gather call.
var proberOverrides = map[string]string{
	"probe_failed_due_to_regex": "http, tcp, unix",
}

// ColWidths holds computed column widths for the markdown table.
type ColWidths struct {
	Name   int
	Type   int
	Prober int
	Labels int
	Help   int
}

// TemplateData is passed to the template engine.
type TemplateData struct {
	Metrics   []MetricDef
	ColWidths ColWidths
}

// discardLogger is used when running probe functions during doc generation
// to suppress diagnostic output that would pollute the generator's stderr.
var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// gatherFromRegistry calls construct with a fresh Prometheus registry, gathers
// all registered metric families from it, and returns them as MetricDefs.
// defaultProber is written into MetricDef.Prober unless proberOverrides has
// an entry for that metric name.
func gatherFromRegistry(defaultProber string, construct func(*prometheus.Registry)) ([]MetricDef, error) {
	reg := prometheus.NewRegistry()
	construct(reg)

	families, err := reg.Gather()
	if err != nil {
		return nil, err
	}

	defs := make([]MetricDef, 0, len(families))
	for _, f := range families {
		p := defaultProber
		if override, ok := proberOverrides[f.GetName()]; ok {
			p = override
		}
		defs = append(defs, MetricDef{
			Name:   f.GetName(),
			Type:   metricTypeName(f.GetType()),
			Prober: p,
			Labels: labelNamesFromFamily(f),
			Help:   f.GetHelp(),
		})
	}
	return defs, nil
}

// metricTypeName converts a dto.MetricType to a lowercase documentation string.
func metricTypeName(t dto.MetricType) string {
	return strings.ToLower(t.String())
}

// labelNamesFromFamily returns the label names from the first metric in the
// family. Sufficient for documentation because all series share the same schema.
func labelNamesFromFamily(f *dto.MetricFamily) []string {
	if len(f.GetMetric()) == 0 {
		return nil
	}
	pairs := f.GetMetric()[0].GetLabel()
	names := make([]string, 0, len(pairs))
	for _, lp := range pairs {
		names = append(names, lp.GetName())
	}
	return names
}

// fromProberRegistry converts prober.MetricRegistry entries to MetricDef.
// Used for metrics whose labels are invisible to Gather() without observed
// values (GaugeVec), or that are only registered on a conditional code path.
func fromProberRegistry() []MetricDef {
	defs := make([]MetricDef, 0, len(prober.MetricRegistry))
	for _, m := range prober.MetricRegistry {
		defs = append(defs, MetricDef{
			Name:   m.Name,
			Type:   m.Type,
			Prober: m.Prober,
			Labels: m.Labels,
			Help:   m.Help,
		})
	}
	return defs
}

// deduplicateByName returns metrics with duplicate names removed, keeping the
// first occurrence. Used to avoid double-listing shared helper metrics that
// appear in multiple gather calls.
func deduplicateByName(metrics []MetricDef) []MetricDef {
	seen := make(map[string]bool, len(metrics))
	out := make([]MetricDef, 0, len(metrics))
	for _, m := range metrics {
		if !seen[m.Name] {
			seen[m.Name] = true
			out = append(out, m)
		}
	}
	return out
}

func sortMetrics(m []MetricDef) {
	sort.Slice(m, func(i, j int) bool {
		if m[i].Prober != m[j].Prober {
			return m[i].Prober < m[j].Prober
		}
		return m[i].Name < m[j].Name
	})
}

func computeColWidths(m []MetricDef) ColWidths {
	return ColWidths{
		Name:   render.MapColumn("Name", m, func(x MetricDef) string { return x.Name }),
		Type:   render.MapColumn("Metric Type", m, func(x MetricDef) string { return x.Type }),
		Prober: render.MapColumn("Prober", m, func(x MetricDef) string { return x.Prober }),
		Labels: render.MapColumn("Labels", m, func(x MetricDef) string { return strings.Join(x.Labels, ", ") }),
		Help:   render.MapColumn("Help", m, func(x MetricDef) string { return x.Help }),
	}
}

func main() {
	var all []MetricDef

	// Metrics from prober.MetricRegistry are appended first so that shared
	// metrics (chooseProtocol, handler) carry their correct Prober annotation
	// and win deduplication over whatever the per-prober gather calls produce.
	all = append(all, fromProberRegistry()...)

	// Config metrics: call the real constructor, gather whatever it registers.
	configDefs, err := gatherFromRegistry("-", func(reg *prometheus.Registry) {
		config.NewSafeConfig(reg)
	})
	if err != nil {
		log.Fatalf("gather config metrics: %v", err)
	}
	all = append(all, configDefs...)

	// DNS metrics: call ProbeDNS with a dummy target. All metrics are registered
	// before any network I/O, so the call failing is expected and harmless.
	// probe_dns_serial is conditional on a successful SOA response and is handled
	// separately via prober.MetricRegistry.
	dnsDefs, err := gatherFromRegistry("dns", func(reg *prometheus.Registry) {
		prober.ProbeDNS(context.Background(), "dummy.invalid", config.Module{}, reg, discardLogger)
	})
	if err != nil {
		log.Fatalf("gather dns metrics: %v", err)
	}
	all = append(all, dnsDefs...)

	// gRPC metrics: the 4 core metrics are registered before chooseProtocol is
	// called, so they appear in Gather() even when the dummy target fails.
	// probe_grpc_healthcheck_response labels and TLS-conditional metrics are
	// handled via prober.MetricRegistry.
	grpcDefs, err := gatherFromRegistry("grpc", func(reg *prometheus.Registry) {
		prober.ProbeGRPC(context.Background(), "dummy.invalid", config.Module{}, reg, discardLogger)
	})
	if err != nil {
		log.Fatalf("gather grpc metrics: %v", err)
	}
	all = append(all, grpcDefs...)

	// ICMP metrics: durationGaugeVec is seeded with all three phases (resolve,
	// setup, rtt) before chooseProtocol is called, so the label is fully
	// visible to Gather() even when the dummy target fails DNS resolution.
	// probe_icmp_reply_hop_limit is conditional on a successful reply and is
	// handled via prober.MetricRegistry.
	icmpDefs, err := gatherFromRegistry("icmp", func(reg *prometheus.Registry) {
		prober.ProbeICMP(context.Background(), "dummy.invalid", config.Module{}, reg, discardLogger)
	})
	if err != nil {
		log.Fatalf("gather icmp metrics: %v", err)
	}
	all = append(all, icmpDefs...)

	// HTTP metrics: 8 core metrics are always registered. With an empty module
	// (no proxy), shouldResolveDNSWithProxy is true so chooseProtocol runs and
	// seeds durationGaugeVec with the "resolve" phase label before failing —
	// making the label name visible to Gather(). Conditional metrics (CEL,
	// Last-Modified, TLS) are handled via prober.MetricRegistry.
	httpDefs, err := gatherFromRegistry("http", func(reg *prometheus.Registry) {
		prober.ProbeHTTP(context.Background(), "http://dummy.invalid", config.Module{}, reg, discardLogger)
	})
	if err != nil {
		log.Fatalf("gather http metrics: %v", err)
	}
	all = append(all, httpDefs...)

	// WebSocket metrics: 3 plain Gauges are registered before chooseProtocol is
	// called and are gatherable. probe_websocket_duration_seconds is registered
	// but its phase label is never seeded before DNS fails on the dummy target,
	// so it is handled via prober.MetricRegistry.
	wsDefs, err := gatherFromRegistry("websocket", func(reg *prometheus.Registry) {
		prober.ProbeWebsocket(context.Background(), "ws://dummy.invalid", config.Module{}, reg, discardLogger)
	})
	if err != nil {
		log.Fatalf("gather websocket metrics: %v", err)
	}
	all = append(all, wsDefs...)

	// (prober.MetricRegistry already appended above)

	all = deduplicateByName(all)
	sortMetrics(all)

	data := TemplateData{
		Metrics:   all,
		ColWidths: computeColWidths(all),
	}

	content, err := render.RenderTemplate(templateFS, "metrics.gotpl", data)
	if err != nil {
		log.Fatalf("render template: %v", err)
	}

	const outputPath = "docs/monitoring/metrics.md"
	if err := render.WriteToFile(outputPath, content); err != nil {
		log.Fatalf("write %s: %v", outputPath, err)
	}

	log.Printf("wrote %s", outputPath)
}

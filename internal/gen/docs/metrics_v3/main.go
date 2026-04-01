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

// main generates a metrics documentation table by intercepting all
// prometheus.Collector registrations via capturingRegistry and calling
// Describe() on each — extracting label names without needing observed values.
//
// Run via: go run ./internal/gen/docs/metrics_v3
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/blackbox_exporter/prober"
)

const outputFile = "internal/gen/docs/metrics_v3/metrics.md"

// descInfo holds what we learn from a single prometheus.Desc.
type descInfo struct {
	fqName         string
	help           string
	variableLabels []string
}

// descPattern matches the output of prometheus.Desc.String():
// Desc{fqName: "name", help: "help", constLabels: {}, variableLabels: {l1,l2}}
var descPattern = regexp.MustCompile(
	`fqName: "([^"]+)".*?help: "([^"]*)".*?variableLabels: \{([^}]*)\}`,
)

func parseDesc(d *prometheus.Desc) (fqName, help string, labels []string) {
	s := d.String()
	m := descPattern.FindStringSubmatch(s)
	if m == nil {
		return "", "", nil
	}
	fqName = m[1]
	help = m[2]
	raw := strings.TrimSpace(m[3])
	if raw == "" {
		return fqName, help, nil
	}
	for _, lbl := range strings.Split(raw, ",") {
		lbl = strings.TrimSpace(lbl)
		if lbl != "" {
			labels = append(labels, lbl)
		}
	}
	return fqName, help, labels
}

// capturingRegistry wraps *prometheus.Registry and intercepts every Register /
// MustRegister call to call Describe() on the collector, extracting label names
// directly from Desc without needing observed metric values.
type capturingRegistry struct {
	real     *prometheus.Registry
	captured map[string]*descInfo // fqName -> descInfo
}

func newCapturingRegistry() *capturingRegistry {
	return &capturingRegistry{
		real:     prometheus.NewRegistry(),
		captured: make(map[string]*descInfo),
	}
}

func (cr *capturingRegistry) Register(c prometheus.Collector) error {
	cr.intercept(c)
	return cr.real.Register(c)
}

func (cr *capturingRegistry) MustRegister(cs ...prometheus.Collector) {
	for _, c := range cs {
		cr.intercept(c)
	}
	cr.real.MustRegister(cs...)
}

func (cr *capturingRegistry) Unregister(c prometheus.Collector) bool {
	return cr.real.Unregister(c)
}

func (cr *capturingRegistry) intercept(c prometheus.Collector) {
	ch := make(chan *prometheus.Desc, 64)
	go func() {
		c.Describe(ch)
		close(ch)
	}()
	for d := range ch {
		fqName, help, labels := parseDesc(d)
		if fqName == "" {
			continue
		}
		if _, exists := cr.captured[fqName]; !exists {
			cr.captured[fqName] = &descInfo{
				fqName:         fqName,
				help:           help,
				variableLabels: labels,
			}
		}
	}
}

// MetricDef is the doc-generation representation of a single metric.
type MetricDef struct {
	Name   string
	Type   string
	Prober string
	Labels []string
	Help   string
}

// proberOverrides maps metric names whose Prober value cannot be inferred from
// a single prober call. Shared metrics registered by chooseProtocol appear
// under whichever prober runs first; override them here.
var proberOverrides = map[string]string{
	"probe_dns_lookup_time_seconds": "dns, grpc, http, icmp, tcp, unix, websocket",
	"probe_ip_addr_hash":            "all",
	"probe_ip_protocol":             "all",
	"probe_failed_due_to_regex":     "http, tcp, unix",
}

// discardLogger suppresses prober diagnostic output during doc generation.
var discardLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// describeFromProber runs construct with a fresh capturingRegistry, then uses
// Gather() on the underlying registry to fill in metric types. Returns
// MetricDefs with defaultProber unless proberOverrides overrides it.
func describeFromProber(defaultProber string, construct func(prometheus.Registerer)) []MetricDef {
	cr := newCapturingRegistry()
	construct(cr)

	families, _ := cr.real.Gather() // errors expected on dummy targets; partial results are fine

	typeByName := make(map[string]string, len(families))
	for _, f := range families {
		typeByName[f.GetName()] = strings.ToLower(f.GetType().String())
	}

	defs := make([]MetricDef, 0, len(cr.captured))
	for name, info := range cr.captured {
		p := defaultProber
		if override, ok := proberOverrides[name]; ok {
			p = override
		}

		metricType := typeByName[name]
		switch {
		case metricType == "" && len(info.variableLabels) > 0:
			metricType = "gaugevec"
		case metricType == "":
			metricType = "gauge"
		case metricType == "gauge" && len(info.variableLabels) > 0:
			metricType = "gaugevec"
		}

		defs = append(defs, MetricDef{
			Name:   name,
			Type:   metricType,
			Prober: p,
			Labels: info.variableLabels,
			Help:   info.help,
		})
	}
	return defs
}

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

// generate collects all metric definitions and renders them as a markdown table.
func generate() string {
	var all []MetricDef

	all = append(all, describeFromProber("-", func(reg prometheus.Registerer) {
		config.NewSafeConfig(reg)
	})...)

	probersToRun := []struct {
		name   string
		target string
		fn     func(context.Context, string, config.Module, prometheus.Registerer, *slog.Logger) bool
	}{
		{"dns", "dummy.invalid", prober.ProbeDNS},
		{"grpc", "dummy.invalid", prober.ProbeGRPC},
		{"http", "http://dummy.invalid", prober.ProbeHTTP},
		{"icmp", "dummy.invalid", prober.ProbeICMP},
		{"websocket", "ws://dummy.invalid", prober.ProbeWebsocket},
	}

	for _, p := range probersToRun {
		target := p.target
		fn := p.fn
		all = append(all, describeFromProber(p.name, func(reg prometheus.Registerer) {
			fn(context.Background(), target, config.Module{}, reg, discardLogger)
		})...)
	}

	// Conditional metrics only reachable on live code paths (TLS handshake,
	// CEL config, SOA response, ICMP reply, query-response config) fall back
	// to MetricRegistry for any name not yet captured above.
	seen := make(map[string]bool, len(all))
	for _, m := range all {
		seen[m.Name] = true
	}
	for _, m := range prober.MetricRegistry {
		if !seen[m.Name] {
			all = append(all, MetricDef{
				Name:   m.Name,
				Type:   m.Type,
				Prober: m.Prober,
				Labels: m.Labels,
				Help:   m.Help,
			})
		}
	}

	all = deduplicateByName(all)
	sortMetrics(all)

	return renderTable(all)
}

func renderTable(metrics []MetricDef) string {
	colName := len("Name")
	colType := len("Metric Type")
	colProber := len("Prober")
	colLabels := len("Labels")
	colHelp := len("Help")

	for _, m := range metrics {
		if l := len(m.Name); l > colName {
			colName = l
		}
		if l := len(m.Type); l > colType {
			colType = l
		}
		if l := len(m.Prober); l > colProber {
			colProber = l
		}
		if l := len(strings.Join(m.Labels, ", ")); l > colLabels {
			colLabels = l
		}
		if l := len(m.Help); l > colHelp {
			colHelp = l
		}
	}

	var sb strings.Builder
	row := func(name, typ, prob, labels, help string) {
		fmt.Fprintf(&sb, "| %-*s | %-*s | %-*s | %-*s | %-*s |\n",
			colName, name,
			colType, typ,
			colProber, prob,
			colLabels, labels,
			colHelp, help,
		)
	}
	sep := func(w int) string { return strings.Repeat("-", w+1) }

	row("Name", "Metric Type", "Prober", "Labels", "Help")
	fmt.Fprintf(&sb, "|:%s|:%s|:%s|:%s|:%s|\n",
		sep(colName), sep(colType), sep(colProber), sep(colLabels), sep(colHelp))

	for _, m := range metrics {
		row(m.Name, m.Type, m.Prober, strings.Join(m.Labels, ", "), m.Help)
	}
	return sb.String()
}

func main() {
	content := generate()
	fmt.Print(content)
	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", outputFile, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", outputFile)
}

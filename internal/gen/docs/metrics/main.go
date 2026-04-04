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

// main generates docs/monitoring/metrics.md from prober.MetricRegistry.
// No network calls, no dummy probes — all metadata is sourced from the
// *_metrics.go files in the prober package.
//
// Run via: go run ./internal/gen/docs/metrics
// Or via:  make generate-metrics-documentation
package main

import (
	"embed"
	"log"
	"sort"
	"strings"

	"github.com/prometheus/blackbox_exporter/internal/gen/docs/render"
	"github.com/prometheus/blackbox_exporter/prober"
)

//go:embed templates
var templateFS embed.FS

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
	Metrics   []prober.MetricDef
	ColWidths ColWidths
}

func sortMetrics(m []prober.MetricDef) {
	sort.Slice(m, func(i, j int) bool {
		if m[i].Prober != m[j].Prober {
			return m[i].Prober < m[j].Prober
		}
		return m[i].Name < m[j].Name
	})
}

func computeColWidths(m []prober.MetricDef) ColWidths {
	return ColWidths{
		Name:   render.MapColumn("Name", m, func(x prober.MetricDef) string { return x.Name }),
		Type:   render.MapColumn("Metric Type", m, func(x prober.MetricDef) string { return x.Type }),
		Prober: render.MapColumn("Prober", m, func(x prober.MetricDef) string { return x.Prober }),
		Labels: render.MapColumn("Labels", m, func(x prober.MetricDef) string { return strings.Join(x.Labels, ", ") }),
		Help:   render.MapColumn("Help", m, func(x prober.MetricDef) string { return x.Help }),
	}
}

func main() {
	metrics := make([]prober.MetricDef, len(prober.MetricRegistry))
	copy(metrics, prober.MetricRegistry)
	sortMetrics(metrics)

	data := TemplateData{
		Metrics:   metrics,
		ColWidths: computeColWidths(metrics),
	}

	content, err := render.RenderTemplate(templateFS, "metrics.gotpl", data)
	if err != nil {
		log.Fatalf("render template: %v", err)
	}

	const outputPath = "docs/monitoring/metrics.md"
	if err := render.WriteToFile(outputPath, content); err != nil {
		log.Fatalf("write %s: %v", outputPath, err)
	}

	log.Printf("wrote %s (%d metrics)", outputPath, len(metrics))
}
